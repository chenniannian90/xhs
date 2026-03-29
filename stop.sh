#!/bin/bash

echo "🛑 停止 NavHub 服务..."

# 停止后端
if [ -f "backend.pid" ]; then
    echo "停止后端服务..."
    kill $(cat backend.pid) 2>/dev/null || true
    rm backend.pid
fi

# 停止前端
if [ -f "frontend.pid" ]; then
    echo "停止前端服务..."
    kill $(cat frontend.pid) 2>/dev/null || true
    rm frontend.pid
fi

# 清理进程
pkill -f "backend/main" 2>/dev/null || true
pkill -f "vite" 2>/dev/null || true

echo "✅ 所有服务已停止"
