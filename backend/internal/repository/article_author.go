package repository

import (
	"context"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
)

type AuthorArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}

type authorArticleRepository struct {
	d dao.AuthorArticleDAO
}

func NewAuthorArticleRepository(d dao.AuthorArticleDAO) AuthorArticleRepository {
	return &authorArticleRepository{
		d: d,
	}
}

func (a *authorArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return a.d.Insert(ctx, dao.Article{
		Tittle:   art.Tittle,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}
func (a *authorArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return a.d.UpdateById(ctx, dao.Article{
		Id:       art.Id,
		Tittle:   art.Tittle,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}
