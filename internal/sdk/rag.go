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

// // EmbedFile 对文件内容进行向量化
// func EmbedFile(ctx context.Context, filePath string) ([][]float32, error) {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	content, err := io.ReadAll(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 按段落分割文件内容
// 	paragraphs := strings.Split(string(content), "\n\n")
// 	return EmbedTexts(ctx, paragraphs)
// }
