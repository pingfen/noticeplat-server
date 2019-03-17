package user

import (
	"context"
	"fmt"
	"github.com/bingbaba/storage"
	"github.com/pingfen/noticeplat-server/pkg/errs"
	pkgstore "github.com/pingfen/noticeplat-server/pkg/storage"
	"github.com/pkg/errors"
	"time"
)

type User struct {
	NickName   string `json:"nickname"`   //用户的昵称
	Sex        int    `json:"sex"`        //用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	City       string `json:"city"`       //用户所在城市
	Country    string `json:"country"`    //	用户所在国家
	Province   string `json:"province	"`  //用户所在省份
	HeadimgUrl string `json:"headimgurl"` //用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	UnionId    string `json:"unionid"`    //	只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。

	WechatProducts map[string]*WechatProduct `json:"wechatProducts"`

	CreateTime int64 `json:"createTime"`
}

type WechatProduct struct {
	Openid string `json:"openid"` //用户的标识，对当前公众号唯一
	Remark string `json:"remark"` //公众号运营者对粉丝的备注，公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
	//TagIdList     []int  `json:"tagid_list"`     //用户被打上的标签ID列表

	// 返回用户关注的渠道来源
	// ADD_SCENE_SEARCH 公众号搜索
	// ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移
	// ADD_SCENE_PROFILE_CARD 名片分享
	// ADD_SCENE_QR_CODE 扫描二维码
	// ADD_SCENEPROFILE LINK 图文页内名称点击
	// ADD_SCENE_PROFILE_ITEM 图文页右上角菜单
	// ADD_SCENE_PAID 支付后关注
	// ADD_SCENE_OTHERS 其他
	SubscribeScene string `json:"subscribe_scene"`
	Subscribe      uint8  `json:"subscribe"`      //用户是否订阅该公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息。
	SubscribeTime  int64  `json:"subscribe_time"` //用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间
}

func getUserKey(uid string) string {
	return fmt.Sprintf("/users/data/%s", uid)
}

func getOpenidKey(product, openid string) string {
	return fmt.Sprintf("/users/.index/%s.%s", product, openid)
}

func Register(ctx context.Context, u *User) error {
	if err := u.valid(); err != nil {
		return err
	}

	// storage instance
	store, err := pkgstore.Get()
	if err != nil {
		return err
	}

	// get
	u_old := new(User)
	err = store.Get(ctx, getUserKey(u.UnionId), u_old)
	if err != nil {
		if storage.IsNotFound(err) {
			u.CreateTime = time.Now().Unix()
		} else {
			return errors.Wrap(err, "check identify failed")
		}
	} else { // has registered
		u.CreateTime = u_old.CreateTime

		if u_old.WechatProducts != nil {
			mergeWechatProducts(u_old.WechatProducts, u.WechatProducts)
		}
	}

	if u.CreateTime <= 0 {
		u.CreateTime = time.Now().Unix()
	}

	// save overwrite
	err = store.Create(ctx, getUserKey(u.UnionId), u, 0)
	if err != nil {
		return errors.Wrap(err, "save failed")
	}

	//index
	err = saveIndex(ctx, u_old, u)
	if err != nil {
		return errors.Wrap(err, "save index failed")
	}

	return nil
}

func Get(ctx context.Context, unionid string) (*User, error) {
	// storage instance
	store, err := pkgstore.Get()
	if err != nil {
		return nil, err
	}

	u := new(User)
	err = store.Get(ctx, getUserKey(unionid), u)
	if err != nil {
		return nil, errors.Wrap(err, "read user info failed")
	}

	return u, nil
}

func Update(ctx context.Context, u *User) error {
	if err := u.valid(); err != nil {
		return err
	}

	u_old, err := Get(ctx, u.UnionId)
	if err != nil {
		return err
	}

	if u_old.WechatProducts != nil {
		mergeWechatProducts(u_old.WechatProducts, u.WechatProducts)
	}

	if u.NickName == "" {
		u.NickName = u_old.NickName
	}
	if u.City == "" {
		u.City = u_old.City
	}
	if u.Country == "" {
		u.Country = u_old.Country
	}
	if u.Province == "" {
		u.Province = u_old.Province
	}
	if u.HeadimgUrl == "" {
		u.HeadimgUrl = u_old.HeadimgUrl
	}
	if u.UnionId == "" {
		u.UnionId = u_old.UnionId
	}

	store, _ := pkgstore.Get()
	err = store.Update(ctx, getUserKey(u.UnionId), 0, u, 0)
	if err != nil {
		return errors.Wrap(err, "update user info failed")
	}

	// index
	err = saveIndex(ctx, u_old, u)
	if err != nil {
		return err
	}
	return nil
}

func Delete(ctx context.Context, unionid string) error {
	u, err := Get(ctx, unionid)
	if err != nil {
		return err
	}

	// delete index
	err = deleteIndex(ctx, u)
	if err != nil {
		return err
	}

	// delete user
	store, _ := pkgstore.Get()
	return store.Delete(ctx, getUserKey(unionid), nil)
}

func mergeWechatProducts(map_old, map_new map[string]*WechatProduct) {
	for product, info := range map_old {
		if _, found := map_new[product]; !found {
			map_new[product] = info
		}
	}
}

func (u *User) valid() error {
	if u.UnionId == "" || u.NickName == "" {
		return errors.Wrap(errs.FIELD_EMPTY, "unionid is empty")
	}

	if u.WechatProducts == nil {
		u.WechatProducts = make(map[string]*WechatProduct)
	}

	return nil
}
