package websocket

import (
	"sync"

	"github.com/hoangminhphuc/goph-chat/common/logger"
)

type Pool struct {
	mu      sync.Mutex
	Clients map[string]*Client
	Register chan *Client
	Unregister chan *Client
	Broadcast chan Message
	logger logger.ZapLogger
}

func NewPool() *Pool {
	return &Pool{
		Clients: make(map[string]*Client),
		Register: make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast: make(chan Message),
		logger: logger.NewZapLogger(),
	}
}

func (p *Pool) Start() {
	for {
		select {
		case client := <- p.Register:
			oldClient, ok := p.Clients[client.ID]

			if ok {
				oldClient.Connection.Close()
			}

			p.logger.Log.Info("New client. Size of connection pool: ", len(p.Clients))

			p.Clients[client.ID] = client
			
			for _, cl := range p.Clients {
				err := cl.Connection.WriteJSON(Message{ ChatUser: client.ID, Body: "New user joined."})
				if err != nil {
					p.logger.Log.Error("Cannot write to client: ", client.ID)
				}
			}

		case client := <- p.Unregister:
			delete(p.Clients, client.ID)	
			p.logger.Log.Info("Client left. Size of connection pool: ", len(p.Clients))
			for _, cl := range p.Clients {
				err := cl.Connection.WriteJSON(Message{ChatUser: client.ID, Body: "User left."})
				if err != nil {
					p.logger.Log.Error("Cannot write to client: ", client.ID)
				}
			}

		case msg := <- p.Broadcast:
			p.logger.Log.Info("Broadcasting message to clients.")
			for _,cl := range p.Clients {
				err := cl.Connection.WriteJSON(msg.Body)
				if err != nil {
					p.logger.Log.Error("Cannot write to client: ", cl.ID)
				}
			}
		}
	}
}

// func (p *Pool) Add(client *Client) {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// 	p.clients[client.ID] = client
// }

// func (p *Pool) Remove(id string) {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// 	delete(p.clients, id)
// }

// func (p *Pool) Broadcast(uid string, message []byte) {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// 	for id, cl := range p.clients {
// 		if id != uid {
// 			cl.Send <- message
// 		}
// 	}
// }