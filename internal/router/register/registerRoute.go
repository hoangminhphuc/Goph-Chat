package register

import (
	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	userRoutes "github.com/hoangminhphuc/goph-chat/module/user/routes"
)

func RegisterAllRoutes(router *gin.RouterGroup, serviceCtx serviceHub.ServiceHub) {
	v1 := router.Group("/v1")

	userRoutes.RegisterRoute(v1.Group("/users"), serviceCtx)
}

