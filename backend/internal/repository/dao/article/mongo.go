package article

import (
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var _ ArticleDAO = &MongoDBArticleDAO{}

type MongoDBArticleDAO struct {
	client  *mongo.Client
	mdb     *mongo.Database
	col     *mongo.Collection
	liveCol *mongo.Collection
	node    *snowflake.Node
}

func InitCollections(mdb *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{"id", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{"author_id", 1}, {"ctime", 1}},
			Options: options.Index(),
		},
	}

	_, err := mdb.Collection("articles").Indexes().CreateMany(ctx, models)
	if err != nil {
		return err
	}
	_, err = mdb.Collection("published_articles").Indexes().CreateMany(ctx, models)
	return err
}

func NewMongoArticleDAO(client *mongo.Client, node *snowflake.Node) *MongoDBArticleDAO {
	mdb := client.Database("webook")
	return &MongoDBArticleDAO{
		client:  client,
		mdb:     mdb,
		col:     mdb.Collection("articles"),
		liveCol: mdb.Collection("published_articles"),
		node:    node,
	}
}

func (m *MongoDBArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Id = int64(m.node.Generate())
	art.Ctime = now
	art.Utime = now
	_, err := m.col.InsertOne(ctx, &art)
	return art.Id, err
}

func (m *MongoDBArticleDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	filter := bson.M{
		"id":        art.Id,
		"author_id": art.AuthorId,
	}
	update := bson.M{
		"$set": bson.M{
			"title":   art.Title,
			"content": art.Content,
			"status":  art.Status,
			"utime":   art.Utime,
		},
	}
	res, err := m.col.UpdateOne(ctx, filter, update)
	if err == nil && res.ModifiedCount != 1 {
		return fmt.Errorf("更新行数错误，更新了%d行", res.ModifiedCount)
	}
	return err
}

func (m *MongoDBArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)

	// 更新制作库
	if art.Id > 0 {
		err = m.UpdateById(ctx, art)
	} else {
		id, err = m.Insert(ctx, art)
		art.Id = id
	}
	if err != nil {
		return 0, err
	}

	// 更新线上库
	err = m.Upsert(ctx, PublishArticle(art))
	return id, err
}

func (m *MongoDBArticleDAO) Upsert(ctx context.Context, art PublishArticle) error {
	now := time.Now().UnixMilli()
	art.Utime = now

	filter := bson.M{"id": art.Id, "author_id": art.AuthorId}
	update := bson.M{"$set": art, "$setOnInsert": bson.M{"ctime": now}}
	_, err := m.liveCol.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (m *MongoDBArticleDAO) SyncStatus(ctx context.Context, id, usrId int64, status uint8) error {
	session, err := m.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessionContext mongo.SessionContext) (any, error) {
		filter := bson.M{"id": id, "author_id": usrId}
		update := bson.M{"$set": bson.M{"status": status}}
		// 更新制作库
		_, err := m.col.UpdateOne(sessionContext, filter, update)
		if err != nil {
			return nil, err
		}
		// 更新线上库
		_, err = m.liveCol.UpdateOne(sessionContext, filter, update)
		return nil, err
	}
	_, err = session.WithTransaction(ctx, callback)
	return err
}
