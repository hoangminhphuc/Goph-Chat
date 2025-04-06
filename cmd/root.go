package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/internal/router"
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

	baseRouter := gin.Default()
	router.InitRoutes(baseRouter, serviceHub)
	baseRouter.Run()

}