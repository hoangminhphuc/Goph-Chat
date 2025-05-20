package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	serviceHub "github.com/hoangminhphuc/goph-chat/boot"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/models"
	"github.com/hoangminhphuc/goph-chat/module/message/business"
	"github.com/hoangminhphuc/goph-chat/module/message/dto"
	"github.com/hoangminhphuc/goph-chat/module/message/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func EditMessage(serviceHub serviceHub.ServiceHub) func(*gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		db := serviceHub.MustGetService(common.PluginDBMain).(*gorm.DB)
		rdb := serviceHub.MustGetService(common.PluginRedisMain).(*redis.Client)
		sqlRepo := repository.NewSQLRepo(db)
		redisRepo := repository.NewRedisRepo(rdb)

		business := business.NewUpdateMessageBusiness(sqlRepo, redisRepo)

		// JSON body
		var data *dto.MessageUpdate
		if err := c.ShouldBindJSON(&data); err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		// Query parameters 
		var roomInfo dto.RoomInfo
		if err := c.ShouldBindQuery(&roomInfo); err != nil {
			common.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		data.RoomID = roomInfo.RoomID

		requester := c.MustGet(common.CurrentUser).(*models.Requester)
		data.UserID = requester.GetUserId()

		if err := business.UpdateMessageByID(c, id, data); err != nil {
			common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

	}
}