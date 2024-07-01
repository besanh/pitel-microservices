package minio

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tel4vn/fins-microservices/common/log"
)

type (
	IMinIO interface {
		GetClient() *minio.Client
		CreateBucket(ctx context.Context, bucketName string) (err error)
		UploadObject(ctx context.Context, bucketName string, objectName string, objectPath string, contentType string) (info minio.UploadInfo, err error)
		PresignedGetObject(ctx context.Context, bucketName string, objectName string, expires time.Duration) (presignedURL string, err error)
	}
	Config struct {
		AccessKeyID     string
		SecretAccessKey string
		Endpoint        string
		UseSSL          bool
		Region          string
	}
	MinIO struct {
		Config Config
		Client *minio.Client
	}
)

var MinIOClient IMinIO

func NewClient(cfg Config) IMinIO {
	s := &MinIO{
		Config: cfg,
	}

	if err := s.Connect(); err != nil {
		log.Fatal(err)
		return nil
	}
	return s
}

func (s *MinIO) Connect() (err error) {
	// Initialize minio client object.
	minioClient, err := minio.New(s.Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s.Config.AccessKeyID, s.Config.SecretAccessKey, ""),
		Secure: s.Config.UseSSL,
	})
	if err != nil {
		return err
	}
	s.Client = minioClient
	return
}

func (s *MinIO) GetClient() *minio.Client {
	return s.Client
}

func (s *MinIO) CreateBucket(ctx context.Context, bucketName string) (err error) {
	err = s.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: s.Config.Region})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := s.Client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Infof("We already own %s", bucketName)
			return nil
		}
	} else {
		log.Infof("Successfully created %s", bucketName)
	}
	return
}

func (s *MinIO) UploadObject(ctx context.Context, bucketName string, objectName string, objectPath string, contentType string) (info minio.UploadInfo, err error) {
	// Upload the test file with FPutObject
	info, err = s.Client.FPutObject(ctx, bucketName, objectName, objectPath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return
	}
	return
}

func (s *MinIO) PresignedGetObject(ctx context.Context, bucketName string, objectName string, expires time.Duration) (presignedURL string, err error) {
	// Generates a presigned url which expires in a day.
	u, err := s.Client.PresignedGetObject(ctx, bucketName, objectName, expires, nil)
	if err != nil {
		return
	}
	presignedURL = u.String()
	return
}
