package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/cache"
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
	GetById(ctx context.Context, id, uid int64) (domain.Article, error)
	GetPubById(ctx context.Context, id int64) (domain.Article, error)
}

type articleRepository struct {
	artDao   article.ArticleDAO
	userRepo UserRepository
	cache    cache.ArticleCache
	log      logger.Logger
}

func NewArticleRepository(d article.ArticleDAO, uRepo UserRepository, c cache.ArticleCache, l logger.Logger) ArticleRepository {
	return &articleRepository{
		artDao:   d,
		userRepo: uRepo,
		cache:    c,
		log:      l,
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
	go func() {
		a.saveArticleCache(ctx, art)
	}()
	return id, err
}

func (a *articleRepository) Update(ctx context.Context, art domain.Article) error {
	err := a.artDao.UpdateById(ctx, a.toEntity(art))
	if err != nil {
		return err
	}
	a.clearCache(ctx, art.Id, art.Author.Id)
	go func() {
		a.saveArticleCache(ctx, art)
	}()
	return nil
}

func (a *articleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	id, err := a.artDao.Sync(ctx, a.toEntity(art))
	if err != nil {
		return 0, err
	}
	a.clearCache(ctx, id, art.Author.Id)
	go func() {
		a.saveArticleCache(ctx, art)

		user, err := a.userRepo.FindById(ctx, art.Author.Id)
		if err != nil {
			a.log.Error("获取用户信息失败",
				logger.Int64("author", art.Author.Id), logger.Error(err))
		}
		art.Author.Name = user.NickName
		a.savePublishedArticleCache(ctx, art)
	}()
	return id, err
}

func (a *articleRepository) SyncStatus(ctx context.Context, id, usrId int64, status domain.ArticleStatus) error {
	err := a.artDao.SyncStatus(ctx, id, usrId, status.ToUint8())
	if err != nil {
		return err
	}
	a.clearCache(ctx, id, usrId)
	return err
}

func (a *articleRepository) List(ctx context.Context, id int64, offset, limit int) ([]domain.Article, error) {
	if offset+limit <= 100 {
		arts, err := a.cache.GetFirstPage(ctx, id)
		if err == nil {
			go func() {
				a.preCache(ctx, arts)
			}()
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

	data := slice.Map[article.Article, domain.Article](res, func(idx int, src article.Article) domain.Article {
		return a.toDomain(src)
	})
	go func() {
		if offset == 0 && limit >= 100 {
			err = a.cache.SetFirstPage(ctx, id, data[:100])
			if err != nil {
				a.log.Error("refresh first page article failed",
					logger.Int64("author", id), logger.Error(err))
			}
		}

		a.preCache(ctx, data)
	}()
	return data, nil
}

func (a *articleRepository) GetById(ctx context.Context, id, uid int64) (domain.Article, error) {
	art, err := a.cache.Get(ctx, id, uid)
	if err == nil {
		return art, err
	}
	dArt, err := a.artDao.FindById(ctx, id, uid)
	if err != nil {
		return domain.Article{}, err
	}

	art = a.toDomain(dArt)
	go func() {
		a.saveArticleCache(ctx, art)
	}()

	return art, nil
}

func (a *articleRepository) GetPubById(ctx context.Context, id int64) (domain.Article, error) {
	// 从缓存中获取数据
	art, err := a.cache.GetPub(ctx, id)
	if err == nil {
		return art, nil
	}

	// 获取文章数据
	pArt, err := a.artDao.FindPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	// 获取作者数据
	user, err := a.userRepo.FindById(ctx, pArt.AuthorId)
	if err != nil {
		return domain.Article{}, err
	}

	art = a.toDomain(article.Article(pArt))
	art.Author.Name = user.NickName

	// 缓存数据
	go func() {
		a.savePublishedArticleCache(ctx, art)
	}()
	return art, nil
}

func (a *articleRepository) preCache(ctx context.Context, arts []domain.Article) {
	if len(arts) > 0 && a.needCache(arts[0]) {
		err := a.cache.Set(context.Background(), arts[0])
		if err != nil {
			a.log.Error("提前预缓存失败", logger.Error(err))
		}
	}
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

func (a *articleRepository) clearCache(ctx context.Context, id, uid int64) {
	err := a.cache.DeleteFirstPage(ctx, uid)
	if err != nil && err != cache.ErrKeyNotExisted {
		a.log.Error("清除第一页缓存失败",
			logger.Int64("author", uid), logger.Error(err))
	}

	err = a.cache.Delete(ctx, id, uid)
	if err != nil && err != cache.ErrKeyNotExisted {
		a.log.Error("清除文章缓存失败",
			logger.Int64("id", id), logger.Int64("author", uid),
			logger.Error(err))
	}

	err = a.cache.DeletePub(ctx, id)
	if err != nil && err != cache.ErrKeyNotExisted {
		a.log.Error("清除发表文章缓存失败",
			logger.Int64("id", id),
			logger.Error(err))
	}
}

func (a *articleRepository) needCache(art domain.Article) bool {
	const CacheDataThreshold = 1024 * 1024
	return len(art.Content) < CacheDataThreshold
}

func (a *articleRepository) saveArticleCache(ctx context.Context, art domain.Article) {
	if !a.needCache(art) {
		return
	}
	err := a.cache.Set(ctx, art)
	if err != nil {
		a.log.Error("缓存制作文章失败",
			logger.Int64("id", art.Id), logger.Error(err))
	}
}

func (a *articleRepository) savePublishedArticleCache(ctx context.Context, art domain.Article) {
	if !a.needCache(art) {
		return
	}

	err := a.cache.SetPub(ctx, art)
	if err != nil {
		a.log.Error("缓存制作文章失败",
			logger.Int64("id", art.Id), logger.Error(err))
	}
}
