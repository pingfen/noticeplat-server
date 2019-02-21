package v1

import (
	"github.com/pingfen/noticeplat-server/server/app"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.RouterGroup = app.Get().Group("/msgpack/api/v1")
)
