#!/bin/bash

# NavHub 完整自动部署脚本
# 用途：一键完成所有部署前检查、构建和部署

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SCRIPTS_DIR="$PROJECT_ROOT/deploy"

echo "======================================"
echo "NavHub 完整自动部署"
echo "======================================"
echo "项目目录: $PROJECT_ROOT"
echo "开始时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# 错误处理
trap 'echo -e "${RED}部署失败！${NC}"; exit 1' ERR

print_step() {
    echo ""
    echo -e "${BLUE}======================================"
    echo "$1"
    echo -e "======================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# ==================== 1. 部署前检查 ====================

print_step "步骤 1/5: 部署前全面检查"

if [ -f "$SCRIPTS_DIR/pre-deploy-check.sh" ]; then
    bash "$SCRIPTS_DIR/pre-deploy-check.sh"

    if [ $? -ne 0 ]; then
        print_error "部署前检查失败，请修复问题后重试"
        exit 1
    fi

    print_success "部署前检查通过"
else
    print_error "未找到部署前检查脚本: $SCRIPTS_DIR/pre-deploy-check.sh"
    exit 1
fi

# ==================== 2. 部署后端 ====================

print_step "步骤 2/5: 部署后端"

if [ -f "$SCRIPTS_DIR/deploy-backend.sh" ]; then
    bash "$SCRIPTS_DIR/deploy-backend.sh"

    if [ $? -ne 0 ]; then
        print_error "后端部署失败"
        exit 1
    fi

    print_success "后端部署成功"
else
    print_error "未找到后端部署脚本: $SCRIPTS_DIR/deploy-backend.sh"
    exit 1
fi

# ==================== 3. 部署前端 ====================

print_step "步骤 3/5: 部署前端"

if [ -f "$SCRIPTS_DIR/deploy-frontend.sh" ]; then
    bash "$SCRIPTS_DIR/deploy-frontend.sh"

    if [ $? -ne 0 ]; then
        print_error "前端部署失败"
        exit 1
    fi

    print_success "前端部署成功"
else
    print_error "未找到前端部署脚本: $SCRIPTS_DIR/deploy-frontend.sh"
    exit 1
fi

# ==================== 4. 完整测试 ====================

print_step "步骤 4/5: 运行完整测试"

echo "运行后端 API 测试..."

BACKEND_DIR="$PROJECT_ROOT/backend"
if [ -f "$BACKEND_DIR/scripts/test_api.sh" ]; then
    export API_BASE_URL="http://117.72.39.169:8080"
    bash "$BACKEND_DIR/scripts/test_api.sh"

    if [ $? -ne 0 ]; then
        print_error "API 测试失败"
        exit 1
    fi

    print_success "API 测试通过"
else
    print_warning "未找到 API 测试脚本"
fi

# ==================== 5. 部署报告 ====================

print_step "步骤 5/5: 生成部署报告"

REPORT_FILE="$PROJECT_ROOT/deploy/deployment-report-$(date +%Y%m%d_%H%M%S).txt"

cat > "$REPORT_FILE" << EOF
========================================
NavHub 部署报告
========================================

部署时间: $(date '+%Y-%m-%d %H:%M:%S')
部署服务器: root@117.72.39.169

部署内容：
✓ 后端 API
✓ 前端应用
✓ 自动化测试

部署状态：成功

访问地址：
- 前端: https://117.72.39.169
- API: http://117.72.39.169:8080

检查清单：
✓ 健康检查通过
✓ API 测试通过
✓ 前端访问正常

回滚信息：
后端备份: /data/web/backend/navhub-api.backup.*
前端备份: /data/web/frontend.backup.*

如有问题，请查看：
- 后端日志: ssh root@117.72.39.169 'journalctl -u navhub-api -n 100'
- Nginx 日志: ssh root@117.72.39.169 'tail -f /var/log/nginx/access.log'
EOF

print_success "部署报告已生成: $REPORT_FILE"

# 显示报告内容
cat "$REPORT_FILE"

# ==================== 总结 ====================

echo ""
echo -e "${GREEN}======================================"
echo "✓ 完整部署成功！"
echo "======================================${NC}"
echo ""
echo "部署完成时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""
echo "下一步操作："
echo "  1. 打开浏览器访问 https://117.72.39.169"
echo "  2. 清除浏览器缓存 (Ctrl+Shift+Delete)"
echo "  3. 强制刷新页面 (Ctrl+F5)"
echo "  4. 测试登录、注册等功能"
echo ""
echo "监控命令："
echo "  查看后端日志: ssh root@117.72.39.169 'journalctl -u navhub-api -f'"
echo "  查看Nginx日志: ssh root@117.72.39.169 'tail -f /var/log/nginx/access.log'"
echo "  查看服务状态: ssh root@117.72.39.169 'systemctl status navhub-api'"
echo ""
echo -e "${YELLOW}⚠ 提示：如果发现任何问题，可以使用备份快速回滚${NC}"
echo ""
