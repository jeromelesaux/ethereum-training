package local

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jeromelesaux/ethereum-training/config"
)

func StoreLocalFile(oldFile, filename, hexa256, email string) error {
	path := filepath.Join(config.MyConfig.GetFilepaths(), hexa256)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create directory [%s] with error :%v\n", path, err)
		return err
	}
	newFile := filepath.Join(path, filename)
	if err := os.Rename(oldFile, newFile); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot move file [%s] to [%s] with error :%v\n", oldFile, newFile, err)
		return err
	}
	mailFile := filepath.Join(path, "mail.txt")
	fw, err := os.Create(mailFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create file [%s] with error :%v\n", mailFile, err)
		return err
	}
	defer fw.Close()
	fw.WriteString(email)
	return nil
}
