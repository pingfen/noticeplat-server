package group

import (
	"context"
	"time"

	"github.com/bingbaba/storage"
)

func Get(ctx context.Context, store storage.Interface, id string) (g *Group, err error) {
	g = new(Group)
	err = store.Get(ctx, GROUP_KEYPREFIX+"/"+id, g)
	if err != nil {
		return
	}

	members, err := store.List(ctx, GROUPMEMBERS_KEYPREFIX, nil, new(Member))
	if err != nil {
		return g, err
	}

	g.Members = make([]Member, len(members))
	for i, m := range members {
		g.Members[i] = m.(Member)
	}

	return
}

func Create(ctx context.Context, store storage.Interface, g *Group) error {
	err := store.Create(ctx, GROUP_KEYPREFIX+"/"+g.OpenGId, g.GroupBase, 0)
	if err != nil {
		return err
	}

	for _, m := range g.Members {
		err = AddMember(ctx, store, g.OpenGId, m.OpenId)
		if err != nil {
			if storage.IsNodeExist(err) || storage.IsConflict(err) {
				continue
			} else {
				return err
			}
		}
	}

	return nil
}

func AddMember(ctx context.Context, store storage.Interface, groupid, openid string) error {
	_, err := Get(ctx, store, groupid)
	if err != nil {
		if storage.IsNotFound(err) {
			err = Create(ctx, store, &Group{GroupBase: &GroupBase{OpenGId: groupid}})
			if err != nil {
				return err
			}
		}
		return err
	}

	return store.Create(ctx, GROUPMEMBERS_KEYPREFIX+"/"+openid, Member{OpenId: openid, TimeStamp: time.Now().Unix()}, 0)
}

func RemoveMember(ctx context.Context, store storage.Interface, opengid, openid string) error {
	_, err := Get(ctx, store, opengid)
	if err != nil {
		return err
	}

	return store.Delete(ctx, GROUPMEMBERS_KEYPREFIX+"/"+openid, &(map[string]string{}))
}
