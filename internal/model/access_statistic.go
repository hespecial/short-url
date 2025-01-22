package model

// AccessStatistic 表示访问统计信息
type AccessStatistic struct {
	Id             int64
	UrlMappingId   int64
	Pv             int
	Uv             int
	LastAccessTime int64
}

// TableName 指定表名
func (AccessStatistic) TableName() string {
	return "access_statistic"
}
