package models

import (
	"fmt"
	"github.com/pingfen/notify/pkg/e"
	"time"
)

const (
	TODOTYPE_TASK  TodoType = "task"
	TODOTYPE_ALARM TodoType = "alert"
	TODOTYPE_EVENT TodoType = "event"

	TODOLEVEL_DEBUG    TodoLevel = 1
	TODOLEVEL_INFO     TodoLevel = 2
	TODOLEVEL_WARNNING TodoLevel = 3
	TODOLEVEL_SERIOUS  TodoLevel = 4
	TODOLEVEL_CRITICAL TodoLevel = 5
)

type TodoType = string
type TodoLevel = int

type Todo struct {
	Owner string    `json:"owner"`
	Type  TodoType  `json:"type"`
	Level TodoLevel `json:"level"`

	ID              string            `json:"id"` // identify
	Subject         string            `json:"subject"`
	Content         string            `json:"content"`
	Labels          map[string]string `json:"labels"`
	Operator        string            `json:"operator"`
	HistoryOperator []string          `json:"historyOperator"`
	StartTime       int64             `json:"startTime"`
	CreateTime      int64             `json:"createTime"`
}

func LevelName(l TodoLevel) string {
	switch l {
	case TODOLEVEL_DEBUG:
		return "调试"
	case TODOLEVEL_INFO:
		return "一般"
	case TODOLEVEL_WARNNING:
		return "警告"
	case TODOLEVEL_SERIOUS:
		return "严重"
	case TODOLEVEL_CRITICAL:
		return "紧急"
	default:
		return "未知"
	}
}
func TypeName(t TodoType) string {
	switch t {
	case TODOTYPE_TASK:
		return "任务"
	case TODOTYPE_ALARM:
		return "告警"
	case TODOTYPE_EVENT:
		return "通知"
	default:
		return "通知"
	}
}

func (t *Todo) Valid() error {
	if t.Subject == "" {
		return fmt.Errorf("主题不可为空%w", e.COMMON_BADREQUEST)
	}
	if t.Content == "" {
		return fmt.Errorf("内容不可为空%w", e.COMMON_PARAM_MISS)
	}
	if t.Owner == "" {
		return fmt.Errorf("负责人不可为空%w", e.COMMON_PARAM_MISS)
	}

	now := time.Now()
	if t.ID == "" {
		t.ID = fmt.Sprintf("%d", now.UnixNano()/1000)
	}
	if t.StartTime == 0 {
		t.StartTime = now.Unix()
	}

	return nil
}

func (t *Todo) LevelName() string {
	return LevelName(t.Level)
}

func (t *Todo) TypeName() string {
	return TypeName(t.Type)
}
