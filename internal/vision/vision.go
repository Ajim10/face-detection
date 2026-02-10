package vision

import (
	"context"

	"cloud.google.com/go/vision/v2/apiv1/visionpb"
)

type Vision interface {
	DetectFaces(ctx context.Context, gcsURI string) ([]*visionpb.FaceAnnotation, error)
}
