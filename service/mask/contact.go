package mask

import (
	"context"
	"fmt"
	"github.com/pingfen/notify/pkg/models"
	pkgstore "github.com/pingfen/notify/storage"
	"time"
)

func SaveContacts(ctx context.Context, uid string, mui *models.MaskUserInfo) error {
	if mui.Id == "" {
		mui.Id = fmt.Sprintf("%d", time.Now().Unix())
	}

	// storage instance
	store, err := pkgstore.Get()
	if err != nil {
		return err
	}

	return store.Create(ctx, getContactKey(uid, mui.Id), mui, 0)
}

func ListContacts(ctx context.Context, uid string) ([]*models.MaskUserInfo, error) {
	// storage instance
	store, err := pkgstore.Get()
	if err != nil {
		return nil, err
	}

	ret, err := store.List(ctx, getContactsDirKey(uid), nil, new(models.MaskUserInfo))
	if err != nil {
		return nil, err
	}

	us := make([]*models.MaskUserInfo, len(ret))
	for i, item := range ret {
		us[i] = item.(*models.MaskUserInfo)
	}

	return us, nil
}

func GetContact(ctx context.Context, user, contactUser string) (*models.MaskUserInfo, error) {
	// storage instance
	store, err := pkgstore.Get()
	if err != nil {
		return nil, err
	}

	u := new(models.MaskUserInfo)
	err = store.Get(ctx, getContactKey(user, contactUser), u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func DeleteContact(ctx context.Context, user, contactUser string) error {
	// storage instance
	store, err := pkgstore.Get()
	if err != nil {
		return err
	}

	u := new(models.MaskUserInfo)
	err = store.Delete(ctx, getContactKey(user, contactUser), u)
	if err != nil {
		return err
	}

	return nil
}

func getContactKey(openid, cid string) string {
	return fmt.Sprintf("/contacts/%s/%s", openid, cid)
}

func getContactsDirKey(openid string) string {
	return fmt.Sprintf("/contacts/%s/", openid)
}
