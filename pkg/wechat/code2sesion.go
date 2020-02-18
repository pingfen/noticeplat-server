package wechat

import (
	"fmt"
	"gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
)

func GetMiniProgramSession(source, code string) (session *oauth2.Session, err error) {
	if endpoint := GetOAuthEndpoint(source); endpoint != nil {
		return oauth2.GetSession(endpoint, code)
	}

	return nil, fmt.Errorf("source invalid")
}
