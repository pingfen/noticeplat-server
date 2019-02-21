package wechat

import (
	"errors"
	"os"

	"gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
)

var (
	e *oauth2.Endpoint
)

func init() {
	appid := os.Getenv("WX_APP_APPID")
	app_secret := os.Getenv("WX_APP_APPSECRET")
	if appid != "" && app_secret != "" {
		e = oauth2.NewEndpoint(appid, app_secret)
	}
}

func GetSession(code string) (*oauth2.Session, error) {
	if e == nil {
		return nil, errors.New("env \"WX_APP_APPID\" and \"WX_APP_APPSECRET\" not found")
	}
	return oauth2.GetSession(e, code)
}
