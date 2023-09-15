package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"sync"
)

const sms_key = "SMS_INFO_DATA"

type SMSCache interface {
	Add(ctx context.Context, info SMSInfo) error
	Take(ctx context.Context, cnt int) ([]SMSInfo, error)
	KeyExists(ctx context.Context) (bool, error)
}

type smsCache struct {
	c  redis.Cmdable
	mu sync.Mutex
}

func NewSMSCache(c redis.Cmdable) SMSCache {
	return &smsCache{
		c:  c,
		mu: sync.Mutex{},
	}
}

func (s *smsCache) Add(ctx context.Context, info SMSInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	jsonData, err := json.Marshal(info)
	if err != nil {
		return err
	}

	return s.c.LPush(ctx, sms_key, jsonData).Err()
}

func (s *smsCache) Take(ctx context.Context, cnt int) ([]SMSInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var res []SMSInfo
	for i := 0; i < cnt; i++ {
		val, err := s.c.RPop(ctx, sms_key).Result()
		if err != nil {
			if err == redis.Nil {
				break
			}
			return nil, err
		}

		var info SMSInfo
		err = json.Unmarshal([]byte(val), &info)
		if err != nil {
			return nil, err
		}
		res = append(res, info)
	}
	return res, nil
}

func (s *smsCache) KeyExists(ctx context.Context) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, err := s.c.Exists(ctx, sms_key).Result()
	if err != nil {
		return false, err
	}

	return val > 0, nil
}

type SMSInfo struct {
	Tpl     string   `json:"tpl"`
	Args    []string `json:"args"`
	Numbers []string `json:"numbers"`
}
