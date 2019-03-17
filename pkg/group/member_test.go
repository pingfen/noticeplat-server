package group

import (
	"context"
	"testing"
)

func TestListDetailMembers(t *testing.T) {
	ms, err := ListDetailMembers(context.Background(), "111111111111111")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if len(ms) == 0 || ms[0].Detail.WechatProducts["noticeplat"].Openid != "openid-xxxxxxxxxxxxxxxxx" {
		t.Fatalf("noticeplat openid fault %+v", ms[0].Detail.WechatProducts["noticeplat"])
	}
}
