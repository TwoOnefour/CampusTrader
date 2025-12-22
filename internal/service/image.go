package service

import (
	"context"
	"io"
)

type ImageService struct {
	storage Storage
}

type Storage interface {
	Save(ctx context.Context, file io.Reader, path string, size int64, contentType string) (string, error)
	Delete(ctx context.Context, path string) error
	GetURL(path string) string
}

func NewImageService(storage Storage) *ImageService {
	return &ImageService{
		storage: storage,
	}
}

func (s *ImageService) Save(ctx context.Context, file io.Reader, path string, size int64, contentType string) (string, error) {
	return s.storage.Save(ctx, file, path, size, contentType)
}
