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

func TestOverwrite(t *testing.T) {
	fillDB(t)
	oldCat := testData["cat4"]

	newCat := ShortenObject{
		Slug:    oldCat.Slug,
		Owner:   "yodie",
		LongURL: "http://imgur.com/t/dogs/8IrIs", // todo find actua dog :-P
	}
	// this currently should overwrite the existing slug
	err := datastore.ShortenURL(newCat.Slug, newCat.LongURL, newCat.Owner, time.Now())
	assert.Nil(t, err, "Shorten to overwrite failed for slug:long '%s':%s, datastore type: %s, err: %s",
		newCat.Slug, newCat.LongURL, datastoreType, err)

	// ensure the old one is gone, new one exists
	l, err := datastore.GetList()
	assert.Nil(t, err, "getlist should not err out: %s", err)
	for _, u := range l {
		if u.Slug == newCat.Slug {
			assert.Equal(t, newCat.LongURL, u.LongURL)
			assert.False(t, newCat.LongURL == oldCat.LongURL, fmt.Sprintf("Old longURL persists for slug: %s", newCat.Slug))
			assert.False(t, newCat.Owner == oldCat.Owner, fmt.Sprintf("Old Owner persists for slug: %s", newCat.Slug))
		}
	}

	// Attempt to get long as well. Somewhat uneccessary, but a precaution against caching, etc
	long, err := datastore.GetLongURL(newCat.Slug)
	assert.Nil(t, err, "Failed to get record '%s' with '%s': %s", newCat.Slug, datastoreType, err)
	assert.Equal(t, newCat.LongURL, long, "Record %s on %s", long, datastoreType)
}

func TestDelete(t *testing.T) {
	fillDB(t)
	s := "cat4"
	cat := testData[s]

	long, err := datastore.GetLongURL(s)
	assert.Nil(t, err, "Failed to get record '%s' with '%s': %s", s, datastoreType, err)
	assert.Equal(t, cat.LongURL, long, "Record %s on %s", long, datastoreType)

	err = datastore.DeleteURL(s)
	assert.Nil(t, err, "Failed to delete '%s' with '%s': %s", s, datastoreType, err)

	newLong, err := datastore.GetLongURL(s)
	assert.Equal(t, ErrNotFound, err,
		fmt.Sprintf("Record %s on %s should have returned not found, instead found: %s with err %s",
			s, datastoreType, newLong, err))
}
