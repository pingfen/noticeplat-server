package wechat

import (
	"github.com/pingfen/noticeplat-server/pkg/todo"
	"testing"
	"time"
)

func TestSendSrvMsg(t *testing.T) {
	openid := "oqwWC1oB2ANQx9cV-vtV7ef_-yWc"

	todo := &todo.Todo{
		Owner: openid,
		Type:  todo.TODOTYPE_TASK,
		Level: todo.TODOLEVEL_CRITICAL,

		ID:      "",
		Subject: "后端开发某某某功能",
		Content: "这是任务描述信息，包含具体功能点与实现方法和注意点",
		Labels: map[string]string{
			"project": "PROJECT-A",
		},
		Operator:  openid,
		StartTime: time.Now().Unix(),
	}

	err := SendTemplateMsg(openid, todo)
	if err != nil {
		t.Fatal(err)
	}
}
