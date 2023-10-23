package repository

import (
	"context"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/dao/article"
)

type ReaderArticleRepository interface {
	Save(ctx context.Context, art domain.Article) error
}

type readerArticleRepository struct {
	r article.ReaderArticleDAO
}

func NewReaderArticleRepository(r article.ReaderArticleDAO) ReaderArticleRepository {
	return &readerArticleRepository{
		r: r,
	}
}

func (r *readerArticleRepository) Save(ctx context.Context, art domain.Article) error {
	return r.r.Upsert(ctx, article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}
