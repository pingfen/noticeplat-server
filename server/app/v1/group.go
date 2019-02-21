package v1

import (
	"github.com/pingfen/noticeplat-server/server/handler"
)

func init() {
	router.GET("/groups/:id", handler.GetGroup)
	router.POST("/groups/:id/members/:openid", handler.AddMember)
	router.DELETE("/groups/:id/members/:openid", handler.RemoveMember)
}
