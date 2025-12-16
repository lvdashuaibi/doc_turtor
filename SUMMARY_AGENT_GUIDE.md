# SummaryAgent 工作流程指南
## 概述

`SummaryAgent` 是一个基于 Eino 框架的文档改写 Agent，采用**循环迭代**的方式，通过 `SummaryAgent` 和 `ReviewerAgent` 的协作，不断改进文档质量，直到满意或达到最大迭代次数。

## 工作流程

```
用户输入：原始文档
    ↓
第 1 次迭代：
├─→ SummaryAgent：改写文档
│   └─ 输出保存到 session: "document_content"
├─→ ReviewerAgent：评审文档
│   ├─ 获取 "document_content" 进行评审
│   ├─ 如果满意 → 调用 satisfied_and_exit 工具 → 退出循环
│   └─ 如果不满意 → 返回改进建议
    ↓
第 2 次迭代（如果需要）：
├─→ SummaryAgent：根据评审意见改进文档
│   └─ 输出保存到 session: "document_content"
├─→ ReviewerAgent：再次评审改进后的文档
│   ├─ 如果满意 → 调用工具 → 退出循环
│   └─ 如果不满意 → 返回新的改进建议
    ↓
... 重复直到满意或达到 MaxIterations (5 次)
    ↓
返回最终改写后的文档
```

## 核心组件

### 1. SummaryAgent（改写 Agent）

**文件**: `agent/summaryAgent.go`

**职责**:
- 接收原始文档或改进建议
- 根据指定规则改写文档
- 将改写结果保存到 session 的 `"document_content"` 键

**关键配置**:
```go
OutputKey: "document_content"  // 将输出自动保存到 session
```

**改写规则**:
1. 保留完整原文原意的情况下用更易懂的语言进行表述
2. 如果有代码希望能够完善代码示例以及优化代码注释和目的
3. 必要的复杂的地方可以使用图表
4. 输出应该只包含改写后的文档内容，不要包含其他说明

### 2. ReviewerAgent（评审 Agent）

**文件**: `agent/reviewerAgent.go`

**职责**:
- 从 session 中读取 `"document_content"`
- 评估文档质量
- 决定是否满意或提供改进建议

**关键工具**: `satisfied_and_exit`
- 当文档满意时调用此工具
- 自动退出循环，停止迭代

**评审标准**:
- 语言是否更易懂
- 代码示例是否完善
- 复杂概念是否有清晰解释
- 整体结构是否合理

### 3. LoopAgent（循环 Agent）

**配置**:
```go
MaxIterations: 5  // 最多迭代 5 次
```

**工作方式**:
- 按顺序执行 SubAgents（SummaryAgent → ReviewerAgent）
- 当 ReviewerAgent 调用 `satisfied_and_exit` 工具时，循环退出
- 如果达到 MaxIterations 仍未满意，循环自动停止

## Session 数据流

### 数据存储和传递

```
第 1 次迭代:
  SummaryAgent 输出 → session["document_content"] = "改写后的文档 v1"
  ↓
  ReviewerAgent 读取 → adk.GetSessionValue(ctx, "document_content")
  ↓
  评审结果 → 作为下一轮 SummaryAgent 的输入

第 2 次迭代:
  SummaryAgent 输入 = 上一轮评审意见
  SummaryAgent 输出 → session["document_content"] = "改写后的文档 v2"
  ↓
  ReviewerAgent 读取 → adk.GetSessionValue(ctx, "document_content")
  ...
```

### 关键 API

**保存到 Session**:
```go
// 在 SummaryAgent 中自动完成（通过 OutputKey 配置）
OutputKey: "document_content"
```

**从 Session 读取**:
```go
// 在 ReviewerAgent 中手动读取
contentToReview, ok := adk.GetSessionValue(ctx, "document_content")
if ok {
    documentContent := contentToReview.(string)
    // 进行评审...
}
```

## 使用示例

### 基础使用

```go
package main

import (
    "context"
    "log"
    myagent "eino_test/agent"
    prints "eino_test/common/utils"
    "github.com/cloudwego/eino/adk"
    "github.com/cloudwego/eino/schema"
)

func main() {
    ctx := context.Background()
    
    // 创建 SummaryAgent
    summaryAgent := myagent.NewSummaryAgent(ctx)
    
    // 创建 Runner
    runner := adk.NewRunner(ctx, adk.RunnerConfig{
        EnableStreaming: true,
        Agent:           summaryAgent,
    })
    
    // 准备输入文档
    originalDocument := `
    # 原始文档
    这是一个需要改写的文档...
    `
    
    // 执行 Agent
    iter := runner.Query(ctx, originalDocument)
    
    // 处理输出
    for {
        event, ok := iter.Next()
        if !ok {
            break
        }
        if event.Err != nil {
            log.Fatal(event.Err)
        }
        prints.Event(event)
    }
}
```

### 与 Supervisor 集成

```go
package main

import (
    "context"
    "log"
    myagent "eino_test/agent"
    "github.com/cloudwego/eino/adk"
    "github.com/cloudwego/eino/adk/prebuilt/supervisor"
    "github.com/cloudwego/eino/schema"
)

func main() {
    ctx := context.Background()
    
    // 创建主 Agent 和 SummaryAgent
    mainAgent := myagent.NewMainAgent(ctx)
    summaryAgent := myagent.NewSummaryAgent(ctx)
    
    // 创建 Supervisor
    supervisorAgent, err := supervisor.New(ctx, &supervisor.Config{
        Supervisor: mainAgent,
        SubAgents:  []adk.Agent{summaryAgent},
    })
    if err != nil {
        panic(err)
    }
    
    // 执行
    iter := supervisorAgent.Run(ctx, &adk.AgentInput{
        Messages: []adk.Message{
            schema.UserMessage("请改写这个文档..."),
        },
    })
    
    for {
        event, ok := iter.Next()
        if !ok {
            break
        }
        if event.Err != nil {
            log.Fatal(event.Err)
        }
        // 处理事件...
    }
}
```

## 关键特性

### 1. 自动 Session 管理

- **SummaryAgent** 通过 `OutputKey` 自动将输出保存到 session
- **ReviewerAgent** 可以直接读取 session 中的数据
- 数据在整个循环过程中自动传递

### 2. 灵活的退出机制

- **工具调用退出**: ReviewerAgent 调用 `satisfied_and_exit` 工具时立即退出
- **迭代次数限制**: 达到 `MaxIterations` 时自动停止
- **错误处理**: 任何错误都会导致循环停止

### 3. 流式处理

- 支持流式输出，实时查看改写进度
- 支持流式评审，实时获取评审意见

## 常见问题

### Q1: 如何修改最大迭代次数？

在 `agent/summaryAgent.go` 中修改：
```go
MaxIterations: 5  // 改为需要的次数
```

### Q2: 如何自定义改写规则？

在 `agent/summaryAgent.go` 中修改 `Instruction` 字段：
```go
Instruction: `你是一个文档改写agent...
// 修改这里的规则
`
```

### Q3: 如何自定义评审标准？

在 `agent/reviewerAgent.go` 中修改 `Instruction` 字段：
```go
Instruction: `你是一个文档评审agent...
// 修改这里的评审标准
`
```

### Q4: 如何获取最终改写后的文档？

最终改写后的文档保存在 session 的 `"document_content"` 中，可以通过以下方式获取：

```go
// 在循环结束后
finalContent, ok := adk.GetSessionValue(ctx, "document_content")
if ok {
    document := finalContent.(string)
    // 使用最终文档...
}
```

### Q5: ReviewerAgent 如何访问改写后的文档？

ReviewerAgent 的 Instruction 中已经说明了如何访问：
```
1、获取 session 中的 "document_content" 来查看改写后的文档
```

Agent 会自动理解这个指令，并在评审时获取相应的数据。

## 技术细节

### 工具调用机制

```go
// 在 ReviewerAgent 中定义的 satisfied_and_exit 工具
satisfiedAndExitTool, err := utils.InferTool(
    "satisfied_and_exit",
    "当文档改写满意时调用此工具退出循环",
    func(ctx context.Context, req *SatisfiedAndExitInput) (string, error) {
        // 发送 BreakLoopAction 来退出循环
        _ = adk.SendToolGenAction(ctx, "satisfied_and_exit", 
            adk.NewBreakLoopAction("reviewerAgent"))
        return req.Summary, nil
    },
)
```

### ReturnDirectly 配置

```go
ReturnDirectly: map[string]bool{
    "satisfied_and_exit": true,
}
```

这个配置确保当调用 `satisfied_and_exit` 工具时，Agent 直接返回结果，不再继续处理。

## 参考资源

- Eino 框架文档: https://github.com/cloudwego/eino
- 相关示例:
  - `eino_example/adk/intro/workflow/loop/` - Loop Agent 示例
  - `eino_example/adk/human-in-the-loop/3_feedback-loop/` - 反馈循环示例
  - `eino_example/adk/intro/session/` - Session 管理示例

## 总结

SummaryAgent 通过以下方式实现了高效的文档改写流程：

1. **循环迭代**: 多次改写和评审，逐步提高文档质量
2. **Session 管理**: 自动保存和传递数据，简化开发
3. **灵活控制**: 支持工具调用退出和迭代次数限制
4. **清晰职责**: SummaryAgent 负责改写，ReviewerAgent 负责评审

这种设计模式可以应用于其他需要迭代改进的场景，如代码审查、文章编辑、数据清洗等。
