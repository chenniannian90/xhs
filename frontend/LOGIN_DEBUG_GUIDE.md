# 登录问题诊断指南

## 🔍 步骤 1：打开浏览器控制台

1. 访问 http://localhost:5173
2. 按 F12 打开开发者工具
3. 切换到 **Console** 标签

## 🔍 步骤 2：清除现有数据

在控制台中运行：
```javascript
localStorage.clear();
location.reload();
```

## 🔍 步骤 3：尝试登录

1. 在登录页面输入你的邮箱和密码
2. 点击"登录"按钮
3. **观察控制台输出**

### 期望看到的日志：

#### 登录成功应该看到：
```
✅ Login successful!
  User: [用户名]
  Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
💾 Stored in localStorage:
  Has user: true
  Has token: true
  Is authenticated: true
```

#### 后续 API 请求应该看到：
```
📡 API Request: GET /api/v1/categories
  Token found: true
  ✅ Authorization header set
```

### 如果看到错误：

#### ❌ 情况 1：登录失败
```
❌ Login failed: [错误信息]
```
**原因：** 邮箱或密码错误
**解决：** 检查邮箱和密码是否正确

#### ❌ 情况 2：Token 未找到
```
  ❌ No token found in auth data
  Auth data structure: {...}
```
**原因：** Token 存储结构不匹配
**需要修复：** 截图控制台输出

#### ❌ 情况 3：401 Unauthorized
```
📡 API Request: GET /api/v1/categories
  Token found: true
  ✅ Authorization header set
[之后出现 Unauthorized request - clearing auth data]
```
**原因：** Token 无效或过期
**解决：** 清除缓存重新登录

## 🔍 步骤 4：使用测试工具

打开测试页面（可选）：
```
file:///Users/mac-new/work/navhub/frontend/test-login.html
```

点击"检查存储"查看当前 localStorage 内容

## 📋 需要提供给诊断的信息：

如果登录仍有问题，请提供：

1. **控制台截图**（显示登录尝试的日志）
2. **Network 标签截图**
   - 点击登录请求
   - 查看 Response 内容
3. **Application 标签截图**
   - Local Storage → navhub-auth
   - Value 内容

## 🛠️ 临时解决方案

如果登录一直失败，尝试：

### 方案 1：重新注册
```
1. 清除所有数据
2. 访问 /register
3. 注册新用户
4. 使用新用户登录
```

### 方案 2：检查现有用户
在数据库中查询：
```bash
docker exec -it navhub-postgres psql -U navhub -d navhub -c "SELECT email, username, email_verified FROM users;"
```

### 方案 3：手动测试 API
在控制台运行：
```javascript
fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: {'Content-Type': 'application/json'},
  body: JSON.stringify({
    email: 'your-email@example.com',
    password: 'your-password'
  })
}).then(r => r.json()).then(console.log)
```

---

**请在浏览器中按照上述步骤操作，并告诉我看到了什么日志输出。**
