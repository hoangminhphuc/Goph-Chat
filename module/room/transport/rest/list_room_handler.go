package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/utils"
	"github.com/hoangminhphuc/goph-chat/module/room/business"
	"github.com/hoangminhphuc/goph-chat/module/room/repository"
	"gorm.io/gorm"
)

func ListRoom(serviceCtx serviceHub.ServiceHub) func(*gin.Context) {
	return func(c *gin.Context) {
		db := serviceCtx.MustGetService(common.PluginDBMain).(*gorm.DB)

		var paging utils.Paging

		if err := c.ShouldBindQuery(&paging); err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		paging.Process()

		repo := repository.NewSQLRepo(db)
		business := business.NewListRoomBusiness(repo)

		rooms, err := business.ListRoom(c.Request.Context(), &paging)

		if err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		for i := range rooms {
			rooms[i].Mask()
		}

		common.SuccessResponse(c, "Rooms retrieved successfully", "paging", paging, "rooms", rooms)
	}
}