package services

import (
	"context"
	"fmt"
	"io"

	storage "github.com/vkhangstack/hexagonal-architecture/internal/adapters/objectStorage"
	"github.com/vkhangstack/hexagonal-architecture/internal/utils"
)

type UploadService struct {
	storage *storage.S3Adapter
}

func NewUploadService(storage *storage.S3Adapter) *UploadService {
	return &UploadService{
		storage: storage,
	}
}

func (u *UploadService) UploadFile(ctx context.Context, fileName string, fileData io.Reader, contentType string) (string, error) {
	extension := utils.GetFileExtension(fileName)
	keyName := utils.UUIDString() + "." + extension
	err := u.storage.Put(ctx, storage.PutInput{
		Key:         keyName,
		Body:        fileData,
		ContentType: contentType,
	})
	fmt.Println("Upload error:", err)

	if err != nil {
		return "", err
	}
	return keyName, nil
}

func (u *UploadService) DeleteFile(ctx context.Context, fileKey string) error {
	return u.storage.Delete(ctx, fileKey)
}

func (u *UploadService) PublicURL(key string) string {
	return u.storage.PublicURL(key)
}
