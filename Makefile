# NavHub Makefile - 快速命令参考

.PHONY: help test deploy check clean

# 默认目标
help:
	@echo "======================================"
	@echo "NavHub 常用命令"
	@echo "======================================"
	@echo ""
	@echo "检查和测试："
	@echo "  make check          - 运行部署前检查"
	@echo "  make test           - 运行所有测试"
	@echo "  make test-api       - 测试后端 API"
	@echo "  make test-frontend  - 测试前端"
	@echo ""
	@echo "部署："
	@echo "  make deploy         - 一键完整部署"
	@echo "  make deploy-backend - 仅部署后端"
	@echo "  make deploy-frontend - 仅部署前端"
	@echo ""
	@echo "开发："
	@echo "  make dev            - 启动前端开发服务器"
	@echo "  make dev-backend    - 启动后端开发服务器"
	@echo "  make build          - 构建前端和后端"
	@echo ""
	@echo "维护："
	@echo "  make logs           - 查看后端日志"
	@echo "  make status         - 查看服务状态"
	@echo "  make clean          - 清理构建文件"
	@echo ""

# ==================== 检查和测试 ====================

check:
	@echo "运行部署前检查..."
	@./deploy/pre-deploy-check.sh

test-api:
	@echo "测试后端 API..."
	@cd backend && API_BASE_URL=http://117.72.39.169:8080 ./scripts/test_api.sh

test-frontend:
	@echo "测试前端..."
	@cd frontend && npm run type-check && npm run lint

test: test-api test-frontend
	@echo "所有测试完成！"

# ==================== 部署 ====================

deploy:
	@echo "开始完整部署..."
	@./deploy/deploy-all.sh

deploy-backend:
	@echo "部署后端..."
	@./deploy/deploy-backend.sh

deploy-frontend:
	@echo "部署前端..."
	@./deploy/deploy-frontend.sh

# ==================== 开发 ====================

dev:
	@echo "启动前端开发服务器..."
	@cd frontend && npm run dev

dev-backend:
	@echo "启动后端开发服务器..."
	@cd backend && go run cmd/server/main.go

build:
	@echo "构建前端..."
	@cd frontend && npm run build
	@echo "构建后端..."
	@cd backend && GOOS=linux GOARCH=amd64 go build -o navhub-api cmd/server/main.go

# ==================== 维护 ====================

logs:
	@echo "查看后端日志..."
	@ssh root@117.72.39.169 'journalctl -u navhub-api -f -n 100'

logs-nginx:
	@echo "查看 Nginx 日志..."
	@ssh root@117.72.39.169 'tail -f /var/log/nginx/access.log'

status:
	@echo "检查服务状态..."
	@ssh root@117.72.39.169 'systemctl status navhub-api --no-pager'

health:
	@echo "检查服务健康状态..."
	@curl -s http://117.72.39.169:8080/health | python3 -m json.tool || echo "健康检查失败"

clean:
	@echo "清理构建文件..."
	@cd frontend && rm -rf dist/
	@cd backend && rm -f navhub-api
	@echo "清理完成"

# ==================== 快捷命令 ====================

# 快速部署（跳过某些检查）
deploy-quick:
	@echo "快速部署（跳过检查）..."
	@./deploy/deploy-backend.sh
	@./deploy/deploy-frontend.sh

# 回滚到上一个版本
rollback:
	@echo "回滚到上一个版本..."
	@echo "警告：这将回滚前端和后端到上一个备份版本"
	@read -p "确认继续？[y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		ssh root@117.72.39.169 \
			'BACKUP=$$(ls -t /data/web/backend/navhub-api.backup.* | head -1); \
			cp $$BACKUP /data/web/backend/navhub-api && \
			systemctl restart navhub-api && \
			echo "后端已回滚" && \
			FRONTEND=$$(ls -td /data/web/frontend.backup.* | head -1); \
			rm -rf /data/web/frontend && \
			cp -r $$FRONTEND /data/web/frontend && \
			echo "前端已回滚"'; \
	fi

# 监控实时状态
monitor:
	@watch -n 2 'bash -c "echo \"=== 服务状态 ===\" && ssh root@117.72.39.169 \"systemctl status navhub-api --no-pager | head -5\" && echo \"\" && echo \"=== 内存使用 ===\" && ssh root@117.72.39.169 \"free -h\" && echo \"\" && echo \"=== 磁盘使用 ===\" && ssh root@117.72.39.169 \"df -h | grep -E '(Filesystem|/data)'\""'
