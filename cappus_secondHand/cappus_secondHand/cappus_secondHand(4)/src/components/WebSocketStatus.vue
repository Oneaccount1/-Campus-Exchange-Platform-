<template>
  <div class="websocket-status" :class="statusClass">
    <el-tooltip :content="tooltipContent" placement="bottom">
      <div class="status-indicator" @click="handleReconnect">
        <div class="status-dot"></div>
        <span v-if="showDetails">{{ status }}</span>
        <span v-if="showLatency && latency" class="latency">{{ latency }}ms</span>
      </div>
    </el-tooltip>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import { ElMessage } from 'element-plus';
import webSocketService from '../utils/websocket';

const props = defineProps({
  showDetails: {
    type: Boolean,
    default: false
  },
  showLatency: {
    type: Boolean,
    default: false
  }
});

// 连接状态
const connected = ref(webSocketService.isConnected);
const latency = ref(null);
const reconnectAttempts = ref(0);
const queuedMessages = ref(0);
const lastStatusChange = ref(new Date().toLocaleTimeString());
const readyStateText = ref(null);

// 状态类
const statusClass = computed(() => {
  return {
    'connected': connected.value,
    'disconnected': !connected.value,
    'with-text': props.showDetails
  };
});

// 状态文本
const status = computed(() => {
  return connected.value ? '已连接' : '已断开';
});

// 提示内容
const tooltipContent = computed(() => {
  const baseContent = connected.value ? 'WebSocket已连接' : 'WebSocket已断开';
  const readyInfo = readyStateText.value ? `，状态: ${readyStateText.value}` : '';
  const latencyInfo = latency.value ? `，延迟: ${latency.value}ms` : '';
  const reconnectInfo = !connected.value && reconnectAttempts.value > 0 
    ? `，重连尝试: ${reconnectAttempts.value}` 
    : '';
  const queueInfo = queuedMessages.value > 0 
    ? `，排队消息: ${queuedMessages.value}` 
    : '';
  const timeInfo = `，最后状态变更: ${lastStatusChange.value}`;
  
  return `${baseContent}${readyInfo}${latencyInfo}${reconnectInfo}${queueInfo}${timeInfo}`;
});

// 处理重连点击
const handleReconnect = () => {
  if (!connected.value) {
    ElMessage.info('正在尝试重新连接...');
    webSocketService.connect();
  } else {
    ElMessage.success('WebSocket连接正常');
  }
};

// 更新状态信息
const updateStatus = () => {
  const info = webSocketService.getConnectionInfo();
  connected.value = info.connected;
  latency.value = info.latency;
  reconnectAttempts.value = info.reconnectAttempts;
  queuedMessages.value = info.queuedMessages;
  
  // 添加WebSocket状态信息
  const readyStates = {
    [-1]: '未初始化',
    [WebSocket.CONNECTING]: '连接中',
    [WebSocket.OPEN]: '已连接',
    [WebSocket.CLOSING]: '关闭中',
    [WebSocket.CLOSED]: '已关闭'
  };
  
  readyStateText.value = readyStates[info.readyState] || '未知';
};

// 连接状态变更回调
const handleConnectionChange = (isConnected) => {
  connected.value = isConnected;
  lastStatusChange.value = new Date().toLocaleTimeString();
  updateStatus();
};

// 启动状态更新定时器
let statusInterval;
onMounted(() => {
  // 注册连接状态监听
  webSocketService.onConnection(handleConnectionChange);
  
  // 每5秒更新一次状态信息
  statusInterval = setInterval(updateStatus, 5000);
  
  // 初始化状态
  updateStatus();
});

// 清理
onBeforeUnmount(() => {
  if (statusInterval) {
    clearInterval(statusInterval);
  }
  // 移除连接状态监听
  webSocketService.offConnection(handleConnectionChange);
});
</script>

<style scoped>
.websocket-status {
  display: inline-flex;
  align-items: center;
  padding: 2px 6px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 5px;
}

.status-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  transition: background-color 0.3s ease;
}

.connected .status-dot {
  background-color: #67C23A;
  box-shadow: 0 0 5px #67C23A;
}

.disconnected .status-dot {
  background-color: #F56C6C;
  box-shadow: 0 0 5px #F56C6C;
}

.connected.with-text {
  color: #67C23A;
}

.disconnected.with-text {
  color: #F56C6C;
}

.latency {
  font-size: 0.8em;
  color: #909399;
  margin-left: 3px;
}

.websocket-status:hover {
  background-color: rgba(0, 0, 0, 0.05);
}
</style> 