package repository

import (
	"context"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/cache"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
	"github.com/johnwongx/webook/backend/pkg/logger"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	IncrLike(ctx context.Context, id int64, biz string, uid int64) error
	DecrLike(ctx context.Context, id int64, biz string, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	AddCollectionItem(ctx context.Context, id int64, biz string, cid, uid int64) error
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
}

type interactiveRepository struct {
	d     dao.InteractiveDAO
	cache cache.InteractiveCache
	l     logger.Logger
}

func NewInteractiveRepository(d dao.InteractiveDAO, cache cache.InteractiveCache, l logger.Logger) InteractiveRepository {
	return &interactiveRepository{
		d:     d,
		cache: cache,
		l:     l,
	}
}

func (i *interactiveRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := i.d.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return i.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}

func (i *interactiveRepository) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	intr, err := i.cache.Get(ctx, biz, bizId)
	if err == nil {
		return intr, nil
	}
	data, err := i.d.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}
	res := i.toDomain(data)
	if er := i.cache.Set(ctx, biz, bizId, res); er != nil {
		i.l.Error("回写缓存失败",
			logger.Int64("bizId", bizId),
			logger.String("biz", biz),
			logger.Error(er))
	}
	return res, nil
}

func (i *interactiveRepository) IncrLike(ctx context.Context, id int64, biz string, uid int64) error {
	err := i.d.IncrLike(ctx, id, biz, uid)
	if err != nil {
		return err
	}
	return i.cache.IncrLikeCntIfPresent(ctx, biz, id)
}

func (i *interactiveRepository) DecrLike(ctx context.Context, id int64, biz string, uid int64) error {
	err := i.d.DecrLike(ctx, id, biz, uid)
	if err != nil {
		return err
	}
	return i.cache.DecrLikeCntIfPresent(ctx, biz, id)
}

func (i *interactiveRepository) AddCollectionItem(ctx context.Context, id int64, biz string, cid, uid int64) error {
	err := i.d.InsertCollectionBiz(ctx, id, biz, cid, uid)
	if err != nil {
		return err
	}
	return i.cache.IncrCollectCntIfPresent(ctx, biz, id)
}

func (i *interactiveRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	info, err := i.d.GetLikeInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return info.Status == 1, nil
	case dao.ErrDataNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (i *interactiveRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := i.d.GetCollectInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrDataNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (i *interactiveRepository) toDomain(data dao.Interactive) domain.Interactive {
	return domain.Interactive{
		BizId:   data.BizId,
		Biz:     data.Biz,
		ReadCnt: data.ReadCnt,
	}
}
