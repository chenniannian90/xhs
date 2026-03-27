// 快速测试脚本：检查 Zustand persist 的实际存储结构
// 在浏览器控制台运行

// 模拟 Zustand persist 的存储
const testStorage = () => {
  console.log('🧪 测试 Zustand persist 存储结构\n');

  // Zustand v4+ persist 中间件的默认存储结构
  console.log('📦 Zustand persist 默认结构:\n');
  console.log('```json');
  console.log(JSON.stringify({
    state: {
      user: { id: '123', username: 'test', email: 'test@example.com' },
      token: 'jwt-token-here',
      isAuthenticated: true
    },
    version: 0
  }, null, 2));
  console.log('```\n');

  // 检查实际存储
  const actual = localStorage.getItem('navhub-auth');
  if (actual) {
    const parsed = JSON.parse(actual);
    console.log('✅ 实际存储结构:\n');
    console.log('```json');
    console.log(JSON.stringify(parsed, null, 2));
    console.log('```\n');

    // 验证 api.ts 能否读取 token
    console.log('🔍 api.ts 读取测试:\n');
    console.log('代码: if (auth.state?.token)');
    if (parsed.state?.token) {
      console.log('✅ 成功: auth.state.token =', parsed.state.token.substring(0, 20) + '...');
    } else if (parsed.token) {
      console.log('❌ 失败: token 在 auth.token，不在 auth.state.token');
      console.log('🔧 需要修复 api.ts 的读取路径');
    } else {
      console.log('❌ 失败: 未找到 token');
    }
  } else {
    console.log('⚠️ 未找到 navhub-auth，请先登录');
  }
};

testStorage();
