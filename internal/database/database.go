package database

import (
	"context"

	"cloud.google.com/go/firestore"
)

type FaceDatabase interface {
	StoreVectors(ctx context.Context, imageURI string, embedding []float32) error
	VectorSearch(ctx context.Context, embedding []float32) ([]*firestore.DocumentSnapshot, error)
}
