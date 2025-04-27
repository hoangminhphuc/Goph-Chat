package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/models"
	"github.com/hoangminhphuc/goph-chat/common/utils"
	"github.com/hoangminhphuc/goph-chat/module/room/business"
	"github.com/hoangminhphuc/goph-chat/module/room/dto"
	"github.com/hoangminhphuc/goph-chat/module/room/repository"
	ws "github.com/hoangminhphuc/goph-chat/internal/server/websocket"
	"gorm.io/gorm"
)

func CreateRoom(serviceCtx serviceHub.ServiceHub) func(*gin.Context) {
	return func(c *gin.Context) {
		db := serviceCtx.MustGetService(common.PluginDBMain).(*gorm.DB)
		wsServer := serviceCtx.MustGetRuntimeService(common.PluginWSMain).(*ws.WebSocketServer)

		var data dto.RoomCreation
		if err := c.ShouldBind(&data); err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		currentUser :=c.MustGet(common.CurrentUser).(*models.Requester)
		data.OwnerID = currentUser.GetUserId()

		repo := repository.NewSQLRepo(db)
		business := business.NewCreateRoomBusiness(repo)

		if err := business.CreateRoom(c, &data); err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		/* 
		! THIS SHOULD NOT BE IN HERE, WILL BE IMPROVE LATER
		! BY INTEGRATING EVENT-DRIVEN ARCHITECTURE (PUB/SUB)
		*/
		if err := wsServer.CreateRoom(data.ID); err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		uuid, _  := utils.EncodeID(uint64(data.ID))

		common.SuccessResponse(c, "Create room successfully", "id", uuid)
		
	}
}