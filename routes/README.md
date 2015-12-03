# routes
--
    import "github.com/Clever/shorty/routes"


## Usage

#### func  DeleteHandler

```go
func DeleteHandler(db database.ShortenBackend) func(http.ResponseWriter, *http.Request)
```
DeleteHandler generates a HTTP handler for deleting URL's given a datastore
backend.

#### func  ListHandler

```go
func ListHandler(db database.ShortenBackend) func(http.ResponseWriter, *http.Request)
```
ListHandler generates a HTTP handler for listing URL's given a datastore
backend.

#### func  MetaHandler

```go
func MetaHandler(protocol, domain string) func(http.ResponseWriter, *http.Request)
```
MetaHandler returns info on the protocol and the domain for our short URLs Seems
a little weird, but follows react patterns

#### func  ReadOnlyHandler

```go
func ReadOnlyHandler() func(http.ResponseWriter, *http.Request)
```
ReadOnlyHandler returns a message to the user saying that they need to find the
write-domain to modify things. This shouldn't be hit too often as the admin
interface itself will hide the user-facing input

#### func  RedirectHandler

```go
func RedirectHandler(db database.ShortenBackend, domain string) func(http.ResponseWriter, *http.Request)
```
RedirectHandler redirects users to their desired location. Not accessed via
Ajax, just by end users

#### func  ShortenHandler

```go
func ShortenHandler(db database.ShortenBackend) func(http.ResponseWriter, *http.Request)
```
ShortenHandler generates a HTTP handler for a shortening URLs given a datastore
backend.
