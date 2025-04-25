package register

import (
	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/internal/middleware"
	roomWebSocketRoutes "github.com/hoangminhphuc/goph-chat/module/room/routes"
	userRoutes "github.com/hoangminhphuc/goph-chat/module/user/routes"
	"gorm.io/gorm"
	userstorage "github.com/hoangminhphuc/goph-chat/module/user/repository"
)

func RegisterAllRoutes(router *gin.RouterGroup, serviceCtx serviceHub.ServiceHub) {
	db := serviceCtx.MustGetService(common.PluginDBMain).(*gorm.DB)
	authStore := userstorage.NewSQLRepo(db) 
	secret := serviceCtx.GetEnvValue("JWT_SECRET")
	middlewareAuth := middleware.RequireAuth(authStore, secret)
	
	v1 := router.Group("/v1")
	userRoutes.RegisterUserRoute(v1, serviceCtx)


	rooms := v1.Group("/rooms", middlewareAuth)
	roomWebSocketRoutes.RegisterWebSocketRoute(rooms, serviceCtx)
}

