package repository

import (
	"context"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
)

type ReaderArticleRepository interface {
	Save(ctx context.Context, art domain.Article) error
}

type readerArticleRepository struct {
	r dao.ReaderArticleDAO
}

func NewReaderArticleRepository(r dao.ReaderArticleDAO) ReaderArticleRepository {
	return &readerArticleRepository{
		r: r,
	}
}

func (r *readerArticleRepository) Save(ctx context.Context, art domain.Article) error {
	return r.r.Upsert(ctx, dao.Article{
		Id:       art.Id,
		Tittle:   art.Tittle,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}
