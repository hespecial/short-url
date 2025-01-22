package dao

import (
	"context"
	"gorm.io/gorm"
	"short-url/internal/model"
)

type ShortUrlDao interface {
	GetUrlMappingByOriginalUrlHash(ctx context.Context, hash string) (*model.UrlMapping, error)
	GetUrlMappingByShortUrlCode(ctx context.Context, shortUrlCode string) (*model.UrlMapping, error)
	CreateUrlMapping(ctx context.Context, urlMapping *model.UrlMapping) error
	CreateAccessLog(ctx context.Context, accessLog *model.AccessLog) error
	UpdateUrlMapping(ctx context.Context, urlMapping *model.UrlMapping) error
}

type shortUrlDao struct {
	db *gorm.DB
}

func NewShortUrlDao(db *gorm.DB) ShortUrlDao {
	return &shortUrlDao{
		db: db,
	}
}

func (s *shortUrlDao) GetUrlMappingByOriginalUrlHash(ctx context.Context, hash string) (*model.UrlMapping, error) {
	urlMapping := model.UrlMapping{}
	return &urlMapping, s.db.WithContext(ctx).
		Where(&model.UrlMapping{
			OriginalUrlHash: hash,
			Deleted:         false,
		}).
		First(&urlMapping).Error
}

func (s *shortUrlDao) GetUrlMappingByShortUrlCode(ctx context.Context, shortUrlCode string) (*model.UrlMapping, error) {
	urlMapping := model.UrlMapping{}
	return &urlMapping, s.db.WithContext(ctx).
		Where(&model.UrlMapping{
			ShortUrlCode: shortUrlCode,
			Deleted:      false,
		}).
		First(&urlMapping).Error
}

func (s *shortUrlDao) CreateUrlMapping(ctx context.Context, urlMapping *model.UrlMapping) error {
	return s.db.WithContext(ctx).Create(urlMapping).Error
}

func (s *shortUrlDao) CreateAccessLog(ctx context.Context, accessLog *model.AccessLog) error {
	return s.db.WithContext(ctx).Create(accessLog).Error
}

func (s *shortUrlDao) UpdateUrlMapping(ctx context.Context, urlMapping *model.UrlMapping) error {
	return s.db.WithContext(ctx).Save(urlMapping).Error
}
