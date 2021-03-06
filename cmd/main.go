package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"

	"github.com/research-pal/backend/api"
	"github.com/research-pal/backend/db"
)

func main() {
	logmode := flag.String("logmode", "", "log mode: DEBUG, ERROR, INFO")
	flag.Parse()
	if strings.EqualFold(*logmode, "DEBUG") {
		// add file path and line number to logs
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("logging in debug mode")
	}

	dbClient := db.NewDBClient()
	defer dbClient.Close()
	api.Init(dbClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting port to %s", port)
	}

	router := mux.NewRouter()

	router.HandleFunc("/notes", jsonHeaders(api.HandleNotesGetFiltered)).Methods("GET")
	router.HandleFunc("/notes/{id}", jsonHeaders(api.HandleNotesGetByID)).Methods("GET")
	router.HandleFunc("/notes/{id}", jsonHeaders(api.HandleNotesPut)).Methods("PUT")
	router.HandleFunc("/notes/{id}", jsonHeaders(api.HandleNotesPatch)).Methods("PATCH")
	router.HandleFunc("/notes/{id}", jsonHeaders(api.HandleNotesDelete)).Methods("DELETE")
	router.HandleFunc("/notes", jsonHeaders(api.HandleNotesPost)).Methods("POST")

	log.Printf("Listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Println("Listening error :", err)
	}
}

func jsonHeaders(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler(w, r)
	}
}
