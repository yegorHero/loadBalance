package algorithm

import (
	"fmt"
	"log"
	"net/url"
	"sync/atomic"
)

type RoundRobin struct {
	servers []*Server
	current atomic.Int32
}

func InitRoundRobin(servers []*Server) *RoundRobin {
	return &RoundRobin{
		servers: servers,
	}
}

func (r *RoundRobin) GetNextServer() (*url.URL, error) {
	index := r.current.Add(1) - 1
	currentIndex := int(index) % len(r.servers)

	for i := 0; i < len(r.servers); i++ {
		idx := (currentIndex + i) % len(r.servers)
		if r.servers[idx].Alive.Load() {
			return r.servers[idx].URL, nil
		}

		log.Printf("server port not alive: %s", r.servers[idx].URL)
	}

	return nil, fmt.Errorf("not alive servers")
}
