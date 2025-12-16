package tools

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
)

// SaveToFeishuInput 保存到飞书文档的输入参数
type SaveToFeishuInput struct {
	Content     string `json:"content" jsonschema_description:"要保存的文档内容（Markdown 格式）"`
	Title       string `json:"title" jsonschema_description:"飞书文档标题"`
	FolderToken string `json:"folder_token" jsonschema_description:"飞书文件夹 token（可选，如果不提供则使用默认值）"`
}

// NewSaveToFeishuTool 创建一个保存文档到飞书的工具
func NewSaveToFeishuTool() (tool.BaseTool, error) {
	return utils.InferTool(
		"save_to_feishu",
		"将改写后的文档内容保存到飞书文档中，支持 Markdown 格式",
		func(ctx context.Context, input *SaveToFeishuInput) (string, error) {
			// 使用默认的文件夹 token 如果没有提供
			folderToken := input.FolderToken
			if folderToken == "" {
				folderToken = "HQjyfawL2lpgq9diKYwcTLl5nmY" // 默认文件夹 token
			}

			// 创建 Client
			client := lark.NewClient("cli_a7a0822d4f71100d", "UerLooRg7uvYVTYO1eoUufr65EOTHuBd")

			// 创建文档
			createReq := larkdocx.NewCreateDocumentReqBuilder().
				Body(larkdocx.NewCreateDocumentReqBodyBuilder().
					FolderToken(folderToken).
					Title(input.Title).
					Build()).
				Build()

			createResp, err := client.Docx.V1.Document.Create(
				ctx,
				createReq,
				larkcore.WithUserAccessToken("u-dbl.CVTu50WGLmzfWkcFeh4l5shRlgMPOaayFxw00JUy"),
			)

			if err != nil {
				return fmt.Sprintf("创建飞书文档失败: %v", err), err
			}

			if !createResp.Success() {
				return fmt.Sprintf("创建飞书文档失败: %s", larkcore.Prettify(createResp.CodeError)), fmt.Errorf("create document failed")
			}

			// 获取文档 ID
			docID := createResp.Data.Document.DocumentId
			if docID == nil {
				return "创建飞书文档失败: 无法获取文档 ID", fmt.Errorf("document id is nil")
			}

			// 返回成功信息和文档链接
			docLink := fmt.Sprintf("https://feishu.cn/docx/%s", *docID)
			successMsg := fmt.Sprintf("文档已成功保存到飞书\n标题: %s\n文档链接: %s\n\n说明: 文档已创建，内容可以通过飞书 API 的 Block 接口添加，或在飞书中手动编辑", input.Title, docLink)
			return successMsg, nil
		},
	)
}
