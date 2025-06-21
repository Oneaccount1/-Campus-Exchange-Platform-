import { defineStore } from 'pinia'
import { useAdminStore } from './admin'
import { useMessageStore } from './message'

export const useUserStore = defineStore('user', {
  state: () => {
    // 尝试从localStorage恢复用户信息
    let userInfo = null;
    let isLoggedIn = false;
    let token = '';
    
    try {
      token = localStorage.getItem('token') || '';
      const savedUserInfo = localStorage.getItem('userInfo');
      
      // 确保token没有双引号
      if (token) {
        token = token.replace(/^"(.*)"$/, '$1');
      }
      
      if (savedUserInfo && token) {
        userInfo = JSON.parse(savedUserInfo);
        isLoggedIn = true;
      }
    } catch (e) {
      console.error('恢复用户信息失败:', e);
      localStorage.removeItem('userInfo');
      localStorage.removeItem('token');
    }
    
    return {
      userInfo,
      isLoggedIn,
      token
    };
  },
  actions: {
    setUserInfo(userInfo) {
      this.userInfo = userInfo;
      this.isLoggedIn = true;
      // 持久化保存用户信息
      localStorage.setItem('userInfo', JSON.stringify(userInfo));
    },
    setToken(token) {
      if (!token) {
        console.error('试图设置空token');
        return;
      }
      
      this.token = token;
      localStorage.setItem('token', token);
    },
    logout() {
      this.userInfo = null;
      this.isLoggedIn = false;
      this.token = '';
      // 清除本地存储
      localStorage.removeItem('token');
      localStorage.removeItem('userInfo');
    },
    // 登录成功后调用
    loginSuccess(userData, token) {
      this.setUserInfo(userData);
      this.setToken(token);
      
      console.log('用户登录成功，用户信息已保存到store');
    },
    
    // 连接WebSocket
    async connectWebSocket() {
      // 导入websocket服务
      const webSocketService = await import('../utils/websocket').then(m => m.default);
      
      // 如果已连接，则不重复连接
      if (webSocketService.isConnected) {
        console.log('WebSocket已连接，无需重复连接');
        return;
      }
      
      // 建立新连接
      webSocketService.connect();
      
      return new Promise((resolve) => {
        // 监听连接成功
        const onConnected = (connected) => {
          if (connected) {
            webSocketService.offConnection(onConnected);
            resolve(true);
          }
        };
        
        webSocketService.onConnection(onConnected);
        
        // 5秒超时
        setTimeout(() => {
          webSocketService.offConnection(onConnected);
          resolve(false);
        }, 5000);
      });
    },
    
    // 登出时调用
    async logout() {
      // 断开WebSocket
      const webSocketService = await import('../utils/websocket').then(m => m.default);
      webSocketService.disconnect();
      
      // 清除数据
      this.userInfo = null;
      this.isLoggedIn = false;
      this.token = '';
      
      // 清除本地存储
      localStorage.removeItem('token');
      localStorage.removeItem('userInfo');
      
      console.log('用户登出，WebSocket连接已关闭');
    }
  }
})

export const useProductStore = defineStore('product', {
  state: () => ({
    products: [],
    categories: [
      { id: 1, name: '电子产品' },
      { id: 2, name: '书籍教材' },
      { id: 3, name: '生活用品' },
      { id: 4, name: '服装鞋帽' },
      { id: 5, name: '运动器材' },
      { id: 6, name: '其他' }
    ],
    currentProduct: null
  }),
  actions: {
    setProducts(products) {
      this.products = products
    },
    setCurrentProduct(product) {
      this.currentProduct = product
    },
    addProduct(product) {
      this.products.unshift(product)
    }
  }
})

export { useAdminStore, useMessageStore }