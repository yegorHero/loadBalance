package main

import (
	"loadBalance/config"
	"loadBalance/internal/algorithm"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

const configPath = "./config.yaml"

func main() {
	cfg := config.Init(configPath)

	alg := algorithm.Init(cfg.AlgorithmType, cfg.BackendAddresses)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Printf("Malformed address: %s", r.RemoteAddr)
		}

		log.Printf("Client IP: %s %s %s", ip, r.Method, r.URL.Path)

		target, err := alg.GetNextServer()
		if err != nil {
			log.Printf("alg.GetNextServer: %v", err)
			http.Error(w, "No backend server are alive", http.StatusServiceUnavailable)
			return
		}

		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = target.Scheme
				req.URL.Host = target.Host
			},
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				http.Error(w, "Backend unavailable: "+err.Error(), http.StatusServiceUnavailable)
			},
		}
		proxy.ServeHTTP(w, r)
	})

	log.Println("Started server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
