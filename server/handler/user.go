package handler

import (
	"net/http"

	"github.com/pingfen/noticeplat-server/pkg/user"
	"github.com/pingfen/noticeplat-server/pkg/wechat"
	"github.com/bingbaba/storage"
	"github.com/gin-gonic/gin"
)

type loginData struct {
	Code string `json:"code"`
}

func UserLogin(ctx *gin.Context) {
	ld := new(loginData)
	err := ctx.BindJSON(ld)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{"ParseJsonFailed", "can't found code: " + err.Error(), nil})
		return
	}

	sess, err := wechat.GetSession(ld.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{"LoginFailed", err.Error(), nil})
		return
	}

	SendWithData(ctx, sess)
}

func GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	var byOpenid bool
	byOpenidStr, found := ctx.GetQuery("openid")
	if found {
		if byOpenidStr == "0" || byOpenidStr == "False" || byOpenidStr == "false" {
			byOpenid = false
		} else {
			byOpenid = true
		}
	} else {
		byOpenid = false
	}

	var u *user.User
	var err error
	if byOpenid {
		u, err = user.GetByOpenId(ctx, store, id)
	} else {
		u, err = user.Get(ctx, store, id)
	}
	if err != nil {
		SendBad(ctx, err)
		return
	}
	SendWithData(ctx, u)
}

func GetUserAvatar(ctx *gin.Context) {
	id := ctx.Param("id")
	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	u, err := user.Get(ctx, store, id)
	if err != nil {
		SendBad(ctx, err)
		return
	}
	SendWithData(ctx, u.AvatarUrl)
}

func UserRegister(ctx *gin.Context) {
	u := new(user.User)
	err := ctx.BindJSON(u)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{"BadRequest", err.Error(), nil})
		return
	}

	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	// check uid
	found, err = user.IdIsExist(ctx, store, u.Id)
	if found {
		ctx.JSON(http.StatusConflict, Response{"UserHasResiter", "the user has register", nil})
		return
	}

	err = user.Add(ctx, store, u)
	if err != nil {
		SendBad(ctx, err)
		return
	}

	SendOK(ctx)
}

func UserModify(ctx *gin.Context) {
	u := new(user.User)
	err := ctx.BindJSON(u)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{"BadRequest", err.Error(), nil})
		return
	}

	store, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}

	err = user.Modify(ctx, store.(storage.Interface), u)
	if err != nil {
		SendBad(ctx, err)
		return
	}

	SendOK(ctx)
}

func UserIdIsExist(ctx *gin.Context) {
	resp := map[string]bool{"exist": true}

	uid := ctx.Param("id")
	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	exist, err := user.IdIsExist(ctx, store, uid)
	if err != nil {
		SendBad(ctx, err)
		return
	}

	resp["exist"] = exist
	SendWithData(ctx, resp)
}

func UserSrvBinding(ctx *gin.Context) {
	uid := ctx.Param("id")
	srvOpenId := ctx.Param("openid")

	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	//u, err := user.Get(ctx, store, uid)
	//if err != nil {
	//	if storage.IsNotFound(err) {
	//		ctx.JSON(http.StatusBadRequest, Response{"BadRequest", err.Error(), nil})
	//		return
	//	} else {
	//		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", err.Error(), nil})
	//		return
	//	}
	//}

	err := user.BindingSrvOpenId(ctx, store, uid, srvOpenId)
	if err != nil {
		SendBad(ctx, err)
		return
	}
	SendOK(ctx)
}
