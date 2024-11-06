package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis为github.com/eachain/360-tuitui-robot/webhook.Cache的redis实现。
// 用于分布式环境。
type Redis struct {
	cli    *redis.Client
	expire time.Duration
}

// NewRedis复用业务所用*redis.Client。expire为nonce过期时间。
func NewRedis(cli *redis.Client, expire time.Duration) *Redis {
	return &Redis{cli: cli, expire: expire}
}

// implement github.com/eachain/360-tuitui-robot/webhook.Cache.
func (rds *Redis) Set(nonce string) bool {
	key := "tuitui:robot:webhook:" + nonce
	ok, err := rds.cli.SetNX(context.Background(), key, nonce, rds.expire).Result()
	if err != nil {
		// redis报错的情况下，允许业务继续执行。
		// 考虑点：重复执行好过不执行。
		return true
	}
	return ok
}
