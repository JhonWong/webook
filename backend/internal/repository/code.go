package repository

import (
	"context"
	"errors"
	"github.com/JhonWong/webook/backend/internal/repository/cache"
)

var (
	ErrSendTooMuch = errors.New("验证码发送过于频繁")
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: c,
	}
}

func (r *CodeRepository) Set(ctx context.Context, biz, code, phone string) error {

}

func (r *CodeRepository) Get(ctx context.Context, biz, phone string) (string, error) {

}
