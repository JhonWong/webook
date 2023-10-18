package repository

import (
	"context"
	"errors"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
)

var _ SMSAsyncRepository = &smsAsyncRepository{}

var (
	ErrorEmptySmsInfo = errors.New("No sms info")
)

type SMSAsyncRepository interface {
	Store(ctx context.Context, info domain.SMSAsyncInfo) error
	Load(ctx context.Context) (domain.SMSAsyncInfo, error)
	UpdateResult(ctx context.Context, id int64, res bool) error
}

type smsAsyncRepository struct {
	d dao.AsyncSMSDAO
}

func NewSmsAsyncRepository(d dao.AsyncSMSDAO) *smsAsyncRepository {
	return &smsAsyncRepository{
		d: d,
	}
}

func (s *smsAsyncRepository) Store(ctx context.Context, info domain.SMSAsyncInfo) error {
	return s.d.Store(ctx, s.toEntity(info))
}

func (s *smsAsyncRepository) Load(ctx context.Context) (domain.SMSAsyncInfo, error) {
	info, err := s.d.Load(ctx)
	if err != nil {
		return domain.SMSAsyncInfo{}, err
	}
	// TODO
	return domain.SMSAsyncInfo{}, nil
}

func (s *smsAsyncRepository) UpdateResult(ctx context.Context, id int64, res bool) error {
	return s.d.UpdateResult(ctx, id, res)
}

func (s *smsAsyncRepository) toEntity(info domain.SMSAsyncInfo) dao.SMSAsyncInfo {

}
