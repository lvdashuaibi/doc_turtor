# 文档改写系统分析报告
## 问题分析

### 1. SupervisorAgent 是否调用了 SummaryAgent？

**答案：否，没有调用 SummaryAgent**

从日志中可以看到：
```
path: [{MainAgent}]
```

这表示只有 `MainAgent` 被调用，`SummaryAgent` 没有被使用。

**原因分析：**

Supervisor 的工作机制是：
- Supervisor 会根据 MainAgent 的输出决定是否需要调用子 Agent
- MainAgent 的 Instruction 中明确指示它要使用 `save_document` 工具保存文档
- MainAgent 直接调用了 `save_document` 工具，完成了任务
- 因此 Supervisor 认为任务已完成，不需要调用 SummaryAgent

**为什么 SummaryAgent 没有被调用：**

1. **Supervisor 的设计目的**：Supervisor 是用来在多个子 Agent 之间进行路由和选择的
2. **当前配置问题**：MainAgent 被配置为可以直接处理文档改写和保存任务
3. **缺少路由逻辑**：MainAgent 没有被指示何时应该将任务转交给 SummaryAgent

### 2. 为什么保存出来的文档不是 Markdown 格式？

**答案：文件确实是 Markdown 格式，但内容格式不够规范**

**观察：**
- 文件名：`改写文档_技术文档.md` ✓ (正确的 .md 扩展名)
- 文件大小：25K ✓ (内容完整)
- 文件内容：纯文本格式 ✓ (可以被 Markdown 解析器识别)

**问题所在：**
- 文件内容缺少 Markdown 格式化元素（如 `#` 标题、`**` 加粗等）
- LLM 生成的内容是纯文本，没有添加 Markdown 语法
- 这是因为 LLM 的改写指令中没有明确要求使用 Markdown 格式

## 解决方案

### 方案 1：修改 MainAgent 的 Instruction 来调用 SummaryAgent

修改 `agent/workflow.go` 中的 MainAgent Instruction：

```go
Instruction: `你是一个文档改写系统的协调者。你的任务是：
1. 接收用户的文档改写请求
2. 将任务转交给 SummaryAgent 进行改写
3. 获取改写结果后，使用 save_document 工具保存到文件

重要提示：
- 必须使用 transfer 操作将任务转交给 SummaryAgent
- 不要自己进行改写，而是让 SummaryAgent 来处理
- 改写完成后，使用 save_document 工具保存文档`,
```

### 方案 2：改进 LLM 的 Instruction 以生成规范的 Markdown 格式

修改 `main.go` 中的 prompt 模板：

```go
schema.SystemMessage(`你是一个专业的技术文档改写专家。你的任务是：
1. 接收用户提供的技术文档
2. 以更易懂、更专业的方式进行重写
3. 使用 save_document 工具将改写后的完整文档保存到 markdown 文件中

重要提示：
- 改写后必须使用 save_document 工具保存文档，不要直接输出
- 确保保存的文件名清晰易识别（例如：改写文档_技术文档.md）
- 保存成功后告知用户文件位置
- 改写后的文档必须使用规范的 Markdown 格式：
  * 使用 # 作为一级标题
  * 使用 ## 作为二级标题
  * 使用 - 或 * 作为列表项
  * 使用 **文本** 进行加粗
  * 使用 \`代码\` 进行代码高亮`),
```

## 当前系统架构

```
用户输入
    ↓
MainAgent (Supervisor)
    ├─ 直接处理文档改写
    ├─ 调用 save_document 工具
    └─ 返回结果
    
SummaryAgent (未被使用)
    ├─ LoopAgent (包含改写和评审循环)
    ├─ SummaryAgent (改写)
    └─ ReviewerAgent (评审)
```

## 建议

1. **如果要使用 SummaryAgent**：
   - 修改 MainAgent 的 Instruction，让它转交任务给 SummaryAgent
   - 使用 `adk.NewTransferToAgentAction()` 进行 Agent 转移

2. **如果要改进 Markdown 格式**：
   - 在 LLM Instruction 中明确指定 Markdown 格式要求
   - 或者在保存前对内容进行后处理，添加 Markdown 格式化

3. **当前最简单的改进**：
   - 只需修改 prompt 模板中的 SystemMessage
   - 添加 Markdown 格式化的具体要求
   - LLM 会自动生成规范的 Markdown 格式内容
