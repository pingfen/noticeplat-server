package user

import (
	"context"
	"github.com/bingbaba/storage"
	pkgstore "github.com/pingfen/noticeplat-server/pkg/storage"
	"github.com/pkg/errors"
)

func GetUnionIdByOpenid(ctx context.Context, product, openid string) (string, error) {
	// storage instance
	store, err := pkgstore.Get()
	if err != nil {
		return "", err
	}

	var unionid string
	key := getOpenidKey(product, openid)
	err = store.Get(ctx, key, &unionid)
	if err != nil {
		return "", errors.Wrap(err, "read key "+key+" failed")
	}

	return unionid, nil
}

func GetByOpenid(ctx context.Context, product, openid string) (*User, error) {
	unionid, err := GetUnionIdByOpenid(ctx, product, openid)
	if err != nil {
		return nil, err
	}

	return Get(ctx, unionid)
}

func deleteIndex(ctx context.Context, u *User) error {
	// storage instance
	store, _ := pkgstore.Get()

	for product, info := range u.WechatProducts {
		err := store.Delete(ctx, getOpenidKey(product, info.Openid), nil)
		if err != nil && !storage.IsNotFound(err) {
			return err
		}
	}

	return nil
}

func saveIndex(ctx context.Context, u_old, u *User) error {
	// storage instance
	store, _ := pkgstore.Get()

	// 0:nochange 1:add 2:modify 3:delete
	operation := make(map[string]int)

	if u_old == nil || u_old.WechatProducts == nil {
		for product := range u.WechatProducts {
			operation[product] = 1
		}
	} else {
		for product, info := range u.WechatProducts {
			if info_old, found := u_old.WechatProducts[product]; found {
				if info.Openid != info_old.Openid {
					operation[product] = 2
				} else {
					operation[product] = 0
				}
			} else {
				operation[product] = 1
			}
		}
		for product := range u_old.WechatProducts {
			if _, found := operation[product]; !found {
				operation[product] = 3
			}
		}
	}

	//index
	//log.Printf("%+v", operation)
	for product, op := range operation {
		switch op {
		case 1:
			info := u.WechatProducts[product]
			err := store.Create(ctx, getOpenidKey(product, info.Openid), u.UnionId, 0)
			if err != nil {
				return errors.Wrap(err, "save new index failed")
			}
		case 2:
			info := u.WechatProducts[product]
			err := store.Create(ctx, getOpenidKey(product, info.Openid), u.UnionId, 0)
			if err != nil {
				return errors.Wrap(err, "save new index failed")
			}

			info_old := u_old.WechatProducts[product]
			err = store.Delete(ctx, getOpenidKey(product, info_old.Openid), nil)
			if err != nil && !storage.IsNotFound(err) {
				return errors.Wrap(err, "delete old index failed")
			}
		case 3:
			info := u_old.WechatProducts[product]
			err := store.Delete(ctx, getOpenidKey(product, info.Openid), nil)
			if err != nil && !storage.IsNotFound(err) {
				return errors.Wrap(err, "delete old index failed")
			}
		default:

		}
	}

	return nil
}
