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
	//r repository.ArticleRepository
	authorRepo repository.AuthorArticleRepository
	readerRepo repository.ReaderArticleRepository
	logger     logger.Logger
}

func NewArticleService(ar repository.AuthorArticleRepository, rr repository.ReaderArticleRepository,
	logger logger.Logger) ArticleService {
	return &articleService{
		authorRepo: ar,
		readerRepo: rr,
		logger:     logger,
	}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.Id > 0 {
		err := a.authorRepo.Update(ctx, art)
		return art.Id, err
	}
	return a.authorRepo.Create(ctx, art)
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if art.Id > 0 {
		err = a.authorRepo.Update(ctx, art)
	} else {
		id, err = a.authorRepo.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	for i := 0; i < 3; i++ {
		err = a.readerRepo.Save(ctx, art)
		if err == nil {
			break
		}
		a.logger.Error("部分失败：保存数据到线上库失败",
			logger.Field{Key: "art_id", Value: id},
			logger.Error(err))
	}
	if err != nil {
		a.logger.Error("部分失败：保存数据到线上库都失败了",
			logger.Field{Key: "art_id", Value: id},
			logger.Error(err))
		return 0, err
	}
	return id, nil
}
