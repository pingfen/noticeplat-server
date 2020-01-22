package todo

import (
	"context"
	"fmt"
	"time"

	pkgstore "github.com/pingfen/noticeplat-server/pkg/storage"
)

func AddUserIndex(ctx context.Context, owner, id, user string) error {
	store, err := pkgstore.Get()
	if err != nil {
		return err
	}

	key := getTodoUserIdxKey(owner, id, user)
	if err := store.Create(ctx, key, time.Now().Unix(), 0); err != nil {
		return fmt.Errorf("save index to %s failed: %w", key, err)
	}

	return nil
}

func getTodoUserIdxKey(owner, id, user string) string {
	return fmt.Sprintf("/todo/current/.index/%s/%s/%s", user, owner, id)
}
