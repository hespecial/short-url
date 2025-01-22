package model

type AccessLog struct {
	Id           int64
	UrlMappingId int64
	Ip           string
	UserAgent    string
	AccessTime   int64
}

func (*AccessLog) TableName() string {
	return "access_log"
}
