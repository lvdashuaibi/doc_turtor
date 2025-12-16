package components

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
)

// LineSplitter 按行数分割文档的自定义分割器
type LineSplitter struct {
	linesPerChunk int
}

// NewLineSplitter 创建一个新的行分割器
// linesPerChunk: 每个分割块包含的行数
func NewLineSplitter(linesPerChunk int) *LineSplitter {
	if linesPerChunk <= 0 {
		linesPerChunk = 2
	}
	return &LineSplitter{
		linesPerChunk: linesPerChunk,
	}
}

// Transform 实现 document.Transformer 接口
func (ls *LineSplitter) Transform(ctx context.Context, docs []*schema.Document, opts ...document.TransformerOption) ([]*schema.Document, error) {
	var result []*schema.Document
	chunkIndex := 1

	for _, doc := range docs {
		// 按换行符分割内容
		lines := strings.Split(doc.Content, "\n")

		// 按指定行数分割，但不在空行处分割
		i := 0
		for i < len(lines) {
			end := i + ls.linesPerChunk
			if end > len(lines) {
				end = len(lines)
			}

			// 如果结束位置不是文档末尾，且下一行是空行，则向前移动分割点
			// 这样可以避免在空行处分割，保持相关内容在一起
			if end < len(lines) && strings.TrimSpace(lines[end]) == "" {
				// 向前查找最后一个非空行
				for end > i && strings.TrimSpace(lines[end-1]) == "" {
					end--
				}
			}

			// 获取这一块的行
			chunk := strings.Join(lines[i:end], "\n")

			// 跳过空块
			if strings.TrimSpace(chunk) == "" {
				i = end + 1
				continue
			}

			// 创建新的文档块
			newDoc := &schema.Document{
				ID:       fmt.Sprintf("%d", chunkIndex),
				Content:  chunk,
				MetaData: doc.MetaData,
			}

			result = append(result, newDoc)
			chunkIndex++

			// 移动到下一个分割点，跳过空行
			i = end
			for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
				i++
			}
		}
	}

	return result, nil
}
