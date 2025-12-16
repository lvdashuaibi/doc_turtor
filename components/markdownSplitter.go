package components

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/components/document"
)

// NewTrans 使用 Markdown 标题分割器
func NewTrans(ctx context.Context) document.Transformer {
	splitter, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"#":      "h1",
			"##":     "h2",
			"###":    "h3",
			"####":   "h4",
			"#####":  "h5",
			"######": "h6",
		},
		TrimHeaders: false,
		IDGenerator: func(ctx context.Context, originalID string, splitIndex int) string {
			return fmt.Sprintf("%d", splitIndex+1)
		},
	})
	if err != nil {
		panic(err)
	}
	return splitter
}

// NewLineTransformer 使用行分割器，按指定行数分割文档
// linesPerChunk: 每个分割块包含的行数（默认为2）
func NewLineTransformer(linesPerChunk int) document.Transformer {
	return NewLineSplitter(linesPerChunk)
}

// NewSmartLineTransformer 使用智能行分割器，能够识别 Markdown 结构
// linesPerChunk: 每个分割块包含的行数（默认为6）
func NewSmartLineTransformer(linesPerChunk int) document.Transformer {
	return NewSmartSplitter(linesPerChunk)
}
