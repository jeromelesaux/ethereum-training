package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jeromelesaux/ethereum-training/config"
	"github.com/jeromelesaux/ethereum-training/storage/amazon"
	"github.com/jeromelesaux/ethereum-training/storage/local"
)

func StoreFile(oldFile, filename, hexa256, email, region, bucket string, storeLocally bool) (err error) {
	if storeLocally {
		return local.StoreLocalFile(oldFile, filename, hexa256, email)
	}
	return amazon.Upload(oldFile, region, bucket)
}

func GetFile(filename, hash, region, bucket string, storeLocally bool) (string, error) {
	filePath := GetFilepath(filename, hash, storeLocally)
	fmt.Fprintf(os.Stdout, "filepath [%s]\n", filePath)
	if storeLocally {
		return filePath, nil
	}
	return amazon.Download(filePath, region, bucket)
}

func GetFilepath(filename, hash string, storeLocally bool) string {
	if storeLocally {
		filePath := filepath.Join(filepath.Join(config.MyConfig.GetFilepaths(), hash), filename)
		return filePath
	}
	return filepath.Join(config.MyConfig.GetFilepaths(), filename)
}
