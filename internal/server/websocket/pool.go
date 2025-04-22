package websocket

import (
	"fmt"
	"sync"

	"github.com/hoangminhphuc/goph-chat/common/logger"
)

// Each Pool is a Room
type Pool struct {
	RoomID 			int
	Clients 		map[int]*Client
	Register 		chan *Client
	Unregister 	chan *Client
	Broadcast 	chan Message
	done 				chan struct{}
	logger 			logger.ZapLogger
	mu 					sync.RWMutex
}

func NewPool(roomID int) *Pool {
	return &Pool{
		RoomID: 			roomID,
		Clients: 			make(map[int]*Client),
		Register: 		make(chan *Client, 10),
		Unregister: 	make(chan *Client, 10),
		Broadcast: 		make(chan Message, 50),
		done:       make(chan struct{}),
		logger: 			logger.NewZapLogger(),
	}
}

func (p *Pool) handleClientRegistration(client *Client) {
	p.mu.Lock()
	// defer p.mu.Unlock()
	oldClient, ok := p.Clients[client.ID]

	if ok {
		oldClient.Connection.Close()
		delete(p.Clients, client.ID) 
	}

	p.Clients[client.ID] = client
	p.logger.Log.Info(fmt.Sprintf("New client joins room %d. Size of connection pool: %d", p.RoomID, len(p.Clients)))

	peers := make([]*Client, 0, len(p.Clients)-1)
	for _, c := range p.Clients {
		if c.ID != client.ID {
			peers = append(peers, c)
		}
	}
	p.mu.Unlock()
	
	// Avoids WriteJSON to blocks other pool operation
	for _, peer := range peers {
		err := peer.Connection.WriteJSON(Message{
			RoomID:   p.RoomID,
			ChatUser: client.ID,
			Body:     "New user joined.",
		})
		if err != nil {
			p.logger.Log.Error("Cannot write to client: ", peer.ID)
		}
	}
}

func (p *Pool) handleClientUnregistration(client *Client) {
	p.mu.Lock()

	delete(p.Clients, client.ID)	
	p.logger.Log.Info("Client left. Size of connection pool: ", len(p.Clients))
	

	peers := make([]*Client, 0, len(p.Clients))
	for _, c := range p.Clients {
		peers = append(peers, c)
	}
	p.mu.Unlock()

	// Avoids WriteJSON to blocks other pool operation
	for _, peer := range peers {
		err := peer.Connection.WriteJSON(Message{
			RoomID:   p.RoomID,
			ChatUser: client.ID,
			Body:     "User left.",
		})
		if err != nil {
			p.logger.Log.Error("Cannot write to client: ", peer.ID)
		}
	}
}

func (p *Pool) broadcastMessage(msg Message) {
	p.mu.RLock()
	peers := make([]*Client, 0, len(p.Clients))
	for _, c := range p.Clients {
		if c.ID != msg.ChatUser {
			peers = append(peers, c)
		}
	}
	p.mu.RUnlock()
	
	p.logger.Log.Info("Broadcasting message to clients.")
	for _, peer := range peers {
		err := peer.Connection.WriteJSON(msg)
		if err != nil {
			p.logger.Log.Error("Cannot write to client: ", peer.ID)
		}
	}
}


func (p *Pool) Start() {
	defer func() {
		if r := recover(); r != nil {
			p.logger.Log.Fatal("pool panic: %v", r)
		}
	}()

	for {
		select {
		case client := <- p.Register:
			p.handleClientRegistration(client)
		case client := <- p.Unregister:
			p.handleClientUnregistration(client)
		case msg := <- p.Broadcast:
			p.broadcastMessage(msg)
		case <-p.done:
			// graceful shutdown
			return
		}
	}
}

func (p *Pool) Stop() {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, client := range p.Clients {
		_ = client.Connection.Close() // Close connection
	}

	close(p.done)
}

// func (p *Pool) ReviveWebsocket() {
// 	if err := recover(); err != nil {
// 			log.Println(
// 				"level", "error",
// 				"err", err,
// 			)
		
// 		go p.Start()
// 	}
// }
