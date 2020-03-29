package storage

import (
	"github.com/go-redis/redis/v7"
	"time"
)

const (
	UrlIdKey = "next.url.id"

	ShortLinkKey = "shortLink:%s:url"

	UrlHashKey = "urlHash:%s:url"

	ShortLinkDetailKey = "shortLink:%s:detail"
)

type RedisClient struct {
	Client *redis.Client
}

type UrlDetail struct {
	Url                 string        `json:"url"`
	CreatedAt           string        `json:"created_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

func NewRedisClient(addr string, pwd string, db int) *RedisClient {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})
	if _, err := c.Ping().Result(); err != nil {
		panic(err)
	}
	return &RedisClient{Client: c}
}