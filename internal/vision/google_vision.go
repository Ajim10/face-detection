package vision

import (
	"context"

	vision "cloud.google.com/go/vision/v2/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
)

type VisionClient struct {
	client *vision.ImageAnnotatorClient
}

func NewVisionClient(ctx context.Context) (*VisionClient, error) {
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}
	return &VisionClient{client: client}, nil
}

func (v *VisionClient) DetectFaces(ctx context.Context, gcsURI string) ([]*visionpb.FaceAnnotation, error) {
	req := &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{
			{
				Image: &visionpb.Image{
					Source: &visionpb.ImageSource{
						GcsImageUri: gcsURI,
					},
				},
				Features: []*visionpb.Feature{
					{
						Type:       visionpb.Feature_FACE_DETECTION,
						MaxResults: 5,
					},
				},
			},
		},
	}

	resp, err := v.client.BatchAnnotateImages(ctx, req)
	if err != nil {
		return nil, err
	}

	annotations := resp.Responses[0].FaceAnnotations
	return annotations, nil
}
