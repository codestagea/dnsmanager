package store

import (
	"github.com/codestagea/bindmgr/internal/model"
	"gorm.io/gorm"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Model struct {
	ID        int64     `json:"id" gorm:"primary_key" example:"1"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:current_timestamp" example:"2018-10-21T16:40:23+08:00"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"default:current_timestamp on update current_timestamp" example:"2018-10-21T16:40:23+08:00"`
}

func Paginate(p *model.PageQuery) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(p.PageSize).Offset(p.Offset).Order("id desc")
	}
}

func truncateStr(s string, size int) string {
	if len(s) < size {
		return s
	}
	return s[0 : size-1]
}

func trimSuffixDot(v string) string {
	ii := len(v) - 1
	for ; ii > 0; ii-- {
		if v[ii] != '.' {
			break
		}
	}
	if ii == 0 {
		return ""
	} else {
		return v[0 : ii+1]
	}
}

const (
	S_Running = "running"
	S_Stopped = "stopped"
)
