package routes

import (
	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/module/message/transport/rest"
)

func RegisterMessageRoute(v1 *gin.RouterGroup, serviceCtx serviceHub.ServiceHub) {

	v1.GET("/:id", rest.GetRecentMessages(serviceCtx))
}