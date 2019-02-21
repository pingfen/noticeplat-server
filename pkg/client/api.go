package client

import (
	"context"
	"time"

	"github.com/bingbaba/storage"
)

func Get(ctx context.Context, store storage.Interface, id string) (c *Client, err error) {
	c = new(Client)
	err = store.Get(ctx, CLIENT_KEYPREFIX+"/"+id, c)
	return
}

func Add(ctx context.Context, store storage.Interface, c *Client) error {
	c.CreateTime = time.Now()
	c.UpdateTime = c.CreateTime

	// user id -> client id
	err := store.Create(ctx, BYUSER_KEYPREIX+"/"+c.User, map[string]string{"id": c.Id}, 0)
	if err != nil {
		return err
	}

	// id -> client info
	err := store.Create(ctx, CLIENT_KEYPREFIX+"/"+c.Id, c, 0)
	if err != nil {
		return err
	}

	return nil
}

func Modify(ctx context.Context, store storage.Interface, c *Client) error {
	old_c, err := Get(ctx, store, c.Id)
	if err != nil {
		return err
	}

	if old_c.User != "" && old_c.User != c.User {
		store.Delete(ctx, BYUSER_KEYPREIX+"/"+c.User, nil)
	}

	c.CreateTime = old_c.CreateTime
	c.UpdateTime = time.Now()
	return Add(ctx, store, c)
}

func Delete(ctx context.Context, store storage.Interface, id string) error {
	err := store.Delete(ctx, BYUSER_KEYPREIX+"/"+id, nil)
	if err != nil {
		return err
	}

	return store.Delete(ctx, CLIENT_KEYPREFIX+"/"+id, make(map[string]interface{}))
}
