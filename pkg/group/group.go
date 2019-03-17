package group

import (
	"context"
	"fmt"
	"github.com/bingbaba/storage"
	"github.com/pingfen/noticeplat-server/pkg/errs"
	pkgstore "github.com/pingfen/noticeplat-server/pkg/storage"
	"github.com/pkg/errors"
)

type Group struct {
	Id         string   `json:"id"`
	OpenGid    string   `json:"openGid"`
	Name       string   `json:"name"`
	HeadimgUrl string   `json:"headimgurl"`
	Managers   []string `json:"managers"`
	Secret     string   `json:"secret,omitempty"`
}

func getGroupKey(gid string) string {
	return fmt.Sprintf("/groups/%s/meta", gid)
}

func Create(ctx context.Context, g *Group, m *Member) error {
	found, err := IsExist(ctx, g.Id)
	if err != nil {
		return err
	}

	if found {
		return errors.Wrap(errs.OBJECT_EXIST, "group has exist")
	}
	g.Managers = []string{m.UnionId}

	// storage instance
	store, _ := pkgstore.Get()
	err = store.Create(ctx, getGroupKey(g.Id), g, 0)
	if err != nil {
		return errors.Wrap(err, "save failed")
	}

	// add manager to member
	err = AddMember(ctx, g.Id, m)
	if err != nil {
		return errors.Wrap(err, "save manager to member failed")
	}

	return nil
}

func Update(ctx context.Context, g *Group) error {
	store, err := pkgstore.Get()
	if err != nil {
		return err
	}

	if g.Managers == nil || len(g.Managers) == 0 {
		g_old := new(Group)
		err = store.Get(ctx, getGroupKey(g.Id), g_old)
		if err != nil {
			return errors.Wrap(err, "group not found")
		}
		g.Managers = g_old.Managers
	}

	err = store.Update(ctx, getGroupKey(g.Id), 0, g, 0)
	if err != nil {
		return errors.Wrap(err, "update group failed")
	}

	return nil
}

func IsExist(ctx context.Context, id string) (bool, error) {
	store, err := pkgstore.Get()
	if err != nil {
		return true, err
	}

	g_old := new(Group)
	err = store.Get(ctx, getGroupKey(id), g_old)
	if err != nil {
		if !storage.IsNotFound(err) {
			return true, errors.Wrap(err, "check identify failed")
		}
	} else {
		return true, nil
	}

	return false, nil
}

func (g *Group) IsManager(unionid string) bool {
	for _, uid := range g.Managers {
		if uid == unionid {
			return true
		}
	}

	return false
}
