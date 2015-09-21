package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/schimmy/shorty/db"
	"github.com/schimmy/shorty/routes"
)

const (
	pgBackend    = "postgres"
	redisBackend = "redis"
)

var (
	port     = flag.String("port", "80", "port to listen on")
	database = flag.String("db", pgBackend, "datastore option to use, one of: ['postgres', 'redis']")
	protocol = flag.String("protocol", "http", "protocol for the short handler - useful to separate for external-facing separate instance")
	domain   = flag.String("domain", "go", "set the domain for the short URL reported to the user")
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
}

func main() {
	var sdb db.ShortenBackend
	switch *database {
	case pgBackend:
		sdb = db.NewPostgresDB()
	case redisBackend:
		sdb = db.NewRedisDB()
	default:
		log.Fatalf("'%s' backend is not offered", *database)
	}

	r := mux.NewRouter()
	r.HandleFunc("/delete", routes.DeleteHandler(sdb)).Methods("POST")
	r.HandleFunc("/shorten", routes.ShortenHandler(sdb)).Methods("POST")
	r.HandleFunc("/list", routes.ListHandler(sdb)).Methods("GET")
	r.HandleFunc("/meta", routes.MetaHandler(*protocol, *domain)).Methods("GET")
	r.PathPrefix("/Shortener.jsx").Handler(http.FileServer(http.Dir("./static")))
	r.PathPrefix("/favicon.png").Handler(http.FileServer(http.Dir("./static")))
	r.HandleFunc("/{slug}", routes.RedirectHandler(sdb, *domain)).Methods("GET")
	r.HandleFunc("/health/check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "STATUS OK")
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		http.ServeFile(w, r, "./static/index.html")
	}).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	http.Handle("/", r)

	fmt.Printf("Starting server on port: %s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
