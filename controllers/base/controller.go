package base

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tautcony/qart/models/response"
)

var (
	AppVer string
)

// JSON response helper
func JSON(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

// Success response
func Success(c *gin.Context, data any, code int) {
	r := &response.BaseResponse{
		Success: true,
		Code:    code,
		Data:    data,
	}
	JSON(c, r)
}

// Fail response
func Fail(c *gin.Context, data any, code int, message string) {
	r := &response.BaseResponse{
		Code:    code,
		Data:    data,
		Message: message,
	}
	JSON(c, r)
}
