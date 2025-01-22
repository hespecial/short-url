package cache

import (
	"context"
	"short-url/common/enum"

	"github.com/redis/go-redis/v9"
)

type ShortUrlCache interface {
	SetUserView(ctx context.Context, key string) (bool, error)
	GetUrlMappingByOriginalUrlHash(ctx context.Context, key string) (string, error)
}

type shortUrlCache struct {
	rdb *redis.Client
}

func NewShortUrlCache(rdb *redis.Client) ShortUrlCache {
	return &shortUrlCache{
		rdb: rdb,
	}
}

func (s *shortUrlCache) SetUserView(ctx context.Context, key string) (bool, error) {
	return s.rdb.SetNX(ctx, key, enum.UserViewExist, enum.UserViewExpire).Result()
}

func (s *shortUrlCache) GetUrlMappingByOriginalUrlHash(ctx context.Context, key string) (string, error) {
	return s.rdb.Get(ctx, key).Result()
}

func (s *shortUrlCache) SetShortUrlCode(ctx context.Context, key string, value string) error {
	return s.rdb.Set(ctx, key, value, 0).Err()
}
