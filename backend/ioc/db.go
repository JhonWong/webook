package ioc

import (
	"fmt"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
	"github.com/johnwongx/webook/backend/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"time"
)

func InitDB(l logger.Logger) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	c := Config{
		DSN: "root:root@tcp(localhost:3306)/mysql",
	}
	err := viper.UnmarshalKey("gorm", &c)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败 %v, 原因 %w", c, err))
	}
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug),
			glogger.Config{
				SlowThreshold: time.Millisecond * 50,
				LogLevel:      glogger.Info,
			}),
	})
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}
