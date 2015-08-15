package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisDB struct {
	RedisPool *redis.Pool
	Host      string
	Port      int
}

func (redisDB *RedisDB) DeleteURL(slug string) error {
	return nil
}

func (redisDB *RedisDB) ShortenURL(slug, longURL, owner string, expires time.Time) error {
	log.Printf("shortening URL: %s to slug: %s", longURL, slug)
	sObj := ShortenObject{
		Slug:     slug,
		LongURL:  longURL,
		Modified: time.Now(),
		Expires:  expires,
	}
	serializedObj, err := json.Marshal(sObj)
	if err != nil {
		return fmt.Errorf("error marshalling json: %s", err.Error())
	}
	redisClient := redisDB.RedisPool.Get()
	defer redisClient.Close()
	//if _, err = redisClient.Do("HSET", redis.Args{}.Add("shorten:"+slug).AddFlat(serializedObj)...); err != nil {
	if _, err = redisClient.Do("SET", "shorten:"+slug, serializedObj); err != nil {
		return err
	}
	return nil
}

func (redisDB *RedisDB) GetLongURL(slug string) (string, error) {
	redisClient := redisDB.RedisPool.Get()
	defer redisClient.Close()
	retVal, err := redisClient.Do("GET", "shorten:"+slug)
	if err != nil {
		return "", err
	}
	if retVal == nil {
		return "", fmt.Errorf("No long url found for slug: %s", slug)
	}
	var retShortObj ShortenObject
	if err = json.Unmarshal(retVal.([]byte), &retShortObj); err != nil {
		return "", fmt.Errorf("Error deserializing slug: %s from redis, err: %s", slug, err)
	}
	return retShortObj.LongURL, nil
}

func (redisDB *RedisDB) GetList() ([]ShortenObject, error) {

	//err := redis.Values(c.Do("SCAN", id))
	//if err != nil {
	//panic(err)
	//}
	returnObj1 := ShortenObject{
		Slug:     "short1",
		LongURL:  "longURL1",
		Modified: time.Now(),
	}
	returnObj2 := ShortenObject{
		Slug:     "short2",
		LongURL:  "longURL2",
		Modified: time.Now(),
	}

	return []ShortenObject{returnObj1, returnObj2}, nil
}

func NewRedisDB() ShortenBackend {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	redisPool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisURL)
		},
	}

	// test initial ping to surface misconfig
	c := redisPool.Get()
	_, err := c.Do("PING")
	c.Close()
	if err != nil {
		log.Fatalf("Error pinging initial redis connection: %s", err)
	}
	return &RedisDB{RedisPool: redisPool}
}
