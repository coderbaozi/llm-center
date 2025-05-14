package sdk

import (
	"context"

	"github.com/openai/openai-go"
)

func EmbeddingTexts(ctx context.Context, texts []string) ([][]float64, error) {
	inputUnion := openai.EmbeddingNewParamsInputUnion(openai.EmbeddingNewParamsInputArrayOfStrings(texts))
	params := openai.EmbeddingNewParams{
		Input: openai.F(inputUnion),
		Model: openai.F("text-embedding-v3"),
	}

	resp, err := QwqClient.Embeddings.New(context.TODO(), params)
	if err != nil {
		return nil, err
	}

	vectors := make([][]float64, len(resp.Data))

	for i, data := range resp.Data {
		vectors[i] = data.Embedding
	}

	return vectors, nil
}
