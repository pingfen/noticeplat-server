package storage

import (
	"fmt"
	"github.com/bingbaba/storage"
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
		return nil, fmt.Errorf("storage not ready")
	}

	return store, nil
}
