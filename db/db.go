package db

import (
	"os"
	"time"
)

// ShortenObject holds the metadata and the mapping for a shortened URL
// I like the NullTime concept from the pq library, so even for other backends
// let's use it instead of checking whether the date is 0001-01-01
type ShortenObject struct {
	Slug     string    `json:"slug",sql:"slug"`
	Owner    string    `json:"owner",sql:"owner"`
	LongURL  string    `json:"long_url",sql:"long_url"`
	Modified time.Time `json:"modified_date",omitempty`
	Expires  time.Time `json:"expire_date",omitempty`
}

// TODO: add tags / who added
// TODO: add delete URL
type ShortenBackend interface {
	DeleteURL(slug string) error
	ShortenURL(slug, longURL, owner string, expires time.Time) error
	GetLongURL(slug string) (string, error)
	GetList() ([]ShortenObject, error)
}

func GetOrDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
