package service

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"time"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

)

type S3Service struct {
	Client     *s3.Client
	BucketName string
}

func NewS3Service() (*S3Service, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Printf("[S3] Failed to load AWS config: %v", err)
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	bucket := os.Getenv("S3_BUCKET")

	log.Printf("[S3] Initialized S3 client | region=%s bucket=%s", cfg.Region, bucket)

	return &S3Service{
		Client:     client,
		BucketName: bucket,
	}, nil
}

func (s *S3Service) UploadFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		log.Printf("[S3] cannot open file: %v", err)
		return "", err
	}
	defer src.Close()

	uploader := manager.NewUploader(s.Client)

	key := fmt.Sprintf("avatars/%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	log.Printf("[S3] 📤 Uploading file to bucket=%s key=%s", s.BucketName, key)

	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   src,
	})
	if err != nil {
		log.Printf("[S3] Upload failed: %v", err)
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.BucketName, key)
	log.Printf("[S3] File uploaded successfully: %s", url)
	return url, nil
}
