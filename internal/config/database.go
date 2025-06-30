package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBS stores the database connections.
var DBS = map[string]*gorm.DB{}

// Config is the global configuration variable.
const (
	MysqlDriver string = "mysql"
)

// DefaultDBName is the default database name.
const (
	DefaultDBName string = "default"
)

// InitDB initializes the database connections based on the configuration.
func InitDB() {
	var dsn string
	for _, v := range Config.DBS {
		if v.ConnectionName == "" {
			v.ConnectionName = DefaultDBName
		}
		if v.Driver == MysqlDriver {
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&loc=%s&parseTime=true",
				v.Username,
				v.Password,
				v.Host,
				v.Port,
				v.Database,
				v.Charset,
				url.QueryEscape("Local"),
			)
		}
		var gormConfig *gorm.Config

		if v.Debug {
			newLogger := logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold: time.Second, // 慢 SQL 阈值
					LogLevel:      logger.Info, // 日志等级
					Colorful:      true,        // 彩色打印
				},
			)
			gormConfig = &gorm.Config{
				SkipDefaultTransaction:                   true,
				PrepareStmt:                              true,
				DisableForeignKeyConstraintWhenMigrating: true,
				Logger:                                   newLogger,
			}
		} else {
			gormConfig = &gorm.Config{}
		}
		db, err := gorm.Open(mysql.Open(dsn), gormConfig)
		DBS[v.ConnectionName] = db
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to database %s: %v", v.ConnectionName, err))
		}
	}

}

func GetDB(name ...string) *gorm.DB {
	var conName string
	if len(name) > 0 {
		conName = name[0]
	} else {
		conName = DefaultDBName
	}
	if db, ok := DBS[conName]; ok {
		return db
	}
	if len(DBS) > 0 {
		for _, db := range DBS {
			return db
		}
	}
	panic("No database connection found")
}
