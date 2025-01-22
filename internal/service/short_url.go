package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"short-url/common/enum"
	"short-url/global"
	"short-url/internal/model"
	"short-url/internal/repo"
	"short-url/pkg/bloom"
	"short-url/util"
	"time"

	"gorm.io/gorm"
)

type ShortUrlService interface {
	RevertToShortUrl(ctx context.Context, url string, priority enum.Priority, comment string) (string, error)
	GetUrlMappingByShortUrlCode(ctx context.Context, shortUrlCode string) (*model.UrlMapping, error)
	LogAccess(ctx context.Context, urlMapping int64, ip, userAgent string) error
}

type shortUrlService struct {
	repo repo.ShortUrlRepo
}

func NewShortUrlService(repo repo.ShortUrlRepo) ShortUrlService {
	return &shortUrlService{
		repo: repo,
	}
}

func getFullShortUrlPath(shortUrlCode string) string {
	return fmt.Sprintf("http://%s:%d/%s", global.Conf.App.Host, global.Conf.App.Port, shortUrlCode)
}

const (
	seed = "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
)

var (
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func (s *shortUrlService) generateShortUrlCode(priority enum.Priority, times int) string {
	if times > 1000 {
		return ""
	}

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

	length := r.Intn(maxLen-minLen+1) + minLen
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = seed[r.Intn(len(seed))]
	}

	if !bloom.Contains(result) {
		bloom.Add(result)
		return string(result)
	}

	return s.generateShortUrlCode(priority, times+1)
}

func (s *shortUrlService) RevertToShortUrl(ctx context.Context, url string, priority enum.Priority, comment string) (string, error) {
	originalUrlHash := util.MD5(url)
	urlMapping, err := s.repo.GetUrlMappingByOriginalUrlHash(ctx, originalUrlHash)
	if err == nil {
		return getFullShortUrlPath(urlMapping.ShortUrlCode), nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}

	urlMapping = &model.UrlMapping{
		ShortUrlCode:    s.generateShortUrlCode(priority, 0),
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

func (s *shortUrlService) GetUrlMappingByShortUrlCode(ctx context.Context, shortUrlCode string) (*model.UrlMapping, error) {
	return s.repo.GetUrlMappingByShortUrlCode(ctx, shortUrlCode)
}

func (s *shortUrlService) LogAccess(ctx context.Context, urlMappingId int64, ip, userAgent string) error {
	accessLog := &model.AccessLog{
		UrlMappingId: urlMappingId,
		Ip:           ip,
		UserAgent:    userAgent,
		AccessTime:   time.Now().Unix(),
	}
	return s.repo.CreateAccessLog(ctx, accessLog)
}

func (s *shortUrlService) ProcessAccess(ctx context.Context, urlMappingId int64) error {
	// accessStatic := &model.AccessStatistic{
	// 	UrlMappingId: urlMappingId,
	// 	LastAccessTime: time.Now().Unix(),
	// }
	panic("implement me")

}
