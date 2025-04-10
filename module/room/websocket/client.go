package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)


type Message struct {
	RoomID 			int
	ChatUser   string `json:"chatUser,omitempty"`
}


type Client struct {
	ID         string
	Connection *websocket.Conn
	Pool   *Pool
}

func NewClient(pool *Pool, conn *websocket.Conn) *Client{
	return &Client{
		Connection: conn,
		Pool:       pool,
	}
}

func (c *Client) Read(bodyChan chan []byte) {
	defer func() {
		c.Pool.Unregister <- c
		c.Connection.Close()
	}()

	for {
		_, p, err := c.Connection.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		var msg Message
		err = json.Unmarshal(p, &msg)
		if err != nil {
			log.Fatal(err)
		}
		// Adds the sender's identity (email) to the message so others know who sent it.
		message := Message {
			RoomID: msg.RoomID,
			ChatUser: c.ID,
		}

		// Sends the message to the Pool
		c.Pool.Broadcast <- message
		log.Println("info:", "Message received: ", message)
	}
}