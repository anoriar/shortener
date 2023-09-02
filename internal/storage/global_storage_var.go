package storage

var GlobalUrlStorage UrlStorageInterface

func init() {
	GlobalUrlStorage = NewUrlStorage()
}
