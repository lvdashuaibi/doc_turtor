# 🔒 安全配置指南
本文档说明如何安全地配置项目，避免将敏感信息上传到 GitHub。

## 📋 目录

1. [敏感信息清单](#敏感信息清单)
2. [环境变量配置](#环境变量配置)
3. [本地开发设置](#本地开发设置)
4. [CI/CD 配置](#cicd-配置)
5. [安全检查清单](#安全检查清单)

---

## 敏感信息清单

以下信息 **绝不应该** 提交到 GitHub：

### 🔑 API 密钥和令牌

- ✅ **OpenAI API Key** - 用于调用 GPT 模型
- ✅ **飞书 App ID 和 Secret** - 用于飞书 API 认证
- ✅ **飞书 Folder Token** - 用于访问飞书文件夹
- ✅ **数据库密码** - Milvus 或其他数据库的认证信息

### 🗝️ 认证信息

- ✅ **OAuth 令牌** - 第三方服务的认证令牌
- ✅ **JWT 密钥** - 用于签名和验证 JWT
- ✅ **会话令牌** - 用户会话相关的令牌

### 📊 配置信息

- ✅ **数据库连接字符串** - 包含密码的连接字符串
- ✅ **服务器地址** - 内部服务器的 IP 或域名
- ✅ **内部 API 端点** - 不应该公开的内部 API

### 📁 文件和数据

- ✅ **生成的文档** - 改写后的文档（可能包含敏感内容）
- ✅ **日志文件** - 可能包含敏感信息的日志
- ✅ **数据库备份** - 包含真实数据的备份文件

---

## 环境变量配置

### 方法 1：使用 .env 文件（推荐用于本地开发）

1. **复制示例文件**
```bash
cp .env.example .env
```

2. **编辑 .env 文件**
```bash
# 使用你喜欢的编辑器打开 .env
vim .env
# 或
nano .env
```

3. **填入实际的值**
```bash
OPENAI_API_KEY=sk-your-actual-key-here
FEISHU_APP_ID=cli_your_actual_id_here
# ... 其他配置
```

4. **确保 .env 在 .gitignore 中**
```bash
# 检查 .gitignore 是否包含 .env
grep "^\.env$" .gitignore
```

### 方法 2：使用系统环境变量

```bash
# 在 shell 配置文件中设置（~/.bashrc, ~/.zshrc 等）
export OPENAI_API_KEY="sk-your-actual-key-here"
export FEISHU_APP_ID="cli_your_actual_id_here"
export FEISHU_APP_SECRET="your_app_secret_here"
export FEISHU_FOLDER_TOKEN="your_folder_token_here"

# 然后重新加载配置
source ~/.bashrc  # 或 source ~/.zshrc
```

### 方法 3：使用 direnv（推荐用于多项目开发）

1. **安装 direnv**
```bash
# macOS
brew install direnv

# Linux
sudo apt-get install direnv

# 或从源码安装
# https://direnv.net/docs/installation.html
```

2. **配置 shell**
```bash
# 对于 bash
echo 'eval "$(direnv hook bash)"' >> ~/.bashrc

# 对于 zsh
echo 'eval "$(direnv hook zsh)"' >> ~/.zshrc

# 重新加载配置
source ~/.bashrc  # 或 source ~/.zshrc
```

3. **创建 .envrc 文件**
```bash
# 在项目根目录创建 .envrc
cat > .envrc << 'EOF'
# 从 .env 文件加载环境变量
dotenv .env
EOF

# 允许 direnv
direnv allow
```

4. **添加 .envrc 到 .gitignore**
```bash
echo ".envrc" >> .gitignore
```

---

## 本地开发设置

### 初始化项目

```bash
# 1. 克隆项目
git clone <repository-url>
cd eino_test

# 2. 复制环境变量示例
cp .env.example .env

# 3. 编辑 .env 文件，填入实际的值
vim .env

# 4. 验证 .env 不会被提交
git status  # 应该看不到 .env

# 5. 安装依赖
go mod download
go mod tidy

# 6. 运行项目
go run main.go
```

### 验证敏感信息不会被提交

```bash
# 检查 .gitignore 配置
cat .gitignore | grep -E "\.env|secrets|credentials|keys"

# 检查暂存区中是否有敏感文件
git diff --cached --name-only | grep -E "\.env|secrets|credentials"

# 检查本地是否有未追踪的敏感文件
git status --porcelain | grep -E "\.env|secrets|credentials"

# 查看 git 会追踪的文件
git ls-files | grep -E "\.env|secrets|credentials"
```

---

## CI/CD 配置

### GitHub Actions 配置

在 `.github/workflows/` 目录中创建工作流文件，使用 GitHub Secrets：

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24.5'
      
      - name: Run tests
        env:
          OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
          FEISHU_APP_ID: ${{ secrets.FEISHU_APP_ID }}
          FEISHU_APP_SECRET: ${{ secrets.FEISHU_APP_SECRET }}
          FEISHU_FOLDER_TOKEN: ${{ secrets.FEISHU_FOLDER_TOKEN }}
        run: go test ./...
```

### 设置 GitHub Secrets

1. 进入 GitHub 仓库设置
2. 选择 **Settings** → **Secrets and variables** → **Actions**
3. 点击 **New repository secret**
4. 添加以下 secrets：
   - `OPENAI_API_KEY`
   - `FEISHU_APP_ID`
   - `FEISHU_APP_SECRET`
   - `FEISHU_FOLDER_TOKEN`

### GitLab CI/CD 配置

```yaml
# .gitlab-ci.yml
stages:
  - test

test:
  stage: test
  image: golang:1.24.5
  script:
    - go test ./...
  variables:
    OPENAI_API_KEY: $OPENAI_API_KEY
    FEISHU_APP_ID: $FEISHU_APP_ID
    FEISHU_APP_SECRET: $FEISHU_APP_SECRET
    FEISHU_FOLDER_TOKEN: $FEISHU_FOLDER_TOKEN
```

---

## 安全检查清单

### 提交前检查

- [ ] 已复制 `.env.example` 为 `.env`
- [ ] 已填入实际的 API 密钥和令牌
- [ ] 已验证 `.env` 在 `.gitignore` 中
- [ ] 已运行 `git status` 确认 `.env` 不会被提交
- [ ] 已检查是否有其他敏感文件会被提交

### 代码审查检查

- [ ] 代码中没有硬编码的 API 密钥
- [ ] 代码中没有硬编码的密码或令牌
- [ ] 所有敏感信息都通过环境变量读取
- [ ] 日志中不会输出敏感信息
- [ ] 错误消息中不会泄露敏感信息

### 仓库维护检查

- [ ] 定期检查 `.gitignore` 是否完整
- [ ] 定期扫描仓库历史中是否有泄露的密钥
- [ ] 如果发现泄露，立即轮换密钥
- [ ] 使用 git-secrets 或类似工具防止密钥泄露

---

## 🚨 如果不小心提交了敏感信息

### 立即行动

1. **轮换所有密钥**
```bash
# 立即轮换 OpenAI API 密钥
# 立即轮换飞书 App Secret
# 立即轮换所有其他敏感信息
```

2. **从 Git 历史中删除**
```bash
# 使用 git-filter-branch（不推荐，可能破坏历史）
git filter-branch --tree-filter 'rm -f .env' HEAD

# 或使用 BFG Repo-Cleaner（推荐）
bfg --delete-files .env
```

3. **强制推送**
```bash
git push origin --force-with-lease
```

4. **通知团队**
- 通知所有团队成员已轮换密钥
- 更新所有相关的配置

### 预防措施

1. **安装 git-secrets**
```bash
# macOS
brew install git-secrets

# Linux
git clone https://github.com/awslabs/git-secrets.git
cd git-secrets
make install

# 配置 git-secrets
git secrets --install
git secrets --register-aws
```

2. **配置 pre-commit hook**
```bash
# 安装 pre-commit
pip install pre-commit

# 创建 .pre-commit-config.yaml
cat > .pre-commit-config.yaml << 'EOF'
repos:
  - repo: https://github.com/awslabs/git-secrets
    rev: master
    hooks:
      - id: git-secrets
EOF

# 安装 hook
pre-commit install
```

---

## 📚 参考资源

- [GitHub - Keeping your account and data secure](https://docs.github.com/en/code-security/getting-started/best-practices-for-repository-security)
- [OWASP - Secrets Management](https://cheatsheetseries.owasp.org/cheatsheets/Secrets_Management_Cheat_Sheet.html)
- [git-secrets](https://github.com/awslabs/git-secrets)
- [direnv](https://direnv.net/)
- [pre-commit](https://pre-commit.com/)

---

## 🤝 问题反馈

如果发现安全问题或有改进建议，请：

1. **不要在 GitHub Issues 中公开讨论**
2. **发送私密邮件** 到项目维护者
3. **使用 GitHub Security Advisory** 报告安全漏洞

---

**最后更新**：2025-12-16

**版本**：1.0.0

**状态**：✅ 安全配置完成
