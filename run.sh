#!/bin/bash
# Eino Demo 一键运行脚本

set -e  # 遇到错误立即退出

echo "🚀 启动 Eino Demo..."
echo ""

# 获取脚本所在目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# 检查可执行文件是否存在
if [ ! -f "eino_demo" ]; then
    echo "⚠️  未找到编译文件，开始编译..."
    bash build.sh
    echo ""
fi

# 检查 .env 文件
if [ ! -f ".env" ]; then
    echo "❌ 错误: 未找到 .env 文件"
    echo "   请创建 .env 文件并配置 DashScope API 密钥"
    exit 1
fi

# 检查 API Key 是否配置
if ! grep -q "DASHSCOPE_API_KEY" .env; then
    echo "❌ 错误: .env 文件中未找到 DASHSCOPE_API_KEY"
    exit 1
fi

# 运行程序
echo "📝 执行程序..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./eino_demo
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "✅ 程序执行完成！"
