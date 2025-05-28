package routes

import (
	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/module/message/transport/rest"
)

func RegisterMessageRoute(message *gin.RouterGroup, serviceCtx serviceHub.ServiceHub) {

	message.GET("/:id/recent", rest.GetRecentMessages(serviceCtx))
	message.PATCH("/:id", rest.EditMessage(serviceCtx))
	message.GET("/:id", rest.ListMessage(serviceCtx))
}