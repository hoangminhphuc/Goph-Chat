package register

import (
	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	roomWebSocketRoutes "github.com/hoangminhphuc/goph-chat/module/room/routes"
	userRoutes "github.com/hoangminhphuc/goph-chat/module/user/routes"
)

func RegisterAllRoutes(router *gin.RouterGroup, serviceCtx serviceHub.ServiceHub) {
	v1 := router.Group("/v1")

	userRoutes.RegisterUserRoute(v1, serviceCtx)
	roomWebSocketRoutes.RegisterWebSocketRoute(v1.Group("/rooms"), serviceCtx)
}

