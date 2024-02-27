package storage

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tel4vn/fins-microservices/common/env"
	"github.com/tel4vn/fins-microservices/common/log"
)

type S3Config struct {
	AccessKey   string
	SecreteKey  string
	Endpoint    string
	UseSSL      bool
	BucketName  string
	Location    string
	ContentType string
}

func (config S3Config) String() string {
	return fmt.Sprintf("AccessKey:%s Endpoint:%s UseSSL:%t BucketName:%s Location:%s ContentType:%s",
		config.AccessKey, config.Endpoint, config.UseSSL, config.BucketName, config.Location, config.ContentType)
}

type InternalConfig struct {
	Location string
}

type StoreInput struct {
	Byte []byte
	Path string
}

func NewStoreInput(byte []byte, path string) *StoreInput {
	return &StoreInput{Byte: byte, Path: path}
}

type RetrieveInput struct {
	Path string
}

func NewRetrieveInput(path string) *RetrieveInput {
	return &RetrieveInput{Path: path}
}

func InitStorage() {
	s3Config := S3Config{
		Endpoint:    env.GetStringENV("STORAGE_ENDPOINT", "localhost"),
		AccessKey:   env.GetStringENV("STORAGE_ACCESS_KEY", ""),
		SecreteKey:  env.GetStringENV("STORAGE_SECRET_KEY", ""),
		ContentType: env.GetStringENV("STORAGE_CONTENT_TYPE", "application/octet-stream"),
		BucketName:  env.GetStringENV("STORAGE_BUCKET_NAME", "dev_fins_document"),
		Location:    env.GetStringENV("STORAGE_LOCATION", "us-east-1"),
		UseSSL:      env.GetBoolENV("STORAGE_USE_SSL", true),
	}
	minioClient, err := minio.New(s3Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3Config.AccessKey, s3Config.SecreteKey, ""),
		Secure: s3Config.UseSSL,
	})
	if err != nil {
		log.Error(err)
	} else {
		log.Debug("Storage client is initialize")
	}
	Instance = NewS3Storage(minioClient, s3Config)
}

var Instance IStorage

type IStorage interface {
	Store(ctx context.Context, input StoreInput) (bool, error)
	Retrieve(ctx context.Context, input RetrieveInput) ([]byte, error)
	RemoveFile(ctx context.Context, input RetrieveInput) (err error)
}

type S3Storage struct {
	Client *minio.Client
	Config S3Config
}

func NewS3Storage(client *minio.Client, config S3Config) *S3Storage {
	return &S3Storage{Client: client, Config: config}
}

func (s *S3Storage) Store(ctx context.Context, input StoreInput) (result bool, err error) {
	cloudStorageConfig := s.Config
	bucketName := cloudStorageConfig.BucketName
	objectName := input.Path
	minioClient := s.Client
	fileBytes := input.Byte
	reader := bytes.NewReader(fileBytes)
	objectSize := reader.Size()
	putObjectOptions := minio.PutObjectOptions{ContentType: cloudStorageConfig.ContentType}

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: cloudStorageConfig.Location})
	if err != nil {
		exists, errBucketExists := s.Client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Infof("Storage connected to bucket: %s", bucketName)
			_, err = minioClient.PutObject(ctx, bucketName, objectName, reader, objectSize, putObjectOptions)
			if err != nil {
				log.Error(err)
				return false, err
			}
		} else {
			log.Error(errBucketExists)
			return false, errBucketExists
		}
	} else {
		log.Infof("Successfully created and connected to bucket %s", bucketName)
		_, err = minioClient.PutObject(ctx, bucketName, objectName, reader, objectSize, putObjectOptions)
		if err != nil {
			log.Error(err)
			return false, err
		}
	}
	return true, nil
}

func (s *S3Storage) Retrieve(ctx context.Context, input RetrieveInput) ([]byte, error) {
	cloudStorageConfig := s.Config
	filePath := input.Path
	bucketName := cloudStorageConfig.BucketName
	objectName := filePath
	minioClient := s.Client

	log.Infof("Query '%s' from cloud storage", objectName)
	object, err := minioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Error(err)
		return nil, err
	} else {
		log.Infof("Successfully get '%s' from cloud storage", objectName)
	}

	// Get the file size
	stat, err := object.Stat()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Read the file into a byte slice
	fileBytes := make([]byte, stat.Size)
	_, err = bufio.NewReader(object).Read(fileBytes)
	if err != nil && err != io.EOF {
		log.Error(err)
		return nil, err
	}

	log.Info("Retrieve file success")
	return fileBytes, nil

}

func (s *S3Storage) RemoveFile(ctx context.Context, input RetrieveInput) (err error) {
	cloudStorageConfig := s.Config
	bucketName := cloudStorageConfig.BucketName
	objectName := input.Path
	minioClient := s.Client
	exists, errBucketExists := s.Client.BucketExists(ctx, bucketName)
	if errBucketExists == nil && exists {
		err = minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
		if err != nil {
			log.Error(err)
			return err
		}
		log.Infof("Successfully remove '%s' from cloud storage", objectName)
		return
	}
	log.Errorf("Bucket '%s' not found", bucketName)
	return errors.New("Bucket " + bucketName + " not found")
}
