package models

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
)

func NewQwenModel(ctx context.Context, modelName string) *openai.ChatModel {
	model, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  os.Getenv("DASHSCOPE_API_KEY"),
		BaseURL: os.Getenv("DASHSCOPE_BASE_URL"),
		Model:   modelName,
	})
	if err != nil {
		panic(err)
	}
	return model
}
