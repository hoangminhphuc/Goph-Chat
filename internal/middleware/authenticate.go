package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/models"
	"github.com/hoangminhphuc/goph-chat/common/utils"
	"github.com/hoangminhphuc/goph-chat/module/user/model"
)

type authenStore interface {
	FindUser(ctx context.Context, conditions map[string]interface{}) (*model.User, error)
}

func extractTokenFromHeaderString(s string) (string, error) {
	parts := strings.Split(s, " ")
	// "Authorization" : "Bearer {token}"

	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", common.NewError("Wrong auth header", http.StatusBadRequest)
	}

	return parts[1], nil
}

func RequireAuth(authenStore authenStore, secret string) func(*gin.Context) {
	return func(c *gin.Context) {
		token, err := extractTokenFromHeaderString(c.GetHeader("Authorization"))

		if err != nil {
			common.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		payload, err := utils.VerifyToken(token, secret)

		if err != nil {
			common.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		user, err := authenStore.FindUser(c.Request.Context(), map[string]interface{}{"id": payload.UserId()})
		
		if err != nil {
			common.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		currentUser := models.NewRequester(user.ID, user.Email, user.Role.String())

		// Set throughout this whole request reponse route (each time access an URL)
		c.Set(common.CurrentUser, currentUser)

		keys := c.Keys
		log.Println(keys) // Check if common.CurrentUser is present in the keys

		c.Next()
	}
}