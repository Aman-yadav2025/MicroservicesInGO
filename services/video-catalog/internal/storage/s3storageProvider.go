package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	minioClient "github.com/minio/minio-go/v7"
	minioCreds "github.com/minio/minio-go/v7/pkg/credentials"
)

type LocalStorage struct {
	BaseDir string
}

type MinioStorage struct {
	client     *minioClient.Client
	bucketName string
	endpoint   string
}

func NewLocalStorage(baseDir string) *LocalStorage {
	os.MkdirAll(baseDir, os.ModePerm)
	return &LocalStorage{BaseDir: baseDir}
}

func (l *LocalStorage) UploadFile(ctx context.Context, file io.Reader, filename string, size int64, contentType string) (string, error) {
    savePath := filepath.Join(l.BaseDir, filename)
    dst, err := os.Create(savePath)
    if err != nil {
        return "", err
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        return "", err
    }
    return fmt.Sprintf("uploads/videos/%s", filename), nil
}

func NewMinioStorage(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*MinioStorage, error) {
	MinioClient, err := minioClient.New(endpoint, &minioClient.Options{
		Creds:  minioCreds.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := MinioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		err := MinioClient.MakeBucket(ctx, bucketName, minioClient.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("Failed to Create Bucket %w", err)
		}
		log.Printf("%s Bucket created in MINIO",bucketName)
	}

	return &MinioStorage{
		client:     MinioClient,
		bucketName: bucketName,
		endpoint:   endpoint,
	}, nil
}

func (m *MinioStorage)UploadFile(ctx context.Context, file io.Reader,fileName string,size int64,contentType string)(string,error){
	_ ,err := m.client.PutObject(ctx,m.bucketName,fileName,file,size,minioClient.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "",fmt.Errorf("Error in Uploading the file: %w",err)
	}

	fileUrl := fmt.Sprintf("http://%s/%s/%s",m.endpoint,m.bucketName,fileName)

	return fileUrl,nil
}
