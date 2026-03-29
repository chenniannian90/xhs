#!/bin/bash

# NavHub 一键部署脚本
# 使用方法: ./deploy.sh [dev|prod] [frontend|backend|all]

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 检查端口是否被占用
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1 ; then
        return 0
    else
        return 1
    fi
}

# 解析参数
MODE=${1:-dev}
COMPONENT=${2:-all}

print_info "🚀 NavHub 部署脚本"
print_info "模式: $MODE | 组件: $COMPONENT"
echo ""

# 启动基础服务
print_info "📦 启动数据库服务..."
docker compose up -d postgres redis mailhog

# 等待数据库就绪
print_info "⏳ 等待数据库就绪..."
sleep 5

# 检查后端
if [ "$COMPONENT" = "backend" ] || [ "$COMPONENT" = "all" ]; then
    print_info "🔧 配置后端..."

    # 检查 .env 文件
    if [ ! -f "backend/.env" ]; then
        print_warning "未找到 backend/.env，从 .env.example 复制..."
        cp backend/.env.example backend/.env
        print_success "已创建 backend/.env"
    fi

    # 安装 Go 依赖
    print_info "📚 安装 Go 依赖..."
    cd backend
    go mod download
    cd ..

    # 启动后端
    if check_port 8080; then
        print_warning "端口 8080 已被占用，尝试停止现有进程..."
        pkill -f "backend/main" || true
        sleep 2
    fi

    print_info "🚀 启动后端服务..."
    cd backend
    if [ "$MODE" = "dev" ]; then
        go run cmd/server/main.go > ../backend.log 2>&1 &
        echo $! > ../backend.pid
    else
        go build -o main cmd/server/main.go
        ./main > ../backend.log 2>&1 &
        echo $! > ../backend.pid
    fi
    cd ..

    sleep 3
    print_success "后端服务已启动"
fi

# 检查前端
if [ "$COMPONENT" = "frontend" ] || [ "$COMPONENT" = "all" ]; then
    print_info "🎨 配置前端..."

    # 检查 node_modules
    if [ ! -d "frontend/node_modules" ]; then
        print_info "📚 安装前端依赖..."
        cd frontend
        npm install
        cd ..
    fi

    # 检查前端环境变量
    if [ ! -f "frontend/.env" ]; then
        print_warning "未找到 frontend/.env，创建默认配置..."
        cat > frontend/.env << 'ENVEOF'
VITE_API_URL=http://localhost:8080
VITE_APP_NAME=NavHub
ENVEOF
        print_success "已创建 frontend/.env"
    fi

    # 启动前端
    if check_port 5173; then
        print_warning "端口 5173 已被占用，尝试停止现有进程..."
        pkill -f "vite" || true
        sleep 2
    fi

    print_info "🚀 启动前端服务..."
    cd frontend
    if [ "$MODE" = "dev" ]; then
        npm run dev > ../frontend.log 2>&1 &
        echo $! > ../frontend.pid
    else
        npm run build
        npm run preview > ../frontend.log 2>&1 &
        echo $! > ../frontend.pid
    fi
    cd ..

    sleep 3
    print_success "前端服务已启动"
fi

echo ""
print_success "🎉 部署完成！"
echo ""
print_info "📊 服务地址："
echo "  - 前端:     http://localhost:5173"
echo "  - 后端 API: http://localhost:8080"
echo "  - MailHog:  http://localhost:8025"
echo ""
print_info "📝 日志文件："
echo "  - 后端: tail -f backend.log"
echo "  - 前端: tail -f frontend.log"
