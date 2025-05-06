package websocket

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/logger"
	"github.com/hoangminhphuc/goph-chat/common/models"
	"github.com/hoangminhphuc/goph-chat/plugin/pubsub"
)


type Message struct {
	RoomID 			int
	ChatUser   	int `json:"chatUser,omitempty"`
	Body 				any
}


type Client struct {
	ID         	int
	Connection 	*websocket.Conn
	Pool 				*Pool
	logger 			logger.ZapLogger
}

func NewClient(id int, conn *websocket.Conn, pool *Pool) *Client {
	return &Client{
		ID:         	id,
		Connection: 	conn,
		Pool: 				pool,
		logger: 			logger.NewZapLogger(),
	}
}

func (c *Client) Read(ctx *gin.Context, bodyChan chan []byte) {
	defer func() {
		select {
		case c.Pool.Unregister <- c:
		default: // avoid panic if Unregister channel is closed or full
		}
		_ = c.Connection.Close()
	}()


	for {
		var msg Message
		if err := c.Connection.ReadJSON(&msg); err != nil {
			if websocket.IsCloseError(err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
			) {
					c.logger.Log.Info("Client closed connection:", c.ID)
			} else {
					c.logger.Log.Error("WebSocket read error:", err)
			}
			break
		}

		currentUser :=ctx.MustGet(common.CurrentUser).(*models.Requester)

		msg.RoomID = c.Pool.RoomID
		msg.ChatUser = currentUser.GetUserId()

		// Sends the message to the Pool
		c.Pool.pubsub.Publish(pubsub.NewMessage(
			pubsub.Topic(fmt.Sprintf("room-%d", c.Pool.RoomID)),
			msg,
		))
	}
}