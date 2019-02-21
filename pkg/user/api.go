package user

import (
	"context"
	"time"

	"fmt"
	"github.com/bingbaba/storage"
)

func Add(ctx context.Context, store storage.Interface, u *User) (err error) {
	u.CreateTime = time.Now()
	u.UpdateTime = u.CreateTime

	defer func() {
		if err != nil {
			store.Delete(ctx, APPOPENID_KEYPREIX+"/"+u.AppOpenId, nil)
			store.Delete(ctx, APPOPENID_KEYPREIX+"/"+u.SrvOpenId, nil)
		}
	}()

	// appOpenid -> id
	if u.AppOpenId != "" {
		err = store.Create(ctx, APPOPENID_KEYPREIX+"/"+u.AppOpenId, map[string]string{"id": u.Id}, 0)
		if err != nil {
			return err
		}
	}

	// srvOpenId -> id
	if u.SrvOpenId != "" {
		err = store.Create(ctx, SRVOPENID_KEYPREIX+"/"+u.SrvOpenId, map[string]string{"id": u.Id}, 0)
		if err != nil {
			return err
		}
	}

	// id -> user info
	err = store.Create(ctx, USER_KEYPREFIX+"/"+u.Id, u, 0)
	if err != nil {
		return err
	}

	return nil
}

func Modify(ctx context.Context, store storage.Interface, u *User) error {
	old_u, err := Get(ctx, store, u.Id)
	if err != nil {
		return err
	}

	if old_u.SrvOpenId != "" && old_u.SrvOpenId != u.SrvOpenId {
		store.Delete(ctx, SRVOPENID_KEYPREIX+"/"+old_u.SrvOpenId, nil)
	}

	if old_u.AppOpenId != "" && old_u.AppOpenId != u.AppOpenId {
		store.Delete(ctx, APPOPENID_KEYPREIX+"/"+old_u.AppOpenId, nil)
	}

	u.UpdateTime = time.Now()
	return Add(ctx, store, u)
}

//
//func BindingSrvOpenId(ctx context.Context, store storage.Interface, uid, srvOpenId string) error {
//	old_u, err := Get(ctx, store, uid)
//	if err != nil {
//		return err
//	}
//
//	if old_u.SrvOpenId != "" {
//		store.Delete(ctx, SRVOPENID_KEYPREIX+"/"+old_u.SrvOpenId, nil)
//	}
//
//	if old_u.AppOpenId != "" {
//		store.Delete(ctx, APPOPENID_KEYPREIX+"/"+old_u.AppOpenId, nil)
//	}
//
//	old_u.BindSrvOpenId(srvOpenId)
//	old_u.UpdateTime = time.Now()
//	return Add(ctx, store, old_u)
//}

func Get(ctx context.Context, store storage.Interface, id string) (u *User, err error) {
	u = new(User)
	err = store.Get(ctx, USER_KEYPREFIX+"/"+id, u)
	return
}

func GetByOpenId(ctx context.Context, store storage.Interface, openid string) (u *User, err error) {
	userid_map := make(map[string]string)
	err = store.Get(ctx, APPOPENID_KEYPREIX+"/"+openid, &userid_map)
	if err != nil {
		return nil, err
	}

	uid, found := userid_map["id"]
	if !found {
		return nil, fmt.Errorf("can't found user id for the openid")
	}

	return Get(ctx, store, uid)
}

func IdIsExist(ctx context.Context, store storage.Interface, id string) (bool, error) {
	_, err := Get(ctx, store, id)
	if err != nil {
		if storage.IsNotFound(err) {
			return false, nil
		} else {
			return true, err
		}
	}

	return true, nil
}

func BindingSrvOpenId(ctx context.Context, store storage.Interface, id, srvOpenId string) error {
	u, err := Get(ctx, store, id)
	if err != nil {
		return err
	}

	if u.SrvOpenId != srvOpenId {
		if u.SrvOpenId != "" {
			store.Delete(ctx, SRVOPENID_KEYPREIX+"/"+u.SrvOpenId, nil)
		}

		u.BindSrvOpenId(srvOpenId)
		u.UpdateTime = time.Now()
		return Add(ctx, store, u)
	}

	return nil
}

func DeleteUser(ctx context.Context, store storage.Interface, id string) error {
	// get user info
	u, err := Get(ctx, store, id)
	if err != nil {
		return err
	}

	// appOpenid -> id
	if u.AppOpenId != "" {
		store.Delete(ctx, APPOPENID_KEYPREIX+"/"+u.AppOpenId, nil)
	}

	// srvOpenId -> id
	if u.SrvOpenId != "" {
		store.Delete(ctx, SRVOPENID_KEYPREIX+"/"+u.SrvOpenId, nil)
	}

	// id -> user info
	err = store.Delete(ctx, USER_KEYPREFIX+"/"+u.Id, nil)
	if err != nil {
		return err
	}

	return nil
}
