package mongo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestMongo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			//fmt.Println(startedEvent.Command)
		},
	}

	opts := options.Client().ApplyURI("mongodb://root:example@localhost:27017").
		SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)

	mdb := client.Database("webook")
	col := mdb.Collection("articles")
	defer func() {
		_, err = col.DeleteMany(ctx, bson.M{})
	}()

	// 新增数据
	res, err := col.InsertOne(ctx, Article{
		Id:      123,
		Title:   "my title",
		Content: "my content",
	})
	assert.NoError(t, err)
	fmt.Printf("id: %s \n", res.InsertedID)

	// 更新数据
	filter := bson.M{"id": 123}
	update := bson.M{"$set": bson.M{"title": "new title"}}
	upres, err := col.UpdateOne(ctx, filter, update)
	assert.NoError(t, err)
	fmt.Println("affected", upres.ModifiedCount)

	upres, err = col.UpdateMany(ctx, filter, bson.M{"$set": Article{Title: "2333"}})
	assert.NoError(t, err)

	// 查找数据
	var art Article
	err = col.FindOne(ctx, filter).Decode(&art)
	assert.NoError(t, err)
	fmt.Printf("%#v \n", art)
}

type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"authorId,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}
