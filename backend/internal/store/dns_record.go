package store

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"github.com/codestagea/bindmgr/internal/model"
	"github.com/codestagea/bindmgr/internal/tools/diff"
	"net"
	"strings"
)

var (
	RT_A     string = "A"
	RT_MX    string = "MX"
	RT_CNAME string = "CNAME"
	RT_NS    string = "NS"
	RT_SOA   string = "SOA"
	RT_PTR   string = "PTR"
	RT_TXT   string = "TXT"
	RT_AAAA  string = "AAAA"
	RT_SVR   string = "SVR"
	RT_URL   string = "URL"

	//RT_SOA 不在记录配置里面
	RECORD_TYPES = []string{RT_A, RT_MX, RT_CNAME, RT_NS, RT_PTR, RT_TXT, RT_AAAA, RT_SVR, RT_URL}
)

type DnsRecord struct {
	ID     int64  `json:"id" gorm:"primary_key" example:"1"`
	ZoneId int64  `gorm:"zone_id" json:"zoneId"`
	Host   string `json:"host" diff:"主机记录"`
	Type   string `gorm:"type" json:"type" diff:"记录类型"`
	Data   string `gorm:"data" json:"data" diff:"记录值"`
	Ttl    int    `gorm:"ttl" json:"ttl" diff:"ttl"`
	Mx     int    `gorm:"mx" json:"mx" diff:"优先级"`
	View   string `gorm:"view" json:"view" diff:"视图"`
	State  string `gorm:"state" json:"state" diff:"状态"`
	Remark string `gorm:"remark" json:"remark" diff:"描述"`
}

func (DnsRecord) TableName() string {
	return "dns_record"
}

func (r *DnsRecord) Validate() error {
	if r.ZoneId <= 0 {
		return errors.New("zone id should be positive number")
	}

	if r.Ttl <= 60 {
		return errors.New("ttl should larger than 60 second")
	}

	if r.Data == "" {
		return errors.New("record data should not null")
	}

	matchType := false
	for _, rt := range RECORD_TYPES {
		if rt == r.Type {
			matchType = true
			break
		}
	}
	if !matchType {
		return fmt.Errorf("record type should be %v", RECORD_TYPES)
	}
	if r.Type == RT_MX {
		if r.Mx <= 0 {
			return errors.New("mx record should have positive mx value")
		}
	} else {
		r.Mx = 0
	}
	if r.Host == "" {
		if r.Type != RT_NS {
			return fmt.Errorf("%s record should have valid host", r.Type)
		}
	} else if r.Host == "." || r.Host[0] == '.' || r.Host[0] == '-' || r.Host[len(r.Host)-1] == '.' || r.Host[len(r.Host)-1] == '-' {
		return errors.New("host should not start/end with . or - mark")
	}
	if r.Type == RT_A && net.ParseIP(r.Data) == nil {
		return errors.New(RT_A + " record's data should be a valid ip")
	}
	if r.Type == RT_CNAME || r.Type == RT_NS {
		r.Data = trimSuffixDot(r.Data)
	}

	//状态
	if r.State == "" {
		r.State = S_Running
	}
	if r.State != S_Running && r.State != S_Stopped {
		return fmt.Errorf("state should be %s", strings.Join([]string{S_Running, S_Stopped}, ", "))

	}
	return nil
}

type DnsRecordStore interface {
	ListPage(zoneId int64, search string, page *model.PageQuery) ([]DnsRecord, int64, error)
	GetById(id int64) (*DnsRecord, error)
	GetByHost(zoneId int64, host, view string) (*DnsRecord, error)
	AddRecord(record *DnsRecord, operator string) error
	UpdateById(record *DnsRecord, operator string) error
}

type dnsRecordStore struct {
	db *gorm.DB
}

func (s *dnsRecordStore) ListPage(zoneId int64, search string, page *model.PageQuery) ([]DnsRecord, int64, error) {
	var r []DnsRecord
	var total int64
	dbQuery := s.db.Table(DnsRecord{}.TableName())
	if search != "" {
		dbQuery = dbQuery.Where("zone_id = ? and `host` like ?", zoneId, "%"+search+"%")
	} else {
		dbQuery = dbQuery.Where("zone_id = ?", zoneId)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		logrus.Errorf("count domain records page fail: %v", err)
		return r, total, err
	}

	if err := dbQuery.Limit(page.PageSize).Offset(page.Offset).Find(&r).Error; err != nil {
		logrus.Errorf("list domain records page fail: %v", err)
		return r, total, err
	}

	return r, total, nil
}

func (s *dnsRecordStore) GetById(id int64) (*DnsRecord, error) {
	domainRecord := DnsRecord{}
	err := s.db.Table(domainRecord.TableName()).Where("id=?", id).Take(&domainRecord).Error
	return &domainRecord, err
}

func (s *dnsRecordStore) GetByHost(zoneId int64, host, view string) (*DnsRecord, error) {
	if view == "" {
		view = "any"
	}
	r := DnsRecord{}
	err := s.db.Table(r.TableName()).Where("zone_id=? and `host`=? and `view`=?", zoneId, host, view).Take(&r).Error
	return &r, err
}

func (s *dnsRecordStore) AddRecord(record *DnsRecord, operator string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		r := tx.Table(record.TableName()).Create(&record)
		if r.Error != nil {
			return r.Error
		}

		diffValues := diff.DiffStructStr(DnsRecord{}, *record)
		log := OperationLog{
			Operator:   operator,
			TargetType: "record",
			TargetId:   record.ID,
			Type:       "add",
			KeyValue:   record.Host,
			Diff:       truncateStr(strings.Join(diffValues, "\n"), 512),
		}
		return tx.Table(log.TableName()).Create(&log).Error
	})
}

func (s *dnsRecordStore) UpdateById(record *DnsRecord, operator string) error {
	var exist DnsRecord

	if err := s.db.Table(record.TableName()).Where("id=?", record.ID).Take(&exist).Error; err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		r := tx.Table(record.TableName()).Where("id = ?", record.ID).Updates(DnsRecord{
			Host:   record.Host,
			Type:   record.Type,
			Data:   record.Data,
			Ttl:    record.Ttl,
			Mx:     record.Mx,
			State:  record.State,
			Remark: record.Remark,
		})
		if r.Error != nil {
			return r.Error
		}

		diffValues := diff.DiffStructStr(exist, *record)
		log := OperationLog{
			Operator:   operator,
			TargetType: "record",
			TargetId:   record.ID,
			Type:       "update",
			KeyValue:   record.Host,
			Diff:       truncateStr(strings.Join(diffValues, "\n"), 512),
		}
		return tx.Table(log.TableName()).Create(&log).Error
	})
}
