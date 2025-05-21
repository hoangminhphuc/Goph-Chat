package websocket

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/logger"
	"github.com/hoangminhphuc/goph-chat/common/models"
	"github.com/hoangminhphuc/goph-chat/internal/cache"
	"github.com/hoangminhphuc/goph-chat/module/message/model"
	"github.com/hoangminhphuc/goph-chat/plugin/pubsub"
)

// type Message struct {
// 	RoomID 			int	`json:"room_id,omitempty"`
// 	ChatUser 		int `json:"user_id,omitempty"`
// 	Body   			string `json:"content"`
// }


type Client struct {
	ID         	int
	Connection 	*websocket.Conn
	Pool 				*Pool
	logger 			logger.ZapLogger
	msgQueue 		*cache.MessageQueue
}

func NewClient(id int, conn *websocket.Conn, pool *Pool, msgQueue *cache.MessageQueue) *Client {
	return &Client{
		ID:         	id,
		Connection: 	conn,
		Pool: 				pool,
		logger: 			logger.NewZapLogger(),
		msgQueue: 		msgQueue,
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
		var msg model.Message
		if err := c.Connection.ReadJSON(&msg); err != nil {
			if websocket.IsCloseError(err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure,
			) {
					c.logger.Log.Info("Client closed connection:", c.ID)
			} else {
					c.logger.Log.Error("WebSocket read error:", err)
			}
			break
		}
		currentUser :=ctx.MustGet(common.CurrentUser).(*models.Requester)

		msg.RoomID = c.Pool.RoomID
		msg.UserID = currentUser.GetUserId()

		data, err := json.Marshal(msg)
		if err != nil {
				c.logger.Log.Error("json marshal error:", err)
				panic(err)
		}

		var msgWithID *model.Message

		if msgWithID, err = c.msgQueue.CacheAndQueue(context.Background(), 
			fmt.Sprint(msg.RoomID), fmt.Sprint(msg.UserID), 
			[]byte(data)); err != nil {
				c.logger.Log.Error("redis cache/queue error:", err)
		}

		// Sends the message to the Pool
		c.Pool.pubsub.Publish(pubsub.NewMessage(
			pubsub.Topic(fmt.Sprintf("room-%d", c.Pool.RoomID)),
			*msgWithID,
		))
	}
}