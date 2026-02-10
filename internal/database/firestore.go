package database

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

type FirestoreClient struct {
	client *firestore.Client
}

type Features struct {
	CreatedAt time.Time `firestore:"createdAt" json:"createdAt"`
	ImageURI  string    `firestore:"imageURI" json:"imageURI"`
	Embedding []float32 `firestore:"embedding" json:"embedding"`
}

func NewFirestoreClient(ctx context.Context, projectID string) (*FirestoreClient, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &FirestoreClient{client: client}, nil
}

func (f *FirestoreClient) StoreVectors(ctx context.Context, imageURI string, embedding []float32) error {
	features := Features{
		CreatedAt: time.Now(),
		ImageURI:  imageURI,
		Embedding: embedding,
	}
	docRef := f.client.Collection("face-detection").Doc(uuid.New().String())
	if _, err := docRef.Set(ctx, features); err != nil {
		return err
	}
	return nil
}

func (f *FirestoreClient) VectorSearch(ctx context.Context, embedding []float32) ([]*firestore.DocumentSnapshot, error) {
	vectorQuery := f.client.Collection("face-detection").FindNearest("embedding", embedding, 10, firestore.DistanceMeasureCosine, &firestore.FindNearestOptions{
		DistanceResultField: "vector_distance",
	})
	docs, err := vectorQuery.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	return docs, nil
}
