package service

import (
	"context"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, id int64, biz string, uid int64) error
	Liked(ctx context.Context, id int64, biz string, uid int64) (bool, error)
	CancelLike(ctx context.Context, id int64, biz string, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Collect(ctx context.Context, id int64, biz string, cid int64, uid int64) error
	Collected(ctx context.Context, id int64, biz string, uid int64) (bool, error)
}

type interactiveService struct {
	r repository.InteractiveRepository
}

func NewInteractiveService(r repository.InteractiveRepository) InteractiveService {
	return &interactiveService{
		r: r,
	}
}

func (i *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return i.r.IncrReadCnt(ctx, biz, bizId)
}

func (i *interactiveService) Like(ctx context.Context, id int64, biz string, uid int64) error {
	return i.r.IncrLike(ctx, id, biz, uid)
}

func (i *interactiveService) CancelLike(ctx context.Context, id int64, biz string, uid int64) error {
	return i.r.DecrLike(ctx, id, biz, uid)
}

func (i *interactiveService) Get(
	ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	return i.r.Get(ctx, biz, bizId)
}

func (i *interactiveService) Collect(ctx context.Context, id int64, biz string, cid, uid int64) error {
	return i.r.AddCollectionItem(ctx, id, biz, cid, uid)
}

func (i *interactiveService) Liked(ctx context.Context, id int64, biz string, uid int64) (bool, error) {
	return i.r.Liked(ctx, biz, id, uid)
}

func (i *interactiveService) Collected(ctx context.Context, id int64, biz string, uid int64) (bool, error) {
	return i.r.Collected(ctx, biz, id, uid)
}
