package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/Clever/pretty-self-hosted-url-shortener/db"
	"github.com/Clever/pretty-self-hosted-url-shortener/routes"
	"github.com/gorilla/mux"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// TODO different backends based on config
	//db := db.NewRedisDB()
	db := db.NewPostgresDB()

	r := mux.NewRouter()
	r.HandleFunc("/shorten", routes.ShortenHandler(db)).Methods("POST")
	r.HandleFunc("/redirect/{slug}", routes.RedirectHandler(db)).Methods("GET")
	r.HandleFunc("/list", routes.ListHandler(db)).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		http.ServeFile(w, r, "./static/index.html")
	}).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	http.Handle("/", r)

	fmt.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
