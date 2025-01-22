package global

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

var Redis *redis.Client

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", Conf.Redis.Host, Conf.Redis.Port),
		DB:   Conf.Redis.Database,
	})
	if err := Redis.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Init redis error: ", err.Error())
	}
}
