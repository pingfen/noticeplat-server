package app

import (
	"github.com/gin-gonic/gin"
	"github.com/pingfen/noticeplat-server/server/middleware/storage"
)

var (
	app *gin.Engine
)

func init() {
	app = gin.New()

	app.Use(gin.Logger())

	app.Use(gin.Recovery())

	app.Use(storage.Storage)
}

func Get()*gin.Engine {
	return app
}
