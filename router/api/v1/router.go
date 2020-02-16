package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/pingfen/notify/router/app"
)

func init() {
	api := app.Get().Group("/api/v2/notify")
	{
		// 登录
		api.POST("/login-by-miniprogram", LoginByMiniprogram)

		authApi := api.Group("", AuthN)
		{
			authApi.GET("/profile", UserProfile)
		}

		// 用户信息
		maskBase := api.Group("/mask")
		{
			maskBase.GET("/captcha", GetCaptcha)
			maskBase.GET("/stores", ListMaskStores)
			maskBase.GET("/stores/:sid", GetMaskStore)
		}

		maskWithAuth := maskBase.Group("", AuthN)
		userMask := maskBase.Group("/users/:uid")
		for _, mask := range []*gin.RouterGroup{maskWithAuth, userMask} {
			mask.POST("/contacts", AddContact)
			mask.GET("/contacts", ListContacts)
			mask.GET("/contacts/:cid", GetContact)
			mask.PUT("/contacts/:cid", AddContact)
			mask.DELETE("/contacts/:cid", DelContact)
			mask.POST("/reserve", ReserveMask)
		}
	}

}
