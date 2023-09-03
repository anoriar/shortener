package storage

import "sync"

var urlStorage URLRepositoryInterface
var urlStorageSyncOnce sync.Once

func GetInstance() URLRepositoryInterface {
	urlStorageSyncOnce.Do(func() {
		urlStorage = newURLStorage(make(map[string]string))
	})
	return urlStorage
}
