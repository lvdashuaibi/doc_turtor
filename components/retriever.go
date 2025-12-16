package components

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	milvusRetriver "github.com/cloudwego/eino-ext/components/retriever/milvus"
)

func NewRetriever(ctx context.Context, embedder *openai.Embedder) *milvusRetriver.Retriever {
	collection := "test3"
	retriever, err := milvusRetriver.NewRetriever(ctx, &milvusRetriver.RetrieverConfig{
		Client:      MilvusCli,
		Collection:  collection,
		VectorField: "vector",
		OutputFields: []string{
			"id",
			"content",
			"metadata",
		},
		TopK:      5,
		Embedding: embedder,
	})
	if err != nil {
		panic(err)
	}
	return retriever
}
