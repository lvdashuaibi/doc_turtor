# 📚 文档改写 Agent 系统
一个基于 Eino 框架的智能文档改写系统，专门为后端开发初学者将复杂的技术文档改写成易懂的版本。系统采用多 Agent 协作架构，支持严格的质量评审、迭代改进和双重保存（本地文件 + 飞书文档）。

## ✨ 核心特性

- 🤖 **多 Agent 协作**：MainAgent、SummaryAgent、ReviewerAgent 三层架构
- 📝 **智能改写**：基于用户背景信息的个性化改写
- ✅ **严格评审**：10 个维度的质量评审标准
- 💾 **双重保存**：同时保存到本地文件和飞书文档
- 🔄 **迭代改进**：不满意时自动返回改进建议进行迭代
- ⚡ **Token 优化**：增量改进模式减少 token 消耗
- 🎨 **Mermaid 图表渲染**：自动将 Mermaid 代码块转换为可渲染的图片 URL

## 🚀 快速开始

### 前置要求

- Go 1.24.5 或更高版本
- Docker 和 Docker Compose（用于 Milvus）
- OpenAI API Key
- 飞书 API 凭证（可选，用于保存到飞书）

### 环境配置

1. **克隆项目**
```bash
git clone <repository-url>
cd eino_test
```

2. **安装依赖**
```bash
go mod download
go mod tidy
```

3. **配置环境变量**

创建 `.env` 文件或设置以下环境变量：

```bash
# OpenAI 配置
OPENAI_API_KEY=your_openai_api_key
OPENAI_BASE_URL=https://api.openai.com/v1  # 可选，默认为官方 URL

# 飞书配置（可选）
FEISHU_APP_ID=your_feishu_app_id
FEISHU_APP_SECRET=your_feishu_app_secret
FEISHU_FOLDER_TOKEN=your_feishu_folder_token

# Milvus 配置
MILVUS_HOST=localhost
MILVUS_PORT=19530
```

4. **启动 Milvus（可选，用于向量存储）**

```bash
# 使用 Docker Compose 启动 Milvus
docker-compose up -d

# 或使用提供的脚本
./start-milvus.sh
```

5. **运行项目**

```bash
# 编译
go build -o eino_test .

# 运行
./eino_test

# 或直接运行
go run main.go
```

## 📖 使用指南

### 基本使用

1. **准备文档**

在项目根目录创建或放置要改写的 Markdown 文档，例如 `text.md`

2. **运行改写**

```bash
go run main.go
```

系统会自动：
- 读取文档
- 根据用户背景信息进行改写
- 进行严格的质量评审
- 如果不满意，自动迭代改进
- 保存到本地文件和飞书文档

3. **查看结果**

改写后的文档会保存到：
- 本地文件：`改写文档_<原文件名>.md`
- 飞书文档：自动创建在指定文件夹中

### 用户背景信息

系统会根据以下用户背景信息进行个性化改写：

```
- 编程语言：掌握 Java、Golang 的简单后端开发
- 中间件经验：熟悉 MySQL、Kafka、Redis 的基本使用
- 架构知识：了解分布式设计的基本概念
- 学习阶段：正在从零开始学习和梳理后端开发框架
```

可在 `agent/summaryAgent.go` 中修改这些信息。

## 🏗️ 项目结构

```
eino_test/
├── agent/                          # Agent 实现
│   ├── frontAgent.go              # 前端 Agent（用户交互）
│   ├── summaryAgent.go            # 改写 Agent（文档改写）
│   ├── reviewerAgent.go           # 评审 Agent（质量评审）
│   ├── Supervisor.go              # 主 Agent（流程协调）
│   └── workflow.go                # 工作流定义
├── components/                     # 核心组件
│   ├── embedder.go                # 向量嵌入
│   ├── indexer.go                 # 向量索引
│   ├── retriever.go               # 向量检索
│   ├── smartSplitter.go           # 智能文本分割
│   ├── markdownSplitter.go        # Markdown 分割
│   ├── lineSplitter.go            # 行分割
│   ├── MilvusCli.go               # Milvus 客户端
│   └── models/
│       └── chat_model.go          # 聊天模型
├── tools/                          # 工具函数
│   ├── file_writer.go             # 文件保存工具
│   ├── feishu_writer.go           # 飞书保存工具
│   ├── mermaid_renderer.go        # Mermaid 渲染工具
│   ├── mermaid_renderer_test.go   # Mermaid 测试
│   └── tool.go                    # 工具注册
├── common/                         # 通用模块
│   ├── constant/
│   │   └── ModelNames.go          # 模型名称常量
│   └── utils/
│       └── prints.go              # 打印工具
├── config/                         # 配置管理
│   └── config.go                  # 配置加载
├── main.go                         # 程序入口
├── go.mod                          # Go 模块定义
├── docker-compose.yml             # Docker 配置
└── README.md                       # 本文件
```

## 🔄 工作流程

### 系统架构

```
用户输入 (文档 + 背景信息)
    ↓
Supervisor (MainAgent)
    ↓
SummaryAgent (改写)
    ↓
ReviewerAgent (评审)
    ↓
满意? ──否→ 返回改进建议 → SummaryAgent (增量改进)
    ↓是
保存到本地文件 + 飞书文档
    ↓
完成
```

### 评审标准

ReviewerAgent 使用 10 个维度的评审标准：

1. **语言易懂性** - 避免生硬学术用语，使用友好语气
2. **内容详细度** - 保留完整结构，充分解释关键概念
3. **举例和代码示例** - 提供完整、可运行的示例
4. **类比学习** - 使用生活中的类比解释复杂概念
5. **对比学习** - 对相似概念进行明确对比
6. **图表辅助** - 使用 Mermaid 图表辅助说明
7. **深度讲解** - 讲解原理、机制、坑点和最佳实践
8. **章节过渡** - 添加过渡性语句连接章节
9. **教学化结构** - 遵循教学大纲和结构化设计
10. **Markdown 格式** - 正确使用标题、列表、代码块等

## 🛠️ 核心组件说明

### SummaryAgent（改写 Agent）

负责根据用户背景信息和改写原则对文档进行改写。

**主要功能：**
- 读取原始文档
- 分析文档结构和内容
- 根据用户背景信息进行个性化改写
- 输出改写版本到 session

**改写原则：**
- 简化语言，避免学术用语
- 添加详细的解释和示例
- 使用类比和对比学习
- 添加 Mermaid 图表
- 遵循教学化结构

### ReviewerAgent（评审 Agent）

负责严格评审改写后的文档，确保质量达到标准。

**主要功能：**
- 检查 10 个评审标准
- 如果满意，保存文档
- 如果不满意，返回具体改进建议
- 支持最多 5 次迭代

**评审流程：**
1. 逐一检查 10 个标准
2. 对于每个标准，判断是否满足
3. 如果有不满足的标准，列出具体改进建议
4. 只有当所有标准都满足时，才能通过评审

### Mermaid 渲染工具

自动将 Markdown 中的 Mermaid 代码块转换为可渲染的图片 URL。

**功能：**
- 检测 Mermaid 代码块（支持两种格式）
- 提取 Mermaid 代码
- 使用 Base64URL 编码
- 生成 mermaid.ink URL
- 替换为 Markdown 图片链接

**支持的 Mermaid 图表类型：**
- 流程图（Flowchart）
- 时序图（Sequence Diagram）
- 类图（Class Diagram）
- 状态图（State Diagram）
- 甘特图（Gantt Chart）

## 📊 数据流转

### Session 数据结构

系统使用 session 在 Agent 之间传递数据：

```go
session := map[string]interface{}{
    "document_path": "path/to/document.md",
    "original_content": "原始文档内容",
    "document_content": "改写后的文档内容",
    "user_background": "用户背景信息",
    "review_feedback": "评审反馈",
    "iteration_count": 1,
}
```

### 数据流转过程

1. **MainAgent** 读取文档和用户背景信息，放入 session
2. **SummaryAgent** 读取 session，进行改写，更新 `document_content`
3. **ReviewerAgent** 读取 session，进行评审
   - 如果满意：调用保存工具
   - 如果不满意：更新 `review_feedback`，返回给 SummaryAgent
4. **SummaryAgent** 根据反馈进行增量改进
5. 重复 3-4 步，直到满意或达到最大迭代次数

## 🔧 配置说明

### config/config.go

```go
// 模型配置
ModelName: "gpt-4o-mini"  // 使用的 LLM 模型

// 向量配置
EmbeddingModel: "text-embedding-3-small"  // 向量模型
EmbeddingDimension: 1536  // 向量维度

// Milvus 配置
MilvusHost: "localhost"
MilvusPort: 19530
```

### 环境变量

| 变量名 | 说明 | 必需 |
|--------|------|------|
| OPENAI_API_KEY | OpenAI API 密钥 | ✅ |
| OPENAI_BASE_URL | OpenAI API 基础 URL | ❌ |
| FEISHU_APP_ID | 飞书应用 ID | ❌ |
| FEISHU_APP_SECRET | 飞书应用密钥 | ❌ |
| FEISHU_FOLDER_TOKEN | 飞书文件夹 Token | ❌ |
| MILVUS_HOST | Milvus 主机 | ❌ |
| MILVUS_PORT | Milvus 端口 | ❌ |

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./tools -v

# 运行 Mermaid 渲染测试
go test -v ./tools -run TestMermaidRendering

# 运行 URL 验证测试
go test -v ./tools -run TestMermaidURLValidity
```

### 测试覆盖

- ✅ Mermaid 代码块检测和转换
- ✅ Base64URL 编码验证
- ✅ 文件保存功能
- ✅ 飞书 API 集成
- ✅ Agent 工作流程

## 📝 常见问题

### Q: 如何修改改写的用户背景信息？

A: 编辑 `agent/summaryAgent.go` 中的 `Instruction` 字段，修改【用户背景信息】部分。

### Q: 如何增加评审标准？

A: 编辑 `agent/reviewerAgent.go` 中的 `Instruction` 字段，在【严格的评审标准】部分添加新的标准。

### Q: 如何禁用飞书保存功能？

A: 在 `agent/reviewerAgent.go` 中，注释掉 `save_to_feishu` 工具的注册。

### Q: 如何修改最大迭代次数？

A: 编辑 `agent/Supervisor.go` 中的 `maxIterations` 配置。

### Q: Mermaid 图表为什么没有渲染？

A: 
1. 检查 Markdown 代码块格式是否正确（必须是 `` ```mermaid ``）
2. 检查 Mermaid 语法是否正确
3. 确保网络连接正常（需要访问 mermaid.ink 服务）

## 🚨 故障排除

### 常见错误

| 错误 | 原因 | 解决方案 |
|------|------|---------|
| `OPENAI_API_KEY not set` | 未设置 OpenAI API 密钥 | 设置环境变量 `OPENAI_API_KEY` |
| `connection refused` | Milvus 未启动 | 运行 `docker-compose up -d` |
| `tool not found` | 工具未注册 | 检查 `tools/tool.go` 中的工具注册 |
| `HTTP 400 on mermaid.ink` | Mermaid 代码格式错误 | 检查 Mermaid 语法和代码块格式 |

### 调试技巧

1. **启用详细日志**
```bash
export LOG_LEVEL=debug
go run main.go
```

2. **检查 Session 内容**
在 Agent 中添加日志输出 session 内容

3. **验证 API 连接**
```bash
curl -X POST https://api.openai.com/v1/chat/completions \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"model":"gpt-4o-mini","messages":[{"role":"user","content":"test"}]}'
```

## 📚 相关文档

- [项目架构详解](./PROJECT_ARCHITECTURE.md) - 完整的系统架构和设计原理
- [SummaryAgent 使用指南](./SUMMARY_AGENT_GUIDE.md) - 改写 Agent 的详细说明
- [Mermaid 渲染修复说明](./MERMAID_RENDERING_FIX.md) - Mermaid 图表渲染的技术细节
- [执行结果分析](./EXECUTION_RESULTS.md) - 系统运行结果和性能分析

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

### 开发流程

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 代码规范

- 遵循 Go 官方代码规范
- 添加必要的注释和文档
- 编写单元测试
- 确保代码通过 `go fmt` 和 `go vet`

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](./LICENSE) 文件。

## 👥 作者

- 项目维护者：[Your Name]
- 贡献者：[Contributors]

## 📞 联系方式

- 提交 Issue：[GitHub Issues]
- 讨论：[GitHub Discussions]
- 邮件：[Your Email]

## 🙏 致谢

感谢以下项目和团队的支持：

- [Eino Framework](https://github.com/cloudwego/eino) - 强大的 Agent 框架
- [Milvus](https://milvus.io/) - 向量数据库
- [OpenAI](https://openai.com/) - LLM 服务
- [飞书 API](https://open.feishu.cn/) - 文档保存服务

---

**最后更新**：2025-12-16

**版本**：1.0.0

**状态**：✅ 生产就绪
