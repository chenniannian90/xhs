#!/bin/bash

# 测试SSH连接并显示服务器信息

# 配置
SERVER_HOST="117.72.39.169"
SERVER_USER="root"

# 颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}测试SSH连接到 ${SERVER_USER}@${SERVER_HOST}...${NC}"
echo ""

# 测试连接
if ssh -o ConnectTimeout=10 ${SERVER_USER}@${SERVER_HOST} "echo '✅ SSH连接成功'" 2>/dev/null; then
    echo -e "${GREEN}✅ SSH连接正常${NC}"
    echo ""

    # 获取服务器信息
    echo "服务器信息:"
    ssh ${SERVER_USER}@${SERVER_HOST} << 'ENDSSH'
echo "=== 系统信息 ==="
hostnamectl | grep -E 'Operating System|Kernel'
echo ""
echo "=== 资源使用 ==="
free -h | head -2
echo ""
df -h | grep -E '^/dev/'
echo ""
echo "=== CPU ==="
top -bn1 | head -5
ENDSSH

    echo ""
    echo -e "${GREEN}✅ 服务器连接测试完成${NC}"
    echo "可以开始部署了！使用: ./deploy-remote.sh"

else
    echo -e "${RED}❌ SSH连接失败${NC}"
    echo ""
    echo "请检查："
    echo "1. 服务器IP地址是否正确: ${SERVER_HOST}"
    echo "2. SSH密钥是否已配置"
    echo "3. 网络连接是否正常"
    echo ""
    echo "配置SSH密钥:"
    echo "  ssh-copy-id ${SERVER_USER}@${SERVER_HOST}"
fi
