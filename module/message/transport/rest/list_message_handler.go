package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/utils"
	"github.com/hoangminhphuc/goph-chat/module/message/business"
	"github.com/hoangminhphuc/goph-chat/module/message/repository"
	"gorm.io/gorm"
)

func ListMessage(serviceCtx serviceHub.ServiceHub) func(*gin.Context) {
	return func(c *gin.Context) {
		db := serviceCtx.MustGetService(common.PluginDBMain).(*gorm.DB)

		roomID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		var paging utils.Paging
		if err := c.ShouldBindQuery(&paging); err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		paging.Process()

		repo := repository.NewSQLRepo(db)
		business := business.NewListMessageBusiness(repo)

		messages, err := business.ListMessage(c.Request.Context(), roomID, &paging)
		if err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		for i := range messages {
			messages[i].Mask()
		}

		common.SuccessResponse(c, "Messages retrieved successfully", "paging", paging, "messages", messages)
	}
}
