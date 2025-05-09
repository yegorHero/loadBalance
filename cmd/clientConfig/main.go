package main

import (
	"log"
	"net/http"
)

func main() {
	mainMux := http.NewServeMux()

	configMux := http.NewServeMux()
	configMux.HandleFunc("/create", nil)
	configMux.HandleFunc("/read", nil)
	configMux.HandleFunc("/delete", nil)
	configMux.HandleFunc("/update", nil)

	mainMux.Handle("/api/config/", http.StripPrefix("/api/config", configMux))

	log.Printf("Start server at :9090")
	log.Fatal(http.ListenAndServe(":9090", mainMux))
}
