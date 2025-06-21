import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建axios实例
const request = axios.create({
  baseURL: 'http://localhost:8080/api/v1', // 后端API的基础URL
  timeout: 5000 // 请求超时时间
})

// 创建管理员API专用的axios实例
const adminRequest = axios.create({
  baseURL: 'http://localhost:8080/api/v1', // 后端API的基础URL
  timeout: 5000 // 请求超时时间
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    // 从localStorage获取token
    const token = localStorage.getItem('token')
    if (token) {
      // 确保不添加双引号
      const cleanToken = token.replace(/^"(.*)"$/, '$1')
      config.headers['Authorization'] = `Bearer ${cleanToken}`
      console.log('请求已附加Token认证头')
    }
    return config
  },
  error => {
    console.error('请求拦截器错误:', error)
    return Promise.reject(error)
  }
)

// 管理员请求拦截器
adminRequest.interceptors.request.use(
  config => {
    // 从localStorage获取管理员token
    const adminToken = localStorage.getItem('admin_token')
    if (adminToken) {
      // 确保不添加双引号
      const cleanToken = adminToken.replace(/^"(.*)"$/, '$1')
      config.headers['Authorization'] = `Bearer ${cleanToken}`
      console.log('管理员请求已附加Token认证头')
    }
    return config
  },
  error => {
    console.error('管理员请求拦截器错误:', error)
    return Promise.reject(error)
  }
)

// 导入pinia store
import { createPinia } from 'pinia'
import { useUserStore, useAdminStore } from '../stores'

// 创建pinia实例
const pinia = createPinia()
// 创建store实例
const _storeTemp = useUserStore(pinia)
const _adminStoreTemp = useAdminStore(pinia)

// 响应拦截器
request.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    // 处理网络错误
    if (!error.response) {
      console.error('网络错误:', error.message)
      ElMessage.error('网络错误，请检查您的连接')
      return Promise.reject(error)
    }

    // 根据不同状态码处理错误
    const status = error.response.status
    
    if (status === 401) {
      console.warn('401 Unauthorized 错误:', error.response.data)
      // 获取当前路径
      const currentPath = window.location.pathname
      
      // 如果已经在登录页，不要重复跳转或显示消息
      if (currentPath !== '/login') {
        // 获取store实例
        const userStore = useUserStore()
        userStore.logout()
        
        // 重定向到登录页面，保留当前路径用于登录后返回
        window.location.href = `/login?redirect=${encodeURIComponent(currentPath)}`
        ElMessage.error('登录已过期，请重新登录')
      }
    } else if (status === 403) {
      ElMessage.error('您没有权限执行此操作')
    } else if (status === 500) {
      ElMessage.error('服务器错误，请稍后再试')
    } else {
      // 通用错误消息
      const message = error.response.data?.message || '请求失败，请稍后再试'
      ElMessage.error(message)
    }
    
    return Promise.reject(error)
  }
)

// 管理员响应拦截器
adminRequest.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    // 处理网络错误
    if (!error.response) {
      console.error('网络错误:', error.message)
      ElMessage.error('网络错误，请检查您的连接')
      return Promise.reject(error)
    }

    // 根据不同状态码处理错误
    const status = error.response.status
    
    if (status === 401) {
      console.warn('401 Unauthorized 错误:', error.response.data)
      // 获取当前路径
      const currentPath = window.location.pathname
      
      // 如果已经在管理员登录页，不要重复跳转或显示消息
      if (currentPath !== '/admin/login') {
        // 获取store实例
        const adminStore = useAdminStore()
        adminStore.logout()
        
        // 重定向到管理员登录页面，保留当前路径用于登录后返回
        window.location.href = `/admin/login?redirect=${encodeURIComponent(currentPath)}`
        ElMessage.error('管理员登录已过期，请重新登录')
      }
    } else if (status === 403) {
      ElMessage.error('您没有权限执行此操作')
    } else if (status === 500) {
      ElMessage.error('服务器错误，请稍后再试')
    } else {
      // 通用错误消息
      const message = error.response.data?.message || '请求失败，请稍后再试'
      ElMessage.error(message)
    }
    
    return Promise.reject(error)
  }
)

export { adminRequest }
export default request 