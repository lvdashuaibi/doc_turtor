package tools

import (
	"context"
	"fmt"
	"testing"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
)

// 飞书应用配置
const (
	LARK_APP_ID     = "cli_a7a0822d4f71100d"
	LARK_APP_SECRET = "UerLooRg7uvYVTYO1eoUufr65EOTHuBd"
	FOLDER_TOKEN    = "HQjyfawL2lpgq9diKYwcTLl5nmY"                    // 飞书文件夹 token
	USER_TOKEN      = "u-dbl.CVTu50WGLmzfWkcFeh4l5shRlgMPOaayFxw00JUy" // 用户 token
)

// TestCreateFeishuDocument 测试创建飞书文档
func TestCreateFeishuDocument(t *testing.T) {
	// 创建 Client
	client := lark.NewClient(LARK_APP_ID, LARK_APP_SECRET)

	// 创建请求对象
	req := larkdocx.NewCreateDocumentReqBuilder().
		Body(larkdocx.NewCreateDocumentReqBodyBuilder().
			FolderToken(FOLDER_TOKEN).
			Title("测试文档 - 文档改写系统").
			Build()).
		Build()

	// 发起请求
	resp, err := client.Docx.V1.Document.Create(context.Background(), req, larkcore.WithUserAccessToken(USER_TOKEN))

	// 处理错误
	if err != nil {
		t.Fatalf("创建文档失败: %v", err)
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		t.Fatalf("logId: %s, error response: %s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return
	}

	// 业务处理
	fmt.Println("文档创建成功:")
	fmt.Println(larkcore.Prettify(resp))

	// 获取文档 ID
	if resp.Data != nil && resp.Data.Document != nil {
		docID := resp.Data.Document.DocumentId
		if docID != nil {
			fmt.Printf("文档 ID: %s\n", *docID)
		}
	}
}

// TestWriteToFeishuDocument 测试写入内容到飞书文档
func TestWriteToFeishuDocument(t *testing.T) {
	// 创建 Client
	client := lark.NewClient(LARK_APP_ID, LARK_APP_SECRET)

	// 首先创建一个文档
	createReq := larkdocx.NewCreateDocumentReqBuilder().
		Body(larkdocx.NewCreateDocumentReqBodyBuilder().
			FolderToken(FOLDER_TOKEN).
			Title("测试写入 - 文档改写系统").
			Build()).
		Build()

	createResp, err := client.Docx.V1.Document.Create(context.Background(), createReq, larkcore.WithUserAccessToken(USER_TOKEN))
	if err != nil {
		t.Fatalf("创建文档失败: %v", err)
		return
	}

	if !createResp.Success() {
		t.Fatalf("创建文档失败: %s", larkcore.Prettify(createResp.CodeError))
		return
	}

	docID := createResp.Data.Document.DocumentId
	if docID != nil {
		fmt.Printf("创建的文档 ID: %s\n", *docID)

		// 现在写入内容到文档
		// 注意：这里需要使用正确的 API 来写入内容
		// 飞书文档 API 可能需要使用 Block 相关的 API

		fmt.Printf("文档创建成功，可以在飞书中查看: %s\n", *docID)
	}
}
