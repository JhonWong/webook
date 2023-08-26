package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type CodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) *CodeCache {
	return &CodeCache{
		client: client,
	}
}

func (c *CodeCache) Set(ctx context.Context, biz, phone, code string) error {

}

func (c *CodeCache) Get(ctx context.Context, biz, phone string) (string, error) {

}
