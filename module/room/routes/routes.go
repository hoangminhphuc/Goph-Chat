package routes

import (
	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/module/room/transport/rest"
	ws "github.com/hoangminhphuc/goph-chat/module/room/transport/websocket"
)

func RegisterWebSocketRoute(rooms *gin.RouterGroup, serviceCtx serviceHub.ServiceHub) {
	rooms.POST("", rest.CreateRoom(serviceCtx))
	rooms.GET("/:id", rest.GetRoomByID(serviceCtx))

	chat := rooms.Group("/ws")
	{
		chat.GET("/:id", ws.HandleWebSocketConnection(serviceCtx))
	}
}
