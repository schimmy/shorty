package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/Clever/shorty/db"
	"github.com/Clever/shorty/routes"
	"github.com/gorilla/mux"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8040"
	}

	// TODO: different backends based on config
	db := db.NewPostgresDB()

	r := mux.NewRouter()
	r.HandleFunc("/delete", routes.DeleteHandler(db)).Methods("POST")
	r.HandleFunc("/shorten", routes.ShortenHandler(db)).Methods("POST")
	r.HandleFunc("/list", routes.ListHandler(db)).Methods("GET")
	r.PathPrefix("/Shortener.jsx").Handler(http.FileServer(http.Dir("./static")))
	r.PathPrefix("/favicon.png").Handler(http.FileServer(http.Dir("./static")))
	r.HandleFunc("/{slug}", routes.RedirectHandler(db)).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		http.ServeFile(w, r, "./static/index.html")
	}).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	http.Handle("/", r)

	fmt.Printf("Starting server on port: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
