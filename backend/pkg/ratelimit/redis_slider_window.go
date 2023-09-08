package ratelimit

import (
	"context"
	_ "embed"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisSliderWindowLimiter struct {
	cmd      redis.Cmdable
	interval time.Duration
	// 阈值
	rate int
}

//go:embed slide_window.lua
var luaRedisSliderWindow string

func NewRedisSliderWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) Limiter {
	return &RedisSliderWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (b *RedisSliderWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	return b.cmd.Eval(ctx, luaRedisSliderWindow, []string{key},
		b.interval.Milliseconds(), b.rate, time.Now().UnixMilli()).Bool()
}
