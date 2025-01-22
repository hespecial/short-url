package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"short-url/common/enum"
	"short-url/internal/model"
	"short-url/internal/repo/cache"
	"short-url/internal/repo/dao"
)

type ShortUrlRepo interface {
	GetUrlMappingByOriginalUrlHash(ctx context.Context, hash string) (*model.UrlMapping, error)
	GetUrlMappingByShortUrlCode(ctx context.Context, shortUrlCode string) (*model.UrlMapping, error)
	CreateUrlMapping(ctx context.Context, urlMapping *model.UrlMapping) error
	CreateAccessLog(ctx context.Context, accessLog *model.AccessLog) error

	SetUserView(ctx context.Context, shortUrlCode, ip string) (bool, error)
}

type shortUrlRepo struct {
	dao   dao.ShortUrlDao
	cache cache.ShortUrlCache
}

func NewShortUrlRepo(dao dao.ShortUrlDao, cache cache.ShortUrlCache) ShortUrlRepo {
	return &shortUrlRepo{
		dao:   dao,
		cache: cache,
	}
}

func (s *shortUrlRepo) GetUrlMappingByOriginalUrlHash(ctx context.Context, hash string) (*model.UrlMapping, error) {
	key := fmt.Sprintf("%s-%s", enum.KeyOriginalUrlHash, hash)
	jsonStr, err := s.cache.GetUrlMappingByOriginalUrlHash(ctx, key)
	if err != nil {
		return nil, err
	}
	if jsonStr == enum.NullCache {
		return s.dao.GetUrlMappingByOriginalUrlHash(ctx, hash)
	}
	var urlMapping model.UrlMapping
	if err = json.Unmarshal([]byte(jsonStr), &urlMapping); err != nil {
		return nil, err
	}
	return &urlMapping, nil
}

func (s *shortUrlRepo) GetUrlMappingByShortUrlCode(ctx context.Context, shortUrlCode string) (*model.UrlMapping, error) {
	return s.dao.GetUrlMappingByShortUrlCode(ctx, shortUrlCode)
}

func (s *shortUrlRepo) CreateUrlMapping(ctx context.Context, urlMapping *model.UrlMapping) error {
	// key := fmt.Sprintf("%s-%s", enum.KeyShortUrlCode, urlMapping.ShortUrlCode)

	return s.dao.CreateUrlMapping(ctx, urlMapping)
}

func (s *shortUrlRepo) CreateAccessLog(ctx context.Context, accessLog *model.AccessLog) error {
	return s.dao.CreateAccessLog(ctx, accessLog)
}

func (s *shortUrlRepo) SetUserView(ctx context.Context, shortUrlCode, ip string) (bool, error) {
	key := fmt.Sprintf("%s-%s-%s", enum.KeyUserView, shortUrlCode, ip)
	return s.cache.SetUserView(ctx, key)
}
