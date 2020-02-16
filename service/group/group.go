package group

import (
	"context"
	"fmt"
	"github.com/bingbaba/storage"
	"github.com/pingfen/notify/pkg/models"
	pkgstore "github.com/pingfen/notify/storage"
)

func getGroupKey(gid string) string {
	return fmt.Sprintf("/groups/%s/meta", gid)
}

func Create(ctx context.Context, g *models.Group, m *Member) error {
	found, err := IsExist(ctx, g.Id)
	if err != nil {
		return err
	}

	if found {
		return fmt.Errorf("该群组未发现%w", errs.OBJECT_EXIST)
	}
	g.Managers = []string{m.UnionId}

	// storage instance
	store, _ := pkgstore.Get()
	err = store.Create(ctx, getGroupKey(g.Id), g, 0)
	if err != nil {
		return fmt.Errorf("创建群组失败%w", err)
	}

	// add manager to member
	err = AddMember(ctx, g.Id, m)
	if err != nil {
		return fmt.Errorf("保存成员信息失败%w", err)
	}

	return nil
}

func Update(ctx context.Context, g *models.Group) error {
	store, err := pkgstore.Get()
	if err != nil {
		return err
	}

	if g.Managers == nil || len(g.Managers) == 0 {
		g_old := new(models.Group)
		err = store.Get(ctx, getGroupKey(g.Id), g_old)
		if err != nil {
			return fmt.Errorf("该群组未发现", errs.OBJECT_EXIST)
		}
		g.Managers = g_old.Managers
	}

	err = store.Update(ctx, getGroupKey(g.Id), 0, g, 0)
	if err != nil {
		return fmt.Errorf("更新群组失败%w", err)
	}

	return nil
}

func IsExist(ctx context.Context, id string) (bool, error) {
	store, err := pkgstore.Get()
	if err != nil {
		return true, err
	}

	g_old := new(models.Group)
	err = store.Get(ctx, getGroupKey(id), g_old)
	if err != nil {
		if !storage.IsNotFound(err) {
			return true, fmt.Errorf("检查群组失败%w", err)
		}
	} else {
		return true, nil
	}

	return false, nil
}
