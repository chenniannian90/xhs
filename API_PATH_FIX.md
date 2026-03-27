# API 路径重复问题修复

## 问题描述

登录请求的 URL 出现路径重复：

```
https://117.72.39.169/api/v1/api/v1/auth/login
                     ^^^^^^^^ ^^^^^^^^
```

**错误原因**：
- API baseURL 配置为 `/api/v1`
- API 调用代码中又包含了 `/api/v1` 前缀
- 导致最终 URL 路径重复

## 修复方案

### 前端配置

**保持** `.env` 配置：
```env
VITE_API_URL=/api/v1
```

**修改**所有 API 调用，去掉 `/api/v1` 前缀：

### 修改的文件

1. **`src/store/authStore.ts`**
   ```typescript
   // 修改前
   api.post('/api/v1/auth/login', ...)
   api.post('/api/v1/auth/register', ...)

   // 修改后
   api.post('/auth/login', ...)
   api.post('/auth/register', ...)
   ```

2. **`src/pages/dashboard/DashboardPage.tsx`**
   ```typescript
   // 修改前
   api.get('/api/v1/categories')
   api.get('/api/v1/sites')

   // 修改后
   api.get('/categories')
   api.get('/sites')
   ```

3. **`src/pages/dashboard/CategoriesPage.tsx`**
   ```typescript
   // 修改前
   api.get('/api/v1/categories')
   api.delete(`/api/v1/categories/${id}`)
   api.delete(`/api/v1/categories/${id}/share`)
   api.post(`/api/v1/categories/${id}/share`)

   // 修改后
   api.get('/categories')
   api.delete(`/categories/${id}`)
   api.delete(`/categories/${id}/share`)
   api.post(`/categories/${id}/share`)
   ```

4. **`src/pages/dashboard/SearchPage.tsx`**
   ```typescript
   // 修改前
   api.get('/api/v1/categories')
   api.get('/api/v1/sites')

   // 修改后
   api.get('/categories')
   api.get('/sites')
   ```

5. **`src/pages/auth/ForgotPasswordPage.tsx`**
   ```typescript
   // 修改前
   api.post('/api/v1/auth/forgot-password', ...)

   // 修改后
   api.post('/auth/forgot-password', ...)
   ```

6. **`src/pages/auth/ResetPasswordPage.tsx`**
   ```typescript
   // 修改前
   api.post('/api/v1/auth/reset-password', ...)

   // 修改后
   api.post('/auth/reset-password', ...)
   ```

## 修复后的 URL 示例

**登录**:
```
修改前: https://117.72.39.169/api/v1/api/v1/auth/login ❌
修改后: https://117.72.39.169/api/v1/auth/login ✅
```

**注册**:
```
修改前: https://117.72.39.169/api/v1/api/v1/auth/register ❌
修改后: https://117.72.39.169/api/v1/auth/register ✅
```

**获取分类**:
```
修改前: https://117.72.39.169/api/v1/api/v1/categories ❌
修改后: https://117.72.39.169/api/v1/categories ✅
```

## 部署信息

- **修复时间**: 2026-03-27
- **前端版本**: index-C2kmMqXp.js
- **服务状态**: ✅ 运行中

## 验证测试

访问 **https://117.72.39.169** 并测试：

1. ✅ 登录功能
2. ✅ 注册功能
3. ✅ 忘记密码
4. ✅ 分类管理
5. ✅ 搜索功能

所有 API 请求现在都应该正确路由到后端。

## 架构说明

**正确的 API 调用模式**：

```
环境变量 (.env):
VITE_API_URL=/api/v1

API 客户端 (api.ts):
baseURL = import.meta.env.VITE_API_URL  // = '/api/v1'

API 调用 (组件中):
api.post('/auth/login')  // 不包含 /api/v1

最终 URL:
/api/v1 + /auth/login = /api/v1/auth/login ✅
```

---

**修复完成时间**: 2026-03-27
**状态**: ✅ 生产环境已部署
