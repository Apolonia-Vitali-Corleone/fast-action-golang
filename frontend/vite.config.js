/**
 * Vite配置文件
 * 配置Element Plus自动导入和OnlyFans主题色
 */
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import path from 'path'

export default defineConfig({
  plugins: [
    vue(),
    // 自动导入Vue和Element Plus的API
    AutoImport({
      imports: ['vue'],
      resolvers: [ElementPlusResolver()],
    }),
    // 自动导入Element Plus组件
    Components({
      resolvers: [ElementPlusResolver()],
    }),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
  },
  server: {
    port: 5173,
    open: true
  }
})
