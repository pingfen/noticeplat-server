package user

import (
	"context"
	"fmt"
	"github.com/bingbaba/storage"
	"github.com/pingfen/notify/pkg/models"
	"github.com/pingfen/notify/pkg/wechat"
	pkgstore "github.com/pingfen/notify/storage"
	"gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
	"time"
)

func LoginByCode(ctx context.Context, source, code string) (*models.MpUser, error) {
	sess, err := wechat.GetMiniProgramSession(source, code)
	if err != nil {
		return nil, fmt.Errorf("登录失败%w", err)
	}

	u, err := GetByOpenid(ctx, source, code)
	if err != nil {
		if storage.IsNotFound(err) {
			// 自动注册
			return RegisterByMpSession(ctx, source, sess)
		} else {
			return nil, fmt.Errorf("查询用户失败%w", err)
		}
	}

	return u, nil
}

func GetByOpenid(ctx context.Context, source, openid string) (*models.MpUser, error) {
	// storage instance
	store, err := pkgstore.Get()
	if err != nil {
		return nil, err
	}

	u := new(models.MpUser)
	err = store.Get(ctx, getKeyByOpenid(source, openid), u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func GetMpByUnionid(ctx context.Context, source, unionid string) (*models.MpUser, error) {
	// storage instance
	store, err := pkgstore.Get()
	if err != nil {
		return nil, err
	}

	var openid string
	err = store.Get(ctx, getIndexKeyByUnionid(source, unionid), &openid)
	if err != nil {
		return nil, err
	}

	return GetByOpenid(ctx, source, openid)
}

func RegisterByMpSession(ctx context.Context, source string, sess *oauth2.Session) (*models.MpUser, error) {
	// storage instance
	store, err := pkgstore.Get()
	if err != nil {
		return nil, err
	}

	u := &models.MpUser{
		OpenId:  sess.OpenId,
		UnionId: sess.UnionId,
	}

	// save overwrite
	u.CreateTime = time.Now().Unix()
	err = store.Create(ctx, getKeyByOpenid(source, u.OpenId), u, 0)
	if err != nil {
		return nil, err
	}

	//index
	err = saveMpUserIndex(ctx, source, u)
	if err != nil {
		return nil, fmt.Errorf("保存Unionid索引数据失败%w", err)
	}

	return u, nil
}

func saveMpUserIndex(ctx context.Context, source string, u *models.MpUser) error {
	store, err := pkgstore.Get()
	if err != nil {
		return err
	}

	if u.UnionId != "" {
		err = store.Create(ctx, getIndexKeyByUnionid(source, u.UnionId), u.OpenId, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func getKeyByOpenid(source, openid string) string {
	return fmt.Sprintf("/users/mp/%s/%s", source, openid)
}

func getIndexKeyByUnionid(source, unionid string) string {
	return fmt.Sprintf("/users/mp/_unionid_/%s/.%s", source, unionid)
}
