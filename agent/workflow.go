package agent

import (
	"context"
	"log"

	"eino_test/common/constant"
	"eino_test/components/models"
	"eino_test/tools"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

func NewMainAgent(ctx context.Context) adk.Agent {
	// 创建保存文档工具
	saveDocumentTool, err := tools.NewSaveDocumentTool()
	if err != nil {
		log.Fatalf("创建保存文档工具失败: %v", err)
	}

	a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "MainAgent",
		Description: "一个负责与用户进行交互的agent，协调文档改写任务",
		Instruction: `你是一个文档改写系统的协调者。你的任务是：
1. 接收用户的文档改写请求
2. 将任务转交给 SummaryAgent 进行改写（SummaryAgent 会进行多轮改写和评审）
3. 获取改写结果后，必须使用 save_document 工具将改写后的完整文档保存到 markdown 文件中
4. 告知用户文档已保存的位置

重要提示：
- 必须将文档改写任务转交给 SummaryAgent，不要自己进行改写
- SummaryAgent 会自动进行改写和评审的循环，直到文档满意为止
- 改写完成后，必须调用 save_document 工具保存文档，不要假设文档已经被保存
- 确保保存的文件名清晰易识别（例如：改写文档_技术文档.md）
- 从改写结果中提取完整的改写后文档内容，然后通过 save_document 工具保存`,
		Model: models.NewQwenModel(ctx, constant.QWEN_TURBO),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{saveDocumentTool},
			},
			ReturnDirectly: map[string]bool{
				"save_document": true,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	return a
}
