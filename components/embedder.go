package components

import (
	"context"
	"os"
	"time"

	openaiEmbedder "github.com/cloudwego/eino-ext/components/embedding/openai"
)

func NewQwenEmbedder(ctx context.Context) *openaiEmbedder.Embedder {
	timeout := 30 * time.Second
	var dimValue = 512
	var dimensions = &dimValue
	embedder, err := openaiEmbedder.NewEmbedder(ctx, &openaiEmbedder.EmbeddingConfig{
		Timeout:    timeout,
		APIKey:     os.Getenv("DASHSCOPE_API_KEY"),
		BaseURL:    os.Getenv("DASHSCOPE_BASE_URL"),
		Model:      "text-embedding-v3",
		Dimensions: dimensions,
	})
	if err != nil {
		panic(err)
	}
	return embedder
}
