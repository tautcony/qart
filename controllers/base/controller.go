package base

import (
	"github.com/gin-gonic/gin"
	"github.com/tautcony/qart/models/response"
	"net/http"
)

var (
	AppVer string
)

// JSON response helper
func JSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// Success response
func Success(c *gin.Context, data interface{}, code int) {
	r := &response.BaseResponse{
		Success: true,
		Code:    code,
		Data:    data,
	}
	JSON(c, r)
}

// Fail response
func Fail(c *gin.Context, data interface{}, code int, message string) {
	r := &response.BaseResponse{
		Code:    code,
		Data:    data,
		Message: message,
	}
	JSON(c, r)
}
