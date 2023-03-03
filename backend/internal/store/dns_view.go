package store

import (
	"errors"
	"gorm.io/gorm"
	"github.com/codestagea/bindmgr/internal/tools/diff"
	"strings"
	"time"
)

type DnsView struct {
	Model
	Name   string `json:"name" gorm:"name" diff:"Name"`
	Remark string `json:"remark" diff:"描述"`
}

func (DnsView) TableName() string {
	return "dns_view"
}

type ViewStore interface {
	List() ([]DnsView, error)
	Add(view *DnsView, operator string) error
	Update(view *DnsView, operator string) error
}

type viewStore struct {
	db *gorm.DB
}

func (s *viewStore) List() ([]DnsView, error) {
	var views []DnsView
	err := s.db.Table(DnsView{}.TableName()).Find(&views).Error
	return views, err
}

func (s *viewStore) Add(view *DnsView, operator string) error {
	view.ID = 0
	var exist DnsView
	err := s.db.Table(view.TableName()).Where("name=?", view.Name).Take(&exist).Error

	if err == nil {
		return errors.New("view already exist")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		r := tx.Table(view.TableName()).Save(view)
		if r.Error != nil {
			return r.Error
		}

		diffValues := diff.DiffStructStr(DnsView{}, *view)
		log := OperationLog{
			Operator:   operator,
			TargetType: "view",
			TargetId:   exist.ID,
			Type:       "add",
			KeyValue:   view.Name,
			Diff:       truncateStr(strings.Join(diffValues, "\n"), 512),
		}
		return tx.Table(log.TableName()).Create(&log).Error
	})
}

func (s *viewStore) Update(view *DnsView, operator string) error {
	var exist DnsView
	if err := s.db.Table(view.TableName()).Where("id=?", view.ID).Take(&exist).Error; err != nil {
		return err
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		r := tx.Table(view.TableName()).Where("id = ?", view.ID).Updates(map[string]interface{}{
			"name":       view.Name,
			"remark":     view.Remark,
			"updated_at": time.Now(),
		})
		if r.Error != nil {
			return r.Error
		}

		diffValues := diff.DiffStructStr(exist, *view)
		log := OperationLog{
			Operator:   operator,
			TargetType: "view",
			TargetId:   exist.ID,
			Type:       "update",
			KeyValue:   view.Name,
			Diff:       truncateStr(strings.Join(diffValues, "\n"), 512),
		}
		return tx.Table(log.TableName()).Create(&log).Error
	})
}
