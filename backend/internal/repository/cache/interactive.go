package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

var (
	//go:embed lua/interactive_incr_cnt.lua
	luaIncrCnt string
)

const (
	fieldReadCnt    = "read_cnt"
	fieldCollectCnt = "collect_cnt"
	fieldLikeCnt    = "like_cnt"
)

type InteractiveCache interface {
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	BatchIncrReadCntIfPresent(ctx context.Context, biz []string, bizId []int64) error
	IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error
}

type RedisInteractiveCache struct {
	client redis.Cmdable
}

func NewRedisInteractiveCache(client redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{
		client: client,
	}
}

func (r *RedisInteractiveCache) IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.key(biz, bizId)},
		fieldLikeCnt, 1).Err()
}

func (r *RedisInteractiveCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.key(biz, bizId)},
		fieldReadCnt, 1).Err()
}

func (r *RedisInteractiveCache) BatchIncrReadCntIfPresent(ctx context.Context, biz []string, bizId []int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisInteractiveCache) DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.key(biz, bizId)},
		fieldLikeCnt, -1).Err()
}

func (r *RedisInteractiveCache) IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.key(biz, bizId)},
		fieldCollectCnt, 1).Err()
}

func (r *RedisInteractiveCache) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	data, err := r.client.HGetAll(ctx, r.key(biz, bizId)).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(data) == 0 {
		return domain.Interactive{}, ErrKeyNotExisted
	}

	collectCnt, _ := strconv.ParseInt(data[fieldCollectCnt], 10, 64)
	likeCnt, _ := strconv.ParseInt(data[fieldLikeCnt], 10, 64)
	readCnt, _ := strconv.ParseInt(data[fieldReadCnt], 10, 64)

	return domain.Interactive{
		BizId:      bizId,
		Biz:        biz,
		ReadCnt:    readCnt,
		LikeCnt:    likeCnt,
		CollectCnt: collectCnt,
	}, nil
}

func (r *RedisInteractiveCache) Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error {
	key := r.key(biz, bizId)
	err := r.client.HMSet(ctx, key,
		fieldLikeCnt, intr.LikeCnt,
		fieldReadCnt, intr.ReadCnt,
		fieldCollectCnt, intr.CollectCnt).Err()
	if err != nil {
		return err
	}
	return r.client.Expire(ctx, key, time.Minute*15).Err()
}

func (r *RedisInteractiveCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}
