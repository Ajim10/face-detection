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

	gcsURI := fmt.Sprintf("gs://%s/%s", cfg.Bucket, cfg.Object)

	storageClient, err := st.NewCloudStorageClient(ctx, cfg.Bucket, cfg.Object)
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

	err = s.ProcessImage(ctx, gcsURI, cfg.Filename, cfg.ProjectID, cfg.Location)
	if err != nil {
		log.Fatal("ProcessImage Error:", err)
	}

}
