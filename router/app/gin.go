package app

import (
	"github.com/gin-gonic/gin"
	"github.com/haier-interx/gogin"
)

var (
	app *gin.Engine
)

func init() {
	app = gin.New()
	app.Use(gogin.Logger())
	app.Use(gogin.Recovery())

}

func Get() *gin.Engine {
	return app
}
