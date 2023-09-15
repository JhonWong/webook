package repository

import (
	"context"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/cache"
)

type SMSRepository interface {
	Put(ctx context.Context, info domain.SMSInfo) error
	Get(ctx context.Context, cnt int) ([]domain.SMSInfo, error)
	IsEmpty(ctx context.Context) bool
}

type smsRepository struct {
	c cache.SMSCache
}

func NewSMSRepository(c cache.SMSCache) SMSRepository {
	return &smsRepository{
		c: c,
	}
}

func (s *smsRepository) Put(ctx context.Context, info domain.SMSInfo) error {
	return s.c.Add(ctx, cache.SMSInfo{
		Tpl:     info.Tpl,
		Args:    info.Args,
		Numbers: info.Numbers,
	})
}

func (s *smsRepository) Get(ctx context.Context, cnt int) ([]domain.SMSInfo, error) {
	infos, err := s.c.Take(ctx, cnt)
	if err != nil {
		return nil, err
	}

	res := make([]domain.SMSInfo, len(infos))
	for i, info := range infos {
		res[i] = domain.SMSInfo{
			Tpl:     info.Tpl,
			Args:    info.Args,
			Numbers: info.Numbers,
		}
	}
	return res, nil
}

func (s *smsRepository) IsEmpty(ctx context.Context) bool {
	res, _ := s.c.KeyExists(ctx)
	return res
}
