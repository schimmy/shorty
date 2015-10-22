package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/Clever/kayvee-go.v2"
)

// msg is a convenience type for kayvee
type msg map[string]interface{}

type ErrNotFound struct{}

func (e ErrNotFound) Error() string {
	return "No items found by the db layer"
}

// ShortenObject holds the metadata and the mapping for a shortened URL
// I like the NullTime concept from the pq library, so even for other backends
// let's use it instead of checking whether the date is 0001-01-01
type ShortenObject struct {
	Slug     string    `json:"slug"`
	Owner    string    `json:"owner"`
	LongURL  string    `json:"long_url"`
	Modified time.Time `json:"modified_date,omitempty"`
	Expires  time.Time `json:"expire_date,omitempty"`
}

// ShortenBackend represents the necessary interface for storing and updating URLs.
type ShortenBackend interface {
	DeleteURL(slug string) error
	ShortenURL(slug, longURL, owner string, expires time.Time) error
	GetLongURL(slug string) (string, error)
	GetList() ([]ShortenObject, error)
	// TODO: add tags / who added
}

// getOrDefault looks for values in the environment and defaults to the provided value
// if it is not found.
func getOrDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	log.Println(kayvee.FormatLog("shorty", kayvee.Info, "configuration", msg{
		"msg": fmt.Sprintf("No value found for '%s', defaulting to '%s'", key, def),
	}))
	return def
}
