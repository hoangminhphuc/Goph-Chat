package websocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/hoangminhphuc/goph-chat/common/logger"
)


type Message struct {
	RoomID 			int
	ChatUser   	string `json:"chatUser,omitempty"`
	Body 				any
}


type Client struct {
	ID         	string
	Connection 	*websocket.Conn
	RoomID     	int
	Pool   			*Pool
	logger 			logger.ZapLogger
}

func (c *Client) Read(bodyChan chan []byte) {
	defer func() {
		c.Pool.Unregister <- c
		c.Connection.Close()
	}()

	defer c.Pool.ReviveWebsocket()


	for {
		_, p, err := c.Connection.ReadMessage()
		if err != nil {
			c.logger.Log.Info("Client disconnected: ", c.ID)
			break
		}


		var msg Message
		err = json.Unmarshal(p, &msg)
		msg.RoomID = c.RoomID
		msg.ChatUser = c.ID

		if err != nil {
			c.logger.Log.Error("Invalid JSON received: ", err)
		}

		// Sends the message to the Pool
		c.Pool.Broadcast <- msg
	}
}