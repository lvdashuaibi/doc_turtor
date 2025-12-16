#!/bin/bash
# Eino Demo 一键编译脚本

set -e  # 遇到错误立即退出

echo "🔨 开始编译 Eino Demo..."
echo ""

# 获取脚本所在目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# 检查 Go 环境
echo "📋 检查 Go 环境..."
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到 Go 命令，请先安装 Go"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo "✓ Go 版本: $GO_VERSION"
echo ""

# 整理依赖
echo "📦 整理依赖..."
go mod tidy
echo "✓ 依赖整理完成"
echo ""

# 编译项目
echo "🔨 编译项目..."
# 自动检测所有 .go 文件
GO_FILES=$(find . -maxdepth 1 -name "*.go" -type f | sort)
if [ -z "$GO_FILES" ]; then
    echo "❌ 错误: 未找到 .go 文件"
    exit 1
fi
echo "   源文件: $(echo $GO_FILES | tr '\n' ' ')"
go build -o eino_demo $GO_FILES
echo "✓ 编译完成"
echo ""

# 显示编译结果
if [ -f "eino_demo" ]; then
    FILE_SIZE=$(ls -lh eino_demo | awk '{print $5}')
    echo "✅ 编译成功！"
    echo "   可执行文件: eino_demo ($FILE_SIZE)"
    echo ""
    echo "💡 运行程序: ./eino_demo"
    echo "   或使用: bash run.sh"
else
    echo "❌ 编译失败"
    exit 1
fi
