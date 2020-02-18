package wechat

import (
	"gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
	"os"

	"gopkg.in/chanxuehong/wechat.v2/mp/core"
)

var (
	// 服务号
	SRV_APPID     string
	SRV_APPSECRET string
	srvClient     *core.Client

	// 小程序
	MINIPROGRAM_APPID        string
	MINIPROGRAM_APPSECRET    string
	miniprogramOAuthEndpoint *oauth2.Endpoint

	// 小程序
	MINIPROGRAM_QINDAOREN_APPID       string
	MINIPROGRAM_QINDAOREN_APPSECRET   string
	miniprogramOAuthEndpointQindaoren *oauth2.Endpoint
)

func init() {
	//// 服务号
	//SRV_APPID = os.Getenv("WX_SRV_APPID")
	//if SRV_APPID == "" {
	//	panic("WX_SRV_APPID was required")
	//}
	//
	//SRV_APPSECRET = os.Getenv("WX_SRV_APPSECRET")
	//if SRV_APPSECRET == "" {
	//	panic("WX_SRV_APPSECRET was required")
	//}
	//
	//http_client := http.DefaultClient
	//ats := core.NewDefaultAccessTokenServer(
	//	SRV_APPID,
	//	SRV_APPSECRET,
	//	http_client,
	//)
	//srvClient = core.NewClient(ats, http_client)

	// 小程序
	MINIPROGRAM_APPID = os.Getenv("WX_MINIPROGRAM_APPID")
	if MINIPROGRAM_APPID == "" {
		panic("WX_MINIPROGRAM_APPID was required")
	}

	MINIPROGRAM_APPSECRET = os.Getenv("WX_MINIPROGRAM_APPSECRET")
	if MINIPROGRAM_APPSECRET == "" {
		panic("WX_MINIPROGRAM_APPSECRET was required")
	}
	miniprogramOAuthEndpoint = oauth2.NewEndpoint(MINIPROGRAM_APPID, MINIPROGRAM_APPSECRET)

	// 小程序
	MINIPROGRAM_QINDAOREN_APPID = os.Getenv("WX_MINIPROGRAM_QINDAOREN_APPID")
	if MINIPROGRAM_QINDAOREN_APPID == "" {
		panic("WX_MINIPROGRAM_QINDAOREN_APPID was required")
	}

	MINIPROGRAM_QINDAOREN_APPSECRET = os.Getenv("WX_MINIPROGRAM_QINDAOREN_APPSECRET")
	if MINIPROGRAM_QINDAOREN_APPSECRET == "" {
		panic("WX_MINIPROGRAM_QINDAOREN_APPSECRET was required")
	}
	miniprogramOAuthEndpointQindaoren = oauth2.NewEndpoint(MINIPROGRAM_QINDAOREN_APPID, MINIPROGRAM_QINDAOREN_APPSECRET)
}

func GetOAuthEndpoint(source string) *oauth2.Endpoint {
	switch source {
	case "notodo":
		return miniprogramOAuthEndpoint
	case "qindaoren":
		return miniprogramOAuthEndpointQindaoren
	default:
		return nil
	}
}

func SourceIsValid(source string) bool {
	switch source {
	case "notodo", "qindaoren":
		return true
	default:
		return false
	}
}
