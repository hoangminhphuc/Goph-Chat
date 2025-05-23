package cmd

import (
	"github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/internal/asyncjob/message"
	"github.com/hoangminhphuc/goph-chat/internal/maintenance"
	rt "github.com/hoangminhphuc/goph-chat/internal/router"
	"github.com/hoangminhphuc/goph-chat/internal/router/register"
	"github.com/hoangminhphuc/goph-chat/internal/server/websocket"
	"github.com/hoangminhphuc/goph-chat/plugin/gorm"
	"github.com/hoangminhphuc/goph-chat/plugin/pubsub"
	"github.com/hoangminhphuc/goph-chat/plugin/redis"
)

func newService() boot.ServiceHub {
	service := boot.NewServiceHub(
		"Goph-Chat",
		boot.RegisterInitService(db.NewGormDB(common.PluginDBMain)),
		boot.RegisterInitService(redis.NewRedisDB()),
		boot.RegisterInitService(pubsub.NewLocalPubSub(common.PluginPubSubMain)),
		boot.RegisterRuntimeService(rt.NewHTTPServer()),
    boot.RegisterRuntimeService(websocket.NewWebSocketServer()),
	)
	return service
}

func Execute() {
	serviceHub := newService()

	logger := serviceHub.GetLogger()

	if err := serviceHub.Init(); err != nil {
		logger.Log.Error(err.Error())
	}

	
	serviceHub.InitializePools(serviceHub.MustGetRuntimeService(common.PluginWSMain).(*websocket.WebSocketServer))
	register.RegisterAllRoutes(serviceHub.MustGetRuntimeService(common.PluginHTTPMain).(*rt.HTTPServer).GetRouter().Group("/"), serviceHub)

	maintenance.StartCleanupService(serviceHub)
	message.RunBackgroundWorkers(serviceHub)

	if err := serviceHub.Start(); err != nil {
		logger.Log.Fatal(err.Error())
	}
}