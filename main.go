package main

import (
	"context"
	myagent "eino_test/agent"
	"eino_test/common/utils"
	"fmt"
	"log"
	"os"

	"eino_test/config"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/supervisor"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()
	config.LoadConfig()

	// 创建 MainAgent 作为 Supervisor
	mainAgent := myagent.NewMainAgent(ctx)

	// 创建 SummaryAgent 作为子 Agent
	summaryAgent := myagent.NewSummaryAgent(ctx)

	// 使用 Supervisor 编排 SummaryAgent
	supervisorAgent, err := supervisor.New(ctx, &supervisor.Config{
		Supervisor: mainAgent,
		SubAgents:  []adk.Agent{summaryAgent},
	})
	if err != nil {
		panic(err)
	}

	// 定义要改写的文档文件路径
	documentPath := "./text.md"

	// 检查文件是否存在
	if _, err := os.Stat(documentPath); err != nil {
		panic(fmt.Sprintf("文档文件不存在: %s", documentPath))
	}

	// 使用 Supervisor 执行文档改写任务
	log.Println("========== 开始执行文档改写任务 ==========")
	log.Println("文档路径:", documentPath)
	log.Println()

	// 创建回调处理器来打印 LLM 和 Agent 的输出
	callbackHandler := utils.NewOutputCallbackHandler()
	// 设置全局回调处理器
	callbacks.InitCallbackHandlers([]callbacks.Handler{callbackHandler})

	// 创建 prompt 模板
	// 注意：详细的改写原则已经在 SummaryAgent 的 Instruction 中定义
	// 这里只需要简单的用户请求即可，传递文件路径而不是文件内容
	promptTemplate := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个文档改写系统的协调者。你的任务是：\n"+
			"1. 接收用户的文档改写请求和文件路径\n"+
			"2. 将任务转交给 SummaryAgent 进行改写\n"+
			"3. SummaryAgent 会使用 read_document 工具读取文件，然后根据用户背景信息和改写原则进行改写\n"+
			"4. 改写完成后，文档会被自动保存到 markdown 文件中"),
		schema.UserMessage("请你帮我改写这份技术文档，文档路径为：{filepath}"),
	)

	// 使用模板格式化消息
	formattedMessages, err := promptTemplate.Format(ctx, map[string]any{
		"filepath": documentPath,
	})
	if err != nil {
		log.Fatalf("格式化 prompt 模板失败: %v", err)
	}

	iter := supervisorAgent.Run(ctx, &adk.AgentInput{
		Messages: formattedMessages,
	})

	// 处理执行过程中的事件
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatal(event.Err)
		}
		utils.Event(event)
	}

	log.Println()
	log.Println("========== 文档改写任务完成 ==========")
}
