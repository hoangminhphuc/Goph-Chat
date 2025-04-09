package websocket

import (
	"sync"
)

type Pool struct {
	mu      sync.Mutex
	clients map[string]*Client
}

func NewPool() *Pool {
	return &Pool{
		clients: make(map[string]*Client),
	}
}

func (p *Pool) Add(client *Client) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clients[client.ID] = client
}

func (p *Pool) Remove(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.clients, id)
}

func (p *Pool) Broadcast(uid string, message []byte) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for id, cl := range p.clients {
		if id != uid {
			cl.Send <- message
		}
	}
}