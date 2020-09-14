package local

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

func StoreLocalFileBuffer(fileBuffer []byte, filename, hexa256, email string) error {
	path := filepath.Join(config.MyConfig.GetFilepaths(), hexa256)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create directory [%s] with error :%v\n", path, err)
		return err
	}
	filePath := filepath.Join(path, filename)

	fw, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create file [%s] with error :%v\n", filePath, err)
		return err
	}
	defer fw.Close()
	fw.Write(fileBuffer)

	mailFile := filepath.Join(path, "mail.txt")
	fmw, err := os.Create(mailFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create file [%s] with error :%v\n", mailFile, err)
		return err
	}
	defer fmw.Close()
	fmw.WriteString(email)
	return nil
}

func GetLocalFile(directoryPath string) (fileName string, err error) {

	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		err = errors.New("no file found")
		return "", err
	}

	for _, v := range files {
		switch v.Name() {
		case "mail.txt":
			break
		default:
			fileName = v.Name()
		}
	}
	return fileName, nil
}

func GetLocalFileBuffer(filePath string) ([]byte, error) {

	fr, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fr.Close()
	buffer := bytes.NewBuffer(nil)
	_, err = io.Copy(buffer, fr)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func getEmail(filePath string) (string, error) {
	fo, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer fo.Close()
	content, err := ioutil.ReadAll(fo)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
