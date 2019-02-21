package handler

import (
	"net/http"

	"github.com/pingfen/noticeplat-server/pkg/group"
	"github.com/bingbaba/storage"
	"github.com/gin-gonic/gin"
)

func GetGroup(ctx *gin.Context) {
	id := ctx.Param("id")
	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	g, err := group.Get(ctx, store, id)
	if err != nil {
		SendBad(ctx, err)
		return
	}
	SendWithData(ctx, g)
}

func AddMember(ctx *gin.Context) {
	id := ctx.Param("id")
	openid := ctx.Param("openid")
	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	err := group.AddMember(ctx, store, id, openid)
	if err != nil {
		SendBad(ctx, err)
		return
	}
	SendWithData(ctx, nil)
}

func RemoveMember(ctx *gin.Context) {
	id := ctx.Param("id")
	openid := ctx.Param("openid")
	si, found := ctx.Get("storage")
	if !found {
		ctx.JSON(http.StatusInternalServerError, Response{"StorageError", "can't found storage from context", nil})
		return
	}
	store := si.(storage.Interface)

	err := group.RemoveMember(ctx, store, id, openid)
	if err != nil {
		SendBad(ctx, err)
		return
	}
	SendWithData(ctx, nil)
}
