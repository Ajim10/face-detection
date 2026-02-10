package facedetection

import (
	"context"
	"fmt"

	"github.com/ajim10/face-detection/internal/database"
	"github.com/ajim10/face-detection/internal/storage"
	"github.com/ajim10/face-detection/internal/vision"
	"github.com/ajim10/face-detection/pkg/imageutil"
	"github.com/ajim10/face-detection/pkg/models"
)

type Service struct {
	db      database.FaceDatabase
	storage storage.Storage
	vision  vision.Vision
}

func NewService(db database.FaceDatabase, storage storage.Storage, vision vision.Vision) *Service {
	return &Service{
		db:      db,
		storage: storage,
		vision:  vision,
	}
}

func (s *Service) ProcessImage(ctx context.Context, gcsURI, filename, projectID, location string) error {
	err := s.storage.UploadImageToCloudStorage(filename)
	if err != nil {
		return err
	}

	annotations, err := s.vision.DetectFaces(ctx, gcsURI)
	if err != nil {
		return err
	}

	face := models.Face{}
	box := models.Box{}
	for i, annotation := range annotations {
		face = models.Face{
			FaceCount: i + 1,
			Anger:     float32(annotation.AngerLikelihood),
			Joy:       float32(annotation.JoyLikelihood),
			Suprise:   float32(annotation.SurpriseLikelihood),
		}
		box = models.Box{
			MinX: int(annotation.BoundingPoly.Vertices[0].X),
			MinY: int(annotation.BoundingPoly.Vertices[0].Y),
			MaxX: int(annotation.BoundingPoly.Vertices[2].X),
			MaxY: int(annotation.BoundingPoly.Vertices[2].Y),
		}
	}

	fmt.Println(face.FaceCount)

	subImg, err := s.storage.SubImage(ctx, gcsURI, box.MinX, box.MinY, box.MaxX, box.MaxY)
	if err != nil {
		return err
	}

	base64Image, err := s.storage.Base64EncodedImage(subImg)
	if err != nil {
		return err
	}

	embedding, err := imageutil.Embedding(projectID, location, base64Image)
	if err != nil {
		return err
	}

	fmt.Println("Len", len(embedding))

	docs, err := s.db.VectorSearch(ctx, embedding)
	if err != nil {
		return err
	}

	for _, doc := range docs {
		fmt.Printf("%v, Distance: %v\n", doc.Data()["imageURI"], (1-doc.Data()["vector_distance"].(float64))*100)
	}

	err = s.db.StoreVectors(ctx, gcsURI, embedding)
	if err != nil {
		return err
	}
	return nil
}
