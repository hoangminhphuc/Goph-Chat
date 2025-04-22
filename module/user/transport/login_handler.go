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

func Login(serviceCtx serviceHub.ServiceHub) func(*gin.Context) {
	return func(c *gin.Context) {
		db := serviceCtx.MustGetService(common.PluginDBMain).(*gorm.DB)
		secret := serviceCtx.GetEnvValue("JWT_SECRET")
		
		var data dto.UserLogin
		if err := c.ShouldBind(&data); err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		repo := repository.NewSQLRepo(db)
		hasher := utils.NewBcryptHash()

		business :=  business.NewLoginBusiness(repo, 
			hasher, 
			utils.AccessTokenExpiredTime, 
			secret)
		
		token, err := business.Login(c.Request.Context(), &data)
		
		if err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		common.SuccessResponse(c, "Login successfully", "token", token)
	}
}