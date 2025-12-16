#!/bin/bash
# Milvus 向量数据库启动脚本

set -e

echo "🚀 启动 Milvus 向量数据库..."
echo ""

# 获取脚本所在目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# 检查 Docker 是否安装
echo "📋 检查 Docker 环境..."
if ! command -v docker &> /dev/null; then
    echo "❌ 错误: 未找到 Docker 命令，请先安装 Docker"
    echo "   访问: https://www.docker.com/products/docker-desktop"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ 错误: 未找到 docker-compose 命令"
    echo "   请确保 Docker Desktop 已安装（包含 docker-compose）"
    exit 1
fi

DOCKER_VERSION=$(docker --version)
echo "✓ $DOCKER_VERSION"
echo ""

# 检查 docker-compose.yml 文件
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ 错误: 未找到 docker-compose.yml 文件"
    exit 1
fi

echo "📦 启动 Docker 容器..."
echo ""

# 启动容器
docker-compose up -d

echo ""
echo "⏳ 等待 Milvus 服务启动..."
sleep 10

# 检查服务状态
echo ""
echo "🔍 检查服务状态..."
docker-compose ps

echo ""
echo "✅ Milvus 启动成功！"
echo ""
echo "📝 服务信息:"
echo "  • Milvus gRPC: localhost:19530"
echo "  • Milvus HTTP: localhost:9091"
echo "  • Attu Web UI: http://localhost:8000"
echo ""
echo "💡 常用命令:"
echo "  • 查看日志: docker-compose logs -f milvus"
echo "  • 停止服务: docker-compose down"
echo "  • 重启服务: docker-compose restart"
echo "  • 删除数据: docker-compose down -v"
echo ""
echo "🌐 打开 Attu 管理界面: http://localhost:8000"
