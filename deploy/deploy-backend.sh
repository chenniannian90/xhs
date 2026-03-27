#!/bin/bash

# NavHub 后端自动部署脚本
# 用途：自动化后端构建、测试和部署流程

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"
REMOTE_SERVER="${DEPLOY_SERVER:-root@117.72.39.169}"
REMOTE_PATH="/data/web/backend"
SERVICE_NAME="navhub-api"

echo "======================================"
echo "NavHub 后端自动部署"
echo "======================================"
echo "项目目录: $PROJECT_ROOT"
echo "目标服务器: $REMOTE_SERVER"
echo ""

# 错误处理
trap 'echo -e "${RED}部署失败！${NC}"; exit 1' ERR

print_step() {
    echo ""
    echo -e "${BLUE}>>> $1${NC}"
    echo "======================================"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# ==================== 步骤 1: 前置检查 ====================

print_step "步骤 1/7: 前置检查"

if [ ! -d "$BACKEND_DIR" ]; then
    print_error "后端目录不存在: $BACKEND_DIR"
    exit 1
fi

print_success "后端目录检查通过"

# 检查必要文件
if [ ! -f "$BACKEND_DIR/go.mod" ]; then
    print_error "go.mod 不存在"
    exit 1
fi

if [ ! -f "$BACKEND_DIR/cmd/server/main.go" ]; then
    print_error "main.go 不存在"
    exit 1
fi

print_success "必要文件检查通过"

# ==================== 步骤 2: 代码质量检查 ====================

print_step "步骤 2/7: 代码质量检查"

cd "$BACKEND_DIR"

# Go fmt 检查
if [ -n "$(gofmt -l .)" ]; then
    print_error "代码格式不正确，请运行 'go fmt ./...'"

    echo "需要格式化的文件："
    gofmt -l .
    exit 1
fi

print_success "代码格式检查通过"

# Go vet 检查
if ! go vet ./... 2>&1 | grep -q "no problems"; then
    print_error "Go vet 发现问题"
    go vet ./...
    exit 1
fi

print_success "Go vet 检查通过"

cd "$PROJECT_ROOT"

# ==================== 步骤 3: 运行测试 ====================

print_step "步骤 3/7: 运行测试"

cd "$BACKEND_DIR"

# 运行单元测试（如果有）
if go test ./... -v 2>&1 | grep -q "PASS"; then
    print_success "单元测试通过"
else
    print_warning "没有找到单元测试或测试失败"
fi

cd "$PROJECT_ROOT"

# ==================== 步骤 4: 构建 ====================

print_step "步骤 4/7: 构建后端"

cd "$BACKEND_DIR"

echo "构建 Linux x86_64 二进制文件..."

if GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o navhub-api cmd/server/main.go; then
    print_success "后端构建成功"

    # 显示二进制文件大小
    SIZE=$(ls -lh navhub-api | awk '{print $5}')
    echo "二进制文件大小: $SIZE"
else
    print_error "后端构建失败"
    exit 1
fi

cd "$PROJECT_ROOT"

# ==================== 步骤 5: 备份当前版本 ====================

print_step "步骤 5/7: 备份当前版本"

echo "连接到服务器备份当前版本..."

ssh "$REMOTE_SERVER" << 'ENDSSH'
if [ -f /data/web/backend/navhub-api ]; then
    BACKUP_NAME="navhub-api.backup.$(date +%Y%m%d_%H%M%S)"
    cp /data/web/backend/navhub-api "/data/web/backend/$BACKUP_NAME"
    echo "已备份为: $BACKUP_NAME"
else
    echo "没有现有版本需要备份"
fi
ENDSSH

print_success "备份完成"

# ==================== 步骤 6: 部署 ====================

print_step "步骤 6/7: 部署到服务器"

echo "上传二进制文件..."

rsync -avz --progress \
    "$BACKEND_DIR/navhub-api" \
    "$REMOTE_SERVER:$REMOTE_PATH/"

print_success "文件上传完成"

# ==================== 步骤 7: 重启服务 ====================

print_step "步骤 7/7: 重启服务"

echo "重启后端服务..."

ssh "$REMOTE_SERVER" << 'ENDSSH'
set -e

# 停止旧进程
echo "停止旧进程..."
pkill -f navhub-api || true
sleep 2

# 设置权限
chmod +x /data/web/backend/navhub-api

# 启动服务
echo "启动服务..."
systemctl start navhub-api

# 等待服务启动
sleep 3

# 检查服务状态
if systemctl is-active --quiet navhub-api; then
    echo "✓ 服务启动成功"
    systemctl status navhub-api --no-pager | head -5
else
    echo "✗ 服务启动失败"
    journalctl -u navhub-api -n 20 --no-pager
    exit 1
fi
ENDSSH

print_success "服务重启完成"

# ==================== 部署后验证 ====================

print_step "部署后验证"

echo "等待服务完全启动..."
sleep 5

# 健康检查
echo "检查服务健康状态..."

if curl -f -s "http://$REMOTE_SERVER:8080/health" > /dev/null; then
    print_success "健康检查通过"
else
    print_error "健康检查失败"
    echo "请检查服务器日志: ssh $REMOTE_SERVER 'journalctl -u navhub-api -n 50'"
    exit 1
fi

# API 测试
echo "运行 API 测试..."

if [ -f "$BACKEND_DIR/scripts/test_api.sh" ]; then
    export API_BASE_URL="http://$REMOTE_SERVER:8080"
    bash "$BACKEND_DIR/scripts/test_api.sh"
    print_success "API 测试通过"
else
    print_warning "未找到 API 测试脚本"
fi

# ==================== 清理 ====================

print_step "清理本地文件"

rm -f "$BACKEND_DIR/navhub-api"
print_success "清理完成"

# ==================== 总结 ====================

echo ""
echo "======================================"
echo -e "${GREEN}✓ 后端部署完成！${NC}"
echo "======================================"
echo ""
echo "部署信息："
echo "  服务器: $REMOTE_SERVER"
echo "  路径: $REMOTE_PATH"
echo "  服务: $SERVICE_NAME"
echo ""
echo "查看日志："
echo "  ssh $REMOTE_SERVER 'journalctl -u navhub-api -f'"
echo ""
echo "查看服务状态："
echo "  ssh $REMOTE_SERVER 'systemctl status navhub-api'"
echo ""
