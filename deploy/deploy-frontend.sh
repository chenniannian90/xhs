#!/bin/bash

# NavHub 前端自动部署脚本
# 用途：自动化前端构建、测试和部署流程

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FRONTEND_DIR="$PROJECT_ROOT/frontend"
REMOTE_SERVER="${DEPLOY_SERVER:-root@117.72.39.169}"
REMOTE_PATH="/data/web/frontend"

echo "======================================"
echo "NavHub 前端自动部署"
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

if [ ! -d "$FRONTEND_DIR" ]; then
    print_error "前端目录不存在: $FRONTEND_DIR"
    exit 1
fi

print_success "前端目录检查通过"

# 检查必要文件
if [ ! -f "$FRONTEND_DIR/package.json" ]; then
    print_error "package.json 不存在"
    exit 1
fi

if [ ! -f "$FRONTEND_DIR/.env" ]; then
    print_error ".env 文件不存在"
    exit 1
fi

print_success "必要文件检查通过"

# 检查依赖
if [ ! -d "$FRONTEND_DIR/node_modules" ]; then
    print_error "node_modules 不存在，请先运行 'npm install'"
    exit 1
fi

print_success "依赖检查通过"

# ==================== 步骤 2: 环境配置检查 ====================

print_step "步骤 2/7: 环境配置检查"

cd "$FRONTEND_DIR"

# 验证 API URL 配置
API_URL=$(grep "^VITE_API_URL=" .env | cut -d'=' -f2)

if [ -z "$API_URL" ]; then
    print_error "VITE_API_URL 未配置"
    exit 1
fi

print_success "VITE_API_URL 配置: $API_URL"

# 检查 API 路径配置问题
if [[ "$API_URL" == *"/api/v1"* ]]; then
    echo "检查 API 路径重复问题..."

    DUPLICATE_COUNT=$(grep -r "api\.\(get\|post\|put\|delete\).*'/api/v1" src/ 2>/dev/null | wc -l)

    if [ "$DUPLICATE_COUNT" -gt 0 ]; then
        print_error "发现 $DUPLICATE_COUNT 处 API 路径重复！"
        echo ""
        echo "问题文件："
        grep -r "api\.\(get\|post\|put\|delete\).*'/api/v1" src/ 2>/dev/null | head -5
        echo ""
        echo "baseURL 已经是 $API_URL"
        echo "代码中不应该再包含 /api/v1 前缀"
        exit 1
    fi

    print_success "没有 API 路径重复问题"
fi

cd "$PROJECT_ROOT"

# ==================== 步骤 3: 代码质量检查 ====================

print_step "步骤 3/7: 代码质量检查"

cd "$FRONTEND_DIR"

# TypeScript 类型检查
echo "运行 TypeScript 类型检查..."

if npx tsc --noEmit 2>&1 | grep -q "error TS"; then
    print_error "TypeScript 类型检查失败"
    echo "详细错误："
    npx tsc --noEmit
    exit 1
fi

print_success "TypeScript 类型检查通过"

# ESLint 检查（如果配置了）
if [ -f ".eslintrc.json" ] || [ -f ".eslintrc.js" ] || [ -f ".eslintrc.cjs" ]; then
    echo "运行 ESLint 检查..."

    if npx eslint src/ --ext .ts,.tsx 2>&1 | grep -q "error"; then
        print_error "ESLint 检查失败"
        npx eslint src/ --ext .ts,.tsx
        exit 1
    fi

    print_success "ESLint 检查通过"
fi

cd "$PROJECT_ROOT"

# ==================== 步骤 4: 运行测试 ====================

print_step "步骤 4/7: 运行测试"

cd "$FRONTEND_DIR"

# 运行测试（如果配置了）
if grep -q '"test"' package.json; then
    echo "运行测试..."

    if npm test 2>&1 | grep -q "passed"; then
        print_success "测试通过"
    else
        print_warning "测试失败或没有配置测试"
    fi
else
    print_warning "未配置测试脚本"
fi

cd "$PROJECT_ROOT"

# ==================== 步骤 5: 构建 ====================

print_step "步骤 5/7: 构建前端"

cd "$FRONTEND_DIR"

echo "清理旧的构建文件..."
rm -rf dist/

echo "构建生产版本..."

if npm run build; then
    print_success "前端构建成功"

    # 显示构建产物
    echo ""
    echo "构建产物："
    ls -lh dist/
    echo ""

    # 检查关键文件
    if [ ! -f "dist/index.html" ]; then
        print_error "index.html 未生成"
        exit 1
    fi

    # 检查 JS 文件
    JS_FILE=$(find dist -name "*.js" -type f | head -1)
    if [ -z "$JS_FILE" ]; then
        print_error "JavaScript 文件未生成"
        exit 1
    fi

    print_success "构建产物验证通过"
else
    print_error "前端构建失败"
    exit 1
fi

cd "$PROJECT_ROOT"

# ==================== 步骤 6: 备份当前版本 ====================

print_step "步骤 6/7: 备份当前版本"

echo "连接到服务器备份当前版本..."

ssh "$REMOTE_SERVER" << 'ENDSSH'
if [ -d /data/web/frontend ]; then
    BACKUP_NAME="frontend.backup.$(date +%Y%m%d_%H%M%S)"
    cp -r /data/web/frontend "/data/web/$BACKUP_NAME"
    echo "已备份为: $BACKUP_NAME"
else
    echo "没有现有版本需要备份"
fi
ENDSSH

print_success "备份完成"

# ==================== 步骤 7: 部署 ====================

print_step "步骤 7/7: 部署到服务器"

echo "上传前端文件..."

rsync -avz --progress --delete \
    "$FRONTEND_DIR/dist/" \
    "$REMOTE_SERVER:$REMOTE_PATH/"

print_success "文件上传完成"

# 验证部署
echo "验证部署..."

# 检查 index.html
if ssh "$REMOTE_SERVER" "test -f $REMOTE_PATH/index.html"; then
    print_success "index.html 已部署"
else
    print_error "index.html 部署失败"
    exit 1
fi

# 检查 JS 文件
if ssh "$REMOTE_SERVER" "test -f $REMOTE_PATH/assets/*.js"; then
    print_success "JavaScript 文件已部署"
else
    print_error "JavaScript 文件部署失败"
    exit 1
fi

# ==================== 部署后验证 ====================

print_step "部署后验证"

echo "测试前端访问..."

# 测试主页
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "http://$REMOTE_SERVER/")

if [ "$HTTP_CODE" = "200" ]; then
    print_success "前端主页可访问 (HTTP 200)"
else
    print_error "前端主页访问失败 (HTTP $HTTP_CODE)"
    exit 1
fi

# 清理浏览器缓存提示
echo ""
echo -e "${YELLOW}⚠ 请清除浏览器缓存或使用 Ctrl+F5 强制刷新${NC}"

# ==================== 清理 ====================

print_step "清理本地文件"

cd "$FRONTEND_DIR"

# 可选：保留 dist 文件用于调试
# rm -rf dist/

print_success "清理完成"

# ==================== 总结 ====================

echo ""
echo "======================================"
echo -e "${GREEN}✓ 前端部署完成！${NC}"
echo "======================================"
echo ""
echo "部署信息："
echo "  服务器: $REMOTE_SERVER"
echo "  路径: $REMOTE_PATH"
echo ""
echo "访问地址："
echo "  http://$REMOTE_SERVER/"
echo "  https://$REMOTE_SERVER/"
echo ""
echo "测试建议："
echo "  1. 清除浏览器缓存"
echo "  2. 打开开发者工具 (F12)"
echo "  3. 访问上述地址"
echo "  4. 测试登录、注册等功能"
echo ""
echo "查看 Nginx 日志："
echo "  ssh $REMOTE_SERVER 'tail -f /var/log/nginx/access.log'"
echo ""
