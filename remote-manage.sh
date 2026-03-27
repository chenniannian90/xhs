#!/bin/bash

# NavHub 远程服务器管理脚本

# 配置
SERVER_HOST="117.72.39.169"
SERVER_USER="root"
SERVER_PATH="/data/web"
DOMAIN="hao246.cn"

# 颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 显示使用说明
show_help() {
    echo -e "${GREEN}NavHub 远程管理工具${NC}"
    echo ""
    echo "使用方法: ./remote-manage.sh [命令]"
    echo ""
    echo "可用命令:"
    echo "  logs-backend      查看后端日志"
    echo "  logs-nginx        查看Nginx日志"
    echo "  logs-all          查看所有日志"
    echo "  status            查看服务状态"
    echo "  restart-backend   重启后端服务"
    echo "  restart-nginx     重启Nginx服务"
    echo "  restart-all       重启所有服务"
    echo "  ssh               连接到服务器"
    echo "  health            健康检查"
    echo ""
}

# 查看后端日志
logs_backend() {
    echo -e "${YELLOW}查看后端日志...${NC}"
    ssh ${SERVER_USER}@${SERVER_HOST} "tail -f ${SERVER_PATH}/logs/backend.log"
}

# 查看Nginx日志
logs_nginx() {
    echo -e "${YELLOW}查看Nginx日志...${NC}"
    ssh ${SERVER_USER}@${SERVER_HOST} "tail -f ${SERVER_PATH}/logs/nginx-error.log"
}

# 查看所有日志
logs_all() {
    echo -e "${YELLOW}查看所有日志...${NC}"
    ssh ${SERVER_USER}@${SERVER_HOST} "multitail ${SERVER_PATH}/logs/backend.log ${SERVER_PATH}/logs/nginx-error.log"
}

# 查看服务状态
status() {
    echo -e "${YELLOW}服务状态:${ ${NC}"
    echo ""
    ssh ${SERVER_USER}@${SERVER_HOST} << ENDSSH
echo "=== 后端服务 ==="
systemctl status navhub-api | head -5
echo ""
echo "=== Nginx 服务 ==="
systemctl status nginx | head -5
echo ""
echo "=== 端口监听 ==="
netstat -tlnp | grep -E '80|8080|443' | head -10
ENDSSH
}

# 重启后端
restart_backend() {
    echo -e "${YELLOW}重启后端服务...${NC}"
    ssh ${SERVER_USER}@${SERVER_HOST} "systemctl restart navhub-api && echo '✅ 后端已重启'"
}

# 重启Nginx
restart_nginx() {
    echo -e "${YELLOW}重启Nginx服务...${NC}"
    ssh ${SERVER_USER}@${SERVER_HOST} "systemctl restart nginx && echo '✅ Nginx已重启'"
}

# 重启所有
restart_all() {
    echo -e "${YELLOW}重启所有服务...${NC}"
    ssh ${SERVER_USER}@${SERVER_HOST} << ENDSSH
systemctl restart navhub-api
systemctl restart nginx
echo "✅ 所有服务已重启"
ENDSSH
}

# SSH连接
ssh_connect() {
    echo -e "${YELLOW}连接到服务器...${NC}"
    ssh ${SERVER_USER}@${SERVER_HOST}
}

# 健康检查
health() {
    echo -e "${YELLOW}健康检查:${NC}"
    echo ""

    # 检查服务器连接
    if ssh -o ConnectTimeout=5 ${SERVER_USER}@${SERVER_HOST} "echo '✅ 服务器连接正常'" 2>/dev/null; then
        echo -e "${GREEN}✅ 服务器连接正常${NC}"
    else
        echo -e "${RED}❌ 服务器连接失败${NC}"
        return 1
    fi

    # 检查后端API
    if curl -s http://${DOMAIN}/health > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 后端API正常${NC}"
        curl -s http://${DOMAIN}/health | python3 -m json.tool 2>/dev/null || echo "  响应格式错误"
    else
        echo -e "${RED}❌ 后端API无响应${NC}"
    fi

    # 检查前端
    if curl -s http://${DOMAIN} > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 前端正常${NC}"
    else
        echo -e "${RED}❌ 前端无响应${NC}"
    fi

    echo ""
    echo "访问地址: http://${DOMAIN}"
}

# 主程序
case "$1" in
    logs-backend)
        logs_backend
        ;;
    logs-nginx)
        logs_nginx
        ;;
    logs-all)
        logs_all
        ;;
    status)
        status
        ;;
    restart-backend)
        restart_backend
        ;;
    restart-nginx)
        restart_nginx
        ;;
    restart-all)
        restart_all
        ;;
    ssh)
        ssh_connect
        ;;
    health)
        health
        ;;
    *)
        show_help
        ;;
esac
