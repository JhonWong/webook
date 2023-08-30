package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var (
	ErrCodeSendTooMany   = errors.New("验证码发送太频繁")
	ErrCodeVerifyTooMany = errors.New("验证次数太多")
	ErrCodeUnknowError   = errors.New("未知错误")
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string, experation time.Duration) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return &LocalCodeCache{
		data: make(map[string]codeInfo),
		lock: &sync.Mutex{},
	}
}

func NewCodeCacheV1(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		client: client,
	}
}

func key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
