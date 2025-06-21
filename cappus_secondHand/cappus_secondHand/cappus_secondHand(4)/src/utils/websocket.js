// WebSocket服务
class WebSocketService {
  constructor() {
    this.socket = null;
    this.isConnected = false;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 10; 
    this.reconnectTimeout = null;
    this.messageCallbacks = [];
    this.connectionCallbacks = [];
    this.enabled = true;
    this.connectionInProgress = false;
    this.reconnectDelay = 2000; // 起始重连延迟，2秒
    this.lastConnectionTime = 0;
    this.connectionStabilityThreshold = 30000; // 连接稳定阈值(30秒)
    this.connectionDebounceTimeout = null;
    this.latency = []; // 存储最近10次延迟数据
    this.messageQueue = []; // 离线消息队列
    this.debug = true; // 是否输出调试日志
    this.healthCheckInterval = null; // 健康检查定时器
    this.lastMessageTime = Date.now(); // 最后一次收到消息的时间
    this.disableReconnect = false; // 禁用自动重连的标志
  }

  // 输出调试日志
  debugLog(...args) {
    if (this.debug) {
      console.log('[WebSocket]', ...args);
    }
  }

  // 获取WebSocket URL
  getWebSocketUrl() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    let host;
    // 在开发环境使用后端服务器地址
    // Vite 环境变量 import.meta.env.DEV 为 true 时表示开发模式
    if (typeof import.meta !== 'undefined' && import.meta.env && import.meta.env.DEV) {
      host = 'localhost:8080';
    } else {
      host = window.location.host;
    }
    return `${protocol}//${host}/api/v1/messages/ws`;
  }

  // 连接WebSocket
  connect() {
    // 检查连接条件
    if (!this.enabled) {
      this.debugLog('WebSocket功能已禁用');
      return;
    }

    // 使用防抖处理，避免短时间内多次连接
    if (this.connectionDebounceTimeout) {
      clearTimeout(this.connectionDebounceTimeout);
    }

    this.connectionDebounceTimeout = setTimeout(() => {
      this._connectImpl();
    }, 300);
  }

  // 实际建立连接的内部方法
  _connectImpl() {
    // 防止重复连接
    if (this.connectionInProgress) {
      this.debugLog('WebSocket连接已在进行中，忽略此次连接请求');
      return;
    }

    // 检查连接频率
    const now = Date.now();
    if (now - this.lastConnectionTime < 2000) {
      this.debugLog(`WebSocket连接过于频繁，将延迟尝试`);
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = setTimeout(() => this.connect(), 2000);
      return;
    }
    this.lastConnectionTime = now;

    // 标记连接进行中
    this.connectionInProgress = true;

    // 关闭可能存在的连接
    if (this.socket) {
      try {
        this.debugLog('关闭已存在的WebSocket连接');
        this.socket.close(1000, "准备重新连接");
      } catch (error) {
        console.error('关闭WebSocket连接失败:', error);
      }
      this.socket = null;
      this.isConnected = false;
    }

    // 获取并验证token
    const token = localStorage.getItem('token');
    if (!token) {
      this.debugLog('未登录，无法建立WebSocket连接');
      this.connectionInProgress = false;
      return;
    }

    // 获取WebSocket URL
    const wsUrl = this.getWebSocketUrl();
    this.debugLog(`尝试连接WebSocket: ${wsUrl}`);
    
    // 创建WebSocket连接
    try {
      const cleanToken = token.replace(/^"(.*)"$/, '$1'); // 移除可能的引号
      const fullUrl = `${wsUrl}?token=${cleanToken}`;
      this.debugLog(`创建WebSocket连接: ${fullUrl}`);
      this.socket = new WebSocket(fullUrl);
      
      // 设置较长的超时时间
      this.socket.timeout = 30000;
    } catch (error) {
      console.error('WebSocket连接创建失败:', error);
      this.connectionInProgress = false;
      return;
    }

    // 连接事件处理
    this.socket.onopen = () => {
      this.debugLog('WebSocket连接已成功建立');
      this.isConnected = true;
      this.reconnectAttempts = 0;
      this.connectionInProgress = false;
      this.lastConnectionTime = Date.now();
      this.lastMessageTime = Date.now();
      this.reconnectDelay = 2000;
      
      // 启动连接健康检查
      this.startHealthCheck();
      
      // 触发连接回调
      this.notifyConnectionStatus(true);
      
      // 处理离线消息队列
      this.processMessageQueue();
    };

    // 接收消息处理
    this.socket.onmessage = this.onmessageHandler.bind(this);

    // 监听协议控制帧 - ping事件
    // 注意：大多数浏览器不会暴露ping/pong控制帧，但我们可以设置一个回调
    // 如果浏览器支持，当收到ping帧时更新最后活动时间
    if (this.socket.onping) {
      this.socket.onping = () => {
        this.lastMessageTime = Date.now();
        this.debugLog('收到WebSocket协议Ping帧');
      };
    }

    // 设置pong回调，更新最后活动时间
    // 注意：大多数浏览器会自动响应ping，但不一定暴露pong事件
    if (this.socket.onpong) {
      this.socket.onpong = () => {
        this.lastMessageTime = Date.now();
        this.debugLog('收到WebSocket协议Pong帧');
      };
    }

    // 关闭连接处理
    this.socket.onclose = (event) => {
      this.debugLog(`WebSocket连接已关闭, code: ${event.code}, reason: ${event.reason}`);
      this.isConnected = false;
      this.connectionInProgress = false;
      
      // 清除定时器
      this.stopHealthCheck();
      
      // 通知连接状态变更
      this.notifyConnectionStatus(false);
      
      // 计算连接持续时间
      const connectionDuration = Date.now() - this.lastConnectionTime;
      this.debugLog(`WebSocket连接持续时间: ${connectionDuration/1000}秒`);
      
      // 只有在未禁用重连并且不是正常关闭的情况下才尝试重连
      if (!this.disableReconnect && event.code !== 1000 && event.code !== 1001) {
        // 根据连接持续时间调整重连策略
        if (connectionDuration < this.connectionStabilityThreshold) {
          this.debugLog(`连接持续时间过短(${connectionDuration/1000}秒)，增加重连延迟`);
          this.reconnectDelay = Math.min(this.reconnectDelay * 1.5, 30000);
        }
        this.attemptReconnect();
      } else {
        this.debugLog('WebSocket正常关闭或已禁用重连，不进行重连');
        this.reconnectDelay = 2000; // 重置重连延迟
      }
    };

    // 错误处理
    this.socket.onerror = (error) => {
      console.error('WebSocket连接错误:', error);
      this.isConnected = false;
      this.connectionInProgress = false;
      
      const errorTime = Date.now();
      const connectionDuration = errorTime - this.lastConnectionTime;
      
      // 连接错误发生得太快，可能是连接问题
      if (connectionDuration < 1000) {
        this.reconnectDelay = Math.min(this.reconnectDelay * 2, 30000);
        this.debugLog(`连接错误发生过快，增加重连延迟至${this.reconnectDelay/1000}秒`);
      }
      
      // 只有在未禁用重连的情况下才尝试重连
      if (!this.disableReconnect) {
        this.attemptReconnect();
      }
    };
  }

  // 通知所有连接状态监听器
  notifyConnectionStatus(isConnected) {
    this.connectionCallbacks.forEach(callback => {
      try {
        callback(isConnected);
      } catch (error) {
        console.error('连接状态回调执行错误:', error);
      }
    });
  }

  // 处理离线消息队列
  processMessageQueue() {
    if (this.messageQueue.length > 0) {
      this.debugLog(`处理${this.messageQueue.length}条离线消息`);
      
      // 复制队列并清空原队列
      const queueCopy = [...this.messageQueue];
      this.messageQueue = [];
      
      // 逐个发送消息
      queueCopy.forEach(message => {
        this.sendMessage(message);
      });
    }
  }

  // 启用WebSocket功能
  enable() {
    this.enabled = true;
    this.debugLog('WebSocket功能已启用');
    this.connect();
  }

  // 禁用WebSocket功能
  disable() {
    this.enabled = false;
    this.debugLog('WebSocket功能已禁用');
    this.disconnect();
  }

  // 断开WebSocket连接
  disconnect() {
    this.debugLog('主动断开WebSocket连接');
    this.stopHealthCheck();
    
    // 设置禁用重连标志
    this.disableReconnect = true;
    
    // 清除所有定时器
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }
    
    if (this.healthCheckInterval) {
      clearInterval(this.healthCheckInterval);
      this.healthCheckInterval = null;
    }
    
    // 关闭连接
    if (this.socket) {
      try {
        this.socket.close(1000, '用户主动断开');
      } catch (error) {
        console.error('关闭WebSocket连接失败:', error);
      }
      this.socket = null;
    }
    
    this.isConnected = false;
    this.notifyConnectionStatus(false);
    
    // 延迟后重置禁用重连标志，允许将来连接时自动重连
    setTimeout(() => {
      this.disableReconnect = false;
    }, 5000);
  }

  // 发送消息
  sendMessage(message) {
    if (!message) {
      console.error('WebSocket发送消息为空');
      return false;
    }
    
    // 转换消息为标准格式
    let standardMessage = message;
    
    // 如果是字符串，假定为JSON字符串
    if (typeof message === 'string') {
      try {
        standardMessage = JSON.parse(message);
      } catch (e) {
        // 如果不是有效的JSON，创建一个原始内容消息
        standardMessage = {
          type: 'raw',
          data: {
            content: message
          }
        };
      }
    } else if (typeof message === 'object') {
      // 标准化旧消息格式
      if (!message.type) {
        standardMessage = {
          type: 'message',
          data: message
        };
      }
      
      // 标准化消息格式，确保所有数据都在data字段中
      if (!standardMessage.data) {
        standardMessage.data = {};
      }
      
      // 将接收者ID放到data中
      if (message.receiver_id && !standardMessage.data.recipientId) {
        standardMessage.data.recipientId = message.receiver_id;
      }
      if (message.receiverId && !standardMessage.data.recipientId) {
        standardMessage.data.recipientId = message.receiverId;
      }
    }
    
    // 序列化消息
    const messageStr = JSON.stringify(standardMessage);
    
    // 连接未就绪，加入队列
    if (!this.isConnected || !this.socket) {
      this.debugLog('WebSocket未连接，消息加入队列:', messageStr);
      this.messageQueue.push(standardMessage);
      
      // 触发连接
      if (!this.connectionInProgress) {
        this.connect();
      }
      return false;
    }
    
    // 发送消息
    try {
      this.socket.send(messageStr);
      this.debugLog('WebSocket消息已发送:', messageStr);
      return true;
    } catch (error) {
      console.error('WebSocket消息发送失败:', error);
      
      // 发送失败，加入队列
      this.messageQueue.push(standardMessage);
      
      // 检查连接状态并尝试重连
      if (this.socket.readyState !== WebSocket.OPEN) {
        this.debugLog('检测到WebSocket连接异常，尝试重新连接');
        this.isConnected = false;
        this.connectionInProgress = false;
        this.connect();
      }
      
      return false;
    }
  }
  
  // 启动健康检查
  startHealthCheck() {
    this.debugLog('启动WebSocket连接健康检查');
    
    // 清理可能存在的旧定时器
    this.stopHealthCheck();
    
    // 设置新的定时器，每30秒检查一次连接状态
    this.healthCheckInterval = setInterval(() => {
      // 只检查连接状态，不输出误导性的日志
      // 浏览器会自动处理ping/pong控制帧，但不会触发message事件
      // 因此lastMessageTime可能长时间不更新，即使连接是活跃的
      
      // 检查连接状态
      if (this.socket && this.socket.readyState !== WebSocket.OPEN) {
        this.debugLog('WebSocket连接已断开，尝试重新连接');
        this.socket.close(4000, '健康检查失败');
      }
    }, 30000);
  }
  
  // 停止健康检查
  stopHealthCheck() {
    if (this.healthCheckInterval) {
      clearInterval(this.healthCheckInterval);
      this.healthCheckInterval = null;
    }
  }

  // 尝试重新连接
  attemptReconnect() {
    if (!this.enabled) {
      this.debugLog('WebSocket功能已禁用，不进行重连');
      return;
    }
    
    // 增加重连次数
    this.reconnectAttempts++;
    
    // 检查是否超过最大重连次数
    if (this.reconnectAttempts > this.maxReconnectAttempts) {
      this.debugLog(`已达到最大重连次数(${this.maxReconnectAttempts})，停止重连`);
      // 通知重连失败
      this.notifyConnectionStatus(false);
      return;
    }
    
    // 计算当前重连延迟
    const delay = Math.min(this.reconnectDelay, 30000);
    
    this.debugLog(`计划${delay/1000}秒后进行第${this.reconnectAttempts}次重连尝试`);
    
    // 设置重连定时器
    clearTimeout(this.reconnectTimeout);
    this.reconnectTimeout = setTimeout(() => {
      this.debugLog(`执行第${this.reconnectAttempts}次重连`);
      this.connect();
    }, delay);
    
    // 指数增长重连延迟，但最多30秒
    this.reconnectDelay = Math.min(this.reconnectDelay * 1.5, 30000);
  }

  // 添加消息监听器
  onMessage(callback) {
    if (typeof callback === 'function' && !this.messageCallbacks.includes(callback)) {
      this.messageCallbacks.push(callback);
    }
  }

  // 移除消息监听器
  offMessage(callback) {
    const index = this.messageCallbacks.indexOf(callback);
    if (index !== -1) {
      this.messageCallbacks.splice(index, 1);
    }
  }

  // 添加连接状态监听器
  onConnection(callback) {
    if (typeof callback === 'function' && !this.connectionCallbacks.includes(callback)) {
      this.connectionCallbacks.push(callback);
      
      // 立即通知当前状态
      if (this.isConnected) {
        callback(true);
      }
    }
  }

  // 移除连接状态监听器
  offConnection(callback) {
    const index = this.connectionCallbacks.indexOf(callback);
    if (index !== -1) {
      this.connectionCallbacks.splice(index, 1);
    }
  }

  // 获取连接信息
  getConnectionInfo() {
    // 检查实际连接状态
    if (this.socket) {
      this.isConnected = this.socket.readyState === WebSocket.OPEN;
    }
    
    return {
      connected: this.isConnected,
      reconnectAttempts: this.reconnectAttempts,
      lastMessageTime: this.lastMessageTime,
      latency: this.getAverageLatency(),
      queuedMessages: this.messageQueue.length,
      readyState: this.socket ? this.socket.readyState : -1
    };
  }
  
  // 获取平均延迟
  getAverageLatency() {
    if (this.latency.length === 0) return null;
    const sum = this.latency.reduce((a, b) => a + b, 0);
    return Math.round(sum / this.latency.length);
  }
  
  // 接收消息处理
  onmessageHandler(event) {
    this.lastMessageTime = Date.now();
    
    let data = event.data;
    let message = null;
    
    try {
      // 尝试解析JSON消息
      if (typeof data === 'string') {
        message = JSON.parse(data);
      } else {
        message = data;
      }
      
      // 计算网络延迟 - 使用服务器时间戳
      if (message && message.timestamp) {
        const latency = Date.now() - message.timestamp;
        if (latency > 0 && latency < 60000) { // 忽略异常值
          this.latency.push(latency);
          if (this.latency.length > 10) {
            this.latency.shift();
          }
        }
      }
      
      // 处理消息 - 标准化消息格式
      if (message && typeof message === 'object') {
        // 标准化消息格式
        if (!message.type && message.sender_id && message.content) {
          message = {
            type: 'message',
            data: message
          };
        }
        
        // 保持向后兼容 - 确保有data字段
        if (!message.data) {
          message.data = {};
        }
        
        // 处理兼容性字段 - 维持旧字段的同时添加新标准字段
        if (message.sender_id && !message.data.senderId) {
          message.data.senderId = message.sender_id;
        }
        if (message.receiver_id && !message.data.recipientId) {
          message.data.recipientId = message.receiver_id;
        }
        if (message.content && !message.data.content) {
          message.data.content = message.content;
        }
      }
    } catch (e) {
      this.debugLog('收到非JSON消息或解析失败:', data);
      // 创建原始消息
      message = {
        type: 'raw',
        data: { content: data }
      };
    }
    
    // 处理消息
    this.messageCallbacks.forEach(callback => {
      try {
        callback(message);
      } catch (error) {
        console.error('WebSocket消息处理错误:', error);
      }
    });
  }
}

// 创建单例
const wsService = new WebSocketService();

export default wsService; 