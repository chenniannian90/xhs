#!/bin/bash

# NavHub 快速部署脚本 - 仅部署前端
# 用于快速更新前端静态文件

set -e

# 颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 配置
SERVER_HOST="117.72.39.169"
SERVER_USER="root"
SERVER_PATH="/data/web"

echo -e "${YELLOW}快速部署前端...${NC}"

# 构建前端
cd /Users/mac-new/work/navhub/frontend
echo "构建前端..."
npm run build

# 上传
echo "上传到服务器..."
rsync -avz --delete \
    -e "ssh -o StrictHostKeyChecking=no" \
    dist/ \
    ${SERVER_USER}@${SERVER_HOST}:${SERVER_PATH}/frontend/

echo -e "${GREEN}✅ 前端部署完成！${NC}"
echo "访问: http://hao246.cn"
