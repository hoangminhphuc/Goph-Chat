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

func GetRoomByID(serviceCtx serviceHub.ServiceHub) func(*gin.Context) {
	return func(c *gin.Context) {
		db := serviceCtx.MustGetService(common.PluginDBMain).(*gorm.DB)

		id, err := utils.DecodeID(c.Param("id"))

		if err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		
		repo := repository.NewSQLRepo(db)
		business := business.NewGetRoomBusiness(repo)

		room, err := business.GetRoom(c.Request.Context(), id)
		if err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		room.Mask()

		common.SuccessResponse(c, "Get room successfully", "room", room)
	}
}