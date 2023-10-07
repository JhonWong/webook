package service

import (
	"context"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
}

type articleService struct {
	r repository.ArticleRepository
}

func NewArticleService(r repository.ArticleRepository) ArticleService {
	return &articleService{
		r: r,
	}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	return a.r.Create(ctx, art)
}
