package register

import (
	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/internal/cache"
	"github.com/hoangminhphuc/goph-chat/internal/middleware"
	roomWebSocketRoutes "github.com/hoangminhphuc/goph-chat/module/room/routes"
	userstorage "github.com/hoangminhphuc/goph-chat/module/user/repository"
	userRoutes "github.com/hoangminhphuc/goph-chat/module/user/routes"
	messageRoutes "github.com/hoangminhphuc/goph-chat/module/message/routes"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterAllRoutes(router *gin.RouterGroup, serviceCtx serviceHub.ServiceHub) {
	db := serviceCtx.MustGetService(common.PluginDBMain).(*gorm.DB)
	redis := serviceCtx.MustGetService(common.PluginRedisMain).(*redis.Client)

	cache.InitDefaultCache(redis)

	authStore := userstorage.NewSQLRepo(db) 
	secret := serviceCtx.GetEnvValue("JWT_SECRET")
	middlewareAuth := middleware.RequireAuth(authStore, secret)

	middlewareRBAC := middleware.RBAC()
	
	v1 := router.Group("/v1")
	userRoutes.RegisterUserRoute(v1, serviceCtx)

	rooms := v1.Group("/rooms", middlewareAuth, middlewareRBAC)
	roomWebSocketRoutes.RegisterWebSocketRoute(rooms, serviceCtx)

	messages := v1.Group("/messages", middlewareAuth, middlewareRBAC)
	messageRoutes.RegisterMessageRoute(messages, serviceCtx)
}

