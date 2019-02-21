package user

import (
	"context"
	"testing"

	"github.com/bingbaba/storage/qcloud-cos"
)

func TestUser(t *testing.T) {
	u := &User{
		Id:        "a",
		AppOpenId: "a-app-openid",
		SrvOpenId: "a-srv-openid",
		Name:      "A",
		Province:  "Shandong",
		Phone:     "",
		Email:     "",
	}
	ctx := context.Background()

	store := cos.NewStorage(cos.NewConfigByEnv())
	err := Add(ctx, store, u)
	if err != nil {
		t.Fatal(err)
	}

	err = SrvBinding(ctx, store, u.Id, "a-srv-openid2")
	if err != nil {
		t.Fatal(err)
	}

	u2, err := Get(ctx, store, u.Id)
	if err != nil {
		t.Fatal(err)
	}

	if u2.SrvOpenID != "a-srv-openid2" {
		t.Fatalf("expect get \"a-srv-openid2\", but get \"%s\"\n", u2.SrvOpenId)
	}
}
