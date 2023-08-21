package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	ErrorCode int         `json:"errorCode"`
	Result    interface{} `json:"result"`
	Message   string      `json:"message"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		200,
		data,
		"ok",
	})
}

func ErrorHandler(c *gin.Context, code int, err string) {
	c.JSON(http.StatusOK, Response{
		code,
		nil,
		err,
	})
}

func ErrorStatusHandler(c *gin.Context, status int, code int, err string) {
	c.JSON(status, Response{
		code,
		nil,
		err,
	})
}
