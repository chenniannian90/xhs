#!/bin/bash

# NavHub 部署前验证脚本
# 用途：在部署前自动检查所有必要条件

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"
FRONTEND_DIR="$PROJECT_ROOT/frontend"

echo "======================================"
echo "NavHub 部署前验证"
echo "======================================"
echo "项目目录: $PROJECT_ROOT"
echo ""

# 错误计数
ERRORS=0
WARNINGS=0

# 打印函数
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
    ERRORS=$((ERRORS + 1))
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
    WARNINGS=$((WARNINGS + 1))
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

print_section() {
    echo ""
    echo "======================================"
    echo "$1"
    echo "======================================"
}

# ==================== 1. 环境检查 ====================

print_section "1. 环境检查"

# 检查 Node.js
if command -v node &> /dev/null; then
    NODE_VERSION=$(node -v)
    print_success "Node.js 已安装: $NODE_VERSION"
else
    print_error "Node.js 未安装"
fi

# 检查 npm
if command -v npm &> /dev/null; then
    NPM_VERSION=$(npm -v)
    print_success "npm 已安装: $NPM_VERSION"
else
    print_error "npm 未安装"
fi

# 检查 Go
if command -v go &> /dev/null; then
    GO_VERSION=$(go version)
    print_success "Go 已安装: $GO_VERSION"
else
    print_error "Go 未安装"
fi

# 检查 rsync
if command -v rsync &> /dev/null; then
    print_success "rsync 已安装"
else
    print_error "rsync 未安装，部署需要 rsync"
fi

# ==================== 2. 后端检查 ====================

print_section "2. 后端检查"

# 检查后端目录
if [ -d "$BACKEND_DIR" ]; then
    print_success "后端目录存在"

    # 检查 go.mod
    if [ -f "$BACKEND_DIR/go.mod" ]; then
        print_success "go.mod 存在"
    else
        print_error "go.mod 不存在"
    fi

    # 检查主要源文件
    if [ -f "$BACKEND_DIR/cmd/server/main.go" ]; then
        print_success "main.go 存在"
    else
        print_error "main.go 不存在"
    fi

    # 检查配置文件
    if [ -f "$BACKEND_DIR/.env" ] || [ -f "$BACKEND_DIR/config.yml" ]; then
        print_success "配置文件存在"
    else
        print_warning "配置文件不存在，将使用环境变量"
    fi
else
    print_error "后端目录不存在: $BACKEND_DIR"
fi

# ==================== 3. 前端检查 ====================

print_section "3. 前端检查"

# 检查前端目录
if [ -d "$FRONTEND_DIR" ]; then
    print_success "前端目录存在"

    # 检查 package.json
    if [ -f "$FRONTEND_DIR/package.json" ]; then
        print_success "package.json 存在"

        # 检查依赖是否已安装
        if [ -d "$FRONTEND_DIR/node_modules" ]; then
            print_success "node_modules 存在"
        else
            print_warning "node_modules 不存在，需要运行 npm install"
        fi
    else
        print_error "package.json 不存在"
    fi

    # 检查 .env 文件
    if [ -f "$FRONTEND_DIR/.env" ]; then
        print_success ".env 文件存在"

        # 验证 VITE_API_URL 配置
        API_URL=$(grep "^VITE_API_URL=" "$FRONTEND_DIR/.env" | cut -d'=' -f2)
        if [ -n "$API_URL" ]; then
            print_success "VITE_API_URL 配置: $API_URL"

            # 检查是否有路径重复问题
            if [[ "$API_URL" == *"/api/v1"* ]]; then
                print_info "检查 API 路径配置..."

                # 检查代码中是否还有 /api/v1 前缀
                DUPLICATE_PATHS=$(grep -r "api\.get.*'/api/v1" "$FRONTEND_DIR/src" 2>/dev/null | wc -l)
                if [ "$DUPLICATE_PATHS" -gt 0 ]; then
                    print_error "发现 $DUPLICATE_PATHS 处 API 路径重复问题！"
                    print_info "baseURL 已经是 $API_URL，代码中不应该再包含 /api/v1"
                    print_info "问题文件："
                    grep -r "api\.get.*'/api/v1" "$FRONTEND_DIR/src" 2>/dev/null | head -5
                else
                    print_success "没有发现 API 路径重复问题"
                fi
            fi
        else
            print_warning "VITE_API_URL 未配置"
        fi
    else
        print_warning ".env 文件不存在"
    fi

    # 检查 TypeScript 配置
    if [ -f "$FRONTEND_DIR/tsconfig.json" ]; then
        print_success "tsconfig.json 存在"
    else
        print_error "tsconfig.json 不存在"
    fi

    # 检查 Vite 配置
    if [ -f "$FRONTEND_DIR/vite.config.ts" ]; then
        print_success "vite.config.ts 存在"
    else
        print_error "vite.config.ts 不存在"
    fi
else
    print_error "前端目录不存在: $FRONTEND_DIR"
fi

# ==================== 4. 代码质量检查 ====================

print_section "4. 代码质量检查"

# 检查 TypeScript 编译
if [ -d "$FRONTEND_DIR" ] && [ -f "$FRONTEND_DIR/package.json" ]; then
    cd "$FRONTEND_DIR"
    print_info "检查 TypeScript 类型..."

    if npx tsc --noEmit 2>/dev/null; then
        print_success "TypeScript 类型检查通过"
    else
        print_error "TypeScript 类型检查失败"
        print_info "运行 'npm run build' 查看详细错误"
    fi
    cd "$PROJECT_ROOT"
fi

# ==================== 5. API 路径配置检查 ====================

print_section "5. API 路径配置检查"

API_PATTERN_CHECKS=0
API_PATTERN_ERRORS=0

# 检查前端 API 调用模式
if [ -d "$FRONTEND_DIR/src" ]; then
    print_info "检查前端 API 调用模式..."

    # 查找所有包含 /api/v1 的 API 调用
    PROBLEMATIC_FILES=$(find "$FRONTEND_DIR/src" -name "*.tsx" -o -name "*.ts" | xargs grep -l "api\.\(get\|post\|put\|delete\).*'/api/v1" 2>/dev/null || true)

    if [ -n "$PROBLEMATIC_FILES" ]; then
        print_error "以下文件包含重复的 /api/v1 路径："
        echo "$PROBLEMATIC_FILES" | while read file; do
            echo "  - $file"
            API_PATTERN_ERRORS=$((API_PATTERN_ERRORS + 1))
        done
    else
        print_success "前端 API 谯径配置正确"
    fi

    API_PATTERN_CHECKS=$((API_PATTERN_CHECKS + 1))
fi

# ==================== 6. 测试检查 ====================

print_section "6. 测试检查"

# 检查是否有测试文件
TEST_COUNT=$(find "$PROJECT_ROOT" -name "*.test.ts" -o -name "*.test.tsx" -o -name "*.spec.ts" 2>/dev/null | wc -l)

if [ "$TEST_COUNT" -gt 0 ]; then
    print_success "发现 $TEST_COUNT 个测试文件"
else
    print_warning "未发现测试文件，建议添加自动化测试"
fi

# ==================== 7. Git 状态检查 ====================

print_section "7. Git 状态检查"

if [ -d "$PROJECT_ROOT/.git" ]; then
    print_success "Git 仓库存在"

    # 检查是否有未提交的更改
    if git -C "$PROJECT_ROOT" diff --quiet 2>/dev/null; then
        print_success "没有未提交的更改"
    else
        CHANGED_FILES=$(git -C "$PROJECT_ROOT" diff --name-only | wc -l)
        print_warning "有 $CHANGED_FILES 个文件未提交"
    fi

    # 检查当前分支
    CURRENT_BRANCH=$(git -C "$PROJECT_ROOT" branch --show-current 2>/dev/null || echo "unknown")
    print_info "当前分支: $CURRENT_BRANCH"

    if [ "$CURRENT_BRANCH" = "main" ] || [ "$CURRENT_BRANCH" = "master" ]; then
        print_warning "您在主分支上，考虑创建功能分支"
    fi
else
    print_warning "不是 Git 仓库"
fi

# ==================== 8. 构建测试 ====================

print_section "8. 构建测试"

# 测试前端构建
if [ -d "$FRONTEND_DIR" ] && [ -f "$FRONTEND_DIR/package.json" ]; then
    print_info "测试前端构建..."
    cd "$FRONTEND_DIR"

    if npm run build 2>&1 | grep -q "built in"; then
        print_success "前端构建成功"
    else
        print_error "前端构建失败"
    fi

    cd "$PROJECT_ROOT"
fi

# 测试后端构建
if [ -d "$BACKEND_DIR" ] && [ -f "$BACKEND_DIR/go.mod" ]; then
    print_info "测试后端构建..."
    cd "$BACKEND_DIR"

    if GOOS=linux GOARCH=amd64 go build -o /tmp/navhub-api-test cmd/server/main.go 2>/dev/null; then
        print_success "后端构建成功 (Linux x86_64)"
        rm -f /tmp/navhub-api-test
    else
        print_error "后端构建失败"
    fi

    cd "$PROJECT_ROOT"
fi

# ==================== 总结 ====================

print_section "检查结果总结"

echo ""
if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✓ 所有检查通过！可以安全部署。${NC}"
    echo ""
    echo "下一步："
    echo "  1. 运行后端测试: cd backend && ./scripts/test_api.sh"
    echo "  2. 部署后端: ./deploy/deploy-backend.sh"
    echo "  3. 部署前端: ./deploy/deploy-frontend.sh"
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo -e "${YELLOW}⚠ 检查通过，但有 $WARNINGS 个警告${NC}"
    echo ""
    echo "建议在部署前查看并解决警告。"
    exit 0
else
    echo -e "${RED}✗ 发现 $ERRORS 个错误，$WARNINGS 个警告${NC}"
    echo ""
    echo "请修复错误后再部署。"
    exit 1
fi
