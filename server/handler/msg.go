package handler

import (
	"net/http"

	"fmt"
	"github.com/pingfen/noticeplat-server/pkg/msg"
	"github.com/pingfen/noticeplat-server/pkg/user"
	"github.com/pingfen/noticeplat-server/pkg/wechat"
	"github.com/bingbaba/storage"
	"github.com/gin-gonic/gin"
	"time"
)

func PostMsg(ctx *gin.Context) {
	clien_type := ctx.Value("ClientType").(msg.ClientType)
	uid_gid := ctx.Param("id")

	// storage
	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	m := new(msg.Message)
	err := ctx.BindJSON(m)
	if err != nil {
		SendBad(ctx, err)
		return
	}

	// check user id
	m.Target = uid_gid
	u, err := user.Get(ctx, store, uid_gid)
	if err != nil {
		if storage.IsNotFound(err) {
			ctx.JSON(http.StatusBadRequest, Response{"UserNotFound", "the user was not register", nil})
			return
		} else {
			SendBad(ctx, err)
			return
		}
	}

	// check message
	err = m.Check()
	if err != nil {
		SendBad(ctx, err)
		return
	}

	// post message
	todo, err := msg.Post(ctx, store, uid_gid, clien_type, m)
	if err != nil {
		SendBad(ctx, err)
		return
	}

	// send notice
	if todo.Alerting {
		err = wechat.SendTemplateMsg(u.SrvOpenId, todo)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, Response{"SendMsgFailed", err.Error(), nil})
			return
		}
	}

	// response
	SendWithData(ctx, todo)
}

func ModifyTodo(ctx *gin.Context) {
	ct := ctx.Value("ClientType").(msg.ClientType)

	// uid
	uid := ctx.Param("id")
	mid := ctx.Param("mid")

	todo := new(msg.Todo)
	err := ctx.BindJSON(todo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{"BadRequest", err.Error(), nil})
		return
	}
	if todo.Message == nil {
		ctx.JSON(http.StatusBadRequest, Response{"BadRequest", "miss message info", nil})
		return
	}

	todo.Owner = uid
	todo.ClientType = ct
	todo.Message.ID = mid

	// storage
	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	err = msg.Modify(ctx, store, todo)
	if err != nil {
		SendBad(ctx, err)
		return
	}

	// response
	SendWithData(ctx, todo)
}

func GetTodo(ctx *gin.Context) {
	ct := ctx.Value("ClientType").(msg.ClientType)

	// uid
	uid := ctx.Param("id")
	mid := ctx.Param("mid")

	// storage
	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	todo, err := msg.GetTodo(ctx, store, uid, mid, ct)
	if err != nil {
		SendBad(ctx, err)
		return
	}

	// response
	SendWithData(ctx, todo)
}

func ListTodo(ctx *gin.Context) {
	ct := ctx.Value("ClientType").(msg.ClientType)

	// uid or gid
	uid_gid := ctx.Param("id")

	// storage
	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	todo_list, err := msg.TodoList(ctx, store, uid_gid, ct)
	if err != nil {
		SendBad(ctx, err)
		return
	}

	// response
	SendWithData(ctx, todo_list)
}

func PostMsgComment(ctx *gin.Context) {
	ct := ctx.Value("ClientType").(msg.ClientType)

	// uid
	uid := ctx.Param("id")
	mid := ctx.Param("mid")

	comment := new(msg.Comments)
	err := ctx.BindJSON(comment)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{"BadRequest", err.Error(), nil})
		return
	}
	if comment.Type == msg.CommentsType_User && comment.Operator == "" {
		ctx.JSON(http.StatusBadRequest, Response{"BadRequest", "operator is required", nil})
		return
	}
	comment.Id = fmt.Sprintf("%d.%s", time.Now().UnixNano(), comment.Operator)
	comment.MsgOwner = uid
	comment.MsgId = mid
	comment.ClientType = ct

	// storage
	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	// save
	err = msg.PostComments(ctx, store, comment)
	if err != nil {
		SendBad(ctx, err)
		return
	}

	// response
	SendWithData(ctx, comment)
}

func ModifyMsgComment(ctx *gin.Context) {
	// uid
	uid := ctx.Param("id")
	mid := ctx.Param("mid")
	cid := ctx.Param("cid")

	comment := new(msg.Comments)
	err := ctx.BindJSON(comment)
	if err != nil {
		SendBad(ctx, err)
		return
	}

	//if comment.MsgId != mid {
	//	ctx.JSON(http.StatusBadRequest, Response{"BadRequest", "msgid inconsistent", nil})
	//	return
	//}
	comment.MsgId = mid
	comment.Id = cid
	comment.MsgOwner = uid

	// storage
	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	err = msg.ModifyComments(ctx, store, comment)
	if err != nil {
		SendBad(ctx, err)
		return
	}

	SendWithData(ctx, comment)
}

func DeleteMsgComment(ctx *gin.Context) {
	ct := ctx.Value("ClientType").(msg.ClientType)

	// uid
	uid := ctx.Param("id")
	mid := ctx.Param("mid")
	cid := ctx.Param("cid")

	// storage
	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	err := msg.DeleteComments(ctx, store, uid, mid, cid, ct)
	if err != nil {
		SendBad(ctx, err)
		return
	}

	SendWithData(ctx, nil)
}

func ListMsgComments(ctx *gin.Context) {
	ct := ctx.Value("ClientType").(msg.ClientType)

	// uid
	uid := ctx.Param("id")
	mid := ctx.Param("mid")

	// storage
	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	cs, err := msg.ListCommentses(ctx, store, uid, mid, ct)
	if err != nil {
		SendBad(ctx, err)
		return
	}

	SendWithData(ctx, cs)
}

func MsgCommentsGood(ctx *gin.Context) {
	ct := ctx.Value("ClientType").(msg.ClientType)

	// uid
	uid := ctx.Param("id")
	mid := ctx.Param("mid")
	cid := ctx.Param("cid")

	// storage
	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	count, err := msg.CommentsGood(ctx, store, uid, mid, cid, ct)
	if err != nil {
		SendBad(ctx, err)
		return
	}

	SendWithData(ctx, count)
}

func MsgCommentsBad(ctx *gin.Context) {
	ct := ctx.Value("ClientType").(msg.ClientType)

	// uid
	uid := ctx.Param("id")
	mid := ctx.Param("mid")
	cid := ctx.Param("cid")

	// storage
	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	count, err := msg.CommentsBad(ctx, store, uid, mid, cid, ct)
	if err != nil {
		SendBad(ctx, err)
		return
	}

	SendWithData(ctx, count)
}
