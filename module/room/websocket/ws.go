package websocket

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool { return true },
}


func ServerWebSocket(c *gin.Context, pool *Pool) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		
		if err != nil {
			return 
		}
		client := NewClient(fmt.Sprintf("%x", time.Now().UnixNano()), pool, conn)
		pool.Register <- client

		requestBody := make(chan []byte) // websocket.Message byte array channel

		go client.Read(requestBody)


}