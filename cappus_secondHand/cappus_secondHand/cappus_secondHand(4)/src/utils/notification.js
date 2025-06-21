import { ElNotification } from 'element-plus'

// 消息通知服务
class NotificationService {
  constructor() {
    this.unreadCount = 0;
    this.unreadCountCallbacks = [];
    this.permission = 'default';
    this.shouldRequestPermission = false;
    
    // 检查并请求通知权限
    this.checkNotificationPermission();
  }

  // 检查并请求通知权限
  checkNotificationPermission() {
    if (!('Notification' in window)) {
      console.log('Browser does not support desktop notifications');
      return;
    }

    this.permission = Notification.permission;
    
    if (this.permission === 'default') {
      // Only request permission on user interaction
      this.shouldRequestPermission = true;
    }
  }

  // 设置未读消息数量
  setUnreadCount(count) {
    this.unreadCount = count;
    
    // 更新页面标题
    this.updatePageTitle();
    
    // 触发回调
    this.unreadCountCallbacks.forEach(callback => callback(count));
  }

  // 增加未读消息数量
  incrementUnreadCount(amount = 1) {
    this.unreadCount += amount;
    
    // 更新页面标题
    this.updatePageTitle();
    
    // 触发回调
    this.unreadCountCallbacks.forEach(callback => callback(this.unreadCount));
  }

  // 清除未读消息数量
  clearUnreadCount() {
    this.unreadCount = 0;
    
    // 更新页面标题
    this.updatePageTitle();
    
    // 触发回调
    this.unreadCountCallbacks.forEach(callback => callback(0));
  }

  // 获取未读消息数量
  getUnreadCount() {
    return this.unreadCount;
  }

  // 监听未读消息数量变化
  onUnreadCountChange(callback) {
    if (typeof callback === 'function') {
      this.unreadCountCallbacks.push(callback);
      
      // 立即触发一次当前值
      callback(this.unreadCount);
    }
  }

  // 取消监听未读消息数量变化
  offUnreadCountChange(callback) {
    this.unreadCountCallbacks = this.unreadCountCallbacks.filter(cb => cb !== callback);
  }

  // 更新页面标题
  updatePageTitle() {
    const originalTitle = '校园二手交易平台';
    
    if (this.unreadCount > 0) {
      document.title = `(${this.unreadCount}) ${originalTitle}`;
    } else {
      document.title = originalTitle;
    }
  }

  // 显示消息通知（Element Plus通知）
  showNotification(options) {
    ElNotification({
      title: options.title || '新消息',
      message: options.message || '',
      type: options.type || 'info',
      position: options.position || 'top-right',
      duration: options.duration || 4500,
      showClose: options.showClose !== undefined ? options.showClose : true,
      onClick: options.onClick
    });
  }

  // 显示桌面通知
  showDesktopNotification(options) {
    if (!('Notification' in window) || this.permission !== 'granted') {
      // 如果不支持桌面通知或没有权限，则使用Element Plus通知
      this.showNotification(options);
      return;
    }
    
    try {
      const notification = new Notification(options.title || '新消息', {
        body: options.message || '',
        icon: options.icon || '/favicon.ico',
        tag: options.tag,
        silent: options.silent !== undefined ? options.silent : false
      });
      
      // 点击通知时的回调
      if (typeof options.onClick === 'function') {
        notification.onclick = options.onClick;
      }
    } catch (error) {
      console.error('显示桌面通知失败:', error);
      
      // 失败时使用Element Plus通知作为备选
      this.showNotification(options);
    }
  }

  // Add new method
  requestPermissionOnUserAction() {
    if (this.shouldRequestPermission) {
      Notification.requestPermission().then(permission => {
        this.permission = permission;
        this.shouldRequestPermission = false;
      });
    }
  }
}

// 创建单例
const notificationService = new NotificationService();

export default notificationService;