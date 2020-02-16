package wechat

import "gopkg.in/chanxuehong/wechat.v2/mp/oauth2"

func GetMiniProgramSession(code string) (session *oauth2.Session, err error) {
	return oauth2.GetSession(miniprogramOAuthEndpoint, code)
}
