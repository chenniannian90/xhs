#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Claude Code 一键安装配置脚本${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 前置条件：阿里云 Token 购买
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}  第一步：购买阿里云 Token${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""
echo "使用 Claude Code 需要 API Token"
echo "国内用户推荐使用阿里云，价格优惠！"
echo ""
echo -e "${GREEN}🔥 推荐套餐：40元包月 = 50万 tokens${NC}"
echo ""
echo "阿里云购买地址："
echo "https://www.aliyun.com/benefit/scene/codingplan"
echo ""
echo "如果已有 API 密钥，可以直接继续"
echo ""
read -p "是否已购买或已有 API 密钥？(y/n): " has_token

if [[ $has_token != "y" && $has_token != "Y" ]]; then
    echo ""
    echo "请先购买阿里云 Token 后再运行此脚本"
    echo ""
    # Mac 使用 open，Linux 使用 xdg-open
    if [[ "$OSTYPE" == "darwin"* ]]; then
        open https://www.aliyun.com/benefit/scene/codingplan
    else
        xdg-open https://www.aliyun.com/benefit/scene/codingplan 2>/dev/null || \
        echo "请手动打开浏览器访问: https://www.aliyun.com/benefit/scene/codingplan"
    fi
    echo ""
    exit 0
fi

echo ""
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}  第二步：安装 Claude Code${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# 检查 Node.js 和 npm
echo -e "${YELLOW}[1/6] 检查 Node.js 和 npm...${NC}"

if ! command -v node &> /dev/null; then
    echo -e "${RED}❌ 未检测到 Node.js${NC}"
    echo ""
    echo "请先安装 Node.js："
    echo "Mac: brew install node"
    echo "Linux: 使用 nvm 安装 (下方提供命令)"
    echo ""
    read -p "是否现在安装 Node.js？(y/n): " install_node

    if [[ $install_node == "y" || $install_node == "Y" ]]; then
        echo ""
        echo "正在安装 nvm..."
        curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

        if [[ $SHELL == "/bin/zsh" ]]; then
            echo 'export NVM_DIR="$HOME/.nvm"' >> ~/.zshrc
            echo '[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"' >> ~/.zshrc
            source ~/.zshrc
        else
            echo 'export NVM_DIR="$HOME/.nvm"' >> ~/.bashrc
            echo '[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"' >> ~/.bashrc
            source ~/.bashrc
        fi

        echo "正在安装 Node.js LTS..."
        nvm install --lts
        nvm use --lts
        nvm alias default lts/*
    else
        echo "请先安装 Node.js 后重新运行此脚本"
        exit 1
    fi
fi

if ! command -v npm &> /dev/null; then
    echo -e "${RED}❌ 未检测到 npm${NC}"
    exit 1
fi

NODE_VERSION=$(node --version)
NPM_VERSION=$(npm --version)

echo -e "${GREEN}✅ 已安装 Node.js: ${NODE_VERSION}${NC}"
echo -e "${GREEN}✅ 已安装 npm: ${NPM_VERSION}${NC}"
echo ""

# 检查 npm 版本
echo -e "${YELLOW}[2/6] 检查 npm 版本...${NC}"

NPM_MAJOR=$(echo $NPM_VERSION | cut -d. -f1)
NPM_MINOR=$(echo $NPM_VERSION | cut -d. -f2)
VERSION_CHECK=$((NPM_MAJOR * 100 + NPM_MINOR))
REQUIRED_VERSION=900

if [ $VERSION_CHECK -lt $REQUIRED_VERSION ]; then
    echo -e "${YELLOW}⚠️  npm 版本过低，正在升级...${NC}"
    echo "当前版本: $NPM_VERSION，需要版本: 9.0.0 或以上"
    echo ""
    echo "正在升级 npm..."
    npm install -g npm
    NPM_VERSION=$(npm --version)
    echo -e "${GREEN}✅ npm 升级完成，新版本: ${NPM_VERSION}${NC}"
    echo ""
else
    echo -e "${GREEN}✅ npm 版本符合要求 (9.0.0+)${NC}"
    echo ""
fi

# 安装 Claude Code
echo -e "${YELLOW}[3/6] 安装 Claude Code...${NC}"
echo ""
echo "正在执行: npm install -g @anthropic-ai/claude-code"
echo ""

if ! npm install -g @anthropic-ai/claude-code; then
    echo ""
    echo -e "${RED}❌ 安装失败！${NC}"
    echo ""
    echo "尝试使用国内镜像源？"
    read -p "是否使用淘宝镜像源？(y/n): " use_mirror

    if [[ $use_mirror == "y" || $use_mirror == "Y" ]]; then
        echo ""
        echo "正在使用淘宝镜像源..."
        npm install -g @anthropic-ai/claude-code --registry=https://registry.npmmirror.com

        if [ $? -ne 0 ]; then
            echo -e "${RED}❌ 镜像源安装也失败了${NC}"
            exit 1
        fi
    else
        exit 1
    fi
fi

echo ""
echo -e "${GREEN}✅ Claude Code 安装成功！${NC}"
echo ""

# 验证安装
echo -e "${YELLOW}[4/6] 验证安装...${NC}"

if ! command -v claude &> /dev/null; then
    echo -e "${YELLOW}⚠️  命令未添加到 PATH${NC}"
    echo ""

    NPM_PREFIX=$(npm config get prefix)
    echo "npm 全局路径: $NPM_PREFIX"
    echo ""

    if [[ "$SHELL" == "/bin/zsh" ]]; then
        echo "请将以下内容添加到 ~/.zshrc："
        echo "export PATH=\"$NPM_PREFIX/bin:\$PATH\""
        echo ""
        echo "执行: source ~/.zshrc"
    else
        echo "请将以下内容添加到 ~/.bashrc："
        echo "export PATH=\"$NPM_PREFIX/bin:\$PATH\""
        echo ""
        echo "执行: source ~/.bashrc"
    fi

    NEED_PATH_FIX=1
else
    echo -e "${GREEN}✅ Claude Code 命令可用${NC}"
    echo ""
fi

# 配置 API 密钥
echo -e "${YELLOW}[5/6] 配置 Claude Code...${NC}"
echo ""

# 创建配置目录
mkdir -p ~/.claude

echo "========================================"
echo "  配置 API 密钥"
echo "========================================"
echo ""

read -p "请输入你的 API 密钥: " api_key
echo ""
echo "阿里云 BASE_URL（默认值，直接回车使用）:"
echo "https://coding.dashscope.aliyuncs.com/apps/anthropic"
echo ""
read -p "请输入 API 基础 URL（非阿里云请填写，直接回车使用阿里云默认值）: " base_url

if [[ -z "$base_url" ]]; then
    echo ""
    echo "使用阿里云默认 BASE_URL..."
    cat > ~/.claude/settings.json << EOF
{
  "env": {
    "ANTHROPIC_API_KEY": "$api_key",
    "ANTHROPIC_BASE_URL": "https://coding.dashscope.aliyuncs.com/apps/anthropic"
  }
}
EOF
else
    echo ""
    echo "使用自定义 BASE_URL..."
    cat > ~/.claude/settings.json << EOF
{
  "env": {
    "ANTHROPIC_API_KEY": "$api_key",
    "ANTHROPIC_BASE_URL": "$base_url"
  }
}
EOF
fi

echo ""
echo -e "${GREEN}✅ 配置文件已创建${NC}"

# 配置跳过新手引导
echo ""
echo -e "${YELLOW}[6/6] 配置跳过新手引导...${NC}"
echo ""
echo "正在配置..."

# 备份现有配置文件（如果存在）
if [[ -f ~/.claude.json ]]; then
    cp ~/.claude.json ~/.claude.json.backup
    echo -e "${GREEN}✅ 已备份现有配置到 ~/.claude.json.backup${NC}"
fi

# 创建新的配置文件
cat > ~/.claude.json << 'EOF'
{
  "hasCompletedOnboarding": true
}
EOF

echo -e "${GREEN}✅ 已配置跳过新手引导${NC}"

# 显示使用说明
echo ""
echo "========================================"
echo -e "${GREEN}  安装配置完成！${NC}"
echo "========================================"
echo ""

if [ ! -z "$NEED_PATH_FIX" ]; then
    echo -e "${YELLOW}⚠️  请先完成 PATH 配置后再使用${NC}"
    echo ""
fi

echo "配置文件："
echo "- API 配置: ~/.claude/settings.json"
echo "- 引导配置: ~/.claude.json"
echo ""

echo "使用方法："
if [ ! -z "$NEED_PATH_FIX" ]; then
    echo "1. 先完成 PATH 配置（见上方说明）"
    echo "2. 执行: source ~/.bashrc 或 source ~/.zshrc"
    echo "3. 进入你的项目目录"
    echo "4. 输入: claude"
else
    echo "1. 进入你的项目目录"
    echo "2. 输入: claude"
fi

echo ""

echo "遇到问题？"
echo "- 查看文档: https://docs.anthropic.com"
echo "- 社区支持: https://github.com/anthropics/claude-code"
echo ""

echo -e "${GREEN}安装成功！祝你使用愉快！🎉${NC}"
echo ""
