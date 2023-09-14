package repository

import (
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/cache"
)

type SMSRepository interface {
	Put(info domain.SMSInfo) error
	Get(cnt int) ([]domain.SMSInfo, error)
	IsEmpty() bool
}

type smsRepository struct {
	c cache.SMSCache
}

func NewSMSRepository(c cache.SMSCache) SMSRepository {
	return &smsRepository{
		c: c,
	}
}

func (s *smsRepository) Put(info domain.SMSInfo) error {
	//TODO implement me
	panic("implement me")
}

func (s *smsRepository) Get(cnt int) ([]domain.SMSInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (s *smsRepository) IsEmpty() bool {
	//TODO implement me
	panic("implement me")
}
