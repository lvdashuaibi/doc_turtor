package components

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
)

// SmartSplitter 智能分割器，能够识别 Markdown 结构
type SmartSplitter struct {
	linesPerChunk int
}

// NewSmartSplitter 创建一个新的智能分割器
func NewSmartSplitter(linesPerChunk int) *SmartSplitter {
	if linesPerChunk <= 0 {
		linesPerChunk = 6
	}
	return &SmartSplitter{
		linesPerChunk: linesPerChunk,
	}
}

// Transform 实现 document.Transformer 接口
// 智能分割策略：
// 1. 按指定行数分割
// 2. 不在列表项中间分割（保持列表项和其子项在一起）
// 3. 跳过空行
func (ss *SmartSplitter) Transform(ctx context.Context, docs []*schema.Document, opts ...document.TransformerOption) ([]*schema.Document, error) {
	var result []*schema.Document
	chunkIndex := 1

	for _, doc := range docs {
		lines := strings.Split(doc.Content, "\n")

		i := 0
		for i < len(lines) {
			// 计算这一块应该包含的行数
			end := i + ss.linesPerChunk
			if end > len(lines) {
				end = len(lines)
			}

			// 如果分割点在列表项中间，向后扩展以包含完整的列表项
			end = ss.extendToCompleteListItem(lines, i, end)

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

// extendToCompleteListItem 扩展分割点以包含完整的列表项
// 如果分割点在列表项中间，向后扩展直到找到下一个同级或更高级的列表项
func (ss *SmartSplitter) extendToCompleteListItem(lines []string, start, end int) int {
	if end >= len(lines) {
		return end
	}

	// 获取分割点前一行的缩进级别
	if end > 0 {
		prevLine := lines[end-1]
		prevIndent := getIndentLevel(prevLine)

		// 检查分割点是否在列表项中间
		// 如果下一行的缩进更深，说明它是当前项的子项，应该包含它
		for end < len(lines) {
			currentLine := lines[end]
			currentIndent := getIndentLevel(currentLine)

			// 如果当前行是空行，跳过
			if strings.TrimSpace(currentLine) == "" {
				end++
				continue
			}

			// 如果当前行的缩进更深（是子项），包含它
			if currentIndent > prevIndent {
				end++
				prevIndent = currentIndent
				continue
			}

			// 如果缩进相同或更浅，停止扩展
			break
		}
	}

	return end
}

// getIndentLevel 获取一行的缩进级别
func getIndentLevel(line string) int {
	count := 0
	for _, ch := range line {
		if ch == ' ' {
			count++
		} else if ch == '\t' {
			count += 4 // 将制表符视为 4 个空格
		} else {
			break
		}
	}
	return count
}
