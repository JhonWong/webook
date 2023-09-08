package ratelimit

import "context"

type Limiter interface {
	//Limit是否有出发限流。key是限流对象
	//bool 代表是否触发限流， true 就是要限流
	//error 限流器本身是否有错误
	Limit(ctx context.Context, key string) (bool, error)
}
