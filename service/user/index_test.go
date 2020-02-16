package user

import (
	"context"
	"testing"
)

func TestGetByOpenid(t *testing.T) {
	u, err := GetByOpenid(context.Background(), "noticeplat", "openid-xxxxxxxxxxxxxxxxx")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	//log.Printf("%+v", u)
	if u.UnionId != "unionid-xxxxxxxxxxxxxx" {
		t.Fatalf("get unionid failed")
	}
}
