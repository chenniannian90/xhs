export default function Test() {
  return (
    <div style={{ padding: '20px', fontFamily: 'sans-serif' }}>
      <h1>🎉 NavHub Frontend Test</h1>
      <p>如果你看到这个页面，说明前端基础功能正常！</p>
      <p>当前时间: {new Date().toLocaleString()}</p>
      <button onClick={() => alert('按钮点击正常！')}>
        点击测试
      </button>
    </div>
  );
}
