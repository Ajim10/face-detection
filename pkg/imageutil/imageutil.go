package imageutil

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	aiplatform "cloud.google.com/go/aiplatform/apiv1beta1"
	"cloud.google.com/go/aiplatform/apiv1beta1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func Embedding(projectID, location, base64Image string) ([]float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	apiEndpoint := fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)
	client, err := aiplatform.NewPredictionClient(ctx, option.WithEndpoint(apiEndpoint))
	if err != nil {
		return nil, fmt.Errorf("aiplatform.NewPredictionClient: %v", err)
	}
	defer client.Close()

	model := os.Getenv("MODEL")
	endpoint := fmt.Sprintf("projects/%s/locations/%s/publishers/google/models/%s", projectID, location, model)

	instance, err := structpb.NewValue(map[string]any{
		"image": map[string]any{
			"bytesBase64Encoded": base64Image,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("structpb.NewValue: %v", err)
	}

	req := &aiplatformpb.PredictRequest{
		Endpoint:  endpoint,
		Instances: []*structpb.Value{instance},
	}

	resp, err := client.Predict(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("client.Predict: %v", err)
	}

	instanceEmbeddingsJson, err := protojson.Marshal(resp.GetPredictions()[0])
	if err != nil {
		return nil, fmt.Errorf("protojson.Marshal: %v", err)
	}

	var instanceEmbeddings struct {
		ImageEmbeddings []float32 `json:"imageEmbedding"`
	}

	if err := json.Unmarshal(instanceEmbeddingsJson, &instanceEmbeddings); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	imageEmbedding := instanceEmbeddings.ImageEmbeddings

	return imageEmbedding, nil
}
