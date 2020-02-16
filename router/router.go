package router

import (
	_ "github.com/pingfen/notify/docs"
	_ "github.com/pingfen/notify/router/api/v1"
	"github.com/pingfen/notify/router/app"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func init() {
	app.Get().GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func Run(addr ...string) error {
	return app.Get().Run(addr...)
}
