package article

import (
	"context"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	Upsert(ctx context.Context, art PublishArticle) error
	SyncStatus(ctx context.Context, id, usrId int64, status uint8) error
	GetByAuthor(ctx context.Context, id int64, offset int, limit int) ([]Article, error)
	FindById(ctx context.Context, id, uid int64) (Article, error)
	FindPubById(ctx context.Context, id int64) (PublishArticle, error)
}
