package storage

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"os"
)

type S3Storage struct {
	client *s3.Client
	bucket string
}

func NewS3Storage(bucket string) *S3Storage {
	return &S3Storage{
		client: s3.New(s3.Options{}),
		bucket: bucket,
	}
}

func (s *S3Storage) Save(ctx context.Context, file io.Reader, path string, size int64, contentType string) (string, error) {
	client := s.client
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   file,
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

func (s *S3Storage) Delete(ctx context.Context, path string) error {
	return nil
}

func (s *S3Storage) GetURL(path string) string {
	return ""
}
