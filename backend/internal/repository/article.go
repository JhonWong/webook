package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/cache"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
	"github.com/johnwongx/webook/backend/internal/repository/dao/article"
	"github.com/johnwongx/webook/backend/pkg/logger"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id, usrId int64, status domain.ArticleStatus) error
	List(ctx context.Context, id int64, offset, limit int) ([]domain.Article, error)
	GetById(ctx *gin.Context, id, uid int64) (domain.Article, error)
}

type articleRepository struct {
	artDao  article.ArticleDAO
	userDao dao.UserDAO
	cache   cache.ArticleCache
	log     logger.Logger
}

func NewArticleRepository(d article.ArticleDAO, ud dao.UserDAO, c cache.ArticleCache, l logger.Logger) ArticleRepository {
	return &articleRepository{
		artDao:  d,
		userDao: ud,
		cache:   c,
		log:     l,
	}
}

func (a *articleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	id, err := a.artDao.Insert(ctx, a.toEntity(art))
	if err != nil {
		return 0, err
	}
	art.Id = id
	uid := art.Author.Id
	err = a.cache.DeleteFirstPage(ctx, uid)
	if err != nil && err != cache.ErrKeyNotExisted {
		a.log.Error("清除第一页缓存失败",
			logger.Int64("author", uid), logger.Error(err))
	}
	return id, err
}

func (a *articleRepository) Update(ctx context.Context, art domain.Article) error {
	err := a.artDao.UpdateById(ctx, a.toEntity(art))
	if err != nil {
		return err
	}
	uid := art.Author.Id
	err = a.cache.DeleteFirstPage(ctx, uid)
	if err != nil && err != cache.ErrKeyNotExisted {
		a.log.Error("清除第一页缓存失败",
			logger.Int64("author", uid), logger.Error(err))
	}
	return nil
}

func (a *articleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	id, err := a.artDao.Sync(ctx, a.toEntity(art))
	if err != nil {
		return 0, err
	}
	art.Id = id
	uid := art.Author.Id
	err = a.cache.DeleteFirstPage(ctx, uid)
	if err != nil && err != cache.ErrKeyNotExisted {
		a.log.Error("清除第一页缓存失败",
			logger.Int64("author", uid), logger.Error(err))
	}
	return id, err
}

func (a *articleRepository) SyncStatus(ctx context.Context, id, usrId int64, status domain.ArticleStatus) error {
	err := a.artDao.SyncStatus(ctx, id, usrId, status.ToUint8())
	if err != nil {
		return err
	}
	err = a.cache.DeleteFirstPage(ctx, usrId)
	if err != nil && err != cache.ErrKeyNotExisted {
		a.log.Error("清除第一页缓存失败",
			logger.Int64("author", usrId), logger.Error(err))
	}
	return err
}

func (a *articleRepository) List(ctx context.Context, id int64, offset, limit int) ([]domain.Article, error) {
	if offset+limit <= 100 {
		arts, err := a.cache.GetFirstPage(ctx, id)
		if err == nil {
			return arts[offset:limit], nil
		}
		if err != cache.ErrKeyNotExisted {
			a.log.Error("Get author article cache form cache failed", logger.Error(err))
		}
	}

	// 慢路径
	res, err := a.artDao.GetByAuthor(ctx, id, offset, limit)
	if err != nil {
		return []domain.Article{}, err
	}

	dres := slice.Map[article.Article, domain.Article](res, func(idx int, src article.Article) domain.Article {
		return a.toDomain(src)
	})
	if offset == 0 && limit >= 100 {
		err = a.cache.SetFirstPage(ctx, id, dres[:100])
		if err != nil {
			a.log.Error("refresh first page article failed",
				logger.Int64("author", id), logger.Error(err))
		}
	}
	return dres, nil
}

func (a *articleRepository) GetById(ctx *gin.Context, id, uid int64) (domain.Article, error) {
	art, err := a.artDao.FindById(ctx, id, uid)
	if err != nil {
		return domain.Article{}, err
	}
	return a.toDomain(art), nil
}

func (a *articleRepository) toDomain(src article.Article) domain.Article {
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
