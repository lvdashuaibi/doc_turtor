package components

import (
	"context"
	"log"

	openaiEmbedder "github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func newQwenIndexer(ctx context.Context, embedder *openaiEmbedder.Embedder) *milvus.Indexer {
	var collection = "test3"
	var fields = []*entity.Field{
		{
			Name:     "id",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "256",
			},
			PrimaryKey: true,
		},
		{
			Name:     "vector", // 确保字段名匹配 - BinaryVector
			DataType: entity.FieldTypeBinaryVector,
			TypeParams: map[string]string{
				"dim": "16384",
			},
		},
		{
			Name:     "content",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "8192",
			},
		},
		{
			Name:     "metadata",
			DataType: entity.FieldTypeJSON,
		},
	}

	indexer, err := milvus.NewIndexer(ctx, &milvus.IndexerConfig{
		Client:     MilvusCli,
		Collection: collection,
		Fields:     fields,
		Embedding:  embedder,
	})
	if err != nil {
		log.Fatalf("Fail to create indexer :%v", err)
	}

	return indexer

}
