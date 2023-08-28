package repository

import (
	"context"
	"time"

	"github.com/JhonWong/webook/backend/internal/repository/cache"
)

var (
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: c,
	}
}

func (r *CodeRepository) Store(ctx context.Context, biz, phone, code string, experation time.Duration) error {
	return r.cache.Set(ctx, biz, phone, code, experation)
}

func (r *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return r.cache.Verify(ctx, biz, phone, code)
}
