import { defineStore } from 'pinia'

export const useAdminStore = defineStore('admin', {
  state: () => {
    // 从localStorage获取管理员信息和登录状态
    const storedAdminInfo = localStorage.getItem('admin_info')
    const storedToken = localStorage.getItem('admin_token')
    
    return {
      adminInfo: storedAdminInfo ? JSON.parse(storedAdminInfo) : null,
      isLoggedIn: !!storedToken,
      token: storedToken || ''
    }
  },
  actions: {
    setAdminInfo(adminInfo) {
      this.adminInfo = adminInfo
      this.isLoggedIn = true
      // 将管理员信息保存到localStorage
      localStorage.setItem('admin_info', JSON.stringify(adminInfo))
    },
    setToken(token) {
      this.token = token
      localStorage.setItem('admin_token', token)
    },
    logout() {
      this.adminInfo = null
      this.isLoggedIn = false
      this.token = ''
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_info')
    }
  }
})