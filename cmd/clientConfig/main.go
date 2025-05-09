package main

import (
	"loadBalance/api"
	"log"
	"net/http"
)

func main() {
	userConfig := api.UserConfig{}
	mainMux := http.NewServeMux()

	configMux := http.NewServeMux()
	configMux.HandleFunc("/create", userConfig.Create)
	configMux.HandleFunc("/read", userConfig.Read)
	configMux.HandleFunc("/update", userConfig.Update)
	configMux.HandleFunc("/delete", userConfig.Delete)

	mainMux.Handle("/api/config/", http.StripPrefix("/api/config", configMux))

	log.Printf("Start server at :9090")
	log.Fatal(http.ListenAndServe(":9090", mainMux))
}
