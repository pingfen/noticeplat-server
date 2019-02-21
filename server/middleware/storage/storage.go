package storage

import (
	"github.com/bingbaba/storage"
	"github.com/bingbaba/storage/qcloud-cos"
	"github.com/gin-gonic/gin"
)

var (
	s storage.Interface
)

func init() {
	s = cos.NewStorage(cos.NewConfigByEnv())
}
func Storage(ctx *gin.Context) {
	ctx.Set("storage", s)
}
