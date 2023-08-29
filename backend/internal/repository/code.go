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

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string, experation time.Duration) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type CachedCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return &CachedCodeRepository{
		cache: c,
	}
}

func (r *CachedCodeRepository) Store(ctx context.Context, biz, phone, code string, experation time.Duration) error {
	return r.cache.Set(ctx, biz, phone, code, experation)
}

func (r *CachedCodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return r.cache.Verify(ctx, biz, phone, code)
}
