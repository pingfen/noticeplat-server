package todo

//
//import (
//	"context"
//	cos "github.com/bingbaba/storage/qcloud-cos"
//	"github.com/pingfen/notify/pkg/storage"
//	"testing"
//	"time"
//)
//
//func init() {
//	storage.Init(cos.NewStorage(cos.NewConfigByEnv()))
//}
//
//func TestTodo(t *testing.T) {
//	todo := &Todo{
//		Owner:   "000000001",
//		Type:    TODOTYPE_ALARM,
//		Level:   TODOLEVEL_INFO,
//		ID:      "",
//		Subject: "this is subject",
//		Content: "this is content...........",
//		Labels: map[string]string{
//			"project": "A",
//		},
//		Operator:        "",
//		HistoryOperator: []string{},
//		StartTime:       time.Now().Unix(),
//		CreateTime:      time.Now().Unix(),
//	}
//
//	// Create
//	err := Create(context.Background(), todo)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	err = Create(context.Background(), todo)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// List
//	todoList, err := List(context.Background(), todo.Owner)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if len(todoList) < 2 {
//		t.Fatalf("expect get 2 todo, but get %d", len(todoList))
//	}
//}
