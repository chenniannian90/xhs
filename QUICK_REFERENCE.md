# NavHub 快速参考卡片

## 🚀 一键部署

```bash
# 完整部署（推荐）
make deploy
# 或
./deploy/deploy-all.sh

# 快速检查
make check
```

---

## 📋 部署前检查

```bash
# 运行所有检查
make check

# 手动检查
./deploy/pre-deploy-check.sh
```

**检查内容**：
- ✅ 环境配置（Node.js, Go, rsync）
- ✅ 后端文件（go.mod, main.go）
- ✅ 前端文件（package.json, .env）
- ✅ **API 路径配置**（防止重复）
- ✅ TypeScript 类型检查
- ✅ 构建测试

---

## 🧪 运行测试

```bash
# 所有测试
make test

# 仅测试 API
make test-api

# 仅测试前端
make test-frontend
```

---

## 🔧 常用命令

| 命令 | 说明 |
|------|------|
| `make help` | 显示所有可用命令 |
| `make check` | 部署前检查 |
| `make deploy` | 一键部署 |
| `make dev` | 启动前端开发服务器 |
| `make logs` | 查看后端日志 |
| `make status` | 查看服务状态 |
| `make health` | 健康检查 |
| `make clean` | 清理构建文件 |

---

## ⚠️ 常见问题

### API 路径重复

**错误**：`/api/v1/api/v1/auth/login`

**原因**：
```typescript
// .env
VITE_API_URL=/api/v1

// 代码中（错误）
api.post('/api/v1/auth/login')  // ❌
```

**解决**：
```typescript
// 代码中（正确）
api.post('/auth/login')  // ✅
```

**检测**：`make check` 会自动检测！

### 部署失败

```bash
# 1. 查看日志
make logs

# 2. 检查服务状态
make status

# 3. 健康检查
make health

# 4. 回滚
make rollback
```

---

## 🔄 回滚流程

```bash
# 自动回滚
make rollback

# 手动回滚后端
ssh root@117.72.39.169
cp /data/web/backend/navhub-api.backup.* /data/web/backend/navhub-api
systemctl restart navhub-api

# 手动回滚前端
ssh root@117.72.39.169
rm -rf /data/web/frontend
cp -r /data/web/frontend.backup.* /data/web/frontend
```

---

## 📊 监控命令

```bash
# 实时监控
make monitor

# 查看后端日志
make logs

# 查看 Nginx 日志
make logs-nginx

# 服务状态
make status
```

---

## 🌐 访问地址

| 服务 | 地址 |
|------|------|
| 前端 | https://117.72.39.169 |
| API | http://117.72.39.169:8080 |
| 健康检查 | http://117.72.39.169:8080/health |

---

## 📁 重要文件

| 文件 | 用途 |
|------|------|
| `deploy/deploy-all.sh` | 一键部署脚本 |
| `deploy/pre-deploy-check.sh` | 部署前检查 |
| `backend/scripts/test_api.sh` | API 测试 |
| `Makefile` | 快捷命令 |

---

## 💡 最佳实践

1. **部署前总是先检查**
   ```bash
   make check
   ```

2. **使用分支开发**
   ```bash
   git checkout -b feature/new-feature
   ```

3. **测试后再部署**
   ```bash
   make test
   make deploy
   ```

4. **定期查看日志**
   ```bash
   make logs
   ```

---

## 📞 获取帮助

```bash
# 查看所有命令
make help

# 查看详细文档
cat AUTOMATED_TESTING_GUIDE.md
```

---

**最后更新**: 2026-03-27
**版本**: 1.0
