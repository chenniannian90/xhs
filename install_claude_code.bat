@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ========================================
echo   Claude Code 一键安装配置脚本
echo ========================================
echo.

:: 前置条件：阿里云 Token 购买
echo ========================================
echo   第一步：购买阿里云 Token
echo ========================================
echo.
echo 使用 Claude Code 需要 API Token
echo 国内用户推荐使用阿里云，价格优惠！
echo.
echo 🔥 推荐套餐：40元包月 = 50万 tokens
echo.
echo 阿里云购买地址：
echo https://www.aliyun.com/benefit/scene/codingplan
echo.
echo 如果已有 API 密钥，可以直接继续
echo.
set /p HAS_TOKEN="是否已购买或已有 API 密钥？(Y/N): "

if /i not "%HAS_TOKEN%"=="Y" (
    echo.
    echo 请先购买阿里云 Token 后再运行此脚本
    echo.
    start https://www.aliyun.com/benefit/scene/codingplan
    echo.
    pause
    exit /b 0
)

echo.
echo ========================================
echo   第二步：安装 Claude Code
echo ========================================
echo.

:: 检查 Node.js 和 npm 是否已安装
echo [1/6] 检查 Node.js 和 npm...
where node >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 未检测到 Node.js
    echo.
    echo 正在打开 Node.js 官网下载页面...
    echo 请下载并安装 LTS 版本的 Node.js
    echo 下载地址: https://nodejs.org/
    start https://nodejs.org/
    echo.
    echo 安装完成后，请重新运行此脚本
    pause
    exit /b 1
)

where npm >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 未检测到 npm
    pause
    exit /b 1
)

for /f "tokens=*" %%i in ('node --version') do set NODE_VERSION=%%i
for /f "tokens=*" %%i in ('npm --version') do set NPM_VERSION=%%i

echo ✅ 已安装 Node.js: %NODE_VERSION%
echo ✅ 已安装 npm: %NPM_VERSION%
echo.

:: 检查 npm 版本
echo [2/6] 检查 npm 版本...
for /f "tokens=1,2 delims=." %%a in ("%NPM_VERSION%") do (
    set MAJOR=%%a
    set MINOR=%%b
)

set /a VERSION_CHECK=100 * %MAJOR% + %MINOR%
set /a REQUIRED_VERSION=900

if %VERSION_CHECK% lss %REQUIRED_VERSION% (
    echo ⚠️  npm 版本过低，需要升级
    echo 当前版本: %NPM_VERSION%，需要版本: 9.0.0 或以上
    echo.
    echo 正在升级 npm...
    call npm install -g npm
    echo ✅ npm 升级完成
    for /f "tokens=*" %%i in ('npm --version') do set NPM_VERSION=%%i
    echo 新版本: %NPM_VERSION%
    echo.
) else (
    echo ✅ npm 版本符合要求 (9.0.0+)
    echo.
)

:: 安装 Claude Code
echo [3/6] 安装 Claude Code...
echo.
echo 正在执行: npm install -g @anthropic-ai/claude-code
echo.
call npm install -g @anthropic-ai/claude-code

if %errorlevel% neq 0 (
    echo.
    echo ❌ 安装失败！
    echo.
    echo 可能的原因：
    echo 1. 网络连接问题
    echo 2. npm 源访问受限
    echo.
    echo 尝试使用国内镜像源？
    set /p USE_MIRROR="是否使用淘宝镜像安装？(Y/N): "
    if /i "%USE_MIRROR%"=="Y" (
        echo.
        echo 正在使用淘宝镜像源...
        call npm install -g @anthropic-ai/claude-code --registry=https://registry.npmmirror.com
        if %errorlevel% neq 0 (
            echo ❌ 镜像源安装也失败了
            pause
            exit /b 1
        )
    ) else (
        pause
        exit /b 1
    )
)

echo.
echo ✅ Claude Code 安装成功！
echo.

:: 验证安装
echo [4/6] 验证安装...
where claude >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  命令未添加到 PATH
    echo.
    echo 需要手动添加 npm 全局路径到系统环境变量
    echo.
    for /f "tokens=*" %%i in ('npm config get prefix') do set NPM_PREFIX=%%i
    echo npm 全局路径: %NPM_PREFIX%
    echo.
    echo 请按照以下步骤操作：
    echo 1. 右键"此电脑" -> 属性 -> 高级系统设置
    echo 2. 点击"环境变量"
    echo 3. 在"系统变量"中找到 Path，点击"编辑"
    echo 4. 点击"新建"，添加: %NPM_PREFIX%
    echo 5. 确定保存，关闭所有终端窗口后重新打开
    echo.
    set NEED_PATH_FIX=1
) else (
    echo ✅ Claude Code 命令可用
    echo.
)

:: 配置 API 密钥
echo [5/6] 配置 Claude Code...
echo.

:: 创建配置目录
if not exist "%USERPROFILE%\.claude" (
    echo 创建配置目录: %USERPROFILE%\.claude
    mkdir "%USERPROFILE%\.claude"
)

:: 询问 API 密钥
echo.
echo ========================================
echo   配置 API 密钥
echo ========================================
echo.
set /p API_KEY="请输入你的 API 密钥: "
echo.
echo 阿里云 BASE_URL（默认值，直接回车使用）:
echo https://coding.dashscope.aliyuncs.com/apps/anthropic
echo.
set /p BASE_URL="请输入 API 基础 URL（非阿里云请填写，直接回车使用阿里云默认值）: "

if "!BASE_URL!"=="" (
    echo.
    echo 使用阿里云默认 BASE_URL...
    (
        echo {
        echo   "env": {
        echo     "ANTHROPIC_API_KEY": "!API_KEY!",
        echo     "ANTHROPIC_BASE_URL": "https://coding.dashscope.aliyuncs.com/apps/anthropic"
        echo   }
        echo }
    ) > "%USERPROFILE%\.claude\settings.json"
) else (
    echo.
    echo 使用自定义 BASE_URL...
    (
        echo {
        echo   "env": {
        echo     "ANTHROPIC_API_KEY": "!API_KEY!",
        echo     "ANTHROPIC_BASE_URL": "!BASE_URL!"
        echo   }
        echo }
    ) > "%USERPROFILE%\.claude\settings.json"
)
echo.
echo ✅ 配置文件已创建

:: 配置跳过新手引导
echo [6/6] 配置跳过新手引导...
echo.
echo 正在配置...

:: 备份现有配置文件（如果存在）
if exist "%USERPROFILE%\.claude.json" (
    copy "%USERPROFILE%\.claude.json" "%USERPROFILE%\.claude.json.backup" >nul
    echo ✅ 已备份现有配置到 %USERPROFILE%\.claude.json.backup
)

:: 创建新的配置文件
(
    echo {
    echo   "hasCompletedOnboarding": true
    echo }
) > "%USERPROFILE%\.claude.json"

echo ✅ 已配置跳过新手引导

:: 显示使用说明
echo.
echo ========================================
echo   安装配置完成！
echo ========================================
echo.

if defined NEED_PATH_FIX (
    echo ⚠️  请先完成环境变量配置后再使用
    echo.
)

echo 配置文件：
echo - API 配置: %USERPROFILE%\.claude\settings.json
echo - 引导配置: %USERPROFILE%\.claude.json
echo.

echo 使用方法：
if defined NEED_PATH_FIX (
    echo 1. 先完成环境变量配置（见上方说明）
    echo 2. 关闭所有终端窗口，重新打开
    echo 3. 进入任意目录
    echo 4. 双击运行 start_claude_code.bat
) else (
    echo 1. 进入任意目录（项目目录或普通文件夹均可）
    echo 2. 双击运行 start_claude_code.bat
    echo 3. 或在终端输入: claude
)

echo.
echo 遇到问题？
echo - 查看文档: https://docs.anthropic.com
echo - 社区支持: https://github.com/anthropics/claude-code
echo.
echo 安装成功！祝你使用愉快！🎉
echo.
pause
