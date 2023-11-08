package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExisted = errors.New("key not existed")
var _ ArticleCache = &RedisArticleCache{}

type ArticleCache interface {
	GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error
	DeleteFirstPage(ctx context.Context, uid int64) error

	Set(ctx context.Context, article domain.Article) error
	Get(ctx context.Context, id, uid int64) (domain.Article, error)
	Delete(ctx context.Context, id, uid int64) error

	SetPub(ctx context.Context, article domain.Article) error
	GetPub(ctx context.Context, id int64) (domain.Article, error)
	DeletePub(ctx context.Context, id int64) error
}

type RedisArticleCache struct {
	client redis.Cmdable
}

func NewRedisArticleCache(client redis.Cmdable) ArticleCache {
	return &RedisArticleCache{
		client: client,
	}
}

func (r *RedisArticleCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	bts, err := r.client.Get(ctx, r.firstPageKey(uid)).Bytes()
	if err == redis.Nil {
		return nil, ErrKeyNotExisted
	} else if err != nil {
		return nil, err
	}
	var arts []domain.Article
	err = json.Unmarshal(bts, &arts)
	return arts, err
}

func (r *RedisArticleCache) SetFirstPage(ctx context.Context, uid int64, arts []domain.Article) error {
	for i := range arts {
		arts[i].Content = arts[i].Abstract()
	}
	bts, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.firstPageKey(uid), bts, time.Minute*10).Err()
}

func (r *RedisArticleCache) DeleteFirstPage(ctx context.Context, uid int64) error {
	err := r.client.Del(ctx, r.firstPageKey(uid)).Err()
	if err == redis.Nil {
		return ErrKeyNotExisted
	}
	return err
}

func (r *RedisArticleCache) Set(ctx context.Context, art domain.Article) error {
	bts, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.articleKey(art.Id, art.Author.Id), bts, time.Minute).Err()
}

func (r *RedisArticleCache) Get(ctx context.Context, id, uid int64) (domain.Article, error) {
	bts, err := r.client.Get(ctx, r.articleKey(id, uid)).Bytes()
	if err == redis.Nil {
		return domain.Article{}, ErrKeyNotExisted
	} else if err != nil {
		return domain.Article{}, err
	}
	var art domain.Article
	err = json.Unmarshal(bts, &art)
	return art, err
}

func (r *RedisArticleCache) Delete(ctx context.Context, id, uid int64) error {
	err := r.client.Del(ctx, r.articleKey(id, uid)).Err()
	if err == redis.Nil {
		return ErrKeyNotExisted
	}
	return err
}

func (r *RedisArticleCache) SetPub(ctx context.Context, art domain.Article) error {
	bts, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.pubArticleKey(art.Id), bts, time.Minute).Err()
}

func (r *RedisArticleCache) GetPub(ctx context.Context, id int64) (domain.Article, error) {
	bts, err := r.client.Get(ctx, r.pubArticleKey(id)).Bytes()
	if err == redis.Nil {
		return domain.Article{}, ErrKeyNotExisted
	} else if err != nil {
		return domain.Article{}, err
	}
	var art domain.Article
	err = json.Unmarshal(bts, &art)
	return art, err
}

func (r *RedisArticleCache) DeletePub(ctx context.Context, id int64) error {
	err := r.client.Del(ctx, r.pubArticleKey(id)).Err()
	if err == redis.Nil {
		return ErrKeyNotExisted
	}
	return err
}

func (r *RedisArticleCache) pubArticleKey(id int64) string {
	return fmt.Sprintf("published_article:%d", id)
}

func (r *RedisArticleCache) articleKey(id, uid int64) string {
	return fmt.Sprintf("article:%d:%d", id, uid)
}

func (r *RedisArticleCache) firstPageKey(uid int64) string {
	return fmt.Sprintf("author_article_list:%d", uid)
}
