package msg

import (
	"context"
	"testing"
	"time"

	"fmt"
	"github.com/bingbaba/storage/qcloud-cos"
)

func TestMsg(t *testing.T) {
	m := &Message{
		ID:        "1",
		MsgType:   MSGTYPE_ALARM_START,
		Level:     LEVEL_WARN,
		Project:   "MyProject",
		Title:     "测试告警",
		Detail:    "这是一条测试告警的详细内容",
		Target:    "a",
		Timestamp: time.Now().Unix(),
	}
	ctx := context.Background()

	store := cos.NewStorage(cos.NewConfigByEnv())
	err := Post(ctx, store, "a", m)
	if err != nil {
		t.Fatal(err)
	}

	todo_list, err := TodoList(ctx, store, "a")
	if err != nil {
		t.Fatal(err)
	}
	if len(todo_list) == 0 {
		t.Fatal("list todo list failed")
	}

	if todo_list[0].ID != "1" {
		t.Fatalf("expect id 1, but get %s", todo_list[0].ID)
	}

	m.MsgType = MSGTYPE_ALARM_RECOVER
	err = Post(ctx, store, "a", m)
	if err != nil {
		t.Fatal(err)
	}

	history_list, err := HistoryList(ctx, store, "a")
	if err != nil {
		t.Fatal(err)
	}
	if len(history_list) == 0 {
		t.Fatal("list history failed")
	}
	for i, h := range history_list {
		fmt.Printf("[%d] %+v\n", i, h)
	}

	if history_list[0].ID != "1" {
		t.Fatalf("expect id 1, but get %s", history_list[0].ID)
	}
}
