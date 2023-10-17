package repository

import (
	"context"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
}

type articleRepository struct {
	d dao.ArticleDAO
}

func NewArticleRepository(d dao.ArticleDAO) ArticleRepository {
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

func (a *articleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Tittle:   art.Tittle,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}
