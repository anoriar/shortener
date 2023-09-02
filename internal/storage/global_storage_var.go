package storage

import "sync"

var urlStorageInstance URLStorageInterface
var once sync.Once
