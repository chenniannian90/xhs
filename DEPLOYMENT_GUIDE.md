# NavHub 远程部署指南

## 📋 概述

NavHub 远程部署工具用于将应用部署到生产服务器：
- **服务器**: root@117.72.39.169
- **部署路径**: /data/web
- **域名**: hao246.cn

## 🚀 部署工具

### 1. 完整部署（推荐首次使用）

```bash
./deploy-remote.sh
```

**功能**:
- ✅ 检查SSH连接
- ✅ 创建远程目录结构
- ✅ 构建��端和后端
- ✅ 上传所有文件
- ✅ 配置systemd服务
- ✅ 配置Nginx反向代理

**部署步骤**:
1. 构建前端生产版本
2. 构建后端Go二进制文件
3. 上传文件到服务器
4. 配置并启动后端服务
5. 配置Nginx

### 2. 快速部署（仅前端）

```bash
./quick-deploy.sh
```

**功能**:
- ✅ 快速构建并上传前端
- ✅ 适合���端代码更新

**使用场景**:
- 修改了前端代码
- 更新了UI界面
- 修复了前端bug

### 3. 远程管理

```bash
./remote-manage.sh [命令]
```

**可用命令**:

| 命令 | 说明 |
|------|------|
| `logs-backend` | 查看后端实时日志 |
| `logs-nginx` | 查看Nginx实时日志 |
| `logs-all` | 查看所有日志 |
| `status` | 查看服务运行状态 |
| `restart-backend` | 重启后端服务 |
| `restart-nginx` | 重启Nginx服务 |
| `restart-all` | 重启所有服务 |
| `ssh` | SSH连接到服务器 |
| `health` | 健康检查 |

**示例**:
```bash
# 查看后端日志
./remote-manage.sh logs-backend

# 查看服务状态
./remote-manage.sh status

# 健康检查
./remote-manage.sh health
```

## 📂 部署结构

```
/data/web/
├── frontend/          # 前端静态文件
│   ├── index.html
│   ├── assets/
│   └── ...
├── backend/           # 后端程序
│   ├── navhub-api     # Go二进制文件
│   ├── internal/      # 代码文件
│   └── .env          # 环境配置
├── logs/              # 日志文件
│   ├── backend.log
│   ├── nginx-access.log
│   └── nginx-error.log
├── ssl/               # SSL证书（可选）
└── scripts/           # 脚本文件
```

## 🔧 配置文件

部署配置位于 `.deploy.config`:

```bash
# 服务器配置
SERVER_HOST=117.72.39.169
SERVER_USER=root
SERVER_PATH=/data/web

# 域名配置
DOMAIN=hao246.cn
```

## 🌐 访问地址

部署完成后，可通过以下地址访问：

- **主站**: http://hao246.cn
- **带WWW**: http://www.hao246.cn
- **API**: http://hao246.cn/api/v1/
- **健康检查**: http://hao246.cn/health

## 🔍 故障排查

### 1. 检查服务状态

```bash
./remote-manage.sh status
```

### 2. 查看日志

```bash
# 后端日志
./remote-manage.sh logs-backend

# Nginx日志
./remote-manage.sh logs-nginx
```

### 3. 重启服务

```bash
# 重启后端
./remote-manage.sh restart-backend

# 重启所有
./remote-manage.sh restart-all
```

### 4. SSH登录服务器

```bash
./remote-manage.sh ssh
```

## 📊 监控和维护

### 日常检查

1. **健康检查**
   ```bash
   ./remote-manage.sh health
   ```

2. **查看日志**
   ```bash
   ./remote-manage.sh logs-backend
   ```

3. **检查磁盘空间**
   ```bash
   ssh root@117.72.39.169 'df -h'
   ```

### 备份

重要数据需要定期备份：

```bash
# 备份数据库
ssh root@117.72.39.169 'pg_dump navhub > backup.sql'

# 备份配置文件
ssh root@117.72.39.169 'tar -czf backup.tar.gz /data/web'
```

## 🔐 安全建议

1. **使用SSH密钥认证**
2. **定期更新系统**
3. **配置防火墙**
4. **启用HTTPS**（推荐使用Let's Encrypt）
5. **定期备份数据**

## 📈 性能优化

1. **启用Gzip压缩**（已配置）
2. **配置静态资源缓存**（已配置）
3. **使用CDN**（可选）
4. **数据库索引优化**

## 🆕 升级部署

当需要更新应用时：

1. **仅前端更新**:
   ```bash
   ./quick-deploy.sh
   ```

2. **完整更新**:
   ```bash
   ./deploy-remote.sh
   ```

3. **更新后重启**:
   ```bash
   ./remote-manage.sh restart-all
   ```

## 📝 注意事项

1. 首次部署需要配置SSH密钥
2. 确保服务器已安装Go和Node.js环境
3. 数据库需要单独配置（PostgreSQL + Redis）
4. 建议在测试环境先验证

## 🆘 获取帮助

如果遇到问题：

1. 查看日志文件
2. 检查服务状态
3. SSH到服务器手动排查
4. 查看Nginx错误日志

---

**最后更新**: 2026-03-26
**维护者**: NavHub Team
