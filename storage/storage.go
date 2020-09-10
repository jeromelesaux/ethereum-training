package storage

import (
	"path/filepath"

	"github.com/jeromelesaux/ethereum-training/storage/amazon"
	"github.com/jeromelesaux/ethereum-training/storage/local"
)

func StoreFile(oldFile, filename, hexa256, email, region, bucket string, storeLocally bool) (err error) {
	if storeLocally {
		return local.StoreLocalFile(oldFile, filename, hexa256, email)
	}
	return amazon.Upload(oldFile, region, bucket)
}

func GetFile(directoryPath, filename, region, bucket string, storeLocally bool) (filePath string, err error) {
	if storeLocally {
		return local.GetLocalFile(directoryPath)
	}
	dir := filepath.Join(directoryPath, filename)
	return amazon.Download(dir, region, bucket)
}
