package store

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"github.com/codestagea/bindmgr/internal/model"
	"strings"
)

type OperationLog struct {
	Model
	Operator   string `gorm:"operator" json:"operator"`
	TargetType string `gorm:"target_type" json:"targetType"`
	TargetId   int64  `gorm:"target_id" json:"targetId"`
	Type       string `gorm:"type" json:"type"`
	KeyValue   string `gorm:"key_value" json:"keyValue"`
	Diff       string `gorm:"diff" json:"diff"`
}

type OperationLogQuery struct {
	Search     string `json:"search"`
	TargetType string `json:"targetType"`
	TargetId   int64  `json:"targetId"`
}

func (OperationLog) TableName() string { return "operation_log" }

type OperationLogStore interface {
	ListPage(q *OperationLogQuery, page *model.PageQuery) ([]OperationLog, int64, error)
}

type operationLogStore struct {
	db *gorm.DB
}

func (s operationLogStore) ListPage(q *OperationLogQuery, page *model.PageQuery) ([]OperationLog, int64, error) {
	var operationLogs []OperationLog
	var total int64
	dbQuery := s.db.Table(OperationLog{}.TableName())
	fields := []string{}
	values := []interface{}{}

	if q.Search != "" {
		fields = append(fields, "key_value like ?")
		values = append(values, "%"+q.Search+"%")
	}
	if q.TargetType != "" {
		fields = append(fields, "target_type = ?")
		values = append(values, q.TargetType)
	}

	if q.TargetId != 0 {
		fields = append(fields, "target_id = ?")
		values = append(values, q.TargetId)
	}
	if len(fields) != 0 {
		dbQuery = dbQuery.Where(strings.Join(fields, " and "), values...)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		logrus.Errorf("count operation logs page fail: %v", err)
		return operationLogs, total, err
	}

	if err := dbQuery.Scopes(Paginate(page)).Find(&operationLogs).Error; err != nil {
		logrus.Errorf("list operation logs page fail: %v", err)
		return operationLogs, total, err
	}

	return operationLogs, total, nil
}
