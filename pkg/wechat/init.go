package wechat

import (
	"net/http"
	"os"

	"gopkg.in/chanxuehong/wechat.v2/mp/core"
)

var (
	client *core.Client
	APPID  string
)

func init() {
	APPID = os.Getenv("WX_SRV_APPID")
	if APPID == "" {
		panic("WX_SRV_APPID was required")
	}

	app_secret := os.Getenv("WX_SRV_APPSECRET")
	if app_secret == "" {
		panic("WX_SRV_APPSECRET was required")
	}

	http_client := http.DefaultClient
	ats := core.NewDefaultAccessTokenServer(
		APPID,
		app_secret,
		http_client,
	)
	client = core.NewClient(ats, http_client)
}
