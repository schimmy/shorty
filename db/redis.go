package db

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	timeFormat = time.RFC3339Nano
	namespace  = "shorty"
)

func ns(slug string) string {
	return namespace + "." + slug
}

type Redis struct {
	c redis.Conn
}

func (r Redis) DeleteURL(slug string) error {
	n, err := redis.Int(r.c.Do("DEL", slug))
	if err != nil {
		return fmt.Errorf("Failed to delete '%s': %s", slug, err)
	} else if n != 1 {
		return fmt.Errorf("Record not found")
	}

	return nil
}

func (r Redis) ShortenURL(slug, longURL, owner string, expires time.Time) error {
	// save as a hash
	_, err := r.c.Do("HMSET", ns(slug),
		"long_url", longURL,
		"owner", owner)
	if err != nil {
		return fmt.Errorf("Failed to save '%s': %s", slug, err)
	}

	// set the expire time
	n, err := redis.Int(r.c.Do("PEXPIREAT", ns(slug), expires.Unix()))
	if err != nil {
		return fmt.Errorf("Failed to set expire time for '%s': %s", slug, err)
	} else if n == 0 {
		return fmt.Errorf("Failed to set expire time for '%s'", slug)
	}

	return nil
}

func (r Redis) GetLongURL(slug string) (string, error) {
	longURL, err := redis.String(r.c.Do("HGET", ns(slug), "long_url"))
	if err != nil {
		return "", fmt.Errorf("Failed to find long_url for '%s': %s", slug, err)
	} else if longURL == "" {
		return "", fmt.Errorf("Failed to find long_url for '%s'", slug)
	}

	return longURL, nil
}

func (r Redis) GetList() ([]ShortenObject, error) {
	// redis.Strings("")
	return []ShortenObject{}, nil
}
