package common

import (
	"fmt"
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

func BuildResponseData(data ...interface{}) (map[string]interface{}, error) {
	if data[0] == nil {
		return nil, nil 
	}


	responseData := make(map[string]interface{})

	for i := 0; i < len(data); i += 2 {
		key := data[i]
		value := data[i+1]

		// Ensure that the key is a string
		if keyStr, ok := key.(string); ok {
			responseData[keyStr] = value
		}  else {
			return nil, fmt.Errorf("invalid key type at position %d: expected string, got %T", i, key)
		}
	}

	return responseData, nil
}

func SuccessResponse(c *gin.Context, message string, data ...interface{}) {
	responseData, err := BuildResponseData(data...)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	apiResponse := APIResponse{
		Code:   http.StatusOK,
		Status: "success",
		Message: message,
	}

	if responseData != nil {
		apiResponse.Data = responseData
	}

	c.JSON(http.StatusOK, apiResponse)
}

func ErrorResponse(c *gin.Context, statusCode int, message string, errs ...APIError) {
	c.JSON(statusCode, APIResponse{
		Code:    statusCode,
		Status:  "error",
		Message: message,
		Errors:  errs,
	})
}