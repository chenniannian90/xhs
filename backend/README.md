# NavHub Backend API

## 📋 项目概述

NavHub 后端是一个基于 Go + Gin + PostgreSQL + GORM 的 RESTful API 服务。

---

## 🏗️ 系统架构

```
┌─────────────┐
│   Client    │
│  (React)   │
└─────┬────┘
       │
       ▼
  ┌─────┐
  │ Gin   │
  │ Router │
  └──────┘
       │
       ▼
  ┌─────┐
  │ GORM  │
  │  ORM   │
  └──────┘
       │
       ▼
  ┌─────┐
  │ PostgreSQL  │
  │ Database │
  └──────┘
       │
```

---

## 📁 目录结构

```
backend/
├── cmd/server/
│   └── main.go                    ✅ 服务器入口
├── internal/
│   ├── config/
│   │   └── config.go             ✅ 配置管理
│   ├── models/
│   │   ├── user.go                 ✅ 用户模型
│   │   ├── oauth_account.go          ✅ OAuth 模型
│   │   ├── email_verification.go    ✅ 邮箱验证模型
│   │   ├── password_reset.go         ✅ 密码重置模型
│   │   ├── category.go              ✅ 分类模型
│   │   ├── site.go                  ✅ 站点模型
│   │   └── models.go              ✅ 模型导出
│   ├── repositories/
│   │   ├── user_repository.go      ✅ 用户数据访问
│   │   ├── category_repository.go  ✅ 分类数据访问
│   │   └── site_repository.go      ✅ 站点数据访问
│   ├── services/
│   │   ├── auth_service.go         ✅ 认证业务逻辑
│   │   ├── email_service.go         ✅ 邮件服务
│   │   ├── category_service.go      ✅ 分类业务逻辑
│   │   └── site_service.go          ✅ 站点业务逻辑
│   ├── handlers/
│   │   ├── auth_handler.go          ✅ 认证 HTTP 处理
│   │   ├── category_handler.go       ✅ 分类 HTTP 处理
│   │   └── site_handler.go          ✅ 站点 HTTP 处理
│   └── middleware/
│       ├── auth.go                  ✅ JWT 认证中间件
│       ├── cors.go                  ✅ CORS 中间件
│       └── logger.go                ✅ 日志中间件
└── pkg/
    ├── jwt/
    │   └── jwt.go                  ✅ JWT 工具
    └── password/
        └── password.go             ✅ 密码工具
├── .env                          ✅ 环境配置
├── go.mod                         ✅ Go 模块
└── README.md                       ✅ 项目文档
```

---

## 🚀 当前状态

### ✅ 已完成

1. **数据库模型**
   - User（用户表）
   - OAuthAccount（OAuth 账户）
   - EmailVerification（邮箱验证）
   - PasswordReset（密码重置）
   - Category（分类表）
   - Site（站点表）

2. **工具包**
   - JWT Manager（令牌生成和验证）
   - Password（哈希和验证）

3. **Repository 层**
   - UserRepository（用户数据操作）
   - CategoryRepository（分类数据操作）
   - SiteRepository（站点数据操作）

4. **Service 层**
   - AuthService（认证业务逻辑）
   - EmailService（邮件发送）
   - CategoryService（分类业务逻辑）
   - SiteService（站点业务逻辑）

5. **Handler 层**
   - AuthHandler（认证 HTTP 端点）
   - CategoryHandler（分类 HTTP 端点）
   - SiteHandler（站点 HTTP 端点）

6. **中间件**
   - Auth（JWT 验证）
   - CORS（跨域支持）
   - Logger（请求日志）

7. **服务器配置**
   - GIN Router（路由设置）
   - 数据库连接（PostgreSQL）

---

## 🔌 API 端点

### 认证端点 (`/api/v1/auth`)

| 方法 | 路径 | 描述 | 保护 |
|------|------|------|--------|
| POST | `/register` | 用户注册 | ❌ |
| POST | `/login` | 用户登录 | ❌ |
| GET | `/me` | 获取当前用户 | ✅ |
| PUT | `/me` | 更新用户信息 | ✅ |
| PUT | `/me/password` | 修改密码 | ✅ |
| PUT | `/me/theme` | 切换主题 | ✅ |
| DELETE | `/me` | 删除账户 | ✅ |
| POST | `/verify-email` | 验证邮箱 | ❌ |
| POST | `/resend-verification` | 重新发送验证邮件 | ✅ |
| POST | `/forgot-password` | 忘记密码 | ❌ |
| POST | `/reset-password` | 重置密码 | ❌ |
| POST | `/refresh-token` | 刷新令牌 | ✅ |
| GET | `/oauth/:provider` | OAuth 登录 | ❌ |
| GET | `/oauth/:provider/callback` | OAuth 回调 | ❌ |

### 分类端点 (`/api/v1/categories`)

| 方法 | 路径 | 描述 | 保护 |
|------|------|------|--------|
| GET | `/categories` | 获取所有分类 | ✅ |
| POST | `/categories` | 创建分类 | ✅ |
| GET | `/categories/:id` | 获取分类详情 | ✅ |
| PUT | `/categories/:id` | 更新分类 | ✅ |
| DELETE | `/categories/:id` | 删除分类 | ✅ |
| POST | `/categories/:id/share` | 生成分享链接 | ✅ |
| DELETE | `/categories/:id/share` | 取消分享 | ✅ |
| GET | `/search` | 搜索分类 | ✅ |
| GET | `/shared/:token` | 访问公开分享 | ❌ |
| POST | `/import` | 导入数据 | ❌ |
| GET | `/export` | 导出数据 | ❌ |
| GET | `/:id/export` | 导出单个分类 | ❌ |

### 站点端点 (`/api/v1/sites`)

| 方法 | 路径 | 描述 | 保护 |
|------|------|------|--------|
| GET | `/sites` | 获取所有站点 | ✅ |
| POST | `/sites` | 创建站点 | ✅ |
| GET | `/sites/:id` | 获取站点详情 | ✅ |
| PUT | `/sites/:id` | 更新站点 | ✅ |
| DELETE | `/sites/:id` | 删除站点 | ✅ |
| PUT | `/sites/:id/move` | 移动站点 | ✅ |
| GET | `/search` | 搜索站点 | ✅ |
| POST | `/batch` | 批量创建站点 | ✅ |

### 搜索端点 (`/api/v1/search`)

| 方法 | 路径 | 描述 | 保护 |
|------|------|------|--------|
| GET | `/sites` | 搜索站点 | ✅ |
| GET | `/categories` | 搜索分类 | ✅ |

### 健康检查

| 方法 | 路径 | 描述 | 保护 |
|------|------|------|--------|
| GET | `/health` | 服务器健康检查 | ❌ |

---

## 🔌 运行中服务器

**状态：** ✅ 运行中
**端口：** 8080
**数据库：** PostgreSQL（容器已启动）

---

## 📝 待实现功能

### 优先级高
- [ ] 邮件验证功能（数据库存储 + 发送）
- [ ] 密码重置功能（数据库存储 + 发送）
- [ ] 数据导入/导出功能

### 优先级中
- [ ] OAuth 登录（Google、GitHub）
- [ ] 公开分享功能（实现访问逻辑）

### 优先级低
- [ ] Redis 缓存（可选）
- [ ] Rate Limiting（防止滥用）
- [ ] 完整的错误处理

---

## 🛠️ 已知问题

1. **数据库连接**
   - 服务器已连接到 PostgreSQL
   - 表已自动创建

2. **端口配置**
   - 确保 Go 服务器和 PostgreSQL 使用相同端口
   - Docker 网络配置正确

---

## 🚀 快速开始

### 安装依赖
```bash
cd backend
go mod tidy
go mod download
```

### 启动数据库
```bash
docker-compose up -d postgres
```

### 运行服务器
```bash
cd backend
go run cmd/server/main.go
```

### 测试 API
```bash
# 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpass123","username":"testuser"}'

# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpass123"}'

# 获取分类
curl -X GET http://localhost:8080/api/v1/categories \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 📊 API 响应格式

### 成功响应
```json
{
  "code": 200,
  "message": "Success",
  "data": { ... }
}
```

### 错误响应
```json
{
  "code": 400,
  "error": "Error message"
}
```

---

**后端 API 基础架构已完成！** 🎉

可以开始测试 API 端点了。
