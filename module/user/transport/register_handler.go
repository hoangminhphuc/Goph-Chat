package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/utils"
	"github.com/hoangminhphuc/goph-chat/module/user/business"
	"github.com/hoangminhphuc/goph-chat/module/user/dto"
	"github.com/hoangminhphuc/goph-chat/module/user/repository"
	"gorm.io/gorm"
)

func Register(serviceCtx serviceHub.ServiceHub) func(*gin.Context) {
	return func(c *gin.Context) {	
		db := serviceCtx.MustGetService(common.PluginDBMain).(*gorm.DB)
		var data dto.UserRegister
		if err := c.ShouldBindJSON(&data); err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		repo := repository.NewSQLRepo(db)
		hasher := utils.NewBcryptHash()

		registerBusiness := business.NewRegisterBusiness(repo, hasher)

		if err := registerBusiness.Register(c.Request.Context(), &data); err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		
		common.SuccessResponse(c, "Register successfully", nil)
	}
}