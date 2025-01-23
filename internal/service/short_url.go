package service

import (
	"context"
	"errors"
	"fmt"
	"short-url/global"
	"short-url/internal/common/enum"
	"short-url/internal/model"
	"short-url/internal/repo"
	"short-url/internal/util"
	"short-url/pkg/bloom"
	"time"

	"gorm.io/gorm"
)

type ShortUrlService struct {
	repo *repo.ShortUrlRepo
}

func NewShortUrlService(repo *repo.ShortUrlRepo) *ShortUrlService {
	return &ShortUrlService{
		repo: repo,
	}
}

func getFullShortUrlPath(shortUrlCode string) string {
	return fmt.Sprintf("http://%s:%d/%s", global.Conf.App.Host, global.Conf.App.Port, shortUrlCode)
}

func (s *ShortUrlService) generateShortUrlCode(priority enum.Priority) string {
	var minLen, maxLen int
	switch priority {
	case enum.PriorityLow:
		minLen, maxLen = 7, 9
	case enum.PriorityMedium:
		minLen, maxLen = 4, 6
	case enum.PriorityHigh:
		minLen, maxLen = 1, 3
	default:
		return ""
	}

	var g func(int) string
	g = func(times int) string {
		if times > 1000 {
			return ""
		}
		random := util.GenerateRandomBytes(minLen, maxLen)
		if !bloom.Contains(random) {
			bloom.Add(random)
			return string(random)
		}
		return g(times + 1)
	}

	return g(0)
}

func (s *ShortUrlService) RevertToShortUrl(ctx context.Context, url string, priority enum.Priority, comment string) (string, error) {
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}
	originalUrlHash := util.MD5(url)
	urlMapping, err := s.repo.GetUrlMappingByOriginalUrlHash(ctx, originalUrlHash)
	if err == nil {
		return getFullShortUrlPath(urlMapping.ShortUrlCode), nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	urlMapping = &model.UrlMapping{
		ShortUrlCode:    s.generateShortUrlCode(priority),
		OriginalUrl:     url,
		OriginalUrlHash: originalUrlHash,
		Priority:        priority,
		CreateTime:      time.Now().Unix(),
		UpdateTime:      time.Now().Unix(),
		Comment:         comment,
	}
	err = s.repo.CreateUrlMapping(ctx, urlMapping)
	if err != nil {
		return "", err
	}
	return getFullShortUrlPath(urlMapping.ShortUrlCode), nil
}

func (s *ShortUrlService) GetUrlMappingByShortUrlCode(ctx context.Context, shortUrlCode string) (*model.UrlMapping, error) {
	return s.repo.GetUrlMappingByShortUrlCode(ctx, shortUrlCode)
}

func (s *ShortUrlService) LogAccess(ctx context.Context, urlMappingId int64, ip, userAgent string) error {
	accessLog := &model.AccessLog{
		UrlMappingId: urlMappingId,
		Ip:           ip,
		UserAgent:    userAgent,
		AccessTime:   time.Now().Unix(),
	}
	return s.repo.CreateAccessLog(ctx, accessLog)
}

func (s *ShortUrlService) ProcessAccess(ctx context.Context, urlMappingId int64, ip string) error {
	accessStatistic, err := s.repo.GetAccessStatisticByUrlMappingId(ctx, urlMappingId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		accessStatistic = &model.AccessStatistic{
			UrlMappingId: urlMappingId,
		}
	}

	accessStatistic.Pv++
	ok, err := s.repo.SetUserView(ctx, urlMappingId, ip)
	if err != nil {
		return err
	}
	if ok {
		accessStatistic.Uv++
	}
	accessStatistic.LastAccessTime = time.Now().Unix()

	return s.repo.SaveAccessStatistic(ctx, accessStatistic)
}
