package sdk

import (
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
