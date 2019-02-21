package v1

import "github.com/pingfen/noticeplat-server/server/handler"

func init() {
	router.POST("/login", handler.UserLogin)

	router.GET("/users/:id", handler.GetUser)
	router.POST("/users", handler.UserRegister)
	router.PUT("/users/:id", handler.UserModify)
	router.GET("/users/:id/avatar", handler.GetUserAvatar)
	router.GET("/users/:id/exist", handler.UserIdIsExist)
	router.POST("/users/:id/service/:openid", handler.UserSrvBinding)
}
