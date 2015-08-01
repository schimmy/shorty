package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"

	pq "github.com/lib/pq"
)

// ShortenObject holds the metadata and the mapping for a shortened URL
// I like the NullTime concept from the pq library, so even for other backends
// let's use it instead of checking whether the date is 0001-01-01
type ShortenObject struct {
	Slug     string      `json:"slug",sql:"slug"`
	LongURL  string      `json:"long_url",sql:"long_url"`
	Modified pq.NullTime `json:"modified_date",sql:"modified"`
	Expires  pq.NullTime `json:"expire_date",sql:"expires"`
}

// TODO: add tags / who added
// TODO: add delete URL
type ShortenBackend interface {
	ShortenURL(slug, longURL string, expires time.Time) error
	GetLongURL(slug string) (string, error)
	GetList() ([]ShortenObject, error)
}

type RedisDB struct {
	RedisPool *redis.Pool
	Host      string
	Port      int
}

type PostgresDB struct {
	c          *sql.DB
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	SchemaName string
	TableName  string
}

func getOrDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func NewPostgresDB() ShortenBackend {
	pgHost := getOrDefault("PG_HOST", "localhost")
	pgPort := getOrDefault("PG_HOST", "5432")
	pgUser := getOrDefault("PG_USER", "url_shortener")
	pgPass := getOrDefault("PG_PASS", "NOPE")
	pgDatabase := getOrDefault("PG_DB", "url_shortener")
	pgSSLMode := getOrDefault("PG_SSL", "disable")

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", pgHost, pgPort, pgUser, pgPass, pgDatabase, pgSSLMode)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	return &PostgresDB{c: db}
}

func (pgDB *PostgresDB) ShortenURL(slug string, longURL string, expires time.Time) error {
	// postgres & redshift don't have an upsert method yet
	// TODO: implement upsert functionality
	existingLong, err := pgDB.GetLongURL(slug)
	if existingLong == "" || err != nil { // TODO figure out what happens on nothing, err?

		q := fmt.Sprintf("INSERT INTO url_shortener.url_shortener(slug, long_url) VALUES ('%s', '%s')", slug, longURL)
		_, err := pgDB.c.Query(q)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("other long url exists: %s", existingLong)
}

func (pgDB *PostgresDB) GetLongURL(slug string) (string, error) {
	//var retObj ShortenObject
	q := fmt.Sprintf("SELECT long_url, expires FROM url_shortener.url_shortener WHERE slug = $1")
	log.Println("query: ", q)
	var long_url string
	var expires pq.NullTime
	err := pgDB.c.QueryRow(q, slug).Scan(&long_url, &expires)
	if err != nil {
		return "", err
	}
	log.Println("long: ", long_url, expires)
	return long_url, nil

}

func (pgDB *PostgresDB) GetList() ([]ShortenObject, error) {
	rows, err := pgDB.c.Query(`SELECT * FROM url_shortener.url_shortener`)
	log.Println("rows: ", rows)
	if err != nil {
		return nil, err
	}
	var retObjs []ShortenObject
	var slug string
	var long_url string
	var expires pq.NullTime
	var modified pq.NullTime
	for rows.Next() {
		err = rows.Scan(&slug, &long_url, &modified, &expires)
		if err != nil {
			return nil, fmt.Errorf("issue scanning row for list: %s", err)
		}
		retObjs = append(retObjs, ShortenObject{
			Slug:     slug,
			LongURL:  long_url,
			Expires:  expires,
			Modified: modified,
		})
	}
	log.Println("ret objs: ", retObjs)
	return retObjs, nil
}

func (redisDB *RedisDB) ShortenURL(slug string, longURL string, expires time.Time) error {
	log.Printf("shortening URL: %s to slug: %s", longURL, slug)
	sObj := ShortenObject{
		Slug:     slug,
		LongURL:  longURL,
		Modified: pq.NullTime{Time: time.Now()},
		Expires:  pq.NullTime{Time: expires},
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
	returnObj1 := ShortenObject{
		Slug:     "short1",
		LongURL:  "longURL1",
		Modified: pq.NullTime{Time: time.Now()},
	}
	returnObj2 := ShortenObject{
		Slug:     "short2",
		LongURL:  "longURL2",
		Modified: pq.NullTime{Time: time.Now()},
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
