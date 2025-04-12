package websocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/hoangminhphuc/goph-chat/common/logger"
)


type Message struct {
	RoomID 			string
	ChatUser   	string `json:"chatUser,omitempty"`
	Body 				any
}


type Client struct {
	ID         	string
	Connection 	*websocket.Conn
	Pool 				*Pool
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
		if err = json.Unmarshal(p, &msg); err != nil {
			c.logger.Log.Error("Invalid JSON received: ", err)
			return
		}

		msg.RoomID = c.Pool.RoomID
		msg.ChatUser = c.ID

		

		// Sends the message to the Pool
		c.Pool.Broadcast <- msg
	}
}