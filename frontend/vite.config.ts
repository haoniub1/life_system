import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    host: '0.0.0.0',
    port: 8085,
    allowedHosts: ['life.cyberniuma.xyz', 'life-test.cyberniuma.xyz', 'localhost'],
    proxy: {
      '/api': {
        target: 'http://localhost:8084',
        changeOrigin: true,
        secure: false,
        configure: (proxy, _options) => {
          proxy.on('proxyRes', (proxyRes, req, res) => {
            // Log cookies for debugging
            const cookies = proxyRes.headers['set-cookie']
            if (cookies) {
              console.log('🍪 Proxy forwarding cookies:', cookies)
            }
          })
        }
      },
      '/uploads': {
        target: 'http://localhost:8084',
        changeOrigin: true
      }
    }
  }
})
