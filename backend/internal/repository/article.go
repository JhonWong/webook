package repository

import (
	"context"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/dao/article"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id, usrId int64, status domain.ArticleStatus) error
}

type articleRepository struct {
	d article.ArticleDAO
}

func NewArticleRepository(d article.ArticleDAO) ArticleRepository {
	return &articleRepository{
		d: d,
	}
}

func (a *articleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return a.d.Insert(ctx, a.toEntity(art))
}

func (a *articleRepository) Update(ctx context.Context, art domain.Article) error {
	return a.d.UpdateById(ctx, a.toEntity(art))
}

func (a *articleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return a.d.Sync(ctx, a.toEntity(art))
}

func (a *articleRepository) SyncStatus(ctx context.Context, id, usrId int64, status domain.ArticleStatus) error {
	return a.d.SyncStatus(ctx, id, usrId, status.ToUint8())
}

func (a *articleRepository) toEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	}
}
