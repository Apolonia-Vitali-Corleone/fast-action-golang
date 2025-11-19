/**
 * Vue应用入口文件
 * 创建并挂载Vue应用，配置Pinia状态管理和路由
 */
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

// 导入Element Plus样式
import 'element-plus/dist/index.css'

// 创建Pinia实例
const pinia = createPinia()

// 创建Vue应用实例
const app = createApp(App)

// 使用Pinia
app.use(pinia)

// 使用路由
app.use(router)

// 挂载到#app元素
app.mount('#app')
