/**
 * Vite配置文件
 * 用于配置开发服务器和构建选项
 */
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// 导出Vite配置
export default defineConfig({
  // 使用Vue插件
  plugins: [vue()],
  // 开发服务器配置
  server: {
    port: 5173,  // 端口号
    open: true   // 自动打开浏览器
  }
})
