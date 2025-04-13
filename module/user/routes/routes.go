package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"gorm.io/gorm"
)

func RegisterUserRoute(v1 *gin.RouterGroup, serviceCtx serviceHub.ServiceHub) {
	db := serviceCtx.MustGetService(common.PluginDBMain).(*gorm.DB)
  v1.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": db.Dialector.Name(),
    })
  })
}