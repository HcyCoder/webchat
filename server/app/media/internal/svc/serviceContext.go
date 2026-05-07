package svc

import (
	"context"
	"log"

	"github.com/team/webchat-server/app/media/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ServiceContext struct {
	Config     config.Config
	MinioClient *minio.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	client, err := minio.New(c.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.Minio.AccessKey, c.Minio.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("minio init: %v", err)
	}

	exists, err := client.BucketExists(context.Background(), c.Minio.Bucket)
	if err == nil && !exists {
		client.MakeBucket(context.Background(), c.Minio.Bucket, minio.MakeBucketOptions{})
	}

	return &ServiceContext{Config: c, MinioClient: client}
}
