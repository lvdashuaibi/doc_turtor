package tools

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// SaveDocumentInput 保存文档的输入参数
type SaveDocumentInput struct {
	Content  string `json:"content" jsonschema_description:"要保存的文档内容"`
	Filename string `json:"filename" jsonschema_description:"保存的文件名（不包含路径，默认保存到当前目录）"`
}

// ReadDocumentInput 读取文档的输入参数
type ReadDocumentInput struct {
	Filepath string `json:"filepath" jsonschema_description:"要读取的 markdown 文件路径（相对路径或绝对路径）"`
}

// NewSaveDocumentTool 创建一个保存文档到 markdown 文件的工具
func NewSaveDocumentTool() (tool.BaseTool, error) {
	return utils.InferTool(
		"save_document",
		"将改写后的文档内容保存到 markdown 文件中，自动将 Mermaid 代码块转换为可渲染的图片 URL",
		func(ctx context.Context, input *SaveDocumentInput) (string, error) {
			// 如果没有指定文件名，使用默认名称
			if input.Filename == "" {
				input.Filename = fmt.Sprintf("改写文档_%s.md", time.Now().Format("20060102_150405"))
			}

			// 确保文件名以 .md 结尾
			if len(input.Filename) < 3 || input.Filename[len(input.Filename)-3:] != ".md" {
				input.Filename += ".md"
			}

			// 将 Mermaid 代码块转换为图片 URL
			processedContent := ConvertMermaidToImageURL(input.Content)

			// 写入文件
			err := os.WriteFile(input.Filename, []byte(processedContent), 0644)
			if err != nil {
				return fmt.Sprintf("保存文档失败: %v", err), err
			}

			// 获取文件的绝对路径
			absPath, err := os.Getwd()
			if err != nil {
				absPath = "当前目录"
			} else {
				absPath = fmt.Sprintf("%s/%s", absPath, input.Filename)
			}

			return fmt.Sprintf("✓ 文档已成功保存到: %s\n\n说明: Mermaid 图表已自动转换为可渲染的图片 URL", absPath), nil
		},
	)
}

// NewReadDocumentTool 创建一个读取 markdown 文件的工具
func NewReadDocumentTool() (tool.BaseTool, error) {
	return utils.InferTool(
		"read_document",
		"从指定的 markdown 文件中读取文档内容",
		func(ctx context.Context, input *ReadDocumentInput) (string, error) {
			// 读取文件
			content, err := os.ReadFile(input.Filepath)
			if err != nil {
				return fmt.Sprintf("读取文档失败: %v", err), err
			}

			// 获取文件的绝对路径
			absPath, err := os.Getwd()
			if err != nil {
				absPath = "当前目录"
			} else {
				absPath = fmt.Sprintf("%s/%s", absPath, input.Filepath)
			}

			return fmt.Sprintf("已成功读取文档 %s，内容如下：\n\n%s", absPath, string(content)), nil
		},
	)
}
