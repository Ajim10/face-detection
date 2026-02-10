package storage

import (
	"context"
	"image"
)

type Storage interface {
	UploadImageToCloudStorage(filename string) error
	LoadImage(ctx context.Context, objectPath string) (image.Image, error)
	Base64EncodedImage(subImg image.Image) (string, error)
	SubImage(ctx context.Context, objectPath string, minX, minY, maxX, maxY int) (image.Image, error)
}
