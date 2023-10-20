package article

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var _ ArticleDAO = &MongoDBArticleDAO{}

type MongoDBArticleDAO struct {
	mdb     *mongo.Database
	col     *mongo.Collection
	liveCol *mongo.Collection
}

func InitCollections(mdb *mongo.Database) error {

}

func NewMongoArticleDAO(mdb *mongo.Database) *MongoDBArticleDAO {
	return &MongoDBArticleDAO{
		mdb:     mdb,
		col:     mdb.Collection("articles"),
		liveCol: mdb.Collection("published_articles"),
	}
}

func (m *MongoDBArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	res, err := m.col.InsertOne(ctx, &art)
}

func (m *MongoDBArticleDAO) UpdateById(ctx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
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
