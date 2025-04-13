package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/module/user/business"
	"github.com/hoangminhphuc/goph-chat/module/user/repository"
	"gorm.io/gorm"
)

func Register(serviceCtx serviceHub.ServiceHub) func(*gin.Context) {
	return func(c *gin.Context) {	
		db := serviceCtx.MustGetService(common.PluginDBMain).(*gorm.DB)

		repo := repository.NewSQLRepo(db)
		hasher := NewBCryptHasher()

		registerBusiness := business.NewRegisterBusiness(repo, hasher)

		if err := registerBusiness.Register(); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		
		c.JSON(http.StatusOK, nil)
	}
}