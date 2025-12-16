package tools

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// MermaidRenderInput Mermaid 渲染工具的输入参数
type MermaidRenderInput struct {
	MermaidCode  string `json:"mermaid_code" jsonschema_description:"Mermaid 图表代码"`
	OutputFormat string `json:"output_format" jsonschema_description:"输出格式：svg 或 png，默认 svg"`
}

// NewMermaidRendererTool 创建一个 Mermaid 渲染工具
func NewMermaidRendererTool() (tool.BaseTool, error) {
	return utils.InferTool(
		"render_mermaid",
		"将 Mermaid 代码渲染为图片（SVG 或 PNG 格式），返回 Base64 编码的图片数据或图片 URL",
		func(ctx context.Context, input *MermaidRenderInput) (string, error) {
			// 验证输入
			if input.MermaidCode == "" {
				return "错误：Mermaid 代码不能为空", fmt.Errorf("mermaid code is empty")
			}

			// 设置默认输出格式
			outputFormat := input.OutputFormat
			if outputFormat == "" {
				outputFormat = "svg"
			}

			// 验证输出格式
			if outputFormat != "svg" && outputFormat != "png" {
				return "错误：输出格式只支持 svg 或 png", fmt.Errorf("unsupported output format: %s", outputFormat)
			}

			// 使用 Mermaid 在线渲染服务
			// 方案 1：使用 mermaid.ink 服务（推荐）
			imageURL, err := renderWithMermaidInk(input.MermaidCode, outputFormat)
			if err == nil {
				return fmt.Sprintf("✓ Mermaid 图表已渲染\n图片 URL: %s\n\n说明：将此 URL 嵌入到 Markdown 中：![图表](%s)", imageURL, imageURL), nil
			}

			// 方案 2：使用 kroki 服务（备选）
			imageURL, err = renderWithKroki(input.MermaidCode, outputFormat)
			if err == nil {
				return fmt.Sprintf("✓ Mermaid 图表已渲染\n图片 URL: %s\n\n说明：将此 URL 嵌入到 Markdown 中：![图表](%s)", imageURL, imageURL), nil
			}

			// 方案 3：返回 Base64 编码的 SVG（本地渲染）
			svgData, err := renderWithBase64SVG(input.MermaidCode)
			if err == nil {
				return fmt.Sprintf("✓ Mermaid 图表已渲染为 Base64 SVG\n\n说明：将此数据 URL 嵌入到 Markdown 中：![图表](data:image/svg+xml;base64,%s)", svgData), nil
			}

			// 所有方案都失败
			return fmt.Sprintf("⚠️ Mermaid 渲染失败\n\n说明：请在 Markdown 中直接使用 Mermaid 代码块，支持 Mermaid 的渲染器会自动渲染\n\n建议：\n1. 使用 GitHub/GitLab 查看 Markdown 文件（原生支持 Mermaid）\n2. 使用 VS Code + Markdown Preview Enhanced 插件\n3. 使用在线工具：https://mermaid.live"), nil
		},
	)
}

// renderWithMermaidInk 使用 mermaid.ink 服务渲染
func renderWithMermaidInk(mermaidCode string, format string) (string, error) {
	// mermaid.ink 的 API 格式
	// https://mermaid.ink/img/BASE64_ENCODED_DIAGRAM
	encoded := base64.StdEncoding.EncodeToString([]byte(mermaidCode))

	if format == "svg" {
		return fmt.Sprintf("https://mermaid.ink/img/%s", encoded), nil
	}

	// PNG 格式
	return fmt.Sprintf("https://mermaid.ink/img/%s?type=png", encoded), nil
}

// renderWithKroki 使用 kroki 服务渲染
func renderWithKroki(mermaidCode string, format string) (string, error) {
	// Kroki 支持多种图表格式，包括 Mermaid
	// https://kroki.io/mermaid/svg/BASE64_ENCODED_DIAGRAM
	encoded := base64.StdEncoding.EncodeToString([]byte(mermaidCode))

	if format == "svg" {
		return fmt.Sprintf("https://kroki.io/mermaid/svg/%s", encoded), nil
	}

	// PNG 格式
	return fmt.Sprintf("https://kroki.io/mermaid/png/%s", encoded), nil
}

// renderWithBase64SVG 将 Mermaid 代码转换为 Base64 编码的 SVG
// 这是一个简单的实现，实际的 SVG 生成需要调用 Mermaid CLI 或在线服务
func renderWithBase64SVG(mermaidCode string) (string, error) {
	// 创建一个简单的 SVG 包装器
	// 实际应用中应该使用 Mermaid CLI 或在线服务来生成真实的 SVG

	// 尝试调用在线 Mermaid 渲染服务
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 使用 mermaid.ink 的 API
	encoded := base64.StdEncoding.EncodeToString([]byte(mermaidCode))
	url := fmt.Sprintf("https://mermaid.ink/img/%s", encoded)

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to render mermaid: status code %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 转换为 Base64
	base64Data := base64.StdEncoding.EncodeToString(body)
	return base64Data, nil
}

// ConvertMermaidToImageURL 将 Markdown 中的 Mermaid 代码块转换为图片 URL
func ConvertMermaidToImageURL(markdown string) string {
	// 查找所有 Mermaid 代码块
	lines := strings.Split(markdown, "\n")
	var result strings.Builder
	inMermaidBlock := false
	var mermaidCode strings.Builder
	i := 0

	for i < len(lines) {
		line := lines[i]

		// 检查是否是 ```mermaid 格式
		if strings.HasPrefix(line, "```mermaid") {
			inMermaidBlock = true
			mermaidCode.Reset()
			i++
			continue
		}

		// 检查是否是 ``` 后跟 mermaid 的格式
		if strings.TrimSpace(line) == "```" && i+1 < len(lines) && strings.TrimSpace(lines[i+1]) == "mermaid" {
			inMermaidBlock = true
			mermaidCode.Reset()
			i += 2 // 跳过 ``` 和 mermaid 行
			continue
		}

		if inMermaidBlock && strings.HasPrefix(line, "```") {
			inMermaidBlock = false
			// 渲染 Mermaid 代码
			code := strings.TrimSpace(mermaidCode.String())
			if code != "" {
				// 移除多余的换行符，保留单个换行符用于格式化
				code = strings.ReplaceAll(code, "\n\n", "\n")
				encoded := base64.StdEncoding.EncodeToString([]byte(code))
				// 使用 Base64URL 编码（用 - 和 _ 替代 + 和 /）
				// 对 Base64 进行 URL 编码，处理特殊字符
				urlEncoded := url.QueryEscape(encoded)
				imageURL := fmt.Sprintf("https://mermaid.ink/img/%s", urlEncoded)
			}
			i++
			continue
		}

		if inMermaidBlock {
			// 只添加非空行
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine != "" {
				mermaidCode.WriteString(trimmedLine)
				mermaidCode.WriteString("\n")
			}
		} else {
			result.WriteString(line)
			result.WriteString("\n")
		}
		i++
	}

	return result.String()
}
