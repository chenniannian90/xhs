# NavHub 完成部署指南

## 📋 当前状态

✅ 已完成：
- 服务器配置：Ubuntu 24.04
- 数据库：PostgreSQL + Redis（已安装并运行）
- 后端API：运行在8080端口
- 前端：已部署到 `/data/web/frontend`
- Nginx：已配置HTTP和HTTPS
- 自签名SSL证书：已配置

⚠️ 待完成：
- 京东云安全组配置（开放80和443端口）
- DNS记录验证
- Let's Encrypt正式SSL证书

---

## 🔥 第一步：配置京东云安全组（必须）

### 操作步骤：

1. **登录京东云控制台**
   - 访问：https://uc.jdcloud.com/
   - 进入：云主机 → 弹性云主机

2. **找到你的服务器**
   - 实例ID: lavm-0unap97y6f
   - IP: 117.72.39.169

3. **配置安全组规则**
   - 点击"安全组" → "配置规则"
   - 添加以下入站规则：

| 协议 | 端口 | 源 | 说明 |
|------|------|-----|------|
| TCP | 80 | 0.0.0.0/0 | HTTP（必须） |
| TCP | 443 | 0.0.0.0/0 | HTTPS（必须） |
| TCP | 22 | 0.0.0.0/0 | SSH（已有） |

4. **保存规则**

### 验证端口是否开放：

```bash
# 在本地执行
telnet hao246.cn 80
telnet hao246.cn 443
```

如果连接成功，说明端口已开放。

---

## 🌐 第二步：验证DNS配置

### 检查DNS记录：

```bash
# 检查A记录
dig hao246.cn
dig www.hao246.cn

# 或使用
nslookup hao246.cn
```

### 正确的DNS配置：

在你的域名注册商（阿里云、腾讯云、GoDaddy等）配置：

```
类型: A
主机记录: @
记录值: 117.72.39.169
TTL: 600

类型: A
主机记录: www
记录值: 117.72.39.169
TTL: 600
```

---

## 🔐 第三步：获取正式SSL证书

**安全组配置完成后**，执行以下命令：

### 方法1：自动获取（推荐）

```bash
ssh root@117.72.39.169

# 停止Nginx避免端口冲突
systemctl stop nginx

# 使用standalone模式获取证书
certbot certonly --standalone \
  -d hao246.cn \
  -d www.hao246.cn \
  --email admin@hao246.cn \
  --agree-tos \
  --non-interactive

# 启动Nginx
systemctl start nginx
```

### 方法2：手动DNS验证（如果方法1失败）

```bash
ssh root@117.72.39.169

# 使用DNS验证
certbot certonly --manual \
  --preferred-challenges dns \
  -d hao246.cn \
  -d www.hao246.cn \
  --email admin@hao246.cn \
  --agree-tos \
  --non-interactive
```

按照提示添加DNS TXT记录，等待生效后继续。

### 证书获取成功后：

```bash
# 配置自动续期
(crontab -l 2>/dev/null | grep -q "certbot renew") || \
  (crontab -l 2>/dev/null; echo "0 0,12 * * * certbot renew --quiet --deploy-hook 'systemctl reload nginx'") | crontab -

# 验证自动续期
certbot renew --dry-run
```

---

## ⚙️ 第四步：更新Nginx配置

使用正式证书后，执行：

```bash
ssh root@117.72.39.169

# 更新Nginx配置
cat > /etc/nginx/conf.d/navhub.conf << 'EOF'
# HTTP 自动跳转 HTTPS
server {
    listen 80;
    server_name hao246.cn www.hao246.cn;

    # Let's Encrypt 验证
    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    # 其他跳转到 HTTPS
    location / {
        return 301 https://$server_name$request_uri;
    }
}

# HTTPS 配置
server {
    listen 443 ssl http2;
    server_name hao246.cn www.hao246.cn;

    # SSL 证书
    ssl_certificate /etc/letsencrypt/live/hao246.cn/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/hao246.cn/privkey.pem;

    # SSL 优化
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    ssl_stapling on;
    ssl_stapling_verify on;

    # 安全头
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # 前端静态文件
    location / {
        root /data/web/frontend;
        try_files $uri $uri/ /index.html;

        # 缓存配置
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # 后端 API
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # 健康检查
    location /health {
        proxy_pass http://127.0.0.1:8080;
        access_log off;
    }

    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/x-javascript application/xml+rss application/json application/javascript;
}
EOF

# 测试配置
nginx -t

# 重启 Nginx
systemctl restart nginx

# 检查服务状态
systemctl status nginx
```

---

## ✅ 第五步：验证部署

### 1. 测试HTTPS访问

```bash
# 测试健康检查
curl https://hao246.cn/health

# 应该返回：
# {"database":"connected","status":"ok"}
```

### 2. 测试网站访问

在浏览器中访问：
- https://hao246.cn
- https://www.hao246.cn

检查：
- ✅ 网站正常加载
- ✅ 浏览器显示小锁图标（HTTPS）
- ✅ 证书有效（点击锁图标查看）

### 3. 测试API

```bash
# 测试API端点
curl https://hao246.cn/api/v1/health
```

### 4. 检查服务状态

```bash
ssh root@117.72.39.169 << 'EOF'
echo "=== 服务状态 ==="
systemctl status postgresql --no-pager -l | head 3
systemctl status redis-server --no-pager -l | head 3
systemctl status navhub-api 2>/dev/null || ps aux | grep navhub-api | grep -v grep
systemctl status nginx --no-pager -l | head 3

echo ""
echo "=== 端口监听 ==="
netstat -tlnp | grep -E '80|443|8080|5432|6379'

echo ""
echo "=== 证书信息 ==="
certbot certificates

echo ""
echo "=== 证书自动续期 ==="
crontab -l | grep certbot
EOF
```

---

## 🎉 完成！

### 部署清单

- [x] 服务器配置（Ubuntu 24.04）
- [x] 数据库安装（PostgreSQL + Redis）
- [x] 后端部署（Go API）
- [x] 前端部署（React静态文件）
- [x] Nginx配置
- [x] 自签名SSL证书（临时）
- [ ] **京东云安全组配置**（需要手动操作）
- [ ] **Let's Encrypt正式证书**（需要安全组配置后）

### 最终访问地址

- 🌐 **网站**: https://hao246.cn
- 🔐 **HTTPS**: https://hao246.cn
- 📡 **API**: https://hao246.cn/api/v1/
- ❤️ **健康检查**: https://hao246.cn/health

---

## 📞 京东云安全组配置帮助

如果找不到安全组配置：

1. **控制台路径**
   - 云主机 → 弹性云主机 → 实例列表
   - 点击实例ID → 更多 → 安全组 → 配置规则

2. **或者使用CLI**
   ```bash
   # 安装京东云CLI
   # 然后执行安全组配置
   ```

3. **联系客服**
   - 京东云技术支持：400-615-1210
   - 提供实例ID: lavm-0unap97y6f
   - 说明需要开放80和443端口

---

## 🔧 故障排查

### 如果HTTPS无法访问

```bash
# 1. 检查安全组
# 在京东云控制台确认80和443端口已开放

# 2. 检查防火墙
ssh root@117.72.39.169 'ufw status'

# 3. 检查Nginx
ssh root@117.72.39.169 'systemctl status nginx'

# 4. 检查端口监听
ssh root@117.72.39.169 'netstat -tlnp | grep -E "80|443"'

# 5. 测试外网访问
curl -I https://hao246.cn
```

### 如果证书获取失败

```bash
# 查看详细日志
ssh root@117.72.39.169 'tail -50 /var/log/letsencrypt/letsencrypt.log'

# 使用standalone模式（需要先停止Nginx）
ssh root@117.72.39.169 << 'EOF'
systemctl stop nginx
certbot certonly --standalone -d hao246.cn -d www.hao246.cn --email admin@hao246.cn --agree-tos --non-interactive
systemctl start nginx
EOF
```

---

## 📱 管理命令

### 查看日志

```bash
# 后端日志
ssh root@117.72.39.169 'tail -f /data/web/logs/backend.log'

# Nginx访问日志
ssh root@117.72.39.169 'tail -f /var/log/nginx/access.log'

# Nginx错误日志
ssh root@117.72.39.169 'tail -f /var/log/nginx/error.log'
```

### 重启服务

```bash
# 重启后端
ssh root@117.72.39.169 'systemctl restart navhub-api'

# 重启Nginx
ssh root@117.72.39.169 'systemctl restart nginx'

# 重启所有
ssh root@117.72.39.169 << 'EOF'
systemctl restart navhub-api
systemctl restart nginx
systemctl restart postgresql
systemctl restart redis-server
EOF
```

### 证书管理

```bash
# 续期证书
ssh root@117.72.39.169 'certbot renew'

# 查看证书
ssh root@117.72.39.169 'certbot certificates'

# 吊销证书
ssh root@117.72.39.169 'certbot revoke --cert-path /etc/letsencrypt/live/hao246.cn/cert.pem'
```

---

**最后更新**: 2026-03-26
**部署状态**: 90% 完成，等待安全组配置
