package interceptor

import (
	"github.com/gin-gonic/gin"
)

// SendSuccessRes sends the success response to the client.
// If statusCode is not provided then "default success status code"(i.e, 200) will be used.
func SendSuccessRes(c *gin.Context, data any, statusCode int) {
	if statusCode == 0 {
		statusCode = DEFAULT_HTTP_SUCCESS_CODE
	}
	response := NewResponse(true, data, "")
	c.AbortWithStatusJSON(statusCode, response)
}

// SendErrRes sends the error response to the client.
// If statusCode is not provided then "default error status code"(i.e, 500) will be used.
// If errMsg is not provided then response with "default error message" will be sent.
func SendErrRes(c *gin.Context, errMsg string, statusCode int) {
	if statusCode == 0 {
		statusCode = DEFAULT_HTTP_ERROR_CODE
	}
	if errMsg == "" {
		errMsg = DEFAULT_ERROR_MSG
	}
	response := NewResponse(false, nil, errMsg)
	c.AbortWithStatusJSON(statusCode, response)
}
