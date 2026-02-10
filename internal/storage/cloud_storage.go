package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/jpeg"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

type CloudStorageClient struct {
	client *storage.Client
	bucket string
	object string
}

func NewCloudStorageClient(ctx context.Context, bucket, object string) (*CloudStorageClient, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &CloudStorageClient{
		client: client,
		bucket: bucket,
		object: object,
	}, nil
}

func (c *CloudStorageClient) LoadImage(ctx context.Context, objectPath string) (image.Image, error) {
	rc, err := c.client.Bucket(c.bucket).Object(objectPath).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	img, _, err := image.Decode(rc)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (c *CloudStorageClient) UploadImageToCloudStorage(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	o := c.client.Bucket(c.bucket).Object(c.object)

	wc := o.NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}

func (c *CloudStorageClient) SubImage(ctx context.Context, objectPath string, minX, minY, maxX, maxY int) (image.Image, error) {
	faceBounds := image.Rect(minX, minY, maxX, maxY)
	img, err := c.LoadImage(ctx, objectPath)
	if err != nil {
		return nil, err
	}
	subImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(faceBounds)

	f, err := os.Create("face.jpg")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = jpeg.Encode(f, subImg, nil)
	if err != nil {
		return nil, err
	}
	return subImg, nil
}

func (c *CloudStorageClient) Base64EncodedImage(subImg image.Image) (string, error) {
	buf := bytes.Buffer{}
	err := jpeg.Encode(&buf, subImg, nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
