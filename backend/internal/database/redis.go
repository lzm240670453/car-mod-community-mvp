package database

import "github.com/redis/go-redis/v9"

func NewRedis(opts redis.Options) *redis.Client {
	return redis.NewClient(&opts)
}
