package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []APIError  `json:"errors,omitempty"` // detailed errors on failure
}

type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Code:   http.StatusOK,
		Status: "success",
		Message: message,
		Data:   data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, errs ...APIError) {
	c.JSON(statusCode, APIResponse{
		Code:    statusCode,
		Status:  "error",
		Message: message,
		Errors:  errs,
	})
}