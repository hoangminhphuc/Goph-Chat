package cmd

import (
	"github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/internal/router/register"
	"github.com/hoangminhphuc/goph-chat/plugin/db"
)

func newService() boot.ServiceHub {
	service := boot.NewServiceHub(
		"Goph-Chat",
		boot.RegisterPlugin(db.NewGormDB(common.PluginDBMain)),
	)
	return service
}

func Execute() {
	serviceHub := newService()

	logger := serviceHub.GetLogger()

	if err := serviceHub.Init(); err != nil {
		logger.Log.Error(err.Error())
	}


	register.RegisterAllRoutes(serviceHub.GetHTTPServer().GetRouter().Group("/"), serviceHub)


	if err := serviceHub.Start(); err != nil {
		logger.Log.Fatal(err.Error())
	}
}