package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

var ctx = context.Background()

var (
	client = &RedisClient{}
)

type RedisClient struct {
	c      *redis.Client
	logger *log.Logger
}
