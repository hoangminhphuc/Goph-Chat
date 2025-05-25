package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/models"
	"github.com/hoangminhphuc/goph-chat/module/user/model"
)

//
var routePermissions = map[string]map[string][]model.UserRole{
	// User Route
    "/v1/register": {
        "POST": {model.RoleAdmin, model.RoleUser},
    },
    "/v1/login": {
        "POST": {model.RoleAdmin, model.RoleUser},
    },

	// Room Route
    "/v1/rooms": {
        "GET":    {model.RoleAdmin, model.RoleUser},
        "POST":   {model.RoleAdmin, model.RoleUser},
    },
		"/v1/rooms/:id": {
        "GET":    {model.RoleAdmin, model.RoleUser},
        "DELETE":   {model.RoleAdmin, model.RoleUser},
    },
		"/v1/rooms/ws/:id": {
        "GET":    {model.RoleAdmin, model.RoleUser},
    },

	// Message Route
    "/v1/messages/:id": {
        "GET":    {model.RoleAdmin, model.RoleUser},
        "PATCH":    {model.RoleAdmin, model.RoleUser},
    },
}

func ParseUserRole(s string) (model.UserRole, error) {
    switch s {
    case "user":
        return model.RoleUser, nil
    case "admin":
        return model.RoleAdmin, nil
    default:
        return -1, fmt.Errorf("invalid role: %s", s)
    }
}


func can(role model.UserRole, path, method string) bool {
    if methods, ok := routePermissions[path]; ok {
        if roles, ok2 := methods[method]; ok2 {
            for _, r := range roles {
                if r == role {
                    return true
                }
            }
        }
    }
    return false
}

func RBAC() gin.HandlerFunc {
    return func(c *gin.Context) {
        // userAny, exists := c.Get("currentUser")
				currentUser := c.MustGet(common.CurrentUser).(*models.Requester)

        path := c.FullPath()        // "/v1/messages/:id"
        method := c.Request.Method  

				userRole, err := ParseUserRole(currentUser.GetRole())
				if err != nil {
					common.ErrorResponse(c, http.StatusBadRequest, err.Error())
					return
				}

        if !can(userRole, path, method) {
            common.ErrorResponse(c, http.StatusForbidden, "permission denied")
            return
        }

        c.Next()
    }
}
