// NavHub API 测试脚本
// 在浏览器控制台运行此脚本进行诊断

console.log('🔍 NavHub API 诊断工具');
console.log('='.repeat(50));

// 测试 1: 检查 localStorage 结构
console.log('\n📦 测试 1: 检查 localStorage');
const authData = localStorage.getItem('navhub-auth');
console.log('navhub-auth 存在:', !!authData);

if (authData) {
  try {
    const parsed = JSON.parse(authData);
    console.log('✅ JSON 解析成功');
    console.log('完整结构:', parsed);

    // 检查不同的 token 路径
    console.log('\n🔑 Token 路径检查:');
    console.log('  - parsed.token:', parsed.token ? '✅ 存在' : '❌ 不存在');
    console.log('  - parsed.state?.token:', parsed.state?.token ? '✅ 存在' : '❌ 不存在');
    console.log('  - parsed.state?.state?.token:', parsed.state?.state?.token ? '✅ 存在' : '❌ 不存在');

    // 检查用户数据
    console.log('\n👤 用户数据检查:');
    console.log('  - parsed.user:', parsed.user ? '✅ 存在' : '❌ 不存在');
    if (parsed.user) {
      console.log('    用户名:', parsed.user.username);
      console.log('    邮箱:', parsed.user.email);
    }

  } catch (e) {
    console.error('❌ JSON 解析失败:', e);
  }
} else {
  console.log('❌ 未找到 navhub-auth 数据');
  console.log('💡 提示: 请先登录或使用 test-zustand-storage.html 模拟认证');
}

// 测试 2: 模拟 API 拦截器逻辑
console.log('\n📡 测试 2: 模拟 API 拦截器');
console.log('api.ts 中的代码:');
console.log('```');
console.log('const token = localStorage.getItem("navhub-auth");');
console.log('if (token) {');
console.log('  const auth = JSON.parse(token);');
console.log('  if (auth.state?.token) {');
console.log('    config.headers.Authorization = `Bearer ${auth.state.token}`;');
console.log('  }');
console.log('}');
console.log('```');

if (authData) {
  const auth = JSON.parse(authData);
  if (auth.state?.token) {
    console.log('✅ Authorization header 会被设置');
    console.log('   Bearer token:', auth.state.token.substring(0, 30) + '...');
  } else if (auth.token) {
    console.log('❌ Authorization header 不会被设置');
    console.log('⚠️ 问题: token 在 auth.token，但 api.ts 读取 auth.state.token');
  } else {
    console.log('❌ 未找到 token');
  }
}

// 测试 3: 实际 API 调用测试
console.log('\n🌐 测试 3: 实际 API 调用');

async function testAPI() {
  const endpoints = [
    { url: '/api/v1/auth/me', method: 'GET', name: '获取当前用户' },
    { url: '/api/v1/categories', method: 'GET', name: '获取分类列表' },
  ];

  for (const endpoint of endpoints) {
    console.log(`\n测试: ${endpoint.name} (${endpoint.url})`);

    try {
      // 使用 fetch 模拟 axios 调用
      const token = localStorage.getItem('navhub-auth');
      let headers = {
        'Content-Type': 'application/json',
      };

      if (token) {
        const auth = JSON.parse(token);
        if (auth.state?.token) {
          headers['Authorization'] = `Bearer ${auth.state.token}`;
        } else if (auth.token) {
          headers['Authorization'] = `Bearer ${auth.token}`;
        }
      }

      console.log('  请求头:', headers);

      const response = await fetch(`http://localhost:8080${endpoint.url}`, {
        method: endpoint.method,
        headers: headers,
      });

      console.log('  状态码:', response.status, response.statusText);

      if (response.ok) {
        const data = await response.json();
        console.log('  ✅ 成功:', data);
      } else {
        const error = await response.json();
        console.log('  ❌ 失败:', error);

        if (response.status === 401) {
          console.log('  ⚠️ 认证失败 - 可能是 token 无效或未设置');
        }
      }
    } catch (error) {
      console.log('  ❌ 网络错误:', error.message);
    }
  }
}

// 运行 API 测试
testAPI().then(() => {
  console.log('\n' + '='.repeat(50));
  console.log('✅ 诊断完成');
  console.log('\n💡 建议:');
  console.log('1. 如果看到 "auth.state.token 不存在"，需要修复 api.ts 的 token 读取路径');
  console.log('2. 如果 API 返回 401，检查 token 是否有效');
  console.log('3. 如果看到 CORS 错误，检查后端 CORS 配置');
});
