# NavHub - 个人导航管理系统

## 🚀 项目简介

NavHub 是一个现代化的个人导航管理系统，支持���后端分离架构，提供优雅的 UI 设计和完善的用户体验。

## ✨ 主要特性

- 🔐 **用户认证**：邮箱注册/登录，OAuth（Google、GitHub）
- 📧 **邮箱验证**：注册后需要验证邮箱
- 🔑 **密码管理**：支持忘记密码功能
- 📁 **分类管理**：创建、编辑、删除分类，支持拖拽排序
- 🔗 **站点管理**：添加、编辑、删除站点，支持图标自定义
- 🔍 **全局搜索**：搜索站点名称、描述、URL
- 📤 **导出功能**：导出为 CSV，支持全部导出或单个分类导出
- 📥 **导入功能**：从 CSV/JSON 导入数据
- 🌐 **公开分享**：分享分类给他人查看
- 🎨 **主题切换**：支持浅色/深色模式
- 📱 **响应式设计**：完美适配桌面和移动设备

## 🏗️ 技术栈

### 前端
- **框架**: React 18 + TypeScript
- **构建工具**: Vite
- **UI 组件**: TailwindCSS + shadcn/ui
- **状态管理**: Zustand
- **数据获取**: TanStack Query (React Query)
- **路由**: React Router v6
- **HTTP 客户端**: Axios

### 后端
- **语言**: Go 1.21+
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL
- **缓存**: Redis
- **认证**: JWT + OAuth2

## 📁 项目结构

```
navhub/
├── frontend/          # React 前端应用
├── backend/           # Go 后端应用
├── docs/             # 项目文档
└── docker-compose.yml # 本地开发环境
```

## 🛠️ 快速开始

### 前置要求

- Node.js 18+
- Go 1.21+
- PostgreSQL 14+
- Redis 7+

### 本地开发

1. **克隆仓库**
   ```bash
   git clone <repository-url>
   cd navhub
   ```

2. **启动数据库**
   ```bash
   docker-compose up -d postgres redis
   ```

3. **启动后端**
   ```bash
   cd backend
   cp .env.example .env
   go mod download
   go run cmd/server/main.go
   ```

4. **启动前端**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

5. **访问应用**
   - 前端: http://localhost:5173
   - 后端: http://localhost:8080

## 📝 开发规范

- 遵循 [Conventional Commits](https://www.conventionalcommits.org/)
- 代码风格：前端使用 ESLint + Prettier，后端使用 gofmt
- 提交前通过测试：`npm test` (前端), `go test ./...` (后端)

## 📄 许可证

MIT License

---

**开发中...** 🚧