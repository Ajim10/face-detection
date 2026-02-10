package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ajim10/face-detection/internal/config"
	"github.com/ajim10/face-detection/internal/database"
	"github.com/ajim10/face-detection/internal/facedetection"
	st "github.com/ajim10/face-detection/internal/storage"
	vi "github.com/ajim10/face-detection/internal/vision"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Config Error:", err)
	}

	bucket := "face-detection"
	object := "test.jpg"
	filename := "profile.jpg"
	gcsURI := fmt.Sprintf("gs://%s/%s", bucket, object)

	storageClient, err := st.NewCloudStorageClient(ctx, bucket, object)
	if err != nil {
		log.Fatal("StorageClient Error:", err)
	}

	firestoreClient, err := database.NewFirestoreClient(ctx, cfg.ProjectID)
	if err != nil {
		log.Fatal("FirestoreClient Error:", err)
	}

	visionClient, err := vi.NewVisionClient(ctx)
	if err != nil {
		log.Fatal("VisionClient Error:", err)
	}

	s := facedetection.NewService(firestoreClient, storageClient, visionClient)

	err = s.ProcessImage(ctx, gcsURI, filename, cfg.ProjectID, cfg.Location)
	if err != nil {
		log.Fatal("ProcessImage Error:", err)
	}

}
