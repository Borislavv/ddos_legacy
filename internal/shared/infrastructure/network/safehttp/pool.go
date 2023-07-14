package safehttp

import (
	"errors"
	"sync"
	"time"
)

type Pool struct {
	mu      *sync.Mutex
	clients []*Client
	timeout time.Duration
}

func NewPool(cap int64, timeout time.Duration) *Pool {
	mu := &sync.Mutex{}

	var clients []*Client
	for i := int64(0); i < cap; i++ {
		clients = append(clients, NewClient(mu))
	}

	return &Pool{
		mu:      &sync.Mutex{},
		clients: clients,
		timeout: timeout,
	}
}

func (p *Pool) GetAvailable() (*Client, error) {
	defer p.mu.Unlock()
	p.mu.Lock()

	start := time.Now()
	for {
		for _, client := range p.clients {
			if client.isAvailable {
				return client, nil
			}

			if time.Since(start) >= p.timeout {
				return nil, errors.New("timeout exceeded, all http clients are unreachable")
			}
		}
	}
}
