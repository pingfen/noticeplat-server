package msg

import (
	"context"
	"fmt"

	"github.com/bingbaba/storage"
)

func Post(ctx context.Context, store storage.Interface, uid_gid string, client_type ClientType, msg *Message) (*Todo, error) {
	switch msg.MsgType {
	case MSGTYPE_ALARM_START, MSGTYPE_TASK:
		todo := &Todo{Message: msg, Owner: uid_gid, ClientType: client_type}
		err := store.Get(ctx, todo.Key(), todo)
		if err != nil {
			todo.GenerateTime = msg.Timestamp
			todo.Operators = []string{msg.Target}
			todo.Alerting = true
		}
		todo.LastOccurTime = msg.Timestamp

		err = saveTodo(ctx, store, todo)
		if err != nil {
			return todo, err
		} else {
			return todo, nil
		}
	case MSGTYPE_ALARM_RECOVER, MSGTYPE_TASK_COMPLATE:
		todo := &Todo{Message: msg, Owner: uid_gid, ClientType: client_type}
		err := store.Get(ctx, todo.Key(), todo)
		if err != nil {
			todo.Message = msg
			todo.GenerateTime = msg.Timestamp
			todo.Operators = []string{msg.Target}
		}
		todo.MsgType = MSGTYPE_ALARM_RECOVER
		todo.RecoveryTime = msg.Timestamp
		todo.Alerting = true

		err = moveToHistory(ctx, store, todo, false)
		if err != nil {
			return todo, err
		} else {
			return todo, nil
		}
	case MSGTYPE_EVENT:
		todo := &Todo{
			Message:      msg,
			Owner:        msg.Target,
			Operators:    []string{msg.Target},
			GenerateTime: msg.Timestamp,
		}
		todo.Alerting = true

		err := moveToHistory(ctx, store, todo, true)
		if err != nil {
			return todo, err
		} else {
			return todo, nil
		}
	default:
		return nil, fmt.Errorf("unknown message type %d", msg.MsgType)
	}
}

func Modify(ctx context.Context, store storage.Interface, todo *Todo) error {
	key := todo.Key()

	old_todo := new(Todo)
	err := store.Get(ctx, key, old_todo)
	if err != nil {
		return err
	}

	if todo.OperatorIsChange(old_todo.Operators) {
		todo.Alerting = true
	}

	new_ho := make([]string, 0, len(todo.HistoryOperators)+len(todo.Operators))
	new_ho = append(new_ho, old_todo.HistoryOperators...)
	new_ho = append(new_ho, old_todo.Operators...)
	todo.HistoryOperators = new_ho

	return store.Update(ctx, key, 0, todo, 0)
}

func TodoList(ctx context.Context, store storage.Interface, uid string, ct ClientType) ([]*Todo, error) {
	ret, err := store.List(ctx, TodoKeyPrefix(uid, ct), nil, new(Todo))
	if err != nil {
		return nil, err
	}

	todo_list := make([]*Todo, len(ret))
	for i, item := range ret {
		todo_list[i] = item.(*Todo)
	}

	return todo_list, nil
}

func GetTodo(ctx context.Context, store storage.Interface, uid, mid string, ct ClientType) (*Todo, error) {
	//key := fmt.Sprintf("%s/%s/%s", KEYPREFIX_TODO_DATA, uid, mid)
	todo := &Todo{Message: &Message{ID: mid}, Owner: uid, ClientType: ct}
	err := store.Get(ctx, todo.Key(), todo)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func HistoryList(ctx context.Context, store storage.Interface, uid string, ct ClientType) ([]*Todo, error) {
	ret, err := store.List(ctx, DoneKeyPrefix(uid, ct), nil, new(Todo))
	if err != nil {
		return nil, err
	}

	todo_list := make([]*Todo, len(ret))
	for i, item := range ret {
		todo_list[i] = item.(*Todo)
	}

	return todo_list, nil
}

func saveTodo(ctx context.Context, store storage.Interface, todo *Todo) error {
	// save data
	key := todo.Key()
	err := store.Create(ctx, key, todo, 0)
	if err != nil {
		return err
	}

	return nil
}

func moveToHistory(ctx context.Context, store storage.Interface, todo *Todo, ignoreDelErr bool) error {
	// save done data
	err := store.Create(ctx, todo.RecoveryKey(), todo, 0)
	if err != nil {
		return err
	}

	// delete todo data
	err = store.Delete(ctx, todo.Key(), nil)
	if !ignoreDelErr && err != nil {
		return err
	}

	return nil
}
