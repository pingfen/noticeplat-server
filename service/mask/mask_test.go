package mask

import (
	"context"
	"github.com/pingfen/notify/pkg/models"
	"log"
	"net/http"
	"testing"
)

func TestMask(t *testing.T) {
	c, err := GetCaptcha(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if c.ContentType != "image/jpeg" {
		t.Fatalf("expect get image/jpeg, but get %s", c.ContentType)
	}

	if len(c.Cookies) == 0 {
		t.Fatalf("expect get captcha cookie, but not")
	}

	ms, err := GetMaskStoreInventory(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(ms) == 0 {
		t.Fatalf("expect get all mask stores inventory, but not")
	}

	for _, m := range ms {
		log.Printf("%+v", m)
	}
}

func TestReserve(t *testing.T) {
	u := &models.MaskUserInfo{"1", "王重胜", "370201199410308213", "13355329830", nil} // 瞎编的数据
	ms := &models.MaskStore{
		Area:    "李沧区",
		Address: "李沧区铜川路216号负一层超市",
		OrgCode: "91370213MA3DQENL2P",
		Name:    "丽达绿城店",
		ID:      132,
	}

	//_, _, cookies, err := GetCaptcha(context.Background())
	//if err != nil {
	//	t.Fatal(err)
	//}

	cookies := []*http.Cookie{
		{
			Name:  "acw_tc",
			Value: "2f61f26715814358833528761e28d0ef1735797764206d29363e62b37c84f4",
			Path:  "/",
		},
		{
			Name:  "_jfinal_captcha",
			Value: "7de663fec320478da76078d751ae5357",
			Path:  "/",
		},
	}

	err := Reserve(context.Background(), u, ms, &models.Captcha{Data: "TEWM", Cookies: cookies})
	if err != nil {
		t.Fatal(err)
	}

}
