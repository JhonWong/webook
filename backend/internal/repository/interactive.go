package repository

import (
	"context"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
}

type interactiveRepository struct {
	d dao.InteractiveDAO
}

func NewInteractiveRepository(d dao.InteractiveDAO) InteractiveRepository {
	return &interactiveRepository{
		d: d,
	}
}

func (i *interactiveRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return i.d.IncrReadCnt(ctx, biz, bizId)
}

func (i *interactiveRepository) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	data, err := i.d.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}
	return i.toDomain(data), nil
}

func (i *interactiveRepository) toDomain(data dao.Interactive) domain.Interactive {
	return domain.Interactive{
		BizId:   data.BizId,
		Biz:     data.Biz,
		ReadCnt: data.ReadCnt,
	}
}
