package oss

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	MinioClient   *minio.Client
	MinioEndpoint string
)

// InitMinioClient 初始化MinIO客户端
func InitMinioClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) error {
	var err error
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return err
	}

	MinioEndpoint = endpoint
	return nil
}
