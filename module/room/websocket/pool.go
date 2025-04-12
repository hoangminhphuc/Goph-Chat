package websocket

import (
	"log"
	"sync"

	"github.com/hoangminhphuc/goph-chat/common/logger"
)

// Each Pool is a Room
type Pool struct {
	RoomID 			string
	Clients 		map[string]*Client
	Register 		chan *Client
	Unregister 	chan *Client
	Broadcast 	chan Message
	logger 			logger.ZapLogger
	mu 					sync.RWMutex
}

func NewPool(roomID string) *Pool {
	return &Pool{
		RoomID: 			roomID,
		Clients: 			make(map[string]*Client),
		Register: 		make(chan *Client),
		Unregister: 	make(chan *Client),
		Broadcast: 		make(chan Message),
		logger: 			logger.NewZapLogger(),
	}
}

func (p *Pool) handleClientRegistration(client *Client) {
	p.mu.Lock()
	defer p.mu.Unlock()
	oldClient, ok := p.Clients[client.ID]

	if ok {
		oldClient.Connection.Close()
		delete(p.Clients, client.ID) 
	}

	p.logger.Log.Info("New client. Size of connection pool: ", len(p.Clients))

	p.Clients[client.ID] = client
	
	for _, cl := range p.Clients {
		if cl.ID != client.ID {
			err := cl.Connection.WriteJSON(Message{ 
				RoomID: p.RoomID,
				ChatUser: client.ID, 
				Body: "New user joined."})
			if err != nil {
				p.logger.Log.Error("Cannot write to client: ", client.ID)
			}
		}
	}
}

func (p *Pool) handleClientUnregistration(client *Client) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.Clients, client.ID)	
	p.logger.Log.Info("Client left. Size of connection pool: ", len(p.Clients))
	
	for _, cl := range p.Clients {
		if cl.ID != client.ID {
			err := cl.Connection.WriteJSON(Message{RoomID: p.RoomID, ChatUser: client.ID, Body: "User left."})
			if err != nil {
				p.logger.Log.Error("Cannot write to client: ", client.ID)
			}
		}
	}

}

func (p *Pool) broadcastMessage(msg Message) {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	p.logger.Log.Info("Broadcasting message to clients.")
			for _,cl := range p.Clients {
				if cl.ID != msg.ChatUser {
					err := cl.Connection.WriteJSON(msg.Body)
					if err != nil {
						p.logger.Log.Error("Cannot write to client: ", cl.ID)
					}
				}
			}
}


func (p *Pool) Start() {
	defer p.ReviveWebsocket()
	for {
		select {
		case client := <- p.Register:
			p.handleClientRegistration(client)
		case client := <- p.Unregister:
			p.handleClientUnregistration(client)
		case msg := <- p.Broadcast:
			p.broadcastMessage(msg)
		}
	}
}

func (p *Pool) ReviveWebsocket() {
	if err := recover(); err != nil {
			log.Println(
				"level", "error",
				"err", err,
			)
		
		go p.Start()
	}
}
