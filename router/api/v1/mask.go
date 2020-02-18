package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pingfen/notify/pkg/e"
	"github.com/pingfen/notify/pkg/models"
	"github.com/pingfen/notify/pkg/resp"
	"github.com/pingfen/notify/router/utils"
	"github.com/pingfen/notify/service/mask"
)

func ListMaskStores(ctx *gin.Context) {
	utils.SendData(ctx, mask.ListMaskStores())
}

func GetMaskStore(ctx *gin.Context) {
	sid := ctx.Param("sid")
	ms, err := mask.GetMaskStoreInventory(ctx)
	if err != nil {
		utils.SendErrResp(ctx, &e.Err{e.COMMON_INTERNAL_ERR.Code, "查询剩余量失败"}, err.Error())
		return
	}

	var maskStore *models.MaskStore
	for _, m := range ms {
		if fmt.Sprintf("%d", m.ID) == sid {
			maskStore = m
			break
		}
	}
	utils.SendData(ctx, maskStore)
}

func AddContact(ctx *gin.Context) {
	uid := GetUID(ctx)
	c := new(models.MaskUserInfo)
	err := ctx.BindJSON(c)
	if err != nil {
		utils.SendErrResp(ctx, e.COMMON_BADREQUEST, err.Error())
		return
	}

	if c.Name == "" || c.CardNo == "" || c.Mobile == "" {
		utils.SendErrResp(ctx, &e.Err{e.COMMON_PARAM_MISS.Code, "填写不能为空"}, "")
		return
	}

	if err := mask.SaveContacts(ctx.Request.Context(), uid, c); err != nil {
		utils.SendErrResp(ctx, &e.Err{e.COMMON_INTERNAL_ERR.Code, "保存用户信息失败"}, err.Error())
		return
	}

	utils.SendData(ctx, nil)
}

func ListContacts(ctx *gin.Context) {
	uid := GetUID(ctx)

	ret, err := mask.ListContacts(ctx.Request.Context(), uid)
	if err != nil {
		utils.SendErrResp(ctx, &e.Err{e.COMMON_INTERNAL_ERR.Code, "查询联系人列表失败"}, err.Error())
		return
	}

	utils.SendData(ctx, ret)
}

func GetContact(ctx *gin.Context) {
	uid := GetUID(ctx)
	cid := ctx.Param("cid")

	ret, err := mask.GetContact(ctx.Request.Context(), uid, cid)
	if err != nil {
		utils.SendErrResp(ctx, &e.Err{e.COMMON_INTERNAL_ERR.Code, "查询联系人失败"}, err.Error())
		return
	}

	utils.SendData(ctx, ret)
}

func DelContact(ctx *gin.Context) {
	uid := GetUID(ctx)
	cid := ctx.Param("cid")

	err := mask.DeleteContact(ctx.Request.Context(), uid, cid)
	if err != nil {
		utils.SendErrResp(ctx, &e.Err{e.COMMON_INTERNAL_ERR.Code, "删除联系人失败"}, err.Error())
		return
	}

	utils.SendData(ctx, nil)
}

func GetCaptcha(ctx *gin.Context) {
	c, err := mask.GetCaptcha(ctx.Request.Context())
	if err != nil {
		utils.SendErrResp(ctx, &e.Err{e.COMMON_INTERNAL_ERR.Code, "获取验证码失败"}, err.Error())
		return
	}

	utils.SendData(ctx, c)
}

func ReserveMask(ctx *gin.Context) {
	uid := GetUID(ctx)

	req := new(models.ReserveRequest)
	if err := ctx.BindJSON(req); err != nil {
		utils.SendErrResp(ctx, &e.Err{e.COMMON_BADREQUEST.Code, "请求格式错误"}, err.Error())
		return
	}

	if req.Contact != "" || req.ContactInfo != nil {
		err := mask.DealReserveRequest(ctx, uid, req)
		if err != nil {
			utils.SendErrResp(ctx, &e.Err{e.COMMON_BADREQUEST.Code, err.Error()}, err.Error())
			return
		}
	} else {
		utils.SendErrResp(ctx, &e.Err{e.COMMON_BADREQUEST.Code, "缺少联系人信息"}, "")
		return
	}

	//utils.SendData(ctx, "预约成功，请等待短信通知")
	utils.SendResp(ctx, &resp.Response{
		&resp.BaseResponse{
			e.COMMON_SUC.Code,
			"预约成功，请等待短信通知",
			"",
		},
		nil,
	})
}
