<template>
  <div class="admin-layout">
    <el-container class="admin-container">
      <el-aside :width="isCollapse ? '64px' : '220px'" class="admin-aside">
        <div class="logo-container">
          <h2 v-show="!isCollapse">后台管理系统</h2>
        </div>
        <el-menu
          :router="true"
          :default-active="activeMenu"
          class="admin-menu"
          :collapse="isCollapse"
          background-color="#001529"
          text-color="#a6adb4"
          active-text-color="#ffffff"
        >
          <el-menu-item index="/admin/dashboard">
            <el-icon><DataBoard /></el-icon>
            <template #title>仪表盘</template>
          </el-menu-item>
          <el-menu-item index="/admin/products">
            <el-icon><ShoppingBag /></el-icon>
            <template #title>商品管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/users">
            <el-icon><User /></el-icon>
            <template #title>用户管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/orders" v-if="hasOrdersModule">
            <el-icon><List /></el-icon>
            <template #title>订单管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/messages" v-if="hasMessagesModule">
            <el-icon><ChatDotRound /></el-icon>
            <template #title>消息管理</template>
          </el-menu-item>
        </el-menu>
      </el-aside>
      
      <el-container>
        <el-header height="60px" class="admin-header">
          <div class="header-left">
            <el-icon class="menu-toggle" @click="toggleCollapse"><Fold v-if="!isCollapse" /><Expand v-else /></el-icon>
            <div class="breadcrumb-container">
              <el-breadcrumb separator="/">
                <el-breadcrumb-item :to="{ path: '/admin/dashboard' }">首页</el-breadcrumb-item>
                <el-breadcrumb-item v-if="currentPageTitle">{{ currentPageTitle }}</el-breadcrumb-item>
              </el-breadcrumb>
            </div>
          </div>
          <div class="header-right">
            <div class="header-actions">
              <el-tooltip content="全屏" placement="bottom">
                <el-icon class="header-icon" @click="toggleFullScreen"><FullScreen /></el-icon>
              </el-tooltip>
              <el-tooltip content="刷新页面" placement="bottom">
                <el-icon class="header-icon" @click="refreshPage"><RefreshRight /></el-icon>
              </el-tooltip>
              <el-tooltip content="前台首页" placement="bottom">
                <el-icon class="header-icon" @click="goToFrontend"><House /></el-icon>
              </el-tooltip>
            </div>
            <el-dropdown trigger="click">
              <div class="admin-info">
                <el-avatar :size="32" class="admin-avatar">
                  {{ adminStore.adminInfo?.username?.charAt(0) || 'A' }}
                </el-avatar>
                <span class="admin-name">{{ adminStore.adminInfo?.username || '管理员' }}</span>
                <el-icon><ArrowDown /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="goToFrontend">
                    <el-icon><House /></el-icon>前台首页
                  </el-dropdown-item>
                  <el-dropdown-item>
                    <el-icon><Setting /></el-icon>个人设置
                  </el-dropdown-item>
                  <el-dropdown-item divided @click="handleLogout">
                    <el-icon><SwitchButton /></el-icon>退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>
        
        <el-main class="admin-main">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </el-main>
        
        <el-footer height="40px" class="admin-footer">
          <div class="footer-content">
            <span>© 2025 校园二手交易平台</span>
            <span>管理员后台系统</span>
          </div>
        </el-footer>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAdminStore } from '../../stores'
import { ElMessageBox, ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()
const adminStore = useAdminStore()

// 侧边栏折叠状态
const isCollapse = ref(false)

// 当前激活的菜单项
const activeMenu = computed(() => {
  return route.path
})

// 当前页面标题
const currentPageTitle = computed(() => {
  const pathMap = {
    '/admin/dashboard': '仪表盘',
    '/admin/products': '商品管理',
    '/admin/users': '用户管理',
    '/admin/orders': '订单管理',
    '/admin/messages': '消息管理'
  }
  return pathMap[route.path]
})

// 功能模块控制
const hasOrdersModule = ref(true)
const hasMessagesModule = ref(true)

// 切换侧边栏折叠状态
const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

// 前往前台首页
const goToFrontend = () => {
  router.push('/')
}

// 刷新页面
const refreshPage = () => {
  window.location.reload()
}

// 切换全屏
const toggleFullScreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
  } else {
    if (document.exitFullscreen) {
      document.exitFullscreen()
    }
  }
}

// 退出登录
const handleLogout = () => {
  ElMessageBox.confirm(
    '确定要退出登录吗？',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    adminStore.logout()
    ElMessage.success('已退出登录')
    router.push('/admin/login')
  }).catch(() => {
    // 取消操作
  })
}

// 在组件挂载时检查管理员登录状态
onMounted(() => {
  // 由于已经在store中实现了持久化，这里只需检查登录状态
  if (!adminStore.isLoggedIn) {
    router.push('/admin/login')
  }
})
</script>

<style scoped>
.admin-layout {
  height: 100vh;
  overflow: hidden;
}

.admin-container {
  height: 100%;
}

.admin-aside {
  background-color: #001529;
  color: #a6adb4;
  overflow: hidden;
  transition: width 0.3s;
  box-shadow: 2px 0 6px rgba(0, 21, 41, 0.35);
  z-index: 10;
}

.logo-container {
  height: 60px;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #fff;
  font-size: 18px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  overflow: hidden;
  transition: all 0.3s;
}

.logo-image {
  height: 32px;
  margin-right: 8px;
}

.admin-menu {
  border-right: none;
  height: calc(100% - 60px);
}

.admin-menu:not(.el-menu--collapse) {
  width: 220px;
}

.admin-header {
  background-color: #fff;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.header-left {
  display: flex;
  align-items: center;
}

.menu-toggle {
  font-size: 18px;
  cursor: pointer;
  margin-right: 20px;
  color: #606266;
  transition: color 0.3s;
}

.menu-toggle:hover {
  color: #409EFF;
}

.breadcrumb-container {
  margin-left: 8px;
}

.header-right {
  display: flex;
  align-items: center;
}

.header-actions {
  display: flex;
  margin-right: 20px;
}

.header-icon {
  font-size: 18px;
  padding: 0 10px;
  cursor: pointer;
  color: #606266;
  transition: color 0.3s;
}

.header-icon:hover {
  color: #409EFF;
}

.admin-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0 8px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.admin-info:hover {
  background-color: #f5f7fa;
}

.admin-avatar {
  margin-right: 8px;
  background-color: #1890ff;
}

.admin-name {
  margin-right: 5px;
  font-size: 14px;
}

.admin-main {
  background-color: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
}

.admin-footer {
  background-color: #fff;
  color: #606266;
  font-size: 12px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-top: 1px solid #f0f0f0;
}

.footer-content {
  display: flex;
  justify-content: space-between;
  width: 100%;
  padding: 0 20px;
}

/* 页面切换动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>