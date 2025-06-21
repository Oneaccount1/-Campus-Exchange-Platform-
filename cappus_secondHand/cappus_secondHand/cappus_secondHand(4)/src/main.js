import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

// 创建Vue应用实例
const app = createApp(App)

// 添加全局错误处理
app.config.errorHandler = (err, instance, info) => {
  console.error('Vue应用错误:', err)
  console.error('错误信息:', info)
  // 可以在这里添加更多的错误处理，如发送到服务器等
}

// 注册Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(createPinia())
app.use(router)
app.use(ElementPlus, { size: 'default', zIndex: 3000 })

// 尝试挂载并捕获任何错误
try {
  app.mount('#app')
  console.log('Vue应用成功挂载')
} catch (error) {
  console.error('Vue应用挂载失败:', error)
}
