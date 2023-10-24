package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/dao/article"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id, usrId int64, status domain.ArticleStatus) error
	List(ctx context.Context, id int64, offset, limit int) ([]domain.Article, error)
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

func (a *articleRepository) List(ctx context.Context, id int64, offset, limit int) ([]domain.Article, error) {
	res, err := a.d.GetByAuthor(ctx, id, offset, limit)
	if err != nil {
		return []domain.Article{}, err
	}

	dres := slice.Map[article.Article, domain.Article](res, func(idx int, src article.Article) domain.Article {
		return domain.Article{
			Id:      src.Id,
			Title:   src.Title,
			Content: src.Content,
			Author: domain.Author{
				Id: src.AuthorId,
			},
			Status: domain.ArticleStatus(src.Status),
			Ctime:  time.UnixMilli(src.Ctime),
			Utime:  time.UnixMilli(src.Utime),
		}
	})
	return dres, nil
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
