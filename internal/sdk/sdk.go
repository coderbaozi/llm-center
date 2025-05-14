package sdk

import (
	"bytes"
	"context" // 添加 context 包
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings" // 添加 strings 包

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var QwqClient *openai.Client

func InitClients() {
	QwqClient = openai.NewClient(
		option.WithAPIKey("sk-00711d2a578e46da9fb882245230a3eb"),
		option.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1/"),
	)
}

func GetQwqClient() *openai.Client {
	return QwqClient
}

func CreateCollection(embeddings [][]float64, collection_id *string) error {
	url := "http://localhost:2333/vector/create_collection"
	client := &http.Client{}

	if len(embeddings) == 0 || len(embeddings[0]) == 0 {
		return fmt.Errorf("embeddings cannot be empty")
	}
	dimension := len(embeddings[0])

	payload := map[string]interface{}{
		"collection_name": collection_id,
		"dimension":       dimension,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func AddEmbeddings(text *string, id *string) error {
	url := "http://localhost:2333/vector/insert"
	client := &http.Client{}

	if text == nil || *text == "" {
		return fmt.Errorf("input text cannot be empty")
	}

	// text 一句话一组按照，, 分割
	processedText := strings.ReplaceAll(*text, ",", "，") // 将英文逗号替换为中文逗号
	texts := strings.Split(processedText, "，")
	if len(texts) == 0 {
		return fmt.Errorf("after splitting, no texts to process")
	}

	embeddings, err := EmbeddingTexts(context.TODO(), texts)
	if err != nil {
		return fmt.Errorf("failed to get embeddings: %w", err)
	}

	if err := CreateCollection(embeddings, id); err != nil {
		// It's possible the collection already exists, so we might not want to fail hard here
		// For now, let's assume CreateCollection handles existing collections or we want to error out
		return fmt.Errorf("failed to create collection: %w", err)
	}

	if len(embeddings) == 0 || len(texts) == 0 || len(embeddings) != len(texts) {
		return fmt.Errorf("embeddings and texts mismatch or are empty. Embeddings count: %d, Texts count: %d", len(embeddings), len(texts))
	}

	dataPayload := make([]map[string]interface{}, len(texts))
	for i, textSnippet := range texts {
		dataPayload[i] = map[string]interface{}{
			"id":     i,
			"vector": embeddings[i],
			"text":   textSnippet,
		}
	}

	payload := map[string]interface{}{
		"collection_name": id,
		"data":            dataPayload,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

type QueryResult struct {
	Data []struct {
		ID       int     `json:"id"`
		Distance float64 `json:"distance"`
		Text     string  `json:"text"`
	} `json:"data"`
}

func QueryEmbeddings(collectionID *string, text *string) ([]string, error) {
	url := "http://localhost:2333/vector/search"
	client := &http.Client{}

	if collectionID == nil || *collectionID == "" {
		return nil, fmt.Errorf("collectionID cannot be empty")
	}
	embeddings, err := EmbeddingTexts(context.TODO(), []string{*text})
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("embedding cannot be empty")
	}

	payload := map[string]interface{}{
		"collection_name": *collectionID,
		"data":            embeddings,
		"output_fields":   []string{"text"},
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var results QueryResult

	if err := json.Unmarshal(bodyBytes, &results); err != nil {
		// Attempt to unmarshal into a simpler structure if the above fails, e.g. just a list of strings
		var simpleResults []string
		if errSimple := json.Unmarshal(bodyBytes, &simpleResults); errSimple == nil {
			return simpleResults, nil
		}
		return nil, fmt.Errorf("failed to unmarshal response body: %w. Body: %s", err, string(bodyBytes))
	}

	var texts []string
	for _, r := range results.Data {
		texts = append(texts, r.Text)
	}

	return texts, nil
}
