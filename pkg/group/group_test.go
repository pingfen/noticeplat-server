package group

import (
	"context"
	"github.com/bingbaba/storage/qcloud-cos"
	"github.com/pingfen/noticeplat-server/pkg/errs"
	"github.com/pingfen/noticeplat-server/pkg/storage"
	"github.com/pkg/errors"
	"testing"
)

func init() {
	storage.Init(cos.NewStorage(cos.NewConfigByEnv()))
}
func TestGroup(t *testing.T) {
	g := &Group{
		Id:      "111111111111111",
		OpenGid: "opengid-xxxxxxxxxxxxxx",
		Name:    "mygroup",
	}
	m := &Member{
		UnionId:   "unionid-xxxxxxxxxxxxxx",
		NickName:  "昵称A",
		ShowPhone: true,
	}

	// create group
	err := Create(context.Background(), g, m)
	if err != nil {
		switch v := errors.Cause(err).(type) {
		case errs.BadRequest:
			if v != errs.OBJECT_EXIST {
				t.Fatalf("%+v", err)
			}
		default:
			t.Fatalf("%+v", err)
		}
	}

	// update group
	g.Name = "群组A"
	err = Update(context.Background(), g)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// list members
	ms, err := ListMembers(context.Background(), g.Id)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if len(ms) != 1 || ms[0].UnionId != "unionid-xxxxxxxxxxxxxx" {
		t.Fatalf("list member failed, %+v", ms[0])
	}
}
