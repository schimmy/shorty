package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/schimmy/shorty/db"
	"gopkg.in/Clever/kayvee-go.v2"
)

const (
	errMsgTemplate  = "<html>Unable to find long URL for slug: <em>'%s'</em>, please consult <a href='http://%s/'>http://%s</a> or the proper admin instance to add this slug.</html>"
	readOnlyMessage = "You have reached the read-only instance of this shortener. If you have access, please visit the admin instance to perform write operations."
)

var (
	reserved = []string{"delete", "shorten", "list", "meta", "Shortener.jsx", "favicon.png"}
)

// msg is a convenience type for kayvee
type msg map[string]interface{}

// returnJSON
func returnJSON(data interface{}, err error, errNum int, w http.ResponseWriter) {
	if err != nil {
		log.Println(kayvee.FormatLog("shorty", kayvee.Error, "internal.error", msg{
			"err": err.Error(),
		}))
		data = msg{"error": err.Error()}
		// if errNum not set
		if errNum == 0 {
			errNum = http.StatusInternalServerError
		}
		w.WriteHeader(errNum)
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println(kayvee.FormatLog("shorty", kayvee.Error, "json.encoding", msg{
			"err": err.Error(),
		}))
	}

	return
}

// ReadOnlyHandler returns a message to the user saying that they need
// to find the write-domain to modify things. This shouldn't be hit too
// often as the admin interface itself will hide the user-facing input
func ReadOnlyHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("In read-only mode, returning 401")
		returnJSON(nil, fmt.Errorf(readOnlyMessage), 401, w)
		return
	}
}

//ShortenHandler generates a HTTP handler for a shortening URLs given a datastore backend.
func ShortenHandler(db db.ShortenBackend) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			returnJSON(nil, fmt.Errorf("couldn't parse form: %s", err.Error()), 500, w)
			return
		}
		slug := r.PostForm.Get("slug")
		longURL := r.PostForm.Get("long_url")
		owner := r.PostForm.Get("owner")

		for _, reserved := range reserved {
			if slug == reserved {
				returnJSON(nil, fmt.Errorf("That slug is reserved: %s", slug), 400, w)
				return
			}
		}

		// for now set expiry to never
		var t time.Time
		err := db.ShortenURL(slug, longURL, owner, t)
		returnJSON(nil, err, 0, w)
		return
	}
}

// DeleteHandler generates a HTTP handler for deleting URL's given a datastore backend.
func DeleteHandler(db db.ShortenBackend) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			returnJSON(nil, fmt.Errorf("couldn't parse form: %s", err.Error()), 500, w)
			return
		}
		slug := r.PostForm.Get("slug")
		err := db.DeleteURL(slug)
		returnJSON(nil, err, 0, w)
	}
}

// MetaHandler returns info on the protocol and the domain for our short URLs
// Seems a little weird, but follows react patterns
func MetaHandler(protocol, domain string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		retObj := map[string]string{
			"protocol": protocol,
			"domain":   domain,
		}
		returnJSON(retObj, nil, 0, w)
	}
}

// ListHandler generates a HTTP handler for listing URL's given a datastore backend.
func ListHandler(db db.ShortenBackend) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sObj, err := db.GetList()
		returnJSON(sObj, err, 0, w)
	}
}

// RedirectHandler redirects users to their desired location.
// Not accessed via Ajax, just by end users
func RedirectHandler(db db.ShortenBackend, domain string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := mux.Vars(r)["slug"]
		long, err := db.GetLongURL(slug)
		if long == "" {
			w.WriteHeader(404)
			fmt.Fprintf(w, fmt.Sprintf(errMsgTemplate, slug, domain, domain))
			return
		}
		if err != nil {
			log.Println(kayvee.FormatLog("shorty", kayvee.Error, "redirect", msg{
				"err": err.Error(),
			}))
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		http.Redirect(w, r, long, 301)
		return
	}
}
