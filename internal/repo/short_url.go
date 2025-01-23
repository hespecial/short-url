package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"short-url/internal/common/enum"
	"short-url/internal/model"
)

type ShortUrlRepo struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewShortUrlRepo(db *gorm.DB, rdb *redis.Client) *ShortUrlRepo {
	return &ShortUrlRepo{
		db:  db,
		rdb: rdb,
	}
}

func (s *ShortUrlRepo) GetUrlMappingByOriginalUrlHash(ctx context.Context, hash string) (*model.UrlMapping, error) {
	key := fmt.Sprintf("%s-%s", enum.KeyOriginalUrlHash, hash)
	jsonStr, err := s.rdb.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		slog.Error(err.Error())
	}
	var urlMapping model.UrlMapping
	if err = json.Unmarshal([]byte(jsonStr), &urlMapping); err == nil {
		return &urlMapping, nil
	}

	urlMapping = model.UrlMapping{}
	err = s.db.WithContext(ctx).
		Where(&model.UrlMapping{
			OriginalUrlHash: hash,
			Deleted:         false,
		}).First(&urlMapping).Error
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(urlMapping)
	if err != nil {
		slog.Error(err.Error())
	}
	if err = s.rdb.Set(ctx, key, string(b), 0).Err(); err != nil {
		slog.Error(err.Error())
	}

	return &urlMapping, nil
}

func (s *ShortUrlRepo) GetUrlMappingByShortUrlCode(ctx context.Context, shortUrlCode string) (*model.UrlMapping, error) {
	key := fmt.Sprintf("%s-%s", enum.KeyShortUrlCode, shortUrlCode)
	jsonStr, err := s.rdb.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		slog.Error(err.Error())
	}
	var urlMapping model.UrlMapping
	if err = json.Unmarshal([]byte(jsonStr), &urlMapping); err == nil {
		return &urlMapping, nil
	}

	urlMapping = model.UrlMapping{}
	err = s.db.WithContext(ctx).
		Where(&model.UrlMapping{
			ShortUrlCode: shortUrlCode,
			Deleted:      false,
		}).
		First(&urlMapping).Error
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(urlMapping)
	if err != nil {
		slog.Error(err.Error())
	}
	if err = s.rdb.Set(ctx, key, string(b), 0).Err(); err != nil {
		slog.Error(err.Error())
	}

	return &urlMapping, nil
}

func (s *ShortUrlRepo) CreateUrlMapping(ctx context.Context, urlMapping *model.UrlMapping) error {
	return s.db.WithContext(ctx).Create(urlMapping).Error
}

func (s *ShortUrlRepo) CreateAccessLog(ctx context.Context, accessLog *model.AccessLog) error {
	return s.db.WithContext(ctx).Create(accessLog).Error
}

func (s *ShortUrlRepo) SetUserView(ctx context.Context, urlMappingId int64, ip string) (bool, error) {
	key := fmt.Sprintf("%s-%d-%s", enum.KeyUserView, urlMappingId, ip)
	return s.rdb.SetNX(ctx, key, enum.UserViewExist, enum.UserViewExpire).Result()
}

func (s *ShortUrlRepo) GetAccessStatisticByUrlMappingId(ctx context.Context, urlMappingId int64) (*model.AccessStatistic, error) {
	var accessStatistic model.AccessStatistic
	return &accessStatistic, s.db.WithContext(ctx).
		First(&accessStatistic, "url_mapping_id = ?", urlMappingId).Error
}

func (s *ShortUrlRepo) SaveAccessStatistic(ctx context.Context, accessStatistic *model.AccessStatistic) error {
	return s.db.WithContext(ctx).Save(accessStatistic).Error
}
