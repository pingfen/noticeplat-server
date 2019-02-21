package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bingbaba/storage"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string      `json:"message"`
	Detail  string      `json:"detail"`
	Data    interface{} `json:"data,omitempty"`
}

func NewResponse(data []byte) (*Response, error) {
	r := new(Response)
	err := json.Unmarshal(data, r)
	return r, err
}

func SendOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{Message: "OK", Detail: ""})
}

func SendBad(ctx *gin.Context, err error) {
	switch err.(type) {
	case *storage.StorageError:
		code, msg, detail := storage.ParseToHttpError(err)
		ctx.JSON(code, Response{Message: msg, Detail: detail})
	default:
		ctx.JSON(500, Response{Message: "InternalServerError", Detail: err.Error()})
	}
}

func SendWithData(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{Message: "OK", Detail: "", Data: data})
}
