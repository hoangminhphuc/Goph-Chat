package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/module/user/transport"
)

func RegisterUserRoute(v1 *gin.RouterGroup, serviceCtx serviceHub.ServiceHub) {
	
	v1.POST("/register", transport.Register(serviceCtx))
	v1.POST("/login", transport.Login(serviceCtx))


  v1.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "sucess",
    })
  })
}