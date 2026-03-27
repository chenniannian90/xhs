# NavHub 自动化测试和部署指南

## 📋 目录

1. [快速开���](#快速开始)
2. [自动化测试](#自动化测试)
3. [部署流程](#部署流程)
4. [常见问题](#常见问题)
5. [最佳实践](#最佳实践)

---

## 🚀 快速开始

### 一键部署（推荐）

```bash
# 从项目根目录执行
./deploy/deploy-all.sh
```

这个脚本会自动：
1. ✅ 运行部署前检查
2. ✅ 部署后端
3. ✅ 部署前端
4. ✅ 运行完整测试
5. ✅ 生成部署报告

### 分步部署

```bash
# 1. 部署前检查
./deploy/pre-deploy-check.sh

# 2. 部署后端
./deploy/deploy-backend.sh

# 3. 部署前端
./deploy/deploy-frontend.sh
```

---

## 🧪 自动化测试

### 后端 API 测试

**测试脚本**：`backend/scripts/test_api.sh`

**运行方式**：

```bash
# 测试本地后端
cd backend
API_BASE_URL=http://localhost:8080 ./scripts/test_api.sh

# 测试生产环境
API_BASE_URL=http://117.72.39.169:8080 ./scripts/test_api.sh
```

**测试内容**：

| 测试项 | 描述 | 验证点 |
|--------|------|--------|
| 健康检查 | GET /health | 返回 200 |
| 用户注册 | POST /api/v1/auth/register | 返回 201 |
| 用户登录 | POST /api/v1/auth/login | 返回 token |
| 错误处理 | ���误密码、不存在用户 | 返回 401 |
| 邮箱已注册 | 重复注册 | 返回 409 |
| 密码规则 | 弱密码验证 | 返回 400 |
| API 路径 | 检查路径重复 | 返回 404 |

**示例输出**：

```
======================================
NavHub Backend API 测试
======================================

测试 1: 健康检查 ... ✓ PASSED (HTTP 200)
测试 2: 注册新用户 ... ✓ PASSED (HTTP 201)
测试 3: 用户登录 ... ✓ PASSED (HTTP 200)
...
======================================
✓ 所有测试通过！
======================================
```

### 前端测试

**TypeScript 类型检查**：

```bash
cd frontend
npm run type-check
```

**代码格式检查**：

```bash
cd frontend
npm run lint
```

**自动修复**：

```bash
cd frontend
npm run lint:fix
```

---

## 🚢 部署流程

### 部署前检查

**脚本**：`deploy/pre-deploy-check.sh`

**检查项**：

1. **环境检查**
   - Node.js / npm
   - Go 编译器
   - rsync 工具

2. **后端检查**
   - go.mod 存在
   - main.go 存在
   - 配置文件

3. **前端检查**
   - package.json 存在
   - .env 配置正确
   - **API 路径配置检查** ⚠️

4. **代码质量**
   - TypeScript 编译
   - Go fmt/vet

5. **Git 状态**
   - 未提交的文件
   - 当前分支

6. **构建测试**
   - 前端构建
   - 后端构建

**运行**：

```bash
./deploy/pre-deploy-check.sh
```

### 后端部署

**脚本**：`deploy/deploy-backend.sh`

**部署步骤**：

1. 前置检查
2. 代码质量检查（gofmt, go vet）
3. 运行单元测试
4. 构建（Linux x86_64）
5. 备份当前版本
6. 上传到服务器
7. 重启服务
8. 健康检查

**运行**：

```bash
./deploy/deploy-backend.sh
```

### 前端部署

**脚本**：`deploy/deploy-frontend.sh`

**部署步骤**：

1. 前置检查
2. 环境配置检查
3. **API 路径重复检查** ⚠️
4. TypeScript 类型检查
5. 运行测试
6. 构建生产版本
7. 备份当前版本
8. 上传到服务器
9. 验证部署

**运行**：

```bash
./deploy/deploy-frontend.sh
```

---

## ❓ 常见问题

### Q1: 为什么频繁出现 API 路径问题？

**原因**：

```
baseURL = /api/v1 (来自 .env)
代码中 = api.post('/api/v1/auth/login')
最终 URL = /api/v1 + /api/v1/auth/login = /api/v1/api/v1/auth/login ❌
```

**解决**：

```typescript
// 正确写法
api.post('/auth/login')  // ✅

// 错误写法
api.post('/api/v1/auth/login')  // ❌
```

**自动检测**：

部署前检查脚本会自动检测这个问题！

### Q2: 如何回滚到上一个版本？

**后端回滚**：

```bash
ssh root@117.72.39.169

# 查看备份
ls -lh /data/web/backend/navhub-api.backup.*

# 恢复备份
cp /data/web/backend/navhub-api.backup.20260327_100000 /data/web/backend/navhub-api

# 重启服务
systemctl restart navhub-api
```

**前端回滚**：

```bash
ssh root@117.72.39.169

# 查看备份
ls -ld /data/web/frontend.backup.*

# 恢复备份
rm -rf /data/web/frontend
cp -r /data/web/frontend.backup.20260327_100000 /data/web/frontend
```

### Q3: 测试失败了怎么办？

1. **查看详细错误**：

```bash
cd backend
API_BASE_URL=http://117.72.39.169:8080 ./scripts/test_api.sh
```

2. **检查服务状态**：

```bash
ssh root@117.72.39.169 'systemctl status navhub-api'
```

3. **查看日志**：

```bash
ssh root@117.72.39.169 'journalctl -u navhub-api -n 50'
```

### Q4: 如何配置不同的部署环境？

创建不同的环境配置文件：

```bash
# 生产环境
export DEPLOY_SERVER=root@117.72.39.169

# 测试环境
export DEPLOY_SERVER=root@test-server.com
```

然后运行部署脚本：

```bash
./deploy/deploy-all.sh
```

---

## ✅ 最佳实践

### 1. 开发流程

```
1. 创建功能分支
   git checkout -b feature/new-feature

2. 本地开发
   npm run dev

3. 运行测试
   npm run type-check
   npm run lint

4. 提交代码
   git add .
   git commit -m "feat: add new feature"

5. 部署前检查
   npm run pre-deploy

6. 部署到生产
   npm run deploy:all
```

### 2. 代码审查清单

- [ ] TypeScript 编译无错误
- [ ] ESLint 检查通过
- [ ] API 测试通过
- [ ] 手动功能测试
- [ ] 更新文档

### 3. 部署前检查清单

- [ ] 运行 `./deploy/pre-deploy-check.sh`
- [ ] 所有检查通过
- [ ] 查看检查报告
- [ ] 确认可以安全部署

### 4. 部署后验证

- [ ] 前端页面可访问
- [ ] 登录功能正常
- [ ] 注册功能正常
- [ ] API 健康检查通过
- [ ] 浏览器控制台无错误

### 5. 监控和维护

```bash
# 定期查看服务状态
ssh root@117.72.39.169 'systemctl status navhub-api'

# 查看错误日志
ssh root@117.72.39.169 'journalctl -u navhub-api -p err -n 50'

# 查看访问日志
ssh root@117.72.39.169 'tail -f /var/log/nginx/access.log'
```

---

## 📚 相关文档

- [问题根因分析](./PROBLEM_ROOT_CAUSE_ANALYSIS.md)
- [API 路径修复](./API_PATH_FIX.md)
- [注册 UX 改进](./REGISTRATION_UX_IMPROVEMENT.md)
- [部署成功文档](./DEPLOYMENT_SUCCESS.md)

---

## 🎯 总结

通过这套自动化测试和部署系统，你可以：

1. ✅ **自动检测问题** - 在部署前发现配置和代码问题
2. ✅ **自动运行测试** - 确保所有 API 端点正常工作
3. ✅ **一键部署** - 简化部署流程，减少人为错误
4. ✅ **快速回滚** - 自动备份，出现问题时快速恢复
5. ✅ **部署报告** - 每次部署都生成详细报告

**开始使用**：

```bash
# 第一次使用，先运行检查
./deploy/pre-deploy-check.sh

# 然后一键部署
./deploy/deploy-all.sh
```

---

**创建时间**: 2026-03-27
**维护者**: NavHub Team
**状态**: ✅ 已实现并测试
