package group

import (
	"context"
	"fmt"
	"github.com/bingbaba/storage"
	pkgstore "github.com/pingfen/noticeplat-server/pkg/storage"
	"github.com/pingfen/noticeplat-server/pkg/user"
	"github.com/pkg/errors"
	"strings"
	"sync"
)

var (
	asyncLimit = make(chan bool, 50)
)

type Member struct {
	UnionId   string `json:"unionid"`
	NickName  string `json:"nickname"`
	ShowPhone bool   `json:"showPhone"`

	Detail *user.User `json:"detail,omitempty"`
}

type DetailMember struct {
	UnionId   string `json:"unionid"`
	NickName  string `json:"nickname"`
	ShowPhone bool   `json:"showPhone"`
}

func getGroupMemberKey(gid, uid string) string {
	return fmt.Sprintf("/groups/%s/members/%s", gid, uid)
}

func AddMember(ctx context.Context, gid string, m *Member) error {
	store, err := pkgstore.Get()
	if err != nil {
		return err
	}

	err = store.Create(ctx, getGroupMemberKey(gid, m.UnionId), m, 0)
	if err != nil {
		return errors.Wrap(err, "add member failed")
	}

	return nil
}

func RemoveMember(ctx context.Context, gid, mid string) error {
	store, err := pkgstore.Get()
	if err != nil {
		return err
	}

	err = store.Delete(ctx, getGroupMemberKey(gid, mid), nil)
	if err != nil {
		return errors.Wrap(err, "remove member failed")
	}

	return nil
}

func ListMembers(ctx context.Context, gid string) ([]*Member, error) {
	store, err := pkgstore.Get()
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("/groups/%s/members", gid)
	ms, err := store.List(ctx, key, nil, new(Member))
	if err != nil {
		return nil, errors.Wrap(err, "list members failed")
	}

	ret := make([]*Member, len(ms))
	for i, m := range ms {
		ret[i] = m.(*Member)
	}

	return ret, nil
}

func ListMemberIds(ctx context.Context, gid string) ([]string, error) {
	store, err := pkgstore.Get()
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("/groups/%s/members", gid)
	item_keys, err := store.List(ctx, key, &storage.SelectionPredicate{KeyOnly: true}, nil)
	if err != nil {
		return nil, errors.Wrap(err, "list members failed")
	}

	ret := make([]string, len(item_keys))
	for i, k := range item_keys {
		ret[i] = strings.TrimPrefix(k.(string), key+"/")
	}

	return ret, nil
}

func ListDetailMembers(ctx context.Context, gid string) ([]*Member, error) {
	ms, err := ListMembers(ctx, gid)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, m := range ms {
		select {
		case <-ctx.Done():
			return ms, errors.Wrap(context.Canceled, "read detail timeout")
		case asyncLimit <- true:
			wg.Add(1)
		}

		go func(m *Member) {
			defer func() {
				wg.Done()
				<-asyncLimit
			}()

			u, err_tmp := user.Get(ctx, m.UnionId)
			if err_tmp != nil {
				err = err_tmp
				return
			}

			m.Detail = u
		}(m)
	}

	wg.Wait()
	return ms, err
}
