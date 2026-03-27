#!/bin/bash

# NavHub 远程部署脚本
# 目标服务器: root@117.72.39.169:/data/web
# 域名: hao246.cn

set -e  # 遇到错误立即退出

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
SERVER_HOST="117.72.39.169"
SERVER_USER="root"
SERVER_PATH="/data/web"
DOMAIN="hao246.cn"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   NavHub 远程部署工具${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "目标服务器: ${SERVER_USER}@${SERVER_HOST}"
echo "部署路径: ${SERVER_PATH}"
echo "域名: ${DOMAIN}"
echo ""

# 检查 SSH 连接
echo -e "${YELLOW}[1/8] 检查 SSH 连接...${NC}"
if ! ssh -o ConnectTimeout=5 ${SERVER_USER}@${SERVER_HOST} "echo 'SSH连接成功'" 2>/dev/null; then
    echo -e "${RED}❌ SSH 连接失败！请确保：${NC}"
    echo "  1. 服务器 IP 地址正确"
    echo "  2. SSH 密钥已配置"
    echo "  3. 网络连接正常"
    exit 1
fi
echo -e "${GREEN}✅ SSH 连接正常${NC}"
echo ""

# 创建远程目录结构
echo -e "${YELLOW}[2/8] 创建远程目录结构...${NC}"
ssh ${SERVER_USER}@${SERVER_HOST} << 'ENDSSH'
set -e

# 创建主目录
mkdir -p /data/web

# 创建子目录
mkdir -p /data/web/frontend
mkdir -p /data/web/backend
mkdir -p /data/web/logs
mkdir -p /data/web/ssl
mkdir -p /data/web/scripts

# 创建日志文件
touch /data/web/logs/frontend.log
touch /data/web/logs/backend.log
touch /data/web/logs/nginx-error.log

echo "目录结构创建完成"
ENDSSH
echo -e "${GREEN}✅ 目录结构创建完成${NC}"
echo ""

# 构建前端
echo -e "${YELLOW}[3/8] 构建前端...${NC}"
cd /Users/mac-new/work/navhub/frontend

# 安装依赖（如果需要）
if [ ! -d "node_modules" ]; then
    echo "安装前端依赖..."
    npm install
fi

# 构建生产版本
echo "构建前端生产版本..."
npm run build

if [ ! -d "dist" ]; then
    echo -e "${RED}❌ 前端构建失败！dist 目录不存在${NC}"
    exit 1
fi
echo -e "${GREEN}✅ 前端构建完成${NC}"
echo ""

# 构建后端
echo -e "${YELLOW}[4/8] 构建后端...${NC}"
cd /Users/mac-new/work/navhub/backend

# 构建 Go 二进制文件
echo "构建后端二进制文件..."
go build -o navhub-api cmd/server/main.go

if [ ! -f "navhub-api" ]; then
    echo -e "${RED}❌ 后端构建失败！navhub-api 文件不存在${NC}"
    exit 1
fi
echo -e "${GREEN}✅ 后端构建完成${NC}"
echo ""

# 上传前端文件
echo -e "${YELLOW}[5/8] 上传前端文件...${NC}"
cd /Users/mac-new/work/navhub/frontend
echo "正在上传前端文件到 ${SERVER_HOST}:${SERVER_PATH}/frontend ..."
rsync -avz --delete \
    -e "ssh -o StrictHostKeyChecking=no" \
    dist/ \
    ${SERVER_USER}@${SERVER_HOST}:${SERVER_PATH}/frontend/

echo -e "${GREEN}✅ 前端文件上传完成${NC}"
echo ""

# 上传后端文件
echo -e "${YELLOW}[6/8] 上传后端文件...${NC}"
cd /Users/mac-new/work/navhub/backend
echo "正在上传后端文件到 ${SERVER_HOST}:${SERVER_PATH}/backend ..."

# 创建临时目录并上传必要文件
ssh ${SERVER_USER}@${SERVER_HOST} "mkdir -p ${SERVER_PATH}/backend-temp"

# 上传二进制文件
rsync -avz \
    -e "ssh -o StrictHostKeyChecking=no" \
    navhub-api \
    ${SERVER_USER}@${SERVER_HOST}:${SERVER_PATH}/backend-temp/

# 上传配置文件和模板
rsync -avz \
    -e "ssh -o StrictHostKeyChecking=no" \
    .env \
    internal/ \
    ${SERVER_USER}@${SERVER_HOST}:${SERVER_PATH}/backend-temp/ \
    --exclude '*.log' \
    --exclude '*.pid'

echo -e "${GREEN}✅ 后端文件上传完成${NC}"
echo ""

# 部署后端（在远程服务器执行）
echo -e "${YELLOW}[7/8] 部署后端服务...${NC}"
ssh ${SERVER_USER}@${SERVER_HOST} << ENDSSH
set -e

SERVER_PATH="/data/web"

# 停止现有服务
echo "停止现有后端服务..."
if [ -f "\${SERVER_PATH}/navhub-api.pid" ]; then
    PID=\$(cat \${SERVER_PATH}/navhub-api.pid)
    if ps -p \$PID > /dev/null 2>&1; then
        kill \$PID
        sleep 2
    fi
fi

# 强制杀死可能存在的进程
pkill -f "navhub-api" || true
sleep 1

# 备份旧版本
if [ -d "\${SERVER_PATH}/backend" ]; then
    mv \${SERVER_PATH}/backend \${SERVER_PATH}/backend-backup-\$(date +%Y%m%d-%H%M%S) || true
fi

# 部署新版本
mv \${SERVER_PATH}/backend-temp \${SERVER_PATH}/backend

# 设置权限
chmod +x \${SERVER_PATH}/backend/navhub-api

# 创建 systemd 服务文件
cat > /etc/systemd/system/navhub-api.service << 'EOF'
[Unit]
Description=NavHub API Service
After=network.target postgresql.service

[Service]
Type=simple
User=root
WorkingDirectory=/data/web/backend
Environment="ENVIRONMENT=production"
Environment="PORT=8080"
ExecStart=/data/web/backend/navhub-api
Restart=always
RestartSec=10
StandardOutput=append:/data/web/logs/backend.log
StandardError=append:/data/web/logs/backend.log

[Install]
WantedBy=multi-user.target
EOF

# 重新加载 systemd
systemctl daemon-reload

# 启动服务
systemctl restart navhub-api
systemctl enable navhub-api

# 等待服务启动
sleep 3

# 检查服务状态
if systemctl is-active --quiet navhub-api; then
    echo "✅ 后端服务启动成功"
else
    echo "❌ 后端服务启动失败"
    journalctl -u navhub-api -n 20 --no-pager
    exit 1
fi

ENDSSH
echo -e "${GREEN}✅ 后端服务部署完成${NC}"
echo ""

# 配置 Nginx
echo -e "${YELLOW}[8/8] 配置 Nginx...${NC}"
ssh ${SERVER_USER}@${SERVER_HOST} << ENDSSH
set -e

DOMAIN="hao246.cn"

# 检查 Nginx 是否安装
if ! command -v nginx &> /dev/null; then
    echo "安装 Nginx..."
    yum install -y nginx || apt-get install -y nginx
fi

# 创建 Nginx 配置
cat > /etc/nginx/conf.d/navhub.conf << EOF
# NavHub 配置文件
# 域名: ${DOMAIN}

# 前端配置
server {
    listen 80;
    server_name ${DOMAIN} www.${DOMAIN};

    # 前端静态文件
    location / {
        root /data/web/frontend;
        try_files \$uri \$uri/ /index.html;

        # 缓存配置
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # 后端 API
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;

        # CORS 配置
        add_header Access-Control-Allow-Origin *;
        add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS";
        add_header Access-Control-Allow-Headers "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization";
        add_header Access-Control-Expose-Headers "Content-Length,Content-Range";

        if (\$request_method = 'OPTIONS') {
            return 204;
        }
    }

    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml+rss application/json application/javascript;

    # 日志
    access_log /data/web/logs/nginx-access.log;
    error_log /data/web/logs/nginx-error.log;
}
EOF

# 测试 Nginx 配置
nginx -t

# 重启 Nginx
systemctl restart nginx
systemctl enable nginx

echo "✅ Nginx 配置完成"
ENDSSH
echo -e "${GREEN}✅ Nginx 配置完成${NC}"
echo ""

# 部署完成
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   🎉 部署完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "🌐 访问地址："
echo -e "   ${GREEN}http://${DOMAIN}${NC}"
echo -e "   ${GREEN}http://www.${DOMAIN}${NC}"
echo ""
echo -e "🔧 管理命令："
echo -e "   查看后端日志: ssh ${SERVER_USER}@${SERVER_HOST} 'tail -f ${SERVER_PATH}/logs/backend.log'"
echo -e "   查看前端日志: ssh ${SERVER_USER}@${SERVER_HOST} 'tail -f ${SERVER_PATH}/logs/nginx-error.log'"
echo -e "   重启后端: ssh ${SERVER_USER}@${SERVER_HOST} 'systemctl restart navhub-api'"
echo -e "   重启Nginx: ssh ${SERVER_USER}@${SERVER_HOST} 'systemctl restart nginx'"
echo ""
echo -e "📊 服务状态："
echo -e "   后端健康检查: ${GREEN}http://${DOMAIN}/api/v1/health${NC}"
echo ""
