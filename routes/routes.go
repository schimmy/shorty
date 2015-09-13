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
	errMsgTemplate = "<html>Unable to find long URL for slug: <em>'%s'</em>, please consult <a href='http://go/'>http://go</a> to add this slug.</html>"
)

var (
	reserved = []string{"delete", "shorten", "list", "Shortener.jsx", "favicon.png"}
)

// msg is a convenience type for kayvee
type msg map[string]interface{}

// returnJSON
func returnJSON(data interface{}, err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		log.Println(kayvee.FormatLog("shorty", kayvee.Error, "internal.error", msg{
			"err": err.Error(),
		}))
		data = msg{"error": err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
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

//ShortenHandler generates a HTTP handler for a shortening URL's given a datastore backend.
func ShortenHandler(db db.ShortenBackend) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			returnJSON(nil, fmt.Errorf("couldn't parse form: %s", err.Error()), w, r)
			return
		}
		slug := r.PostForm.Get("slug")
		longURL := r.PostForm.Get("long_url")
		owner := r.PostForm.Get("owner")

		for _, reserved := range reserved {
			if slug == reserved {
				returnJSON(nil, fmt.Errorf("That slug is reserved: %s", slug), w, r)
				return
			}
		}

		// for now set expiry to never
		var t time.Time
		err := db.ShortenURL(slug, longURL, owner, t)
		returnJSON(nil, err, w, r)
		return
	}
}

// DeleteHandler generates a HTTP handler for deleting URL's given a datastore backend.
func DeleteHandler(db db.ShortenBackend) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			returnJSON(nil, fmt.Errorf("couldn't parse form: %s", err.Error()), w, r)
			return
		}
		slug := r.PostForm.Get("slug")
		err := db.DeleteURL(slug)
		returnJSON(nil, err, w, r)
	}
}

// ListHandler generates a HTTP handler for listing URL's given a datastore backend.
func ListHandler(db db.ShortenBackend) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sObj, err := db.GetList()
		returnJSON(sObj, err, w, r)
	}
}

// RedirectHandler redirects users to their desired location.
// Not accessed via Ajax, just by end users
func RedirectHandler(db db.ShortenBackend) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := mux.Vars(r)["slug"]
		long, err := db.GetLongURL(slug)
		if long == "" {
			w.WriteHeader(404)
			fmt.Fprintf(w, fmt.Sprintf(errMsgTemplate, slug))
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
