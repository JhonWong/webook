package cache

import (
	"context"
	_ "embed"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewRedisCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		client: client,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string, experation time.Duration) error {
	//TODO:将过期时间设置到lua脚本中
	res, err := c.client.Eval(ctx, luaSetCode, []string{key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}

	switch res {
	case 0:
		return nil
	case -1:
		return ErrCodeSendTooMany
	default:
		return errors.New("系统错误")
	}
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	//TODO:将过期时间设置到lua脚本中
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{key(biz, phone)}, code).Int()
	if err != nil {
		return false, err
	}

	switch res {
	case 0:
		return true, nil
	case -1:
		return false, ErrCodeVerifyTooMany
	case -2:
		return false, nil
	}
	return false, ErrCodeUnknowError
}
