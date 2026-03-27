# NavHub 部署命令快速参考

## 🚀 部署命令

### 完整部署
```bash
./deploy-remote.sh
```
首次部署或完整更新，包括前后端和配置

### 快速部署（仅前端）
```bash
./quick-deploy.sh
```
仅更新前端静态文件

## 🛠️ 管理命令

### 查看日志
```bash
./remote-manage.sh logs-backend      # 后端日志
./remote-manage.sh logs-nginx        # Nginx日志
./remote-manage.sh logs-all          # 所有日志
```

### 服务管理
```bash
./remote-manage.sh status            # 查看状态
./remote-manage.sh restart-backend   # 重启后端
./remote-manage.sh restart-nginx     # 重启Nginx
./remote-manage.sh restart-all       # 重启所有
```

### 诊断
```bash
./remote-manage.sh health            # 健康检查
./remote-manage.sh ssh               # SSH连接
./test-connection.sh                 # 测试连接
```

## 🌐 访问地址

- 主站: http://hao246.cn
- API: http://hao246.cn/api/v1/
- 健康检查: http://hao246.cn/health

## 📊 服务器信息

- IP: 117.72.39.169
- 用户: root
- 路径: /data/web
- 域名: hao246.cn

## 📁 文件结构

```
/data/web/
├── frontend/    # 前端静态文件
├── backend/     # 后端程序
└── logs/        # 日志文件
```

## ⚡ 快速操作

### 更新前端代码
```bash
cd frontend
# 修改代码...
npm run build
../quick-deploy.sh
```

### 查看后端错误
```bash
./remote-manage.sh logs-backend
# 按 Ctrl+C 退出
```

### 重启所有服务
```bash
./remote-manage.sh restart-all
```

## 🆘 故障排查

### 网站无法访问
```bash
./remote-manage.sh health
./remote-manage.sh status
```

### 查看错误日志
```bash
./remote-manage.sh logs-backend
./remote-manage.sh logs-nginx
```

### SSH到服务器
```bash
./remote-manage.sh ssh
# 或
ssh root@117.72.39.169
```

## 📚 详细文档

- 完整部署指南: `DEPLOYMENT_GUIDE.md`
- 工具说明: `DEPLOY_TOOLS_README.md`
- 配置文件: `.deploy.config`
