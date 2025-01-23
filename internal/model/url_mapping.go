package model

import (
	"short-url/internal/common/enum"
)

type UrlMapping struct {
	Id              int64
	ShortUrlCode    string
	OriginalUrl     string
	OriginalUrlHash string
	Priority        enum.Priority
	Disable         bool
	CreateTime      int64
	UpdateTime      int64
	ExpireTime      int64
	Deleted         bool
	Comment         string
}

func (*UrlMapping) TableName() string {
	return "url_mapping"
}
