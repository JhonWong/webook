package cache

import (
	"context"
	"errors"
	"github.com/johnwongx/webook/backend/internal/repository/cache/lru"
	"time"
)

type LocalCodeCache struct {
	c *lru.LRUCache[string]
}

type codeInfo struct {
	code       string
	create     time.Time
	expiration time.Time
	verifyCnt  int
}

func NewLocalCodeCache(c *lru.LRUCache[string]) CodeCache {
	return &LocalCodeCache{
		c: c,
	}
}

func (l *LocalCodeCache) Set(ctx context.Context, biz, phone, code string, experation time.Duration) error {
	curKey := key(biz, phone)
	now := time.Now()
	if ok, val := l.c.Get(curKey); ok {
		//存在并且创建时间距离当前时间小于1分钟
		cInfo, ok := val.(codeInfo)
		if !ok {
			return errors.New("系统错误")
		}
		if now.Sub(cInfo.create) < time.Minute {
			return ErrCodeSendTooMany
		}
	}

	info := codeInfo{
		code:       code,
		create:     now,
		expiration: now.Add(experation),
		verifyCnt:  3,
	}
	l.c.Put(curKey, info)
	return nil
}

func (l *LocalCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	curKey := key(biz, phone)
	ok, val := l.c.Get(curKey)
	if !ok {
		return false, ErrCodeUnknowError
	}

	now := time.Now()
	cInfo, ok := val.(codeInfo)
	if !ok {
		return false, errors.New("系统错误")
	}

	if cInfo.expiration.Before(now) {
		return false, errors.New("验证码过期")
	}

	if cInfo.verifyCnt < 0 {
		return false, ErrCodeVerifyTooMany
	}

	if cInfo.code != code {
		cInfo.verifyCnt -= 1
		l.c.Put(curKey, cInfo)
		return false, nil

	}

	cInfo.verifyCnt = -1
	l.c.Put(curKey, cInfo)
	return true, nil
}
