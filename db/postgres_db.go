package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	c          *sql.DB
	SchemaName string
	TableName  string
}

func NewPostgresDB() db.ShortenBackend {
	pgHost := db.GetOrDefault("PG_HOST", "localhost")
	pgPort := db.GetOrDefault("PG_HOST", "5432")
	pgUser := db.GetOrDefault("PG_USER", "shortener")
	pgPass := db.GetOrDefault("PG_PASS", "NOPE")
	pgDatabase := db.GetOrDefault("PG_DB", "shortener")
	pgSchema := db.GetOrDefault("PG_SCHEMA", "shortener")
	pgTable := db.GetOrDefault("PG_TABLE", "shortener")
	pgSSLMode := db.GetOrDefault("PG_SSL", "disable")

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", pgHost, pgPort, pgUser, pgPass, pgDatabase, pgSSLMode)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	return &PostgresDB{
		c:          db,
		SchemaName: pgSchema,
		TableName:  pgTable,
	}
}

func (pgDB *PostgresDB) DeleteURL(slug string) error {
	_, err := pgDB.c.Query(fmt.Sprintf("DELETE FROM %s.%s WHERE slug=$1", pgDB.SchemaName, pgDB.TableName), slug)
	if err != nil {
		return err
	}
	log.Printf("Successfully deleted slug: %s", slug)
	return nil

}
func (pgDB *PostgresDB) ShortenURL(slug, longURL, owner string, tags []string, expires time.Time) error {
	// postgres & redshift don't have an upsert method yet
	existingLong, err := pgDB.GetLongURL(slug)
	if existingLong == "" || err != nil { // TODO figure out what happens on nothing, err?
		//q := fmt.Sprintf("INSERT INTO %s.%s(slug, long_url, expires, modified) VALUES($1, $2, $3, $4)")
		q := fmt.Sprintf("INSERT INTO %s.%s(slug, long_url, owner) VALUES($1, $2, $3)", pgDB.SchemaName, pgDB.TableName)
		_, err := pgDB.c.Query(q, slug, longURL, owner)
		if err != nil {
			return fmt.Errorf("Issue inserting new row for slug: %s, err is: %s", slug, err)
		}
		return nil
	}
	// Otherwise, upsert
	q := fmt.Sprintf("UPDATE %s.%s SET long_url=$2, owner=$3 WHERE slug=$1", pgDB.SchemaName, pgDB.TableName)
	_, err = pgDB.c.Query(q, slug, longURL, owner)
	if err != nil {
		return err
	}
	log.Printf("Successfully updated slug: %s", slug)
	return nil
}

func (pgDB *PostgresDB) GetLongURL(slug string) (string, error) {
	//var retObj ShortenObject
	q := fmt.Sprintf("SELECT long_url FROM %s.%s WHERE slug = $1", pgDB.SchemaName, pgDB.TableName)
	var long_url string
	err := pgDB.c.QueryRow(q, slug).Scan(&long_url)
	if err != nil {
		return "", err
	}
	log.Println("long: ", long_url)
	return long_url, nil

}

func (pgDB *PostgresDB) GetList() ([]db.ShortenObject, error) {
	rows, err := pgDB.c.Query(fmt.Sprintf("SELECT slug, long_url, owner FROM %s.%s", pgDB.SchemaName, pgDB.TableName))
	if err != nil {
		return nil, err
	}
	var retObjs []db.ShortenObject
	var slug string
	var long_url string
	var owner string
	//tags :=  []string{}
	//var expires time.Time
	//var modified time.Time
	for rows.Next() {
		err = rows.Scan(&slug, &long_url, &owner) //, &modified, &expires)
		if err != nil {
			return nil, fmt.Errorf("issue scanning row for list: %s", err)
		}
		retObjs = append(retObjs, db.ShortenObject{
			Slug:    slug,
			LongURL: long_url,
			Owner:   owner,
			//Expires:  expires,
			//Modified: modified,
		})
	}
	return retObjs, nil
}
