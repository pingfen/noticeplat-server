package storage

import (
	"github.com/bingbaba/storage"
	"github.com/pingfen/noticeplat-server/pkg/errs"
	"github.com/pkg/errors"
	"sync"
)

var (
	store storage.Interface
	once  sync.Once
)

func Init(s storage.Interface) {
	if s == nil {
		return
	}

	once.Do(func() {
		store = s
	})
}

func Get() (storage.Interface, error) {

	if store == nil {
		return nil, errors.Wrap(errs.STORAGE_ERROR, "storage not ready")
	}

	return store, nil
}
