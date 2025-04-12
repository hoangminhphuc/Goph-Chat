package routes

import (
	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	ws "github.com/hoangminhphuc/goph-chat/module/room/websocket"
)

func RegisterWebSocketRoute(v1 *gin.RouterGroup, serviceCtx serviceHub.ServiceHub) {
	roomManager := ws.GetRoomCenter()
	chat := v1.Group("/ws")

	chat.GET("/:id", func(c *gin.Context) {
		ws.ServerWebSocket(c, roomManager)
	})
}
