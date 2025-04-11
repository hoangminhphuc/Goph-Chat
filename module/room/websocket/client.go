package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)


type Message struct {
	RoomID 			int
	ChatUser   	string `json:"chatUser,omitempty"`
	Body 				any
}


type Client struct {
	ID         string
	Connection *websocket.Conn
	Pool   *Pool
}

func NewClient(userID string, pool *Pool, conn *websocket.Conn) *Client{
	return &Client{
		Connection: conn,
		Pool:       pool,
		ID:         userID,
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
		msg.ChatUser = c.ID
		log.Println(msg)
		if err != nil {
			log.Fatal(err)
		}

		// Sends the message to the Pool
		c.Pool.Broadcast <- msg
	}
}