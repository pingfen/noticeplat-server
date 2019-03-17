package user

import (
	"context"
	"github.com/bingbaba/storage/qcloud-cos"
	"github.com/pingfen/noticeplat-server/pkg/storage"
	"testing"
	"time"
)

func init() {
	storage.Init(cos.NewStorage(cos.NewConfigByEnv()))
}

func TestRegister(t *testing.T) {
	u := &User{
		NickName:   "aaa",
		Sex:        1,
		City:       "Qingdao",
		Country:    "中国",
		Province:   "Shandong",
		HeadimgUrl: "",
		UnionId:    "unionid-xxxxxxxxxxxxxx",

		WechatProducts: map[string]*WechatProduct{
			"noticeplat": {
				Openid:         "openid-xxxxxxxxxxxxxxxxx",
				SubscribeScene: "ADD_SCENE_SEARCH",
				Subscribe:      1,
				SubscribeTime:  time.Now().Unix(),
			},
			"noticeplat-mp": {
				Openid:         "openid-yyyyyyyyyyyyyyyyy",
				SubscribeScene: "ADD_SCENE_SEARCH",
				Subscribe:      1,
				SubscribeTime:  time.Now().Unix(),
			},
		},
	}

	err := Register(context.Background(), u)
	if err != nil {
		t.Fatalf("%+v", err)
	}
}

func TestUpdate(t *testing.T) {

	u := &User{
		UnionId:  "unionid-xxxxxxxxxxxxxx",
		NickName: "bbb",
		WechatProducts: map[string]*WechatProduct{
			"noticeplat-mp": {
				Openid:         "openid-zzzzzzzzzzzzzzzzzz",
				SubscribeScene: "ADD_SCENE_SEARCH",
				Subscribe:      1,
				SubscribeTime:  time.Now().Unix(),
			},
		},
	}

	err := Update(context.Background(), u)
	if err != nil {
		t.Fatalf("%+v", err)
	}
}
