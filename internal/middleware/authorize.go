package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hoangminhphuc/goph-chat/common"
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
			c.JSON(http.StatusUnauthorized, err.Error())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		_, err = utils.VerifyToken(token, secret)

		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}