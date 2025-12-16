# 使用 Eino ADK 构建您的第一个 AI 智能体：从 Excel Agent 实战开始

本文旨在指导您如何利用 Eino ADK（Agent Development Kit）构建一个功能强大的多智能体系统。

## 什么是 Excel Agent？

Excel Agent 是一种能够“理解用户的问题、读懂表格内容并编写和执行代码”的智能助手。它将复杂的 Excel 处理任务分解为清晰的步骤，通过自动规划、工具调用与结果验证，稳定地完成各种 Excel 数据处理任务。

Excel Agent 是基于 Eino ADK 构建的 Multi-Agent 系统，其架构如下图所示：

### Excel Agent 内部包含的几个 Agent 功能分别为：

- **Planner**：分析用户的输入，将问题拆解为可执行的计划。
- **Executor**：正确执行当前计划中的第一步。
- **CodeAgent**：接收 Executor 的指令，调用多种工具（如读写文件、运行 Python 代码等）完成任务。
- **WebSearchAgent**：接收 Executor 的指令，进行网络搜索。
- **Replanner**：根据 Executor 的执行结果和现有计划，决定继续执行、调整计划或结束执行。
- **ReportAgent**：根据运行过程和结果生成总结性质的报告。

## Excel Agent 的典型使用场景

在实际业务中，您可以将 Excel Agent 视为一位“Excel 专家 + 自动化工程师”。当您提供一个原始表格和目标描述时，它会提出可行的解决方案并完成执行：

- **数据清理与格式化**：从一个包含大量数据的 Excel 文件中完成去重、空值处理、日期格式标准化操作。
- **数据分析与报告生成**：从销售数据中提取每月的销售总额，聚合统计、透视，并最终生成并导出图表报告。
- **自动化预算计算**：根据不同部门的预算申请，自动计算总预算并生成部门预算分配表。
- **数据匹配与合并**：将多个不同来源的客户信息表进行匹配合并，生成完整的客户信息数据库。

## Excel Agent 的完整运行动线

[图片]

## 核心收益

- **减少人工操作**：将复杂繁琐的 Excel 处理工作交给 Agent 自动完成。
- **提高产出质量**：通过“规划—执行—反思”闭环减少漏项与错误。
- **增强可扩展性**：各 Agent 独立构建，低耦合利于迭代更新。

Excel Agent 既可以单独使用，也可以作为子 Agent，集成在一个复合的多专家系统中，由外部路由到此 Agent 上，解决 Excel 领域相关的问题。

## ChatModelAgent：与 LLM 交互的基石

ChatModelAgent 是 Eino ADK 中的一个核心预构建 Agent，内部使用了 ReAct 模式（一种让模型“思考-行动-观察”的链式推理模式）。

ChatModelAgent 的主要目的是让 ChatModel 进行显式的、一步一步的“思考”，结合思考过程驱动行动，观测历史思考过程与行动结果继续进行下一步的思考与行动，最终解决复杂问题。

### ReAct 模式的工作流程如下：

1. 调用 ChatModel（Reason）
2. LLM 返回工具调用请求（Action）
3. ChatModelAgent 执行工具（Act）
4. 将工具结果返回给 LLM（Observation），结合之前的上下文继续生成，直到模型判断不需要调用工具后结束

在 Excel Agent 中，每个 Agent 的核心都是这样一个 ChatModelAgent。以 Executor 运行“读取用户输入表格的头信息”这个步骤为例，我们可以通过观察完整的运行过程来理解 ReAct 模式在 ChatModelAgent 中的表现：

- Executor：经过判断，将任务转交给 CodeAgent 运行
- CodeAgent：接收到任务“读取用户输入表格的头信息”
- Think-1：上下文未提供工作目录下的所有文件，需要查看
- Act-1: 调用 Bash 工具，ls 查看工作目录下的所有文件
- Think-2: 找到了用户输入的文件，判断需要编写 Python 代码读取 xlsx 表格的首行
- Act-2: 调用 PythonRunner 工具，书写代码并运行，获取运行结果
- Think-3: 获取到了 xlsx 首行，判断任务完成

运行完成，将表格头信息返回给 Executor

## Plan-Execute Agent：基于「规划-执行-反思」的多智能体协作框架

Plan-Execute Agent 是 Eino ADK 中一种基于「规划-执行-反思」范式的多智能体协作框架，旨在解决复杂任务的分步拆解、执行与动态调整问题。它通过 Planner（规划器）、Executor（执行器）和 Replanner（重规划器）三个核心智能体的协同工作，实现任务的结构化规划、工具调用执行、进度评估与动态 replanning，最终达成用户目标。

### 示例代码：

```go
// NewPlanner creates a new planner agent based on the provided configuration.
func NewPlanner(_ context.Context, cfg *PlannerConfig) (adk.Agent, error)

// NewExecutor creates a new executor agent.
func NewExecutor(ctx context.Context, cfg *ExecutorConfig) (adk.Agent, error)

// NewReplanner creates a new replanner agent.
func NewReplanner(_ context.Context, cfg *ReplannerConfig) (adk.Agent, error)

// New creates a new plan-execute-replan agent with the given configuration.
func New(ctx context.Context, cfg *Config) (adk.Agent, error)
```

而 Excel Agent 的核心能力恰好为【解决用户在 Excel 领域的问题】，与该智能体协作框架定位一致：

- **规划者（Planner）**：明确目标，自动拆解可执行步骤
- **执行者（Executor）**：调用工具（Excel 读取、系统命令、Python 代码）完成规划中的每一个详细步骤
- **反思者（Replanner）**：根据执行进度决定继续、调整规划或结束

Planner 和 Replanner 会将用户模糊的指令拆解为清晰的、可执行的步骤清单，即包含多个步骤（Step）的计划（Plan），Eino ADK 为此提供了灵活的 Plan 接口定义，支持用户自定义 Plan 结构与细节：

```go
type Plan interface {
    // FirstStep returns the first step to be executed in the plan.
    FirstStep() string
    // Marshaler serializes the Plan into JSON.
    // The resulting JSON can be used in prompt templates.
    json.Marshaler
    // Unmarshaler deserializes JSON content into the Plan.
    // This processes output from structured chat models or tool calls into the Plan structure.
    json.Unmarshaler
}
```

默认情况下，框架会使用内置的 Plan 结构作为兜底配置，例如下面就是 Excel Agent 产生的一个完整运行计划：

### 任务计划
- [x] 1. Read the contents of '模拟出题.csv' from the working directory into a pandas DataFrame.
- [x] 2. Identify the question type (e.g., multiple-choice, short-answer) for each row in the DataFrame.
- [x] 3. For non-short-answer questions, restructure the data to place question, answer, explanation, and options in the same row.
- [x] 4. For short-answer questions, merge the answer content into the explanation column and ensure question and merged explanation are in the same row.
- [x] 5. Verify that all processed rows have question, answer (where applicable), explanation, and options (where applicable) in a single row with consistent formatting.
- [x] 6. Generate a cleaned report presenting the formatted questions with all relevant components (question, answer, explanation, options) in unified rows.

## Workflow Agents：可控的多 Agent 运行流水线

在 Excel Agent 中，存在一些需要按照特定顺序运行 agent 的情况：

- **顺序运行**：先运行 Planner，再运行 Executor 和 Replanner；Planner 只运行一次。
- **循环运行**：Executor 和 Replanner 需要按需循环运行多次，每次循环运行都是先运行 Executor 后运行 Replanner。
- **顺序运行**：Plan-Executor 整体运行完后，固定运行一次 ReportAgent 进行总结。

对于这些拥有固定执行流程的场景，Eino ADK 提供了三种流程编排方式，协助用户快速搭建可控的工作流：

### SequentialAgent：依次执行一系列子 Agent

```go
import github.com/cloudwego/eino/adk

// 依次执行 制定研究计划 -> 搜索资料 -> 撰写报告
sequential := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
    Name: "research_pipeline",
    SubAgents: []adk.Agent{
        planAgent,    // 制定研究计划
        searchAgent,  // 搜索资料
        writeAgent,   // 撰写报告
    },
})
```

### LoopAgent：重复执行配置的子 Agent 序列

```go
import github.com/cloudwego/eino/adk

// 循环执行 5 次，每次顺序为：分析当前状态 -> 提出改进方案 -> 验证改进效果
loop := adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{
    Name: "iterative_optimization",
    SubAgents: []adk.Agent{
        analyzeAgent,  // 分析当前状态
        improveAgent,  // 提出改进方案
        validateAgent, // 验证改进效果
    },
    MaxIterations: 5,
})
```

### ParallelAgent：允许多个子 Agent 并发执行

```go
import github.com/cloudwego/eino/adk

// 并发执行 情感分析 + 关键词提取 + 内容摘要
parallel := adk.NewParallelAgent(ctx, &adk.ParallelAgentConfig{
    Name: "multi_analysis",
    SubAgents: []adk.Agent{
        sentimentAgent,  // 情感分析
        keywordAgent,    // 关键词提取
        summaryAgent,    // 内容摘要
    },
})
```

## Agent 抽象：灵活定义 Agent 的基础

Eino ADK 的核心是一个简洁而强大的 Agent 接口，每个 Agent 都有明确的身份（Name）、清晰的职责（Description）和标准化的执行方式（Run），为 Agent 之间的发现与调用提供了基础。

### 统一的 Agent 抽象

ADK 提供的预构建 Agent（ChatModelAgent，Plan-Execute Agent，Workflow Agents）都遵循该接口定义。您也可以基于该接口，书写自定义 Agent，完成定制化需求。

```go
type Agent interface {
    Name(ctx context.Context) string
    Description(ctx context.Context) string
    Run(ctx context.Context, input *AgentInput, options ...AgentRunOption) *AsyncIterator[*AgentEvent]
}
```

### 标准化输入

Agent 通常以 LLM 为核心，因此 Eino ADK 定义的 Agent 的输入与 LLM 接收的输入一致：

```go
type AgentInput struct {
    Messages        []Message
    EnableStreaming bool
}

type Message = *schema.Message // *schema.Message 是模型输入输出的结构定义
```

### 异步事件驱动输出

Agent 的输出是一个 AgentEvent 的异步迭代器，其中的 AgentEvent 表示 Agent 在其运行过程中产生的核心事件数据。其中包含了 Agent 的元信息、输出、行为和报错信息：

```go
type AgentEvent struct {
    AgentName string    // 产生 Event 的 Agent 名称（框架自动填充）

    RunPath []RunStep   // 到达当前 Agent 的完整运行轨迹（框架自动填充）

    Output *AgentOutput // Agent 输出消息内容

    Action *AgentAction // Agent 动作事件内容

    Err error           // Agent 报错
}

type AgentOutput struct {
    MessageOutput *MessageVariant // 模型消息输出内容

    CustomizedOutput any          // 自定义输出内容
}

type MessageVariant struct {
    IsStreaming bool            // 是否为流式输出
    
    Message       Message        // 非流式消息输出
    MessageStream MessageStream  // 流式消息输出

    Role schema.RoleType         // 消息角色
    ToolName string              // 工具名称
}

type AgentAction struct {
    Exit bool                               // Agent 退出

    Interrupted *InterruptInfo              // Agent 中断

    TransferToAgent *TransferToAgentAction  // Agent 跳转

    CustomizedAction any                    // 自定义 Agent 动作
}
```

异步迭代器允许 Agent 在运行过程中的任意时刻向迭代器发送消息（Agent 调用模型结果、工具运行结果、中间状态等等），同时调用方以一种有序、阻塞的方式消费这一系列事件：

```go
iter := myAgent.Run(ctx, "hello") // get AsyncIterator

for {
    event, ok := iter.Next()
    if !ok {
        break
    }
    // handle event
}
```

## Agent 协作：隐藏在 Agent 后的数据传递

Excel Agent 架构图中的节点代表每个具体的 Agent，边代表了数据流通与任务转移。在构建多 Agent 系统时，让不同 Agent 之间高效、准确地共享信息至关重要。

这些信息不仅包含 Agent 的输入输出，还有全局的、部分可见的种种额外信息，例如：

- Executor 执行需要从 Planner / Replanner 拿到一个结构化的、可被拆分为详细步骤（Step）的计划（Plan），而非一段非结构化的 LLM 原始输出消息。
- ReportAgent 需要拿到完整的运行计划、运行过程与运行产物才能正确产生报告。

Eino ADK 包含两种基础的数据传递机制：

### History

每一个 Agent 产生的 AgentEvent 都会被保存到这个隐藏的 History 中，调用一个新 Agent 时 History 中的 AgentEvent 会被转换并拼接到 AgentInput 中。默认情况下，其他 Agent 的 Assistant 或 Tool Message，被转换为 User Message，这相当于在告诉当前的 LLM：“刚才， Agent_A 调用了 some_tool ，返回了 some_result 。现在，轮到你来决策了。”。 通过这种方式，其他 Agent 的行为被当作了提供给当前 Agent 的“外部信息”或“事实陈述”，而不是它自己的行为，从而避免了 LLM 的上下文混乱。

### 共享 Session

单次运行过程中持续存在的 KV 存储，用于支持跨 Agent 的状态管理和数据共享，一次运行中的任何 Agent 可以在任何时间读写 SessionValues。以 Plan-Execute Agent 模式为例，Planner 生成首个计划并写入 Session；Executor 从 Session 读取计划并执行；Replanner 从 Session 读取当前计划后，结合运行结果，将更新后的计划写回 Session 覆盖当前的计划。

```go
// Agent 内获取全部 SessionValues
func GetSessionValues(ctx context.Context) map[string]any

// Agent 内指定 key 获取 SessionValues 中的值
func GetSessionValue(ctx context.Context, key string) (any, bool)

// Agent 内添加 SessionValues
func AddSessionValue(ctx context.Context, key string, value any)

// Agent 内批量添加 SessionValues
func AddSessionValues(ctx context.Context, kvs map[string]any)

// WithSessionValues 在 Agent 运行前由外部注入 SessionValues
func WithSessionValues(v map[string]any) AgentRunOption
```

除了完善的 Agent 间数据传递机制，Eino ADK 从实践出发，提供了多种 Agent 协作模式：

### 预设 Agent 运行顺序（Workflow）

以代码中预设好的流程运行，Agent 的执行顺序是事先确定、可预测的。对应 Workflow Agents 章节提到的三种范式。

### 移交运行（Transfer）

携带本 Agent 输出结果上下文，将任务移交至子 Agent 继续处理。适用于智能体功能可以清晰的划分边界与层级的场景，常结合 ChatModelAgent 使用，通过 LLM 的生成结果进行动态路由。结构上，以此方式进行协作的两个 Agent 称为父子 Agent：

```go
// 设置父子 Agent 关系
func SetSubAgents(ctx context.Context, agent Agent, subAgents []Agent) (Agent, error)

// 指定目标 Agent 名称，构造 Transfer Event
func NewTransferToAgentAction(destAgentName string) *AgentAction
```

### 显式调用（ToolCall）

将 Agent 视为工具进行调用，适用于 Agent 运行仅需要明确清晰的参数而非完整运行上下文的场景。常结合 ChatModelAgent，将 Agent 作为工具运行后将结果返回给 ChatModel 继续处理。除此之外，ToolCall 同样支持调用符合工具接口构造的、不含 Agent 的普通工具。

```go
// 将 Agent 转换为 Tool
func NewAgentTool(_ context.Context, agent Agent, options ...AgentToolOption) tool.BaseTool
```

## Excel Agent 示例运行

### 配置环境与输入输出路径

- **环境变量**：Excel Agent 运行依赖的完整环境变量可参考项目 README。
- **运行输入**：包括一段用户需求描述和待处理的一系列文件，其中：
  - `main.go` 中首行表示用户输入的需求描述，可自行修改：
    ```go
    func main() {
        // query := schema.UserMessage("统计附件文件中推荐的小说名称及推荐次数，并将结果写到文件中。凡是带有《》内容都是小说名称，形成表格，表头为小说名称和推荐次数，同名小说只列一行，推荐次数相加")
        // query := schema.UserMessage("读取模拟出题.csv 中的内容，规范格式将题目、答案、解析、选项放在同一行，简答题只把答案写入解析即可")
        query := schema.UserMessage("请帮我将 question.csv 表格中的第一列提取到一个新的 csv 中")
    }
    ```
- **adk/multiagent/integration-excel-agent/playground/input** 为默认的附件输入路径，附件输入路径支持配置，参考 README。
- **adk/multiagent/integration-excel-agent/playground/test_data** 路径下提供了几个示例文件，您可以将文件复制到附件输入路径下来进行测试运行。

### 运行输出

Excel Agent 输入的附件、运行的中间产物与最终结果都会放置在工作路径下：`adk/multiagent/integration-excel-agent/playground/${uuid}`，输出路径支持配置，参考 README。

### 查看运行结果

Excel Agent 单次运行会在输出路径下创建一个新的工作目录，并在该目录下完成任务，运行时产生的中间产物与最终结果都会写到该目录下。

以“请帮我将 question.csv 表格中的第一列提取到一个新的 csv 中”这个任务为例，运行完成后在工作目录下包含的内容有：

- **原始输入**：从输入路径获取到的 `question.csv`
- **Planner / Replanner 给出的运行计划**：`plan.md`
- **Executor 中的 CodeAgent 书写的代码**：`$uuid.py`
- **运行中间产物**：`extracted_first_column.csv` 和 `first_column.csv`
- **最终报告**：`final_report.json`

### 运行过程输出

Excel Agent 会将每个步骤的运行结果输出到日志中。下面仍以“请帮我将 question.csv 表格中的第一列提取到一个新的 csv 中”这个任务为例，向您展示 Excel Agent 在运行过程中的几个关键步骤及其输出，并通过对步骤的解释，直观地呈现 Agent 的运行流程及其强大能力。

#### Planner 生成 JSON 格式的初始计划

```json
{
    "steps": [
        {
            "index": 1,
            "desc": "Read the 'questions.csv' file into a pandas DataFrame."
        },
        {
            "index": 2,
            "desc": "Extract the first column from the DataFrame."
        },
        {
            "index": 3,
            "desc": "Save the extracted first column to a new CSV file."
        }
    ]
}
```

#### Executor 将 CodeAgent 作为工具进行调用，执行计划中的首个步骤

```json
{
    "name": "Executor",
    "tool name": "CodeAgent",
    "arguments": {
        "request": "Read the 'questions.csv' file into a pandas DataFrame using pandas. Use the pandas.read_csv function and store the result in a variable named df."
    }
}
```

#### CodeAgent 使用 PythonRunner 工具运行代码，并使用 ReAct 模式自动纠错，修正代码中的错误

```json
{
    "name": "Executor",
    "tool name": "PythonRunner",
    "arguments": {
        "code": "```python\nfirst_column = df.iloc[:, 0]\n```"
    }
}
```

#### PythonRunner 代码运行报错

```json
{
    "name": "Executor",
    "tool response": "Traceback (most recent call last):\n  File \"/User/user/go/src/github.com/cloudwego/eino-examples/adk/multiagent/integration-excel-agent/playground/00f118af-4bd8-42f7-8d11-71f2801218bd/00f118af-4bd8-42f7-8d11-71f2801218bd.py\", line 1, in <module>\n    first_column = df.iloc[:, 0]\n                   ^^\nNameError: name 'df' is not defined"
}
```

#### ReAct 模式自动纠错，修正无法运行的代码

```json
{
    "name": "Executor",
    "answer": "The error occurs because the DataFrame `df` is not defined. We need to first load the data from the existing CSV file `questions.csv` into `df`. Here's the corrected code:",
    "tool name": "PythonRunner",
    "arguments": {
        "code": "```python\nimport pandas as pd\ndf = pd.read_csv('questions.csv')\nfirst_column = df.iloc[:, 0]\nprint(first_column.head())  # Verify the result\n```"
    }
}
```

#### 代码运行成功，返回运行结果

```json
{
    "name": "Executor",
    "path": [{SequentialAgent} {plan_execute_replan} {Planner} {execute_replan} {Executor}],
    "tool response": 
    0    multiple-choice
    1    multiple-choice
    2    multiple-choice
    3    multiple-choice
    4    multiple-choice
    Name: type, dtype: object
}
```

#### Replanner 判断计划完成，提交运行结果至 ReportAgent

```json
{
    "name": "Replanner",
    "answer": {
        "is_success": true,
        "result": "已成功将'questions.csv'表格中的第一列提取到新的CSV文件'extracted_first_column.csv'中。",
        "files": [
            {
                "desc": "包含原表格第一列数据的新CSV文件",
                "path": "extracted_first_column.csv"
            }
        ]
    }
}
```

#### ReportAgent 进行总结，结束执行

```json
{
    "name": "Report",
    "tool name": "SubmitResult",
    "arguments": {
        "is_success": true,
        "result": "Successfully extracted the first column from questions.csv and saved it to first_column.csv.",
        "files": [
            {
                "path": "/User/user/go/src/github.com/cloudwego/eino-examples/adk/multiagent/integration-excel-agent/playground/00f118af-4bd8-42f7-8d11-71f2801218bd/first_column.csv",
                "desc": "A CSV file containing only the first column data from the original questions.csv."
            }
        ]
    }
}
```

## 总结

Excel Agent 所呈现的并非“单一智能体”的技巧，而是一套以 Eino ADK 为底座的 Multi-Agent 系统工程化方法论：

- 以 ChatModelAgent 的 ReAct 能力为基石，让模型“可思考、会调用”。
- 以 WorkflowAgents 的编排能力，让 Multi-Agent 系统中的每个 Agent 以用户预期的顺序运行。
- 以 Planner–Executor–Replanner 的闭环，让复杂任务“可拆解、能纠错”。
- 以 History / Session 的数据传递机制，让多 Agent “能协作、可回放”。

💡立即开始你的智能体开发之旅

⌨️ 查看 Excel Agent 源码：integration-excel-agent [2]
📚 查看更多文档：Eino ADK 文档[3]
🛠️ 浏览 ADK 源码：Eino ADK 源码[4]
💡 探索 ADK 全部示例：Eino ADK Examples[5]
🤝 加入开发者社区：与其他开发者交流经验和最佳实践

Eino ADK，让智能体开发变得简单而强大！

## 参考资料

[1] ReAct: https://react-lm.github.io/

[2] Eino 示例代码: https://github.com/cloudwego/eino-examples/tree/main/adk/multiagent/integration-excel-agent

[3] Eino ADK 文档: https://www.cloudwego.io/zh/docs/eino/core_modules/eino_adk/

[4] Eino ADK 源码: https://github.com/cloudwego/eino/tree/main/adk

[5] Eino ADK Examples: https://github.com/cloudwego/eino-examples/tree/main/adk