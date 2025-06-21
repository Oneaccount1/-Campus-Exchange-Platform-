<template>
  <el-container class="app-container">
    <!-- 当不是管理员路由时才显示前台头部 -->
    <el-header height="60px" v-if="!isAdminRoute">
      <el-menu
        :router="true"
        mode="horizontal"
        :ellipsis="false"
        class="nav-menu"
      >
        <el-menu-item index="/">
          <el-icon><HomeFilled /></el-icon>
          校园二手交易平台
        </el-menu-item>
        <div class="flex-grow" />
        <el-menu-item index="/products">
          <el-icon><ShoppingBag /></el-icon>
          商品列表
        </el-menu-item>
        <el-menu-item index="/publish">
          <el-icon><Plus /></el-icon>
          发布商品
        </el-menu-item>
        <el-menu-item index="/messages">
          <el-badge :value="messageStore.totalUnread" :hidden="messageStore.totalUnread === 0" :max="99">
            <el-icon><Message /></el-icon>
            消息中心
          </el-badge>
        </el-menu-item>
        <template v-if="userStore.isLoggedIn">
          <el-sub-menu index="user-menu">
            <template #title>
              <el-avatar :size="32" :src="userStore.userInfo?.avatar || ''">
                {{ userStore.userInfo?.username?.charAt(0) }}
              </el-avatar>
            </template>
            <el-menu-item index="/user">
              <el-icon><User /></el-icon>
              个人中心
            </el-menu-item>
            <el-menu-item @click="handleLogout">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-menu-item>
          </el-sub-menu>
        </template>
        <template v-else>
          <el-menu-item index="/login">
            <el-icon><Key /></el-icon>
            登录
          </el-menu-item>
          <el-menu-item index="/register">
            <el-icon><UserFilled /></el-icon>
            注册
          </el-menu-item>
        </template>
      </el-menu>
    </el-header>
    
    <el-main :class="{ 'admin-main': isAdminRoute }">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </el-main>

    <!-- 当不是管理员路由时才显示前台底部 -->
    <el-footer height="60px" v-if="!isAdminRoute">
      <div class="footer-content">
        © 2025 校园二手交易平台 - 让闲置物品流转起来
      </div>
    </el-footer>
  </el-container>
</template>

<script setup>
import { onMounted, onBeforeUnmount, ref, computed } from 'vue'
import { useUserStore, useMessageStore } from './stores'
import { ElMessage } from 'element-plus'
import { HomeFilled, ShoppingBag, Plus, Message, User, UserFilled, Key, SwitchButton } from '@element-plus/icons-vue'
import webSocketService from './utils/websocket'
import * as messageApi from './api/message'
import { getUserInfo } from './api/user'
import { useRouter, useRoute } from 'vue-router'
import WebSocketStatus from './components/WebSocketStatus.vue'

const userStore = useUserStore()
const messageStore = useMessageStore()
const router = useRouter()
const route = useRoute()

// 判断当前是否为管理员路由
const isAdminRoute = computed(() => {
  return route.path.startsWith('/admin')
})

const handleLogout = () => {
  // 断开WebSocket连接
  webSocketService.disconnect()
  
  // 清空消息数据
  messageStore.clearAll()
  
  // 退出登录
  userStore.logout()
  ElMessage.success('退出登录成功')
}

// 初始化获取未读消息数量
const initUnreadCount = async () => {
  if (!userStore.isLoggedIn) return
  
  try {
    // 启用API调用，后端已实现
    const response = await messageApi.getUnreadCount()
    if (response.data) {
      messageStore.setTotalUnread(response.data.count || 0)
    }
  } catch (error) {
    console.error('获取未读消息数量失败:', error)
    // 如果API调用失败，设置默认值
    messageStore.setTotalUnread(0)
  }
}

// WebSocket初始化标志
let wsInitialized = false

// 初始化WebSocket连接
const initWebSocket = () => {
  if (!userStore.isLoggedIn) return
  
  console.log('正在初始化全局WebSocket连接...')
  
  // 确保先断开任何现有连接
  webSocketService.disconnect()
  
  // 建立连接
  webSocketService.connect()
  
  // 监听连接状态变化
  webSocketService.onConnection((connected) => {
    const status = connected ? '已连接' : '已断开'
    console.log('全局WebSocket连接状态变更:', status)
    
    // 更新应用状态
    appStatus.value.wsConnected = connected
    appStatus.value.lastWsStatusChange = new Date().toLocaleTimeString()
    
    // 如果连接断开且用户已登录，不再需要额外处理
    // WebSocketService内部会自动重连
  })
  
  // 设置全局消息处理器
  webSocketService.onMessage(handleGlobalWebSocketMessage)
  
  // 定期更新连接状态信息
  startConnectionInfoUpdates()
  
  wsInitialized = true
}

// 定期更新连接状态信息
const startConnectionInfoUpdates = () => {
  // 创建应用状态对象，用于监控
  const appStatus = ref({
    wsConnected: webSocketService.isConnected,
    wsLatency: null,
    wsReconnectAttempts: 0,
    wsQueuedMessages: 0,
    lastWsStatusChange: new Date().toLocaleTimeString()
  })
  
  // 每15秒更新一次连接状态信息
  const infoInterval = setInterval(() => {
    if (!userStore.isLoggedIn) {
      clearInterval(infoInterval)
      return
    }
    
    // 获取连接信息
    const info = webSocketService.getConnectionInfo()
    
    // 更新状态
    appStatus.value.wsConnected = info.connected
    appStatus.value.wsLatency = info.latency
    appStatus.value.wsReconnectAttempts = info.reconnectAttempts
    appStatus.value.wsQueuedMessages = info.queuedMessages
    
    // 输出调试信息
    if (process.env.NODE_ENV !== 'production') {
      console.log('WebSocket状态:', {
        connected: info.connected,
        latency: info.latency ? `${info.latency}ms` : 'N/A',
        reconnectAttempts: info.reconnectAttempts,
        queuedMessages: info.queuedMessages
      })
    }
  }, 15000)
  
  // 组件卸载时清理
  onBeforeUnmount(() => {
    clearInterval(infoInterval)
  })
  
  return appStatus
}

// 创建应用状态对象
const appStatus = startConnectionInfoUpdates()

// 全局WebSocket消息处理
const handleGlobalWebSocketMessage = (message) => {
  if (!message) return
  
  const type = message.type || 'unknown'
  const data = message.data || {}
  
  // 仅在开发环境输出调试日志
  if (process.env.NODE_ENV !== 'production') {
    console.log('App收到WebSocket消息类型:', type)
  }
  
  // 处理特定类型的消息
  switch (type) {
    case 'message':
      // 收到新消息时更新未读消息数
      initUnreadCount()
      // 触发桌面通知
      if (data.senderId !== userStore.userInfo?.id) {
        notifyNewMessage(data)
      }
      break
    case 'online_status':
      // 用户在线状态变更，可以在这里处理
      break
    case 'notification':
      // 显示系统通知
      if (data && data.message) {
        ElMessage({
          message: data.message,
          type: data.type || 'info'
        })
      }
      break
  }
}

// 触发新消息通知
const notifyNewMessage = (messageData) => {
  import('./utils/notification').then(module => {
    const service = module.default
    service.showDesktopNotification({
      title: '新消息通知',
      message: messageData.content || '您有一条新消息',
      onClick: () => {
        router.push('/messages')
      }
    })
  })
}

// 检查并恢复用户会话
const restoreUserSession = async () => {
  // 检查是否已存在token
  const token = localStorage.getItem('token')
  if (!token) {
    console.log('未找到登录凭证，用户未登录')
    return false
  }
  
  // 检查store状态
  if (userStore.isLoggedIn && userStore.userInfo) {
    console.log('用户已登录且信息完整')
    
    // 确保WebSocket已连接
    if (!webSocketService.isConnected) {
      initWebSocket()
    }
    
    return true
  }
  
  // 尝试从API获取用户信息
  try {
    console.log('尝试从服务器重新获取用户会话...')
    const response = await getUserInfo()
    
    if (response.code === 200 && response.data) {
      console.log('用户会话成功恢复:', response.data.username)
      userStore.setUserInfo(response.data)
      userStore.setToken(token)
      
      // 初始化WebSocket
      initWebSocket()
      
      return true
    } else {
      console.warn('API返回成功但用户数据无效，静默失败:', response)
      // 不要立即登出，静默失败
      return false
    }
  } catch (error) {
    console.error('恢复用户会话失败，可能是token已过期:', error)
    // 不要立即登出或显示消息，让请求拦截器处理认证错误
    return false
  }
}

onMounted(async () => {
  // 尝试恢复用户会话
  const sessionRestored = await restoreUserSession()
  
  // 初始化
  if (sessionRestored) {
    // 先获取未读消息
    await initUnreadCount()
  }
})

onBeforeUnmount(() => {
  // 移除全局消息处理器
  webSocketService.offMessage(handleGlobalWebSocketMessage)
  
  // 清理资源，使用1000正常关闭码
  if (webSocketService.isConnected) {
    webSocketService.disconnect()
  }
})
</script>

<style>
.app-container {
  min-height: 100vh;
}

.nav-menu {
  width: 100%;
  display: flex;
  align-items: center;
}

.flex-grow {
  flex-grow: 1;
}

.el-header {
  padding: 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.el-main {
  padding: 20px;
  background-color: #f5f7fa;
}

.el-footer {
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #fff;
  border-top: 1px solid #e4e7ed;
}

.footer-content {
  color: #909399;
  font-size: 14px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 为WebSocket状态指示器添加样式 */
.user-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}
</style>
