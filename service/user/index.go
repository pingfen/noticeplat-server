package user

//func GetIdByOpenid(ctx context.Context, product, openid string) (string, error) {
//	// storage instance
//	store, err := pkgstore.Get()
//	if err != nil {
//		return "", err
//	}
//
//	var id string
//	key := getOpenidKey(product, openid)
//	err = store.Get(ctx, key, &id)
//	if err != nil {
//		return "", fmt.Errorf("read key %s failed %w", key, err)
//	}
//
//	return id, nil
//}
//
////func GetByOpenid(ctx context.Context, product, openid string) (*User, error) {
////	id, err := GetIdByOpenid(ctx, product, openid)
////	if err != nil {
////		return nil, err
////	}
////
////	return Get(ctx, id)
////}
//
//func deleteIndex(ctx context.Context, u *User) error {
//	// storage instance
//	store, _ := pkgstore.Get()
//
//	for product, info := range u.WechatProducts {
//		err := store.Delete(ctx, getOpenidKey(product, info.Openid), nil)
//		if err != nil && !storage.IsNotFound(err) {
//			return err
//		}
//	}
//	if u.UnionId != "" {
//		err := store.Delete(ctx, getUnionidKey(u.UnionId), nil)
//		if err != nil && !storage.IsNotFound(err) {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func saveIndex(ctx context.Context, u_old, u *User) error {
//	// storage instance
//	store, _ := pkgstore.Get()
//
//	// unionid
//	if u_old == nil || (u_old.UnionId == "" && u.UnionId != "") {
//		err := store.Create(ctx, getUnionidKey(u.UnionId), u.Id, 0)
//		if err != nil {
//			return fmt.Errorf("save new index failed: %w", err)
//		}
//	}
//
//	// 0:nochange 1:add 2:modify 3:delete
//	operation := make(map[string]int)
//
//	if u_old == nil || u_old.WechatProducts == nil {
//		for product := range u.WechatProducts {
//			operation[product] = 1
//		}
//	} else {
//		for product, info := range u.WechatProducts {
//			if info_old, found := u_old.WechatProducts[product]; found {
//				if info.Openid != info_old.Openid {
//					operation[product] = 2
//				} else {
//					operation[product] = 0
//				}
//			} else {
//				operation[product] = 1
//			}
//		}
//		for product := range u_old.WechatProducts {
//			if _, found := operation[product]; !found {
//				operation[product] = 3
//			}
//		}
//	}
//
//	//index
//	//log.Printf("%+v", operation)
//	for product, op := range operation {
//		switch op {
//		case 1:
//			info := u.WechatProducts[product]
//			err := store.Create(ctx, getOpenidKey(product, info.Openid), u.Id, 0)
//			if err != nil {
//				return fmt.Errorf("save new index failed: %w", err)
//			}
//		case 2:
//			info := u.WechatProducts[product]
//			err := store.Create(ctx, getOpenidKey(product, info.Openid), u.Id, 0)
//			if err != nil {
//				return fmt.Errorf("save new index failed: %w", err)
//			}
//
//			info_old := u_old.WechatProducts[product]
//			err = store.Delete(ctx, getOpenidKey(product, info_old.Openid), nil)
//			if err != nil && !storage.IsNotFound(err) {
//				return fmt.Errorf("delete old index failed: %w", err)
//			}
//		case 3:
//			info := u_old.WechatProducts[product]
//			err := store.Delete(ctx, getOpenidKey(product, info.Openid), nil)
//			if err != nil && !storage.IsNotFound(err) {
//				return fmt.Errorf("delete old index failed: %w", err)
//			}
//		default:
//
//		}
//	}
//
//	return nil
//}
