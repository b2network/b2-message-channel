package middlewares

import (
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

type customResponse struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data,omitempty"`
	Reason string      `json:"reason,omitempty"`
}

type ErrorResponse struct {
	Status int `json:"-"`
	customResponse
}

func (e ErrorResponse) Error() string {
	return e.Msg
}

func ErrorHandler(c *gin.Context, err any) {
	resErr := err.(ErrorResponse)
	sentry.CaptureException(resErr)
	resErr.Reason = ""
	c.AbortWithStatusJSON(resErr.Status, resErr.customResponse)
}
