package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

type LocalCodeCache struct {
	data map[string]codeInfo
	lock sync.Locker
}

type codeInfo struct {
	code       string
	create     time.Time
	expiration time.Time
	verifyCnt  int
}

func (l *LocalCodeCache) Set(ctx context.Context, biz, phone, code string, experation time.Duration) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	curKey := key(biz, phone)
	if val, ok := l.data[curKey]; ok {
		//存在并且创建时间距离当前时间小于1分钟
		if val.expiration.Sub(time.Now()) < time.Minute {
			return ErrCodeSendTooMany
		}
	}

	info := codeInfo{
		code:       code,
		create:     time.Now(),
		expiration: time.Now().Add(experation),
		verifyCnt:  3,
	}
	l.data[curKey] = info
	return nil
}

func (l *LocalCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	curKey := key(biz, phone)
	val, ok := l.data[curKey]
	if !ok {
		return false, ErrCodeUnknowError
	}

	now := time.Now()
	if val.expiration.Before(now) {
		return false, errors.New("验证码过期")
	}

	if val.verifyCnt < 0 {
		return false, ErrCodeVerifyTooMany
	}

	if val.code != code {
		val.verifyCnt -= 1
		l.data[curKey] = val
		return false, nil

	}

	val.verifyCnt = -1
	l.data[curKey] = val
	return true, nil
}
