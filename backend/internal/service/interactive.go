package service

import (
	"context"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
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

func (i *interactiveService) Get(
	ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	return i.r.Get(ctx, biz, bizId)
}
