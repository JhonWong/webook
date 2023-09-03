package ioc

import "github.com/johnwongx/webook/backend/internal/repository/cache/lru"

func InitLRUCache() *lru.LRUCache[string] {
	return lru.NewLRUCache[string](3)
}
