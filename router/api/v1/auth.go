package v1

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pingfen/notify/pkg/e"
	"github.com/pingfen/notify/router/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	CONTEXT_OPENID_TAG = "OpenID"
	CONTEXT_UNION_TAG  = "UnionID"
)

func AuthN(ctx *gin.Context) {
	uid, err := getUidFromJwt(ctx)
	if err != nil {
		utils.SendErrResp(ctx, e.AUTH_NOTLOGIN, err.Error())
		ctx.Abort()
		return
	}

	//u, err := storage.GlobalStore.GetUser(uid)
	//if err != nil {
	//	utils.SendErrResp(ctx, http.StatusUnauthorized, "read user info failed", err.Error())
	//	return
	//}

	ctx.Set(CONTEXT_OPENID_TAG, uid)
	//ctx.Set(CONTEXT_OPENID_TAG, u.OpenId)
	//ctx.Set(CONTEXT_UNION_TAG, u.UnionId)

	ctx.Next()
}

func getUidFromJwt(ctx *gin.Context) (uid string, err error) {
	ss := getToken(ctx)
	if ss == "" {
		return "", fmt.Errorf("token not found")
	}

	token, err := jwt.ParseWithClaims(ss, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("parse jwt failed %w", err)
	}
	claims := token.Claims.(*jwt.StandardClaims)
	err = claims.Valid()
	if err != nil {
		return "", fmt.Errorf("jwt invalid %w", err)
	}

	if claims.Issuer != "notodo" {
		return "", fmt.Errorf("please login again")
	}

	ctx.Set("claims", claims)
	return claims.Subject, nil
}

func getToken(ctx *gin.Context) string {
	auth_info := ctx.GetHeader("Authorization")
	if auth_info != "" {
		if strings.HasPrefix(auth_info, "Bearer ") {
			return strings.TrimPrefix(auth_info, "Bearer ")
		}
	}
	return ""
}

//func GetOpenID(ctx *gin.Context) string {
//	return ctx.GetString(CONTEXT_OPENID_TAG)
//}
//
//func GetUnionID(ctx *gin.Context) string {
//	return ctx.GetString(CONTEXT_UNION_TAG)
//}

func GetUID(ctx *gin.Context) string {
	return getUID(ctx)
}

func GetClaims(ctx *gin.Context) *jwt.StandardClaims {
	ret, found := ctx.Get("claims")
	if !found {
		return nil
	}
	return ret.(*jwt.StandardClaims)
}

func getUID(ctx *gin.Context) (uid string) {
	// uid = ctx.GetString(CONTEXT_UNION_TAG)
	if ctx.GetString(CONTEXT_OPENID_TAG) != "" {
		return ctx.GetString(CONTEXT_OPENID_TAG)
	}
	return ctx.Param("uid")
}
