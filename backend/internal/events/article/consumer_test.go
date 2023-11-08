//go:build e2e

package article

import (
	"github.com/IBM/sarama"
	"github.com/johnwongx/webook/backend/internal/repository"
	"github.com/johnwongx/webook/backend/internal/repository/cache"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
	"github.com/johnwongx/webook/backend/pkg/logger"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	client := InitKafka()
	logger := logger.NewNopLogger()
	db := InitDB(logger)
	cmdable := InitRedis()
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	interactiveRepository := repository.NewInteractiveRepository(interactiveDAO, interactiveCache, logger)
	kafkaConsumer := NewKafkaConsumer(client, interactiveRepository, logger)
	kafkaConsumer.Start()
	time.Sleep(time.Minute * 10)
}

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.addr"),
	})
	return redisClient
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}

func InitDB(l logger.Logger) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	c := Config{
		DSN: "root:root@tcp(localhost:13316)/webook",
	}
	//err := viper.UnmarshalKey("gorm", &c)
	//if err != nil {
	//	panic(fmt.Errorf("初始化配置失败 %v, 原因 %w", c, err))
	//}
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

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	cfg := Config{
		Addrs: []string{"localhost:9094"},
	}
	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}
