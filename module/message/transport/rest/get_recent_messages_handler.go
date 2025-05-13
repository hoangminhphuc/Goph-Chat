package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/models"
	"github.com/hoangminhphuc/goph-chat/module/message/business"
	"github.com/hoangminhphuc/goph-chat/module/message/repository"
	"github.com/redis/go-redis/v9"
)

func GetRecentMessages(serviceCtx serviceHub.ServiceHub) func(*gin.Context) {
	return func(c *gin.Context) {
		roomID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		currentUser := c.MustGet(common.CurrentUser).(*models.Requester)
		rdb := serviceCtx.MustGetService(common.PluginRedisMain).(*redis.Client)
		userID := currentUser.GetUserId()

		repo := repository.NewRedisRepo(rdb)
		business := business.NewGetRecentMessagesBusiness(repo)

		messages, err := business.GetRecentMessages(c, roomID, userID)
		if err != nil {
			common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		common.SuccessResponse(c, "Get recent messages successfully", "messages", messages)
	}
}