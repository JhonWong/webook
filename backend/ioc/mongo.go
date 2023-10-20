package ioc

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func InitMongoDB() *mongo.Database {
	type Config struct {
		DSN string `yaml:"dsn"`
		DB  string `yaml:"db"`
	}
	c := Config{
		DSN: "mongodb://root:root@localhost:27017",
	}
	err := viper.UnmarshalKey("db", &c)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败 %v, 原因 %w", c, err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			//fmt.Println(startedEvent.Command)
		},
	}

	opts := options.Client().ApplyURI(c.DSN).
		SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败 %v, 原因 %w", c, err))
	}

	mdb := client.Database("webook")
	return mdb
}
