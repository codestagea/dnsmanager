package store

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"github.com/codestagea/bindmgr/internal/model"
	"github.com/codestagea/bindmgr/internal/tools/diff"
	"strings"
)

type DnsZone struct {
	ID         int64  `json:"id" gorm:"primary_key" example:"1"`
	Zone       string `json:"zone" diff:"Zone"`
	Refresh    int    `json:"refresh" diff:"Refresh"`
	Retry      int    `json:"retry" diff:"Retry"`
	Expire     int    `json:"expire" diff:"Expire"`
	Minimum    int    `json:"minimum" diff:"Minimum"`
	Serial     int64  `json:"serial" diff:"Serial"`
	HostMaster string `json:"hostMaster" diff:"hostMaster"`
	PrimaryNs  string `json:"primaryNs" diff:"PrimaryNs"`
	State      string `json:"state" diff:"状态"`
	Remark     string `json:"remark" diff:"描述"`
}

func (DnsZone) TableName() string {
	return "dns_zone"
}
func (z *DnsZone) Validate() error {
	errArr := make([]string, 0)
	if z.Zone == "" {
		errArr = append(errArr, "zone name should not be empty")
	}
	if z.Refresh <= 60 {
		errArr = append(errArr, "refresh should longer than 60 seconds")
	}
	if z.Retry <= 60 {
		errArr = append(errArr, "retry should longer than 60 seconds")
	}

	if z.Expire <= 60 {
		errArr = append(errArr, "expire should longer than 60 seconds")
	}

	if z.Minimum <= 60 {
		errArr = append(errArr, "minimum should longer than 60 seconds")
	}

	hostMaster := trimSuffixDot(z.HostMaster)
	if hostMaster == "" {
		errArr = append(errArr, "host master email should not be empty")
	}
	z.HostMaster = hostMaster + "."

	primaryNs := trimSuffixDot(z.PrimaryNs)
	if primaryNs == "" {
		errArr = append(errArr, "primary ns should not be empty")
	}
	z.PrimaryNs = primaryNs + "."

	if z.State == "" {
		z.State = S_Running
	}
	if z.State != S_Running && z.State != S_Stopped {
		errArr = append(errArr, fmt.Sprintf("state should be %s", strings.Join([]string{S_Running, S_Stopped}, ", ")))
	}

	if len(errArr) > 0 {
		return errors.New(strings.Join(errArr, "\n"))
	} else {
		return nil
	}
}

type DnsZoneStore interface {
	ListPage(search string, page *model.PageQuery) ([]DnsZone, int64, error)
	GetById(id int64) (*DnsZone, error)
	GetByName(name string) (*DnsZone, error)
	AddZone(d *DnsZone, operator string) error
	UpdateById(d *DnsZone, operator string) error
}

type dnsZoneStore struct {
	db *gorm.DB
}

var _ DnsZoneStore = &dnsZoneStore{}

func (s *dnsZoneStore) ListPage(search string, page *model.PageQuery) ([]DnsZone, int64, error) {
	var zones []DnsZone
	var total int64
	dbQuery := s.db.Table(DnsZone{}.TableName())
	if search != "" {
		dbQuery = dbQuery.Where("name like ?", "%"+search+"%").Or("code like ?", "%"+search+"%")
	}
	if err := dbQuery.Count(&total).Error; err != nil {
		logrus.Errorf("count zone page fail: %v", err)
		return zones, total, err
	}
	if err := dbQuery.Limit(page.PageSize).Offset(page.Offset).Find(&zones).Error; err != nil {
		logrus.Errorf("list zones page fail: %v", err)
		return zones, total, err
	}

	return zones, total, nil
}

func (s *dnsZoneStore) GetById(id int64) (*DnsZone, error) {
	zone := DnsZone{}
	err := s.db.Table(zone.TableName()).Where("id=?", id).Take(&zone).Error
	return &zone, err
}

func (s *dnsZoneStore) GetByName(name string) (*DnsZone, error) {
	zone := DnsZone{}
	err := s.db.Table(zone.TableName()).Where("zone=?", name).Take(&zone).Error
	return &zone, err
}

func (s *dnsZoneStore) AddZone(d *DnsZone, operator string) error {
	var exist DnsZone

	err := s.db.Table(d.TableName()).Where("zone=?", d.Zone).Take(&exist).Error

	if err == nil {
		return errors.New("zone already exist")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("get zone %s fail: %v", d.Zone, err)
		return errors.New("get zone by name fail")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		r := tx.Table(d.TableName()).Create(&d)
		if r.Error != nil {
			return r.Error
		}

		diffValues := diff.DiffStructStr(DnsZone{}, *d)
		log := OperationLog{
			Operator:   operator,
			TargetType: "zone",
			TargetId:   d.ID,
			Type:       "add",
			KeyValue:   d.Zone,
			Diff:       truncateStr(strings.Join(diffValues, "\n"), 512),
		}
		return tx.Table(log.TableName()).Create(&log).Error
	})
}

func (s *dnsZoneStore) UpdateById(d *DnsZone, operator string) error {
	var exist DnsZone
	if err := s.db.Table(exist.TableName()).Where("id=?", d.ID).Take(&exist).Error; err != nil {
		return err
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		//不允许修改名字
		r := tx.Table(d.TableName()).Model(&d).Where("id = ?", d.ID).Updates(map[string]interface{}{
			"refresh":     d.Refresh,
			"retry":       d.Retry,
			"expire":      d.Expire,
			"minimum":     d.Minimum,
			"serial":      gorm.Expr("serial + ?", 1),
			"host_master": d.HostMaster,
			"primary_ns":  d.PrimaryNs,
			"remark":      d.Remark,
			"state":       d.State,
		})
		if r.Error != nil {
			return r.Error
		}

		diffValues := diff.DiffStructStr(exist, *d)
		log := OperationLog{
			Operator:   operator,
			TargetType: "zone",
			TargetId:   d.ID,
			Type:       "update",
			KeyValue:   d.Zone,
			Diff:       truncateStr(strings.Join(diffValues, "\n"), 512),
		}
		return tx.Table(log.TableName()).Create(&log).Error
	})
}
