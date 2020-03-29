package storage

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/mattheath/base62"
	"net/http"
	er "shorterUrl/src/error"
	"time"
)

func (c *RedisClient) Shorten(url string, exp int64) (string, error) {
	hash := toSha1(url)
	s, err := c.Client.Get(fmt.Sprintf(UrlHashKey, hash)).Result()
	if err == redis.Nil {
		// value not exist
	} else if err != nil {
		return "", nil
	} else {
		if s != "{}" {
			return s, nil
		}
	}

	// increase the global counter
	err = c.Client.Incr(UrlIdKey).Err()
	if err != nil {
		return "", nil
	}

	// encode global counter to base62
	id, err := c.Client.Get(UrlIdKey).Int64()
	if err != nil {
		return "", nil
	}
	eid := base62.EncodeInt64(id)

	// store the url against this encoded id
	err = c.Client.Set(
		fmt.Sprintf(ShortLinkKey, eid),
		url,
		time.Duration(exp)*time.Minute).Err()
	if err != nil {
		return "", nil
	}

	detail, err := json.Marshal(
		&UrlDetail{
			Url:                 url,
			CreatedAt:           time.Now().String(),
			ExpirationInMinutes: time.Duration(exp),
		})
	if err != nil {
		return "", nil
	}

	// store the url detail against this encoded id
	err = c.Client.Set(
		fmt.Sprintf(ShortLinkDetailKey, eid),
		detail,
		time.Duration(exp)*time.Minute).Err()
	if err != nil {
		return "", nil
	}
	return eid, nil
}

func (c *RedisClient) ShortLinkInfo(eid string) (interface{}, error) {
	res, err := c.Client.Get(fmt.Sprintf(ShortLinkDetailKey, eid)).Result()
	if err == redis.Nil {
		return "", er.StatusError{
			Code: http.StatusNotFound,
			Err:  errors.New("UnKnown short url"),
		}
	} else if err != nil {
		return "", err
	} else {
		return res, nil
	}
}

func (c *RedisClient) UnShorten(eid string) (string, error) {
	res, err := c.Client.Get(fmt.Sprintf(ShortLinkKey, eid)).Result()
	if err == redis.Nil {
		return "", er.StatusError{
			Code: http.StatusNotFound,
			Err:  errors.New("not exist eid to find url"),
		}
	} else if err != nil {
		return "", err
	} else {
		return res, nil
	}
}

func toSha1(str string) string {
	sh := sha1.New()
	return string(sh.Sum([]byte(str)))
}
