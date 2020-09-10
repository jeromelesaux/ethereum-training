package storage

import (
	"github.com/jeromelesaux/ethereum-training/storage/amazon"
	"github.com/jeromelesaux/ethereum-training/storage/local"
)

func StoreFile(oldFile, filename, hexa256, email, region, bucket string, storeLocally bool) (err error) {
	if storeLocally {
		return local.StoreLocalFile(oldFile, filename, hexa256, email)
	}
	return amazon.Upload(oldFile, region, bucket)
}

func GetFile(directoryPath, region, bucket string, storeLocally bool) (filePath string, err error) {
	if storeLocally {
		return local.GetLocalFile(directoryPath)
	}
	return amazon.Download(directoryPath, region, bucket)
}
