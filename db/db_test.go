package db

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	datastore     ShortenBackend
	datastoreType string

	testData = map[string]ShortenObject{
		"cat1": {
			Slug:    "cat1",
			Owner:   "garfield",
			LongURL: "http://i.imgur.com/D30DR9V.png",
		},
		"cat2": {
			Slug:    "cat2",
			Owner:   "john",
			LongURL: "http://imgur.com/t/cats/H4YWLta",
		},
		"cat3": {
			Slug:    "cat3",
			Owner:   "nate",
			LongURL: "http://imgur.com/t/cats/Uzabwxi",
		},
		"cat4": {
			Slug:    "cat4",
			Owner:   "schimmy",
			LongURL: "http://imgur.com/t/cats/8IrIs",
		},
	}
)

func clearDB() {
	switch datastoreType {
	case "redis":
		c := datastore.(Redis).p.Get()
		defer c.Close()
		_, err := c.Do("FLUSHDB")
		if err != nil {
			log.Fatalf("Failed to clear redis db: %s", err)
		}
	case "postgres":
		pg := datastore.(*PostgresDB)
		_, err := pg.c.Exec(fmt.Sprintf("DELETE FROM %s.%s", pg.SchemaName, pg.TableName))
		if err != nil {
			log.Fatalf("Failed to clear postgres db: %s", err)
		}
	}
}

func fillDB(t *testing.T) {
	clearDB()
	for i, u := range testData {
		err := datastore.ShortenURL(u.Slug, u.LongURL, u.Owner, time.Time{})
		assert.Nil(t, err, "Failed to insert record '%s' to '%s': %s", i, datastoreType, err)
	}
}

func TestMain(m *testing.M) {
	datastoreType = "redis"
	datastore = NewRedisDB()
	redisResult := m.Run()

	datastoreType = "postgres"
	datastore = NewPostgresDB()
	postgresResult := m.Run()

	os.Exit(redisResult | postgresResult)
}

func TestRetrieve(t *testing.T) {
	fillDB(t)

	for i, u := range testData {
		long, err := datastore.GetLongURL(u.Slug)
		assert.Nil(t, err, "Failed to get record %d ('%s') with '%s': %s",
			i, u.Slug, datastoreType, err)
		assert.Equal(t, u.LongURL, long, "Record %d on %s", i, datastoreType)
	}
}

func TestGetList(t *testing.T) {
	fillDB(t)

	l, err := datastore.GetList()
	assert.Nil(t, err, "getlist should not err out: %s", err)
	for _, u := range l {
		assert.Equal(t, testData[u.Slug], u)
	}
}
