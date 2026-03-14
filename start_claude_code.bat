@echo off
chcp 65001 >nul

echo ========================================
echo   启动 Claude Code
echo ========================================
echo.

:: 显示当前目录信息
echo 📁 当前目录: %cd%
echo.
echo 💡 提示: Claude Code 可以在任何文件夹中使用！
echo    - 项目目录 ✅
echo    - 普通文件夹 ✅
echo    - 文档目录 ✅
echo.

:: 检查 Claude Code 是否已安装
where claude >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Claude Code 未安装
    echo.
    echo 请先运行 install_claude_code.bat 进行安装
    echo.
    pause
    exit /b 1
)

:: 检查配置文件是否存在
if not exist "%USERPROFILE%\.claude\settings.json" (
    echo ⚠️  未检测到配置文件
    echo.
    echo 请先配置 API 密钥
    echo 配置文件位置: %USERPROFILE%\.claude\settings.json
    echo.
    echo 或运行安装脚本进行配置
    echo.
    pause
    exit /b 1
)

:: 启动 Claude Code
echo 正在启动 Claude Code...
echo.
claude

pause
