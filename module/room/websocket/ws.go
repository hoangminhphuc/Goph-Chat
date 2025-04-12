package websocket

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hoangminhphuc/goph-chat/common/logger"
)

var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool { return true },
}


func ServerWebSocket(c *gin.Context, pool *Pool) {
	roomID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}


	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		
		if err != nil {
			return 
		}
		client := &Client {
			ID: fmt.Sprintf("%d", time.Now().Unix()),
			Connection: conn,
			RoomID: roomID,
			Pool: pool,
			logger: logger.NewZapLogger(),
		}
		pool.Register <- client

		requestBody := make(chan []byte) // websocket.Message byte array channel

		go client.Read(requestBody)


}