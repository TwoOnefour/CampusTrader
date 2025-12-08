import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
// web/vite.config.ts
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      // 当你请求 /api/v1/user 时，Vite 会自动转发给 http://localhost:8080/api/v1/user
      '/api': {
        target: 'http://localhost:8080', // 你的 Go 后端地址
        changeOrigin: true,
      }
    }
  }
})
