package models

import (
	"time"

	"cloud.google.com/go/firestore"
)

type Features struct {
	CreatedAt time.Time          `firestore:"createdAt" json:"createdAt"`
	ImageURI  string             `firestore:"imageURI" json:"imageURI"`
	Embedding firestore.Vector32 `firestore:"embedding" json:"embedding"`
}

type Box struct {
	MinX int
	MinY int
	MaxX int
	MaxY int
}

type Face struct {
	FaceCount int
	Anger     float32
	Joy       float32
	Suprise   float32
}
