package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gopkg.in/Clever/kayvee-go.v2"

	// this is standard procedure for registering a drive with db/sql
	_ "github.com/lib/pq"
)

const (
	connStr = "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"
)

// PostgresDB represents a connection to a Postgres database.
type PostgresDB struct {
	c          *sql.DB
	SchemaName string
	TableName  string
}

// NewPostgresDB configures and connects to a postgres database.
func NewPostgresDB() ShortenBackend {
	pgHost := getOrDefault("PG_HOST", "localhost")
	pgPort := getOrDefault("PG_PORT", "5432")
	pgUser := getOrDefault("PG_USER", "shortener")
	pgPass := getOrDefault("PG_PASSWORD", "NOPE")
	pgDatabase := getOrDefault("PG_DATABASE", "shortener")
	pgSchema := getOrDefault("PG_SCHEMA", "shortener")
	pgTable := getOrDefault("PG_TABLE", "shortener")
	pgSSLMode := getOrDefault("PG_SSL", "disable")

	connString := fmt.Sprintf(connStr, pgHost, pgPort, pgUser, pgPass, pgDatabase, pgSSLMode)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %s", err)
	}

	return &PostgresDB{
		c:          db,
		SchemaName: pgSchema,
		TableName:  pgTable,
	}
}

// DeleteURL removes a URL from the database.
func (pgDB *PostgresDB) DeleteURL(slug string) error {
	_, err := pgDB.c.Query(fmt.Sprintf("DELETE FROM %s.%s WHERE slug=$1", pgDB.SchemaName, pgDB.TableName), slug)
	if err != nil {
		return err
	}
	log.Println(kayvee.FormatLog("shorty", kayvee.Info, "slug.delete", msg{
		"name": slug,
	}))
	return nil

}

// ShortenURL creates a new record for a shortened URL.
func (pgDB *PostgresDB) ShortenURL(slug, longURL, owner string, expires time.Time) error {
	// postgres doesn't have an upsert method yet (coming in 9.5)
	existingLong, err := pgDB.GetLongURL(slug)
	if existingLong == "" || err != nil {
		q := fmt.Sprintf("INSERT INTO %s.%s(slug, long_url, owner) VALUES($1, $2, $3)", pgDB.SchemaName, pgDB.TableName)
		_, err := pgDB.c.Query(q, slug, longURL, owner)
		if err != nil {
			return fmt.Errorf("Issue inserting new row for slug: %s, err is: %s", slug, err)
		}
		log.Println(kayvee.FormatLog("shorty", kayvee.Info, "slug.new", msg{
			"name":     slug,
			"long_url": longURL,
			"owner":    owner,
		}))
		return nil
	}

	// Otherwise, upsert
	q := fmt.Sprintf("UPDATE %s.%s SET long_url=$2, owner=$3 WHERE slug=$1", pgDB.SchemaName, pgDB.TableName)
	_, err = pgDB.c.Query(q, slug, longURL, owner)
	if err != nil {
		return err
	}
	log.Println(kayvee.FormatLog("shorty", kayvee.Info, "slug.update", msg{
		"name":     slug,
		"long_url": longURL,
		"owner":    owner,
	}))
	return nil
}

// GetLongURL searches for the short URL reference in order to return the long url.
func (pgDB *PostgresDB) GetLongURL(slug string) (string, error) {
	q := fmt.Sprintf("SELECT long_url FROM %s.%s WHERE slug = $1", pgDB.SchemaName, pgDB.TableName)
	var longURL string
	err := pgDB.c.QueryRow(q, slug).Scan(&longURL)
	if err != nil {
		return "", err
	}
	return longURL, nil
}

// GetList lists all shortened URLs.
func (pgDB *PostgresDB) GetList() ([]ShortenObject, error) {
	rows, err := pgDB.c.Query(fmt.Sprintf("SELECT slug, long_url, owner FROM %s.%s", pgDB.SchemaName, pgDB.TableName))
	if err != nil {
		return nil, err
	}
	var retObjs []ShortenObject
	for rows.Next() {
		var so ShortenObject
		err = rows.Scan(&so.Slug, &so.LongURL, &so.Owner)
		if err != nil {
			return nil, fmt.Errorf("issue scanning row for list: %s", err)
		}
		retObjs = append(retObjs, so)
	}
	return retObjs, nil
}
