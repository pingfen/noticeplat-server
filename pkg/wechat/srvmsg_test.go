package wechat

import (
	"github.com/pingfen/noticeplat-server/pkg/msg"
	"testing"
	"time"
)

func TestSendSrvMsg(t *testing.T) {
	todo := &msg.Todo{
		Message: &msg.Message{
			ID:        "00000001",
			MsgType:   msg.MSGTYPE_ALARM_START,
			Level:     msg.LEVEL_ERROR,
			Project:   "监控系统",
			Title:     "机器CPU利用率过高",
			Content:   "机器 127.0.0.1 CPU利用率过高,达到95%",
			Timestamp: time.Now().Unix(),
			Target:    "a",
		},
	}

	openid := "o3HH903-YBeNkhLmdomElPrmxLZI"
	err := SendTemplateMsg(openid, todo)
	if err != nil {
		t.Fatal(err)
	}
}
