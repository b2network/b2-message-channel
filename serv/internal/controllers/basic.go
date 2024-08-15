package controllers

import (
	"bsquared.network/b2-message-channel-serv/internal/middlewares"
	"bsquared.network/b2-message-channel-serv/internal/utils/sso"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type BasicController struct{}

func (c *BasicController) Authorization(ctx *gin.Context, secret string) (*sso.UserClaims, error) {
	authorization := ctx.GetHeader("Authorization")
	if strings.HasPrefix(authorization, "Bearer ") {
		authorization = strings.ReplaceAll(authorization, "Bearer ", "")
	}
	claims, err := sso.VerifyToken(secret, authorization)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (c *BasicController) Return(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, data)
}

func (c *BasicController) Success(ctx *gin.Context, data ...interface{}) {
	var d interface{}
	if len(data) > 0 {
		d = data[0]
	}
	res := struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data,omitempty"`
	}{
		Code: 0,
		Msg:  "success",
		Data: d,
	}
	ctx.AbortWithStatusJSON(http.StatusOK, res)
}

func (c *BasicController) Error(ctx *gin.Context, statusCode int, message string, key string, code int, reason error) {
	log.Errorf("[%s]code: %d, reason: %v\n", key, code, reason)
	err := middlewares.ErrorResponse{Status: statusCode}
	err.Msg = fmt.Sprintf("%s[%d]", message, code)
	err.Reason = fmt.Sprintf("[%s]code: %d, reason: %v\n", key, code, reason)
	panic(err)
}
