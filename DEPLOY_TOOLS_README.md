# NavHub 部署工具集

完整的远程部署解决方案，用于将 NavHub 部署到生产服务器。

## 🎯 快速开始

### 1. 测试连接

```bash
./test-connection.sh
```

### 2. 完整部署

```bash
./deploy-remote.sh
```

### 3. 访问网站

```
http://hao246.cn
```

## 📦 部署工具

### 完整部署 - `deploy-remote.sh`

一键完成整个部署流程，包括：
- ✅ 构建前端和后端
- ✅ 上传所有文件
- ✅ 配置systemd服务
- ✅ 配置Nginx反向代理
- ✅ 启动所有服务

**使用**:
```bash
./deploy-remote.sh
```

### 快速部署 - `quick-deploy.sh`

仅部署前端，适合快速更新：

**使用**:
```bash
./quick-deploy.sh
```

**场景**:
- 前端代码修改
- UI更新
- 前端bug修复

### 远程管理 - `remote-manage.sh`

服务器运维管理工具：

**使用**:
```bash
./remote-manage.sh [命令]
```

**可用命令**:
```bash
./remote-manage.sh logs-backend      # 查看后端日志
./remote-manage.sh logs-nginx        # 查看Nginx日志
./remote-manage.sh status            # 查看服务状态
./remote-manage.sh restart-backend   # 重启后端
./remote-manage.sh restart-nginx     # 重启Nginx
./remote-manage.sh health            # 健康检查
./remote-manage.sh ssh               # SSH连接
```

### 连接测试 - `test-connection.sh`

测试SSH连接和服务器状态：

**使用**:
```bash
./test-connection.sh
```

## 🗂️ 文件结构

```
navhub/
├── deploy-remote.sh          # 完整部署脚本
├── quick-deploy.sh           # 快速部署脚本
├── remote-manage.sh          # 远程管理脚本
├── test-connection.sh        # 连接测试脚本
├── .deploy.config            # 部署配置文件
├── DEPLOYMENT_GUIDE.md       # 详细部署指南
└── DEPLOY_TOOLS_README.md    # 本文件
```

## ⚙️ 配置

### 服务器信息

- **服务器**: 117.72.39.169
- **用户**: root
- **路径**: /data/web
- **域名**: hao246.cn

### 修改配置

如需修改部署目标，编辑 `.deploy.config`:

```bash
SERVER_HOST=117.72.39.169
SERVER_USER=root
SERVER_PATH=/data/web
DOMAIN=hao246.cn
```

## 📊 部署后结构

```
/data/web/
├── frontend/              # 前端静态文件
├── backend/               # 后端程序
│   ├── navhub-api        # Go二进制
│   └── internal/         # 代码文件
├── logs/                  # 日志文件
│   ├── backend.log
│   ├── nginx-access.log
│   └── nginx-error.log
└── ssl/                   # SSL证书
```

## 🌐 访问地址

- **主站**: http://hao246.cn
- **API**: http://hao246.cn/api/v1/
- **健康检查**: http://hao246.cn/health

## 🔧 常用操作

### 查看日志

```bash
# 后端日志
./remote-manage.sh logs-backend

# Nginx日志
./remote-manage.sh logs-nginx
```

### 重启服务

```bash
# 重启后端
./remote-manage.sh restart-backend

# 重启Nginx
./remote-manage.sh restart-nginx

# 重启所有
./remote-manage.sh restart-all
```

### 检查状态

```bash
# 服务状态
./remote-manage.sh status

# 健康检查
./remote-manage.sh health
```

## 🚨 故障排查

### 服务无法启动

1. 查看日志: `./remote-manage.sh logs-backend`
2. 检查端口: `ssh root@117.72.39.169 'netstat -tlnp'`
3. 检查配置: `ssh root@117.72.39.169 'systemctl status navhub-api'`

### 前端无法访问

1. 检查Nginx: `./remote-manage.sh restart-nginx`
2. 查看错误: `./remote-manage.sh logs-nginx`
3. 检查文件: `ssh root@117.72.39.169 'ls -la /data/web/frontend'`

### 网站无法访问

1. 健康检查: `./remote-manage.sh health`
2. 检查防火墙: `ssh root@117.72.39.169 'firewall-cmd --list-all'`
3. 重启所有: `./remote-manage.sh restart-all`

## 📈 性能优化

### 已启用

- ✅ Nginx Gzip压缩
- ✅ 静态资源缓存
- ✅ HTTP/2支持（需配置SSL）

### 建议优化

- 🔧 配置HTTPS（Let's Encrypt）
- 🔧 使用CDN加速
- 🔧 数据库连接池优化
- 🔧 启用Redis缓存

## 🔐 安全建议

1. **SSH密钥认证**
   ```bash
   ssh-copy-id root@117.72.39.169
   ```

2. **配置防火墙**
   ```bash
   ssh root@117.72.39.169 'firewall-cmd --add-port=80/tcp --permanent'
   ```

3. **启用HTTPS**
   ```bash
   # 使用Let's Encrypt
   ssh root@117.72.39.169 'certbot --nginx -d hao246.cn'
   ```

4. **定期备份**
   ```bash
   ssh root@117.72.39.169 'pg_dump navhub > backup.sql'
   ```

## 📝 维护计划

### 每日

- 检查服务状态: `./remote-manage.sh health`
- 查看错误日志: `./remote-manage.sh logs-backend`

### 每周

- 检查磁盘空间
- 更新系统安全补丁
- 备份数据库

### 每月

- 检查访问日志
- 优化数据库性能
- 更新依赖包

## 🆘 获取帮助

### 查看详细文档

```bash
cat DEPLOYMENT_GUIDE.md
```

### 查看管理工具帮助

```bash
./remote-manage.sh
```

### SSH连接服务器

```bash
./remote-manage.sh ssh
```

## 📞 支持

如遇到问题：

1. 查看日志文件
2. 运行健康检查
3. SSH到服务器手动排查
4. 查阅详细部署指南

---

**创建时间**: 2026-03-26
**最后更新**: 2026-03-26
**版本**: 1.0.0
