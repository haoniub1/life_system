const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const path = require('path');

const app = express();
const PORT = 8082;

// 代理 API 请求到后端
app.use('/api', createProxyMiddleware({
  target: 'http://localhost:8084',
  changeOrigin: true,
  logLevel: 'info'
}));

// 代理上传文件
app.use('/uploads', createProxyMiddleware({
  target: 'http://localhost:8084',
  changeOrigin: true
}));

// 静态文件服务 - 优先使用构建后的文件
const distPath = path.join(__dirname, 'frontend/dist');
app.use(express.static(distPath));

// SPA 路由处理  
app.get('/*', (req, res) => {
  res.sendFile(path.join(distPath, 'index.html'));
});

app.listen(PORT, () => {
  console.log(`🚀 Life System Proxy Server running on http://localhost:${PORT}`);
  console.log(`📁 Serving static files from: ${distPath}`);
  console.log(`🔗 API proxy: /api/* -> http://localhost:8084`);
  console.log(`📤 Uploads proxy: /uploads/* -> http://localhost:8084`);
});