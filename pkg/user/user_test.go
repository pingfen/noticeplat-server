package user

import (
	"context"
	"fmt"
	"github.com/bingbaba/storage/qcloud-cos"
	"github.com/pingfen/noticeplat-server/pkg/storage"
	"sync"
	"testing"
	"time"
)

func init() {
	storage.Init(cos.NewStorage(cos.NewConfigByEnv()))
}

func newTestUser() *User {
	return &User{
		NickName:   "aaa",
		Sex:        1,
		City:       "Qingdao",
		Country:    "中国",
		Province:   "Shandong",
		HeadimgUrl: "",
		UnionId:    "",

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
}
func TestRegister(t *testing.T) {

	asyncLimit := make(chan bool, 20)
	var wg sync.WaitGroup
	for i := 1; i <= 200; i++ {
		asyncLimit <- true
		wg.Add(1)
		start := time.Now()

		go func(i int) {
			defer func() {
				wg.Done()
				<-asyncLimit
			}()

			u := newTestUser()
			u.UnionId = fmt.Sprintf("%09d", i)
			u.NickName = fmt.Sprintf("%09d", i)
			u.WechatProducts["noticeplat"].Openid = fmt.Sprintf("openid-noticeplat-%09d", i)
			u.WechatProducts["noticeplat-mp"].Openid = fmt.Sprintf("openid-noticeplat-mp-%09d", i)
			err := Register(context.Background(), u)
			if err != nil {
				t.Fatalf("%+v", err)
			}
			fmt.Printf("[%d] %s\n", i, time.Now().Sub(start).String())
		}(i)
	}

	wg.Wait()

}

//func TestUpdate(t *testing.T) {
//
//	u := &User{
//		UnionId:  "unionid-xxxxxxxxxxxxxx",
//		NickName: "bbb",
//		WechatProducts: map[string]*WechatProduct{
//			"noticeplat-mp": {
//				Openid:         "openid-zzzzzzzzzzzzzzzzzz",
//				SubscribeScene: "ADD_SCENE_SEARCH",
//				Subscribe:      1,
//				SubscribeTime:  time.Now().Unix(),
//			},
//		},
//	}
//
//	err := Update(context.Background(), u)
//	if err != nil {
//		t.Fatalf("%+v", err)
//	}
//}
