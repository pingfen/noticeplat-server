package todo

import (
	"context"
	"fmt"
	"github.com/bingbaba/storage"
	"github.com/pingfen/noticeplat-server/pkg/errs"
	pkgstore "github.com/pingfen/noticeplat-server/pkg/storage"
	"github.com/pingfen/noticeplat-server/pkg/utils"
	"github.com/pkg/errors"
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

func getTodoKey(owner, id string) string {
	return fmt.Sprintf("/todo/current/%s/%s", owner, id)
}

func Create(ctx context.Context, t *Todo) error {
	if err := t.valid(); err != nil {
		return err
	}

	// storage instance
	store, err := pkgstore.Get()
	if err != nil {
		return err
	}

	key := getTodoKey(t.Owner, t.ID)

	// check if exist
	t_old := new(Todo)
	err = store.Get(ctx, key, t_old)
	if err != nil {
		if !storage.IsNotFound(err) {
			return errors.Wrap(err, "check repeat failed")
		}
	} else {
		t.CreateTime = t_old.CreateTime
		t.StartTime = t_old.StartTime
	}

	err = store.Create(ctx, key, t, 0)
	if err != nil {
		return errors.Wrap(err, "save failed")
	}

	// user -> todo
	if t.Operator != "" {
		if err = AddUserIndex(ctx, t.Owner, t.ID, t.Operator); err != nil {
			return err
		}
	}

	return nil
}

func Update(ctx context.Context, t *Todo) error {
	return Create(ctx, t)
}

func List(ctx context.Context, owner string) ([]*Todo, error) {
	store, err := pkgstore.Get()
	if err != nil {
		return nil, err
	}

	key := getTodoKey(owner, "")
	ret, err := store.List(ctx, key, nil, &Todo{})
	if err != nil {
		return nil, errors.Wrap(err, "list failed")
	}

	todo_list := make([]*Todo, len(ret))
	for i, item := range ret {
		todo_list[i] = item.(*Todo)
	}

	return todo_list, nil
}

func update(ctx context.Context, t *Todo) error {
	store, err := pkgstore.Get()
	if err != nil {
		return err
	}

	key := getTodoKey(t.Owner, t.ID)
	err = store.Update(ctx, key, 0, t, 0)
	if err != nil {
		return errors.Wrap(err, "update to storage failed")
	}
	return nil
}

func Get(ctx context.Context, owner, id string) (*Todo, error) {
	// storage instance
	store, err := pkgstore.Get()
	if err != nil {
		return nil, err
	}

	t := new(Todo)
	key := getTodoKey(owner, id)
	err = store.Get(ctx, key, t)
	if err != nil {
		return nil, errors.Wrap(err, "get from storage failed")
	}

	return t, nil
}

func ModifyOperator(ctx context.Context, owner, id, operator string) error {
	if operator == "" {
		return errors.New("operator is empty")
	}

	t, err := Get(ctx, owner, id)
	if err != nil {
		return errors.Wrap(err, "read from storage failed")
	}

	if t.Operator == operator {
		return nil
	}

	newOperator := false
	if !utils.InSlice(t.Operator, t.HistoryOperator) {
		newOperator = true

		t.HistoryOperator = append(t.HistoryOperator, t.Operator)

		// the todo comment

	}
	t.Operator = operator

	if err = update(ctx, t); err != nil {
		return err
	}

	// new operator
	if newOperator {
		if err = AddUserIndex(ctx, t.Owner, t.ID, t.Operator); err != nil {
			return err
		}
	}

	return nil
}

func Finish(ctx context.Context, t *Todo) error {

	// move to history
	return nil
}

func (t *Todo) valid() error {
	if t.Subject == "" {
		return errors.Wrap(errs.FIELD_EMPTY, "subject is empty")
	}
	if t.Content == "" {
		return errors.Wrap(errs.FIELD_EMPTY, "content is empty")
	}
	if t.Owner == "" {
		return errors.Wrap(errs.FIELD_EMPTY, "owner is empty")
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
