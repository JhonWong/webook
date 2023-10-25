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

func NewMongoArticleDAO(mdb *mongo.Database, node *snowflake.Node) *MongoDBArticleDAO {
	return &MongoDBArticleDAO{
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
	now := time.Now().UnixMilli()
	filter := bson.M{"id": id, "author_id": usrId}
	update := bson.M{"$set": bson.D{{"status", status}, {"utime", now}}}
	res, err := m.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return fmt.Errorf("update article collection failed, id or author_id err, id:%d, author_id:%d \n", id, usrId)
	}
	res, err = m.liveCol.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return fmt.Errorf("update article collection failed, id or author_id err, id:%d, author_id:%d \n", id, usrId)
	}
	return err
}

func (m *MongoDBArticleDAO) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error) {
	filter := bson.M{"author_id": uid}
	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)).SetSort(bson.M{"utime": -1})
	cur, err := m.liveCol.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var arts []Article
	err = cur.All(ctx, arts)
	if err != nil {
		return nil, err
	}
	return arts, nil
}

func (m *MongoDBArticleDAO) FindById(ctx context.Context, id, uid int64) (Article, error) {
	filter := bson.M{"id": id, "author_id": uid}
	var art Article
	err := m.col.FindOne(ctx, filter).Decode(art)
	if err != nil {
		return Article{}, err
	}
	return art, nil
}
