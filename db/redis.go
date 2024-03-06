package db

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

const RAPORT_REDIS_QUEUE_BUILDY_KEY string = "profile:progress-raport:build-list"
const RAPORT_REDIS_QUEUE_FAILED_BUILD_KEY string = "profile:progress-raport:failed-delivery-list"

const RAPORT_RESULT_REDIS_QUEUE_BUILDY_KEY string = "profile:raport-result:build-list"
const RAPORT_RESULT_REDIS_QUEUE_FAILED_BUILD_KEY string = "profile:raport-result:failed-delivery-list"

var RedisContext context.Context
var Redis *redis.ClusterClient

func NewRedisCluster() *redis.ClusterClient {
	if Redis != nil {
		return Redis
	}
	env := os.Getenv("REDIS_CLUSTERS")
	if env == "" {
		panic("REDIS_CLUSTERS is not set")
	}

	cluters := strings.Split(env, ",")

	c := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    cluters,
		Password: os.Getenv("REDIS_CLUSTERS_PASSWORD"),
	})

	Redis = c
	RedisContext = context.Background()
	return c
	//return Redis
}

func GetBuildRaportWorkerProcessInSecond() int64 {
	env := os.Getenv("RAPORT_WORKER_REFRESH_SECOND")

	if env != "" {
		n, err := strconv.Atoi(env)
		if err != nil {
			panic("error parse RAPORT_WORKER_REFRESH_SECOND environment variable is not integer")
		}

		return int64(n)
	}

	return 1
}

func GetBuildRaportResultWorkerProcessInSecond() int64 {
	env := os.Getenv("RAPORT_RESULT_WORKER_REFRESH_SECOND")

	if env != "" {
		n, err := strconv.Atoi(env)
		if err != nil {
			panic("error parse RAPORT_WORKER_REFRESH_SECOND environment variable is not integer")
		}

		return int64(n)
	}

	return 1
}
