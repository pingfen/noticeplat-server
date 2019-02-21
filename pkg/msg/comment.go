package msg

import (
	"context"
	"fmt"
	"time"

	"github.com/bingbaba/storage"
)

type CommentsType int

const (
	CommentsType_User   CommentsType = 0
	CommentsType_System CommentsType = 1

	KEYPREFIX_COMMENT = "/comments/data"
)

type Comments struct {
	Id         string       `json:"id"`
	MsgId      string       `json:"msgid"`
	MsgOwner   string       `json:"msgOwner"`
	Content    string       `json:"content"`
	Type       CommentsType `json:"type"`
	TimeStamp  int64        `json:"timestamp"`
	Operator   string       `json:"operator"`
	ClientType ClientType   `json:"clientType"`

	Good int `json:"good"`
	Bad  int `json:"bad"`
}

func (c *Comments) Key() string {
	return getCommentsKey(c.MsgOwner, c.MsgId, c.Id, c.ClientType)
}

func getCommentsKey(uid, mid, cid string, ct ClientType) string {
	return fmt.Sprintf("%s/%s.%s.%s/%s", KEYPREFIX_COMMENT, ct, uid, mid, cid)
}

func PostComments(ctx context.Context, store storage.Interface, c *Comments) error {
	// use timestamp default
	if c.Id == "" {
		c.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	if c.TimeStamp == 0 {
		c.TimeStamp = time.Now().Unix()
	}

	err := store.Create(ctx, c.Key(), c, 0)
	if err != nil {
		return err
	}
	return nil
}

func ModifyComments(ctx context.Context, store storage.Interface, c *Comments) error {
	if c.Id == "" {
		return fmt.Errorf("comment id is required")
	}

	err := store.Update(ctx, c.Key(), 0, c, 0)
	if err != nil {
		return err
	}
	return nil
}

func DeleteComments(ctx context.Context, store storage.Interface, uid, mid, cid string, ct ClientType) error {
	err := store.Delete(ctx, getCommentsKey(uid, mid, cid, ct), nil)
	if err != nil {
		return err
	}
	return nil
}

func ListCommentses(ctx context.Context, store storage.Interface, uid, mid string, ct ClientType) ([]*Comments, error) {
	ret, err := store.List(ctx, getCommentsKey(uid, mid, "", ct), nil, new(Comments))
	if err != nil {
		return nil, err
	}

	cs := make([]*Comments, len(ret))
	for i, item := range ret {
		cs[i] = item.(*Comments)
	}

	return cs, nil
}

func CommentsGood(ctx context.Context, store storage.Interface, uid, mid, cid string, ct ClientType) (int, error) {
	key := getCommentsKey(uid, mid, cid, ct)
	c := new(Comments)
	err := store.Get(ctx, key, c)
	if err != nil {
		return 0, err
	}
	c.Good++
	err = ModifyComments(ctx, store, c)
	return c.Good, err
}

func CommentsBad(ctx context.Context, store storage.Interface, uid, mid, cid string, ct ClientType) (int, error) {
	key := getCommentsKey(uid, mid, cid, ct)
	c := new(Comments)
	err := store.Get(ctx, key, c)
	if err != nil {
		return 0, err
	}
	c.Bad++
	err = ModifyComments(ctx, store, c)
	return c.Bad, err
}
