package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	ws "github.com/hoangminhphuc/goph-chat/module/room/websocket"
)

var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool { return true },
}

func RegisterWebSocketRoute(v1 *gin.RouterGroup, serviceCtx serviceHub.ServiceHub) {
	pool := ws.NewPool()
	go pool.Start()
	v1.GET("/ws", func(c *gin.Context) {
		ws.ServerWebSocket(c, pool)

		// go func() {
		// 	defer func() {
		// 		pool.Remove(client.ID)
		// 		client.Connection.Close()
		// 	}()
		// 	for {
		// 		// Read message from the client (frontend)
		// 		_, msg, err := client.Connection.ReadMessage()
		// 		log.Printf("📥 Received from client %d: %s", client.ID, string(msg))
		// 		if err != nil {
		// 			break
		// 		}
		// 	// Broadcast to all send channel of clients so it can sends to the frontend
		// 		pool.Broadcast(client.ID, msg)
		// 	}
		// }()

		// go func() {
		// 	for msg := range client.Send {
		// 		// Send message to the client (frontend) if receive message
		// 		if err := client.Connection.WriteMessage(websocket.TextMessage, msg); err != nil {
		// 			break
		// 		}
		// 		log.Printf("📤 Sent to client %d: %s", client.ID, string(msg))
		// 	}
		// }()

	})
}
