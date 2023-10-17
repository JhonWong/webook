package service

import (
	"context"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository"
	"github.com/johnwongx/webook/backend/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
}

type articleService struct {
	r      repository.ArticleRepository
	logger logger.Logger
}

func NewArticleService(r repository.ArticleRepository, logger logger.Logger) ArticleService {
	return &articleService{
		r:      r,
		logger: logger,
	}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.Id > 0 {
		err := a.r.Update(ctx, art)
		return art.Id, err
	}
	return a.r.Create(ctx, art)
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	return a.r.Sync(ctx, art)
}
