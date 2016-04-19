package db

import (
	"fmt"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"gopkg.in/Clever/kayvee-go.v3/logger"
)

const (
	timeFormat = time.RFC3339Nano
	namespace  = "shorty"
)

func ns(slug string) string {
	return namespace + "." + slug
}

// Redis is a wrapper over a Redis pool connection.
type Redis struct {
	p *redis.Pool
}

// NewRedisDB connects to Redis and pools connections.
func NewRedisDB() ShortenBackend {
	redisURL := getOrDefault("REDIS_URL", "localhost:6379")
	pool := newPool(redisURL)
	conn, err := pool.Dial()
	if err != nil {
		lg.ErrorD("redis.connection.failure", logger.M{
			"msg": fmt.Sprintf("Failed to connect to redis at '%s': %s", redisURL, err)})
		os.Exit(1)
	}
	conn.Close()
	return Redis{p: pool}
}

// DeleteURL removes all reference to a URL.
func (r Redis) DeleteURL(slug string) error {
	c := r.p.Get()
	defer c.Close()

	n, err := redis.Int(c.Do("DEL", ns(slug)))
	if err != nil {
		return fmt.Errorf("Failed to delete '%s': %s", slug, err)
	} else if n < 1 {
		return ErrNotFound
	}

	lg.InfoD("slug.delete", logger.M{"name": slug})
	return nil
}

// ShortenURL adds a URL object to the db.
func (r Redis) ShortenURL(slug, longURL, owner string, expires time.Time) error {
	c := r.p.Get()
	defer c.Close()

	// save as a hash
	_, err := c.Do("HMSET", ns(slug),
		"long_url", longURL,
		"owner", owner)
	if err != nil {
		return fmt.Errorf("Failed to save '%s': %s", slug, err)
	}

	// set the expire time if it is a valid time
	if !expires.Before(time.Now()) {
		n, err := redis.Int(c.Do("PEXPIREAT", ns(slug), expires.Unix()))
		if err != nil {
			return fmt.Errorf("Failed to set expire time for '%s': %s", slug, err)
		} else if n == 0 {
			return fmt.Errorf("Failed to set expire time for '%s'", slug)
		}
	}

	return nil
}

// GetLongURL searches for the short URL reference in order to return the long url.
func (r Redis) GetLongURL(slug string) (string, error) {
	c := r.p.Get()
	defer c.Close()

	longURL, err := redis.String(c.Do("HGET", ns(slug), "long_url"))

	if longURL == "" || err == redis.ErrNil {
		return "", ErrNotFound
	} else if err != nil {
		return "", fmt.Errorf("Failed to find long_url for '%s': %s", slug, err)
	}

	return longURL, nil
}

// GetList lists all shortened URLs.
func (r Redis) GetList() ([]ShortenObject, error) {
	c := r.p.Get()
	defer c.Close()

	// get the list of every URL hash
	urlKeys, err := redis.Strings(c.Do("KEYS", ns("*")))
	if err != nil {
		return []ShortenObject{}, fmt.Errorf("Failed to scan to find all records: %s", err)
	}

	// pipeline retrieving every
	for _, key := range urlKeys {
		c.Send("HMGET", key, "slug", "owner", "long_url")
	}

	// flush all commands
	err = c.Flush()
	if err != nil {
		return []ShortenObject{}, fmt.Errorf("Failed to flush query for list: %s", err)
	}

	objs := make([]ShortenObject, len(urlKeys))
	for i := range objs {
		url, err := redis.Values(c.Receive())
		if err != nil {
			return []ShortenObject{}, fmt.Errorf("Failed to recieve value for list: %s", err)
		}

		u := objs[i]
		_, err = redis.Scan(url, &u.Slug, &u.Owner, &u.LongURL)
		if err != nil {
			return []ShortenObject{}, fmt.Errorf("Failed to marshal value: %s", err)
		}
	}

	return []ShortenObject{}, nil
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
