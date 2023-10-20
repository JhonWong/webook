package article

import (
	"context"
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
	filter := bson.M{
		"id":art.Id,
		"author_id":art.AuthorId,
	}
	update := bson.M{
		"title":
	}
	m.col.UpdateOne(ctx, filter,)
}

func (m *MongoDBArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBArticleDAO) Upsert(ctx context.Context, art PublishArticle) error {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBArticleDAO) SyncStatus(ctx context.Context, id, usrId int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}
