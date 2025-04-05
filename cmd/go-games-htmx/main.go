package main

import (
	"go-games-htmx/api"
	sqlite "go-games-htmx/database"
	"go-games-htmx/handlers"
	"log"
	"net/http"
	"time"
)

func main() {
	// Check if API key is defined
	err := api.CheckApiKey()
	if err != nil {
		log.Fatal(err)
	}

	// Init the database once
	err = sqlite.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Assets
	http.Handle("GET /assets/", http.FileServer(http.Dir(".")))

	// Pages
	http.HandleFunc("GET /{$}", handlers.HandleGETIndex)
	http.HandleFunc("GET /search", handlers.HandleGETSearch)
	http.HandleFunc("GET /game", handlers.HandleGETGame)

	// 404 Not Found
	http.HandleFunc("GET /", handlers.Handle404)

	s := &http.Server{
		Addr:                         ":80",
		Handler:                      http.DefaultServeMux,
		DisableGeneralOptionsHandler: true,
		ReadTimeout:                  10 * time.Second,
		ReadHeaderTimeout:            10 * time.Second,
		WriteTimeout:                 10 * time.Second,
		IdleTimeout:                  10 * time.Second,
		MaxHeaderBytes:               http.DefaultMaxHeaderBytes,
		ErrorLog:                     log.Default(),
	}

	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
