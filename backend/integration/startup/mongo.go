package startup

import (
	"context"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var mdb *mongo.Database

func InitTestMongoDB() *mongo.Database {
	if mdb == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		monitor := &event.CommandMonitor{
			Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
				//fmt.Println(startedEvent.Command)
			},
		}

		dsn := "mongodb://root:example@localhost:27017"
		opts := options.Client().ApplyURI(dsn).
			SetMonitor(monitor)
		client, err := mongo.Connect(ctx, opts)
		if err != nil {
			panic(err)
		}
		mdb = client.Database("webook")
	}
	return mdb
}
