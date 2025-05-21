package websocket

import (
	"fmt"
	"sync"

	"github.com/hoangminhphuc/goph-chat/common/logger"
	"github.com/hoangminhphuc/goph-chat/module/message/model"
	"github.com/hoangminhphuc/goph-chat/plugin/pubsub"
)

// Each Pool is a Room
type Pool struct {
	RoomID 			int
	Clients 		map[int]*Client
	Register 		chan *Client
	Unregister 	chan *Client
	// Broadcast 	chan Message
	done 				chan struct{}
	pubsub    	*pubsub.LocalPubSub
	subChans     	[]<-chan *pubsub.Message
	unsubs     	[]func()
	logger 			logger.ZapLogger
	mu 					sync.RWMutex
}

func NewPool(roomID int, ps *pubsub.LocalPubSub) *Pool {
	p := &Pool{
		RoomID: 			roomID,
		Clients: 			make(map[int]*Client),
		Register: 		make(chan *Client, 10),
		Unregister: 	make(chan *Client, 10),
		done:       make(chan struct{}),
		pubsub:     ps,
		logger: 			logger.NewZapLogger(),
	}

	p = RegisterTopics(p, 
		pubsub.Topic(fmt.Sprintf("room-%d", roomID)), 
		pubsub.Topic(fmt.Sprintf("room-%d:updated-msg", roomID)))

	return p
}

func RegisterTopics(p *Pool, topics ...pubsub.Topic) *Pool {
    for _, topic := range topics {
        ch, unsub := p.pubsub.Subscribe(topic)
        p.subChans = append(p.subChans, ch)
        p.unsubs = append(p.unsubs, unsub)
    }
    return p
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
	
	msg := model.Message{
			RoomID:  p.RoomID,       
			UserID:  client.ID,  
			Content: "New user joined.", 
		}
	msg.Mask()
	// Avoids WriteJSON to blocks other pool operation
	for _, peer := range peers {
		err := peer.Connection.WriteJSON(msg)
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

	msg := model.Message{
			RoomID: p.RoomID,
			UserID: client.ID,
			Content: "User left.",
		}

	msg.Mask()

	// Avoids WriteJSON to blocks other pool operation
	for _, peer := range peers {
		err := peer.Connection.WriteJSON(msg)
		if err != nil {
			p.logger.Log.Error("Cannot write to client: ", peer.ID)
		}
	}
}

func (p *Pool) broadcastMessage(msg model.Message) {
	p.mu.RLock()
	peers := make([]*Client, 0, len(p.Clients))
	for _, c := range p.Clients {
		if c.ID != msg.UserID {
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

	mergedChannel := p.merge()

	for {
		select {
		case client := <- p.Register:
			p.handleClientRegistration(client)
		case client := <- p.Unregister:
			p.handleClientUnregistration(client)
		case msg, ok := <- mergedChannel:
			if !ok {
				return
			}
			
			message := msg.GetData().(model.Message)
			message.Mask()

			p.broadcastMessage(message)
		case <-p.done:
			// graceful shutdown
			return
		}
	}
}

func (p *Pool) merge() <-chan *pubsub.Message {
    out := make(chan *pubsub.Message)

    // For each subscription, start a goroutine that forwards its messages
    for _, ch := range p.subChans {
        ch := ch // capture loop var
        go func() {
            for msg := range ch {
                out <- msg
            }
        }()
    }

    return out
}


func (p *Pool) Stop() {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, client := range p.Clients {
		_ = client.Connection.Close() // Close connection
	}

	close(p.done)
	for _, unsub := range p.unsubs {
        unsub()
  }
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
