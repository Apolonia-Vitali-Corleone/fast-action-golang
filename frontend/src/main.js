/**
 * Vue应用入口文件
 * 创建并挂载Vue应用，配置Pinia状态管理
 */
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'

// 创建Pinia实例
const pinia = createPinia()

// 创建Vue应用实例
const app = createApp(App)

// 使用Pinia
app.use(pinia)

// 挂载到#app元素
app.mount('#app')
