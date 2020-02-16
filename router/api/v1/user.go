package v1

import (
	"github.com/pingfen/notify/pkg/e"
	"github.com/pingfen/notify/router/utils"
	"github.com/pingfen/notify/service/user"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type LoginReq struct {
	Code string `json:"code"`
}

type LoginResp struct {
	Session string `json:"session"`
	UserID  string `json:"userId"`
}

var (
	signingKey = []byte("notify_wx_2019")
)

func LoginByMiniprogram(ctx *gin.Context) {
	lr := new(LoginReq)

	// request
	err := ctx.BindJSON(lr)
	if err != nil {
		utils.SendErrResp(ctx, e.COMMON_BADREQUEST, err.Error())
		return
	}
	if lr.Code == "" {
		utils.SendErrResp(ctx, &e.Err{e.COMMON_BADREQUEST.Code, "code不可为空"}, "")
		return
	}

	// openid
	userInfo, err := user.LoginByCode(ctx.Request.Context(), "notodo", lr.Code)
	if err != nil {
		utils.SendErrResp(ctx, &e.Err{e.COMMON_BADREQUEST.Code, err.Error()}, err.Error())
		return
	}

	// user
	user_id := userInfo.OpenId

	// create session
	now := time.Now()
	claims := &jwt.StandardClaims{
		Audience:  "miniprogram",
		ExpiresAt: now.Add(time.Hour * 24 * 30).Unix(),
		Issuer:    "notodo",
		Subject:   user_id,
		IssuedAt:  now.Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signingKey)
	if err != nil {
		utils.SendErrResp(ctx, &e.Err{e.COMMON_INTERNAL_ERR.Code, "签名失败，请稍后重试"}, err.Error())
		return
	}

	utils.SendData(ctx, token)
}

func UserProfile(ctx *gin.Context) {
	if claims := GetClaims(ctx); claims != nil {
		utils.SendData(ctx, claims)
		return
	}

	utils.SendErrResp(ctx, e.AUTH_NOTLOGIN, "")
}
