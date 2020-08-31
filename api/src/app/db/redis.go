package db

import (
	"api/config"
	"fmt"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"os"
)

var Redis *redis.Client

func RedisInit() {
	addr := fmt.Sprintf("%s:%d", config.Config.Redis.Host, config.Config.Redis.Port)

	Redis = redis.NewClient(&redis.Options{
		Addr: addr,
		Password: "",
	})

	pong, err := Redis.Ping().Result()

	if err == nil && pong == "PONG" {
		log.Infof("Connected to Redis: %s", addr)
	} else {
		log.Error("Error while connecting to Redis")
		log.Error(err)
		os.Exit(1)
	}
}

