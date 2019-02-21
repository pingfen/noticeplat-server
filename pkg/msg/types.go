package msg

import (
	"errors"
	"fmt"
	"time"
)

type MsgType = int
type Level = int
type ClientType = string

const (
	MSGTYPE_ALARM_START   MsgType = 0
	MSGTYPE_ALARM_RECOVER MsgType = 1
	MSGTYPE_EVENT         MsgType = 2
	MSGTYPE_TASK          MsgType = 3
	MSGTYPE_TASK_COMPLATE MsgType = 4
	MSGTYPE_MAX           MsgType = MSGTYPE_TASK_COMPLATE

	LEVEL_DEBUG Level = 0
	LEVEL_INFO  Level = 1
	LEVEL_WARN  Level = 2
	LEVEL_ERROR Level = 3

	CLIENTTYPE_USER  ClientType = "user"
	CLIENTTYPE_GROUP ClientType = "group"

	KEYPREFIX_TODO_DATA  = "/todo/data"   // key /todo/data/$owner/$mid
	KEYPREFIX_TODO_INDEX = "/todo/.index" // key /todo/.index/$uid/$owner.$mid

	KEYPREFIX_DONE_DATA  = "/done/data"   // key /done/data/$recoveryTime.$owner.$mid
	KEYPREFIX_DONE_INDEX = "/done/.index" // key /done/.index/$uid/$recoveryTime.$owner.$mid
)

func LevelName(level Level) string {
	switch level {
	case LEVEL_DEBUG:
		return "调试"
	case LEVEL_INFO:
		return "一般"
	case LEVEL_WARN:
		return "警告"
	case LEVEL_ERROR:
		return "紧急"
	default:
		return "未知"
	}
}

func MsgTypeName(t MsgType) string {
	switch t {
	case MSGTYPE_ALARM_START:
		return "告警"
	case MSGTYPE_ALARM_RECOVER:
		return "已恢复"
	case MSGTYPE_EVENT:
		return "事件"
	case MSGTYPE_TASK:
		return "任务"
	case MSGTYPE_TASK_COMPLATE:
		return "已完成"
	default:
		return "未知"
	}
}

type Message struct {
	ID      string  `json:"id"`
	MsgType MsgType `json:"type"`
	Level   Level   `json:"level"`

	Project   string `json:"project"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp,omitempty"`

	Target string `json:"target"`
}

func (m *Message) Check() error {
	if m.ID == "" {
		return errors.New("id was required")
	} else if m.MsgType < 0 || m.MsgType > MSGTYPE_MAX {
		return errors.New("type was required")
	} else if m.Target == "" {
		return errors.New("target was required")
	} else if m.Title == "" {
		return errors.New("title was required")
	} else if m.Content == "" {
		return errors.New("content was required")
	}

	if m.Timestamp == 0 {
		m.Timestamp = time.Now().Unix()
	}
	return nil
}

type Todo struct {
	*Message

	Owner            string     `json:"owner"`
	ClientType       ClientType `json:"clientType"`
	Operators        []string   `json:"operators"`
	HistoryOperators []string   `json:"historyOperators"`

	GenerateTime  int64 `json:"generateTime"`
	LastOccurTime int64 `json:"lastOccurTime,omitempty"`
	RecoveryTime  int64 `json:"recoveryTime,omitempty"`

	Alerting bool `json:"-"`
}

func (todo *Todo) OperatorIsChange(old []string) bool {
	new := todo.Operators
	if len(new) != len(old) {
		return true
	}
	for i, item := range new {
		if item != old[i] {
			return true
		}
	}
	return false
}

type History struct {
	*Message
	RecoverTime int64 `json:"recoverTime"`
}

func TodoKeyPrefix(uid_gid string, ct ClientType) string {
	return fmt.Sprintf("%s/%s.%s", KEYPREFIX_TODO_DATA, ct, uid_gid)
}

func (todo *Todo) Key() string {
	return fmt.Sprintf("%s/%s.%s/%s", KEYPREFIX_TODO_DATA, todo.ClientType, todo.Owner, todo.Message.ID)
}

func (todo *Todo) RecoveryKey() string {
	return fmt.Sprintf("%s/%s.%s/%d.%s", KEYPREFIX_DONE_DATA, todo.ClientType, todo.Owner, todo.RecoveryTime, todo.Message.ID)
}

func DoneKeyPrefix(uid_gid string, ct ClientType) string {
	return fmt.Sprintf("%s/%s.%s", KEYPREFIX_DONE_DATA, ct, uid_gid)
}
