// Package storage internal/common/storage/local.go
package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	BaseDir string // 本地存放的根目录，例如 "./static/uploads"
	BaseURL string // 对应的访问前缀，例如 "/static/uploads"
}

func NewLocalStorage(baseDir, baseURL string) *LocalStorage {
	// 确保存储目录存在
	_ = os.MkdirAll(baseDir, os.ModePerm)
	return &LocalStorage{
		BaseDir: baseDir,
		BaseURL: baseURL,
	}
}

func (s *LocalStorage) Save(ctx context.Context, file io.Reader, path string, size int64, contentType string) (string, error) {
	fullPath := filepath.Join(s.BaseDir, path)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return s.GetURL(path), nil
}

func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.BaseDir, path)
	return os.Remove(fullPath)
}

func (s *LocalStorage) GetURL(path string) string {
	// 简单拼接，注意处理路径分隔符
	return fmt.Sprintf("%s/%s", s.BaseURL, path)
}
