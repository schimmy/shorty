package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Clever/pretty-self-hosted-url-shortener/db"
	"github.com/gorilla/mux"
)

func returnJson(data interface{}, err error, w http.ResponseWriter, r *http.Request) {
	retCode := 200
	var content interface{}
	if err != nil {
		log.Printf("500 Err: %s", err)
		data = map[string]interface{}{"error": err.Error()}
		retCode = 500
	}
	content, err = json.Marshal(data)
	if err != nil {
		log.Fatalf("error marshalling json: %s", err.Error())
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(retCode)
	w.Write(content.([]byte))
	return
}

func ShortenHandler(db db.ShortenBackend) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			returnJson(nil, fmt.Errorf("couldn't parse form: %s", err.Error()), w, r)
			return
		}
		slug := r.PostForm.Get("slug")
		longURL := r.PostForm.Get("long_url")

		// for now set expiry to never
		var t time.Time
		err := db.ShortenURL(slug, longURL, t)
		returnJson(nil, err, w, r)
		return
	}
}

func DeleteHandler(db db.ShortenBackend) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			returnJson(nil, fmt.Errorf("couldn't parse form: %s", err.Error()), w, r)
			return
		}
		slug := r.PostForm.Get("slug")
		err := db.DeleteURL(slug)
		returnJson(nil, err, w, r)
	}
}
func ListHandler(db db.ShortenBackend) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sObj, err := db.GetList()
		returnJson(sObj, err, w, r)
	}
}

// Not accessed via Ajax, just by end users
func RedirectHandler(db db.ShortenBackend) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := mux.Vars(r)["slug"]
		log.Println("slug: ", slug)
		long, err := db.GetLongURL(slug)
		if err != nil {
			log.Printf("Error in redirect: %s", err)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		http.Redirect(w, r, long, 301)
		return
	}
}
