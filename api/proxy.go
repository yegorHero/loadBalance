package api

import (
	"loadBalance/internal/algorithm"
	"log"
	"net/http"
	"net/http/httputil"
)

type ProxyHandler struct {
	selector algorithm.ServerSelector
}

func NewProxyHandler(selector algorithm.ServerSelector) *ProxyHandler {
	return &ProxyHandler{
		selector: selector,
	}
}

// инплементация интерфейса хэндлера с проксированием поступающего запроса
func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target, err := p.selector.GetNextServer()
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
}
