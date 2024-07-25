package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type ResponseJson struct {
	Status int    `json:"-"`
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Data   any    `json:"data,omitempty"`
}

// IsEmpty 判断结构体是否为空
func (r ResponseJson) IsEmpty() bool {
	return reflect.DeepEqual(r, ResponseJson{})
}

// 构建状态码 ，如果 传入的ResponseJson没有Status 就使用默认的状态码
func buildStatus(resp ResponseJson, defaultStatus int) int {
	if resp.Status == 0 {
		return defaultStatus
	}
	return resp.Status
}

func HttpResponse(ctx *gin.Context, status int, resp ResponseJson) {
	if resp.IsEmpty() {
		ctx.AbortWithStatus(status)
		return
	}
	ctx.AbortWithStatusJSON(status, resp)
}

func Success(ctx *gin.Context, data any) {
	resp := ResponseJson{
		Code: 0,
		Msg:  "",
		Data: data,
	}
	HttpResponse(ctx, buildStatus(resp, http.StatusOK), resp)
}

func Fail(ctx *gin.Context, resp ResponseJson) {
	HttpResponse(ctx, buildStatus(resp, http.StatusBadRequest), resp)
}

func ServerFail(ctx *gin.Context, resp ResponseJson) {
	HttpResponse(ctx, buildStatus(resp, http.StatusInternalServerError), resp)

}
