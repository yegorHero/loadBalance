package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var servers = []string{
	"8081",
	"8082",
	"8083",
}

func main() {

	var wg sync.WaitGroup

	for _, port := range servers {
		wg.Add(1)
		go func(port string) {
			defer wg.Done()

			mux := http.NewServeMux()

			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				hostname, _ := os.Hostname()
				fmt.Fprintf(w, "Backend server on port %s, host: %s, Request path: %s\n", port, hostname, r.URL.Path)
			})

			log.Printf("Backend started at :%s\n", port)
			if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
				log.Fatal(err)
			}
		}(port)
	}

	wg.Wait()
}
