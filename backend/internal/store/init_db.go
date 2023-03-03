package store

import (
	"github.com/codestagea/bindmgr/config"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var Users UserStore
var Views ViewStore
var OperationLogs OperationLogStore
var DnsRecords DnsRecordStore
var DnsZones DnsZoneStore

func InitDb(cfg config.Database) (*gorm.DB, error) {
	// m := config.GVA_CONFIG.Mysql
	mysqlConfig := mysql.Config{
		DSN: cfg.Dsn, // DSN data source name
	}

	gormCfg := &gorm.Config{}
	if cfg.Debug {
		gormCfg.Logger = gormlogger.New(
			logrus.New(),
			gormlogger.Config{
				SlowThreshold:             time.Second,     // Slow SQL threshold
				LogLevel:                  gormlogger.Info, // Log level
				IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,           // Disable color
			},
		)
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), gormCfg); err != nil {
		return nil, err
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+cfg.Driver)
		Users = &userStore{db: db}
		Views = &viewStore{db: db}
		OperationLogs = &operationLogStore{db: db}
		DnsRecords = &dnsRecordStore{db: db}
		DnsZones = &dnsZoneStore{db: db}
		return db, nil
	}

}
