package todo

import (
	"context"
	"fmt"
	"github.com/bingbaba/storage"
	"github.com/pingfen/notify/pkg/e"
	"github.com/pingfen/notify/pkg/models"
	"github.com/pingfen/notify/pkg/utils"
	pkgstore "github.com/pingfen/notify/storage"
)

func getTodoKey(owner, id string) string {
	return fmt.Sprintf("/todo/current/%s/%s", owner, id)
}

func Create(ctx context.Context, t *models.Todo) error {
	if err := t.Valid(); err != nil {
		return err
	}

	// storage instance
	store, err := pkgstore.Get()
	if err != nil {
		return err
	}

	key := getTodoKey(t.Owner, t.ID)

	// check if exist
	t_old := new(models.Todo)
	err = store.Get(ctx, key, t_old)
	if err != nil {
		if !storage.IsNotFound(err) {
			return fmt.Errorf("提交内容已存在%w", e.COMMON_CONFILCT)
		}
	} else {
		t.CreateTime = t_old.CreateTime
		t.StartTime = t_old.StartTime
	}

	err = store.Create(ctx, key, t, 0)
	if err != nil {
		return fmt.Errorf("save failed %w", err)
	}

	// user -> todo
	if t.Operator != "" {
		if err = AddUserIndex(ctx, t.Owner, t.ID, t.Operator); err != nil {
			return err
		}
	}

	return nil
}

func Update(ctx context.Context, t *models.Todo) error {
	return Create(ctx, t)
}

func List(ctx context.Context, owner string) ([]*models.Todo, error) {
	store, err := pkgstore.Get()
	if err != nil {
		return nil, err
	}

	key := getTodoKey(owner, "")
	ret, err := store.List(ctx, key, nil, &models.Todo{})
	if err != nil {
		return nil, fmt.Errorf("list failed %w", err)
	}

	todo_list := make([]*models.Todo, len(ret))
	for i, item := range ret {
		todo_list[i] = item.(*models.Todo)
	}

	return todo_list, nil
}

func update(ctx context.Context, t *models.Todo) error {
	store, err := pkgstore.Get()
	if err != nil {
		return err
	}

	key := getTodoKey(t.Owner, t.ID)
	err = store.Update(ctx, key, 0, t, 0)
	if err != nil {
		return fmt.Errorf("update to storage failed %w", err)
	}
	return nil
}

func Get(ctx context.Context, owner, id string) (*models.Todo, error) {
	// storage instance
	store, err := pkgstore.Get()
	if err != nil {
		return nil, err
	}

	t := new(models.Todo)
	key := getTodoKey(owner, id)
	err = store.Get(ctx, key, t)
	if err != nil {
		return nil, fmt.Errorf("get from storage failed %w", err)
	}

	return t, nil
}

func ModifyOperator(ctx context.Context, owner, id, operator string) error {
	if operator == "" {
		return fmt.Errorf("operator不可为空%w", e.COMMON_PARAM_MISS)
	}

	t, err := Get(ctx, owner, id)
	if err != nil {
		return fmt.Errorf("read from storage failed %w", err)
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

func Finish(ctx context.Context, t *models.Todo) error {

	// move to history
	return nil
}
