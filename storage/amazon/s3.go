package amazon

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	awsSessionLoader sync.Once
	awsSession       *session.Session
)

func UploadBuffer(fileBuffer []byte, filePath, region, bucket string) error {
	if err := getSession(region); err != nil {
		return err
	}

	fileSize := int64(len(fileBuffer))
	_, err := s3.New(awsSession).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(filePath),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

func Upload(filePath, region, bucket string) error {
	if err := getSession(region); err != nil {
		return err
	}
	upFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer upFile.Close()
	upFileInfo, err := upFile.Stat()
	if err != nil {
		return err
	}
	fileSize := upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	_, err = s3.New(awsSession).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(filePath),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

func getSession(region string) (err error) {
	awsSessionLoader.Do(func() {
		awsSession, err = session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
	})
	return err
}

func Download(filePath, region, bucket string) (string, error) {
	if err := getSession(region); err != nil {
		return "", err
	}
	downFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer downFile.Close()

	downloader := s3manager.NewDownloader(awsSession)

	_, err = downloader.Download(downFile,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filePath),
		})
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func DownloadBuffer(filePath, region, bucket string) ([]byte, error) {
	if err := getSession(region); err != nil {
		return nil, err
	}
	buffer := make([]byte, 0)
	wbuffer := aws.NewWriteAtBuffer(buffer)
	downloader := s3manager.NewDownloader(awsSession)

	n, err := downloader.Download(wbuffer,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(filePath),
		})
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return buffer, fmt.Errorf("Expected content file superior to 0 for filepath %s", filePath)
	}
	return buffer, nil
}
