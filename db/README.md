# db
--
    import "github.com/Clever/shorty/db"


## Usage

#### type ErrNotFound

```go
type ErrNotFound struct{}
```


#### func (ErrNotFound) Error

```go
func (e ErrNotFound) Error() string
```

#### type PostgresDB

```go
type PostgresDB struct {
	SchemaName string
	TableName  string
}
```

PostgresDB represents a connection to a Postgres database.

#### func (*PostgresDB) DeleteURL

```go
func (pgDB *PostgresDB) DeleteURL(slug string) error
```
DeleteURL removes a URL from the database.

#### func (*PostgresDB) GetList

```go
func (pgDB *PostgresDB) GetList() ([]ShortenObject, error)
```
GetList lists all shortened URLs.

#### func (*PostgresDB) GetLongURL

```go
func (pgDB *PostgresDB) GetLongURL(slug string) (string, error)
```
GetLongURL searches for the short URL reference in order to return the long url.

#### func (*PostgresDB) ShortenURL

```go
func (pgDB *PostgresDB) ShortenURL(slug, longURL, owner string, expires time.Time) error
```
ShortenURL creates a new record for a shortened URL.

#### type Redis

```go
type Redis struct {
}
```

Redis is a wrapper over a Redis pool connection.

#### func (Redis) DeleteURL

```go
func (r Redis) DeleteURL(slug string) error
```
DeleteURL removes all reference to a URL.

#### func (Redis) GetList

```go
func (r Redis) GetList() ([]ShortenObject, error)
```
GetList lists all shortened URLs.

#### func (Redis) GetLongURL

```go
func (r Redis) GetLongURL(slug string) (string, error)
```
GetLongURL searches for the short URL reference in order to return the long url.

#### func (Redis) ShortenURL

```go
func (r Redis) ShortenURL(slug, longURL, owner string, expires time.Time) error
```
ShortenURL adds a URL object to the db.

#### type ShortenBackend

```go
type ShortenBackend interface {
	DeleteURL(slug string) error
	ShortenURL(slug, longURL, owner string, expires time.Time) error
	GetLongURL(slug string) (string, error)
	GetList() ([]ShortenObject, error)
}
```

ShortenBackend represents the necessary interface for storing and updating URLs.

#### func  NewPostgresDB

```go
func NewPostgresDB() ShortenBackend
```
NewPostgresDB configures and connects to a postgres database.

#### func  NewRedisDB

```go
func NewRedisDB() ShortenBackend
```
NewRedsDB connects to Redis and pools connections.

#### type ShortenObject

```go
type ShortenObject struct {
	Slug     string    `json:"slug"`
	Owner    string    `json:"owner"`
	LongURL  string    `json:"long_url"`
	Modified time.Time `json:"modified_date,omitempty"`
	Expires  time.Time `json:"expire_date,omitempty"`
}
```

ShortenObject holds the metadata and the mapping for a shortened URL I like the
NullTime concept from the pq library, so even for other backends let's use it
instead of checking whether the date is 0001-01-01
