package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	resp2 "github.com/pingfen/notify/pkg/resp"
)

func SendData(ctx *gin.Context, v interface{}) {
	SendResp(ctx, resp2.NewSucResponse(v))
}

func SendErrResp(ctx *gin.Context, err error, detail string) {
	SendResp(ctx, resp2.NewErrResponse(err, detail))
}

func SendResp(ctx *gin.Context, resp *resp2.Response) {
	if resp.Code != 10000 {
		if resp.Detail != "" {
			PushCtxErr(ctx, fmt.Errorf("%d :: %s :: %s", resp.Code, resp.Message, resp.Detail))
		} else {
			PushCtxErr(ctx, fmt.Errorf("%d :: %s", resp.Code, resp.Message))
		}
		ctx.Abort()
	}

	//if !GetQueryParamBool(ctx, "debug") {
	//	resp.Detail = ""
	//}

	if GetQueryParamBool(ctx, "pretty") || GetQueryParamBool(ctx, "debug") {
		ctx.IndentedJSON(200, resp)
	} else {
		ctx.JSON(200, resp)
	}
}

func PushCtxErr(ctx *gin.Context, err error) {
	ctx.Error(gin.Error{Err: err, Type: gin.ErrorTypePrivate})
}
