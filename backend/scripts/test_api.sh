#!/bin/bash

# NavHub Backend API 自动化测试
# 用途：部署前验证所有 API 端点是否正常工作

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
BASE_URL="${API_BASE_URL:-http://localhost:8080}"
TEST_EMAIL="test_$(date +%s)@example.com"
TEST_USERNAME="testuser_$(date +%s)"
TEST_PASSWORD="Test123!@"

echo "======================================"
echo "NavHub Backend API 测试"
echo "======================================"
echo "测试服务器: $BASE_URL"
echo ""

# 测试计数器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 测试辅助函数
test_api() {
    local test_name="$1"
    local method="$2"
    local endpoint="$3"
    local data="$4"
    local expected_code="$5"

    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    echo -n "测试 $TOTAL_TESTS: $test_name ... "

    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL$endpoint" \
            -H 'Content-Type: application/json')
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint" \
            -H 'Content-Type: application/json' \
            -d "$data")
    fi

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" = "$expected_code" ]; then
        echo -e "${GREEN}✓ PASSED${NC} (HTTP $http_code)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        echo -e "${RED}✗ FAILED${NC}"
        echo "  期望: HTTP $expected_code"
        echo "  实际: HTTP $http_code"
        echo "  响应: $body"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
}

echo "======================================"
echo "1. 健康检查"
echo "======================================"

test_api \
    "健康检查" \
    "GET" \
    "/health" \
    "" \
    "200"

echo ""
echo "======================================"
echo "2. 用户注册"
echo "======================================"

test_api \
    "注册新用户" \
    "POST" \
    "/api/v1/auth/register" \
    "{\"username\":\"$TEST_USERNAME\",\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}" \
    "201"

# 保存 token 用于后续测试
echo ""
echo "获取登录 token..."
login_response=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
    -H 'Content-Type: application/json' \
    -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}")

TOKEN=$(echo "$login_response" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}✗ 无法获取登录 token${NC}"
    exit 1
else
    echo -e "${GREEN}✓ 成功获取 token${NC}"
fi

echo ""
echo "======================================"
echo "3. 用户登录"
echo "======================================"

test_api \
    "使用错误密码登录" \
    "POST" \
    "/api/v1/auth/login" \
    "{\"email\":\"$TEST_EMAIL\",\"password\":\"WrongPassword123!\"}" \
    "401"

test_api \
    "使用不存在的用户登录" \
    "POST" \
    "/api/v1/auth/login" \
    "{\"email\":\"notexist@example.com\",\"password\":\"$TEST_PASSWORD\"}" \
    "401"

echo ""
echo "======================================"
echo "4. 邮箱已注册测试"
echo "======================================"

test_api \
    "重复注册已存在的邮箱" \
    "POST" \
    "/api/v1/auth/register" \
    "{\"username\":\"anotheruser\",\"email\":\"$TEST_EMAIL\",\"password\":\"Test123!@\"}" \
    "409"

echo ""
echo "======================================"
echo "5. 需要认证的 API"
echo "======================================"

test_api \
    "无 token 访问分类" \
    "GET" \
    "/api/v1/categories" \
    "" \
    "401"

echo ""
echo "======================================"
echo "6. API 路径验证"
echo "======================================"

# 检查是否有重复路径
echo "检查 API 路径重复问题..."

test_result=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/api/v1/api/v1/auth/login" \
    -H 'Content-Type: application/json' \
    -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}")

if [ "$test_result" = "404" ]; then
    echo -e "${GREEN}✓ 没有路径重复问题 (返回 404)${NC}"
else
    echo -e "${RED}✗ 可能存在路径重复问题 (返回 $test_result)${NC}"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi

echo ""
echo "======================================"
echo "7. 密码验证规则"
echo "======================================"

test_api \
    "密码太短" \
    "POST" \
    "/api/v1/auth/register" \
    "{\"username\":\"shortpwd\",\"email\":\"shortpwd@test.com\",\"password\":\"A1!a\"}" \
    "400"

test_api \
    "密码缺少大写字母" \
    "POST" \
    "/api/v1/auth/register" \
    "{\"username\":\"noupper\",\"email\":\"noupper@test.com\",\"password\":\"test123!@\"}" \
    "400"

test_api \
    "密码缺少数字" \
    "POST" \
    "/api/v1/auth/register" \
    "{\"username\":\"nodigit\",\"email\":\"nodigit@test.com\",\"password\":\"TestAbc!@\"}" \
    "400"

test_api \
    "密码缺少特殊字符" \
    "POST" \
    "/api/v1/auth/register" \
    "{\"username\":\"nospecial\",\"email\":\"nospecial@test.com\",\"password\":\"Test123abc\"}" \
    "400"

echo ""
echo "======================================"
echo "测试结果汇总"
echo "======================================"

echo "总测试数: $TOTAL_TESTS"
echo -e "通过: ${GREEN}$PASSED_TESTS${NC}"
echo -e "失败: ${RED}$FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo ""
    echo -e "${GREEN}======================================"
    echo "✓ 所有测试通过！"
    echo "======================================${NC}"
    exit 0
else
    echo ""
    echo -e "${RED}======================================"
    echo "✗ 有 $FAILED_TESTS 个测试失败"
    echo "======================================${NC}"
    exit 1
fi
