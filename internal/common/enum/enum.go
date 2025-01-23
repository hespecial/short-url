package enum

import "time"

const (
	ConfigFile = "./config/config.yaml"

	//NullCache          = "null"
	KeyUserView        = "uv"
	KeyOriginalUrlHash = "ouh"
	KeyShortUrlCode    = "suc"

	UserViewExpire = 24 * time.Hour
	UserViewExist  = true
)

type RevertType uint8

const (
	RevertToShort RevertType = iota
	RevertToOrigin
)

type Priority uint8

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
)
