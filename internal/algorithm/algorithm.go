package algorithm

import (
	"log"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

type ServerSelector interface {
	GetNextServer() (*url.URL, error)
}

type Server struct {
	URL   *url.URL
	Alive atomic.Bool
}

func Init(algorithmType string, addrServers []string) ServerSelector {
	servers := make([]*Server, 0, len(addrServers))
	for _, addr := range addrServers {
		u, err := url.Parse(addr)
		if err != nil {
			log.Fatalf("error parse url server: %v", err)
		}

		servers = append(servers, &Server{
			URL: u,
		})
	}

	startHealthCheck(servers, 5*time.Second)

	switch algorithmType {
	case "round-robin":
		return InitRoundRobin(servers)
	case "weighted-round-robin":
	case "least-response-time":
	case "resource-based":

	}

	return InitRoundRobin(servers)
}

func (s *Server) CheckHealth(timeout time.Duration) {
	client := http.Client{Timeout: timeout}
	resp, err := client.Get(s.URL.String())

	if err != nil || resp.StatusCode != http.StatusOK {
		s.Alive.Store(false)
		if err != nil {
			log.Printf("server %s i DOWN: %v", s.URL, err)
		} else {
			log.Printf("Server %s returned: %d", s.URL, resp.StatusCode)
		}
		return
	}

	resp.Body.Close()
	s.Alive.Store(true)
}

func startHealthCheck(servers []*Server, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			for _, s := range servers {
				s.CheckHealth(2 * time.Second)
			}
		}
	}()
}
