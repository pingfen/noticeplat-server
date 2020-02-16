package group

import (
	"context"
	"fmt"
	"github.com/bingbaba/storage"
	"github.com/pingfen/notify/pkg/e"
	"github.com/pingfen/notify/pkg/models"
	"github.com/pingfen/notify/service/user"
	pkgstore "github.com/pingfen/notify/storage"
	"strings"
	"sync"
)

var (
	asyncLimit = make(chan bool, 50)
)

type Member struct {
	OpenId    string `json:"openid"`
	NickName  string `json:"nickname"`
	ShowPhone bool   `json:"showPhone"`

	Detail *models.MpUser `json:"detail,omitempty"`
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

	err = store.Create(ctx, getGroupMemberKey(gid, m.OpenId), m, 0)
	if err != nil {
		return fmt.Errorf("添加成员失败%w", err)
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
		return fmt.Errorf("删除成员失败%w", err)
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
		return nil, fmt.Errorf("查询成员列表失败%w", err)
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
		return nil, fmt.Errorf("查询成员列表失败%w", err)
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
			return ms, fmt.Errorf("执行超时%w", e.COMMON_INTERNAL_CALLING_TIMEOUT)
		case asyncLimit <- true:
			wg.Add(1)
		}

		go func(m *Member) {
			defer func() {
				wg.Done()
				<-asyncLimit
			}()

			u, err_tmp := user.GetByOpenid(ctx, "notodo", m.OpenId)
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
