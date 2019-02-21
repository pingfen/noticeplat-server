package v1

import (
	"github.com/pingfen/noticeplat-server/pkg/msg"
	"github.com/pingfen/noticeplat-server/server/handler"
	"github.com/gin-gonic/gin"
)

func init() {
	user_router := router.Group("/users/:id",
		func(ctx *gin.Context) {
			ctx.Set("ClientType", msg.CLIENTTYPE_USER)
		},
	)

	// message
	user_router.POST("/messages", handler.PostMsg)

	// todo_info
	user_router.GET("/todo", handler.ListTodo)
	user_router.GET("/todo/:mid", handler.GetTodo)
	user_router.PUT("/todo/:mid", handler.ModifyTodo)

	// comment
	user_router.POST("/todo/:mid/comments", handler.PostMsgComment)
	user_router.GET("/todo/:mid/comments", handler.ListMsgComments)
	user_router.PUT("/todo/:mid/comments/:cid", handler.ModifyMsgComment)
	user_router.POST("/todo/:mid/comments/:cid/good", handler.MsgCommentsGood)
	user_router.POST("/todo/:mid/comments/:cid/bad", handler.MsgCommentsBad)
	user_router.DELETE("/todo/:mid/comments/:cid", handler.DeleteMsgComment)
}
