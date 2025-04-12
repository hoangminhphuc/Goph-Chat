package websocket

import (
	"log"

	"github.com/hoangminhphuc/goph-chat/common/logger"
)

type Pool struct {
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

func (p *Pool) handleClientRegistration(client *Client) {
	oldClient, ok := p.Clients[client.ID]

	if ok {
		oldClient.Connection.Close()
		delete(p.Clients, client.ID) 
	}

	p.logger.Log.Info("New client. Size of connection pool: ", len(p.Clients))

	p.Clients[client.ID] = client
	
	for _, cl := range p.Clients {
		if cl.RoomID == client.RoomID && cl.ID != client.ID {
			err := cl.Connection.WriteJSON(Message{ 
				RoomID: client.RoomID,
				ChatUser: client.ID, 
				Body: "New user joined."})
			if err != nil {
				p.logger.Log.Error("Cannot write to client: ", client.ID)
			}
		}
	}
}

func (p *Pool) handleClientUnregistration(client *Client) {
	delete(p.Clients, client.ID)	
	p.logger.Log.Info("Client left. Size of connection pool: ", len(p.Clients))
	
	for _, cl := range p.Clients {
		if cl.RoomID == client.RoomID && cl.ID != client.ID {
			err := cl.Connection.WriteJSON(Message{RoomID: client.RoomID, ChatUser: client.ID, Body: "User left."})
			if err != nil {
				p.logger.Log.Error("Cannot write to client: ", client.ID)
			}
		}
	}

}

func (p *Pool) broadcastMessage(msg Message) {
	p.logger.Log.Info("Broadcasting message to clients.")
			for _,cl := range p.Clients {
				if cl.RoomID == msg.RoomID && cl.ID != msg.ChatUser {
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
