<template>
  <div class="admin-messages-container">
    <div class="page-header">
      <h2 class="page-title">消息管理</h2>
      <el-button type="primary" @click="sendSystemMessage">
        <el-icon><Message /></el-icon>发送系统消息
      </el-button>
    </div>

    <el-tabs v-model="activeTab" class="message-tabs">
      <el-tab-pane label="消息列表" name="messages">
        <el-card shadow="hover" class="filter-card">
          <div class="filter-container">
            <div class="filter-item">
              <el-input
                v-model="searchQuery"
                placeholder="搜索内容/用户"
                class="search-input"
                clearable
                @clear="handleSearch"
                @input="handleSearch"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
            </div>
            <div class="filter-item">
              <el-select v-model="typeFilter" placeholder="消息类型" clearable @change="handleSearch">
                <el-option label="用户消息" value="user" />
                <el-option label="系统消息" value="system" />
              </el-select>
            </div>
            <div class="filter-item">
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                @change="handleSearch"
              />
            </div>
            <div class="filter-actions">
              <el-button type="primary" @click="handleSearch">
                <el-icon><Search /></el-icon>搜索
              </el-button>
              <el-button @click="resetFilters">
                <el-icon><RefreshLeft /></el-icon>重置
              </el-button>
            </div>
          </div>
        </el-card>
          
        <el-card shadow="hover" class="data-card">
          <div class="table-toolbar">
            <div class="toolbar-left">
              <el-tag type="success">共 {{ totalMessages }} 条记录</el-tag>
            </div>
            <div class="toolbar-right">
              <el-tooltip content="刷新" placement="top">
                <el-button circle @click="refreshTable">
                  <el-icon><Refresh /></el-icon>
                </el-button>
              </el-tooltip>
            </div>
          </div>
          
          <el-table 
            :data="messagesList" 
            style="width: 100%" 
            v-loading="loading"
            :header-cell-style="{ background: '#f5f7fa' }"
            border
          >
            <el-table-column label="消息类型" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.type === 'system' ? 'danger' : 'primary'">
                  {{ scope.row.type === 'system' ? '系统消息' : '用户消息' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="发送者" width="120">
              <template #default="scope">
                <div class="user-info">
                  {{ scope.row.sender }}
                  <el-tag v-if="scope.row.type === 'system'" size="small" effect="plain" type="danger">系统</el-tag>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="接收者" width="120">
              <template #default="scope">
                {{ scope.row.receiver }}
              </template>
            </el-table-column>
            <el-table-column prop="content" label="消息内容" show-overflow-tooltip />
            <el-table-column label="发送时间" width="170">
              <template #default="scope">
                {{ formatDateTime(scope.row.createTime) }}
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === '已读' ? 'success' : 'info'">
                  {{ scope.row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="scope">
                <el-button type="primary" size="small" @click="viewMessage(scope.row)">
                  <el-icon><View /></el-icon>查看
                </el-button>
                <el-button type="danger" size="small" @click="deleteMessage(scope.row.id)">
                  <el-icon><Delete /></el-icon>删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <div class="pagination-container">
            <el-pagination
              :current-page="currentPage"
              :page-size="pageSize"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="totalMessages"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
              background
            />
          </div>
        </el-card>
      </el-tab-pane>
      
      <el-tab-pane label="会话列表" name="conversations">
        <el-card shadow="hover" class="filter-card">
          <div class="filter-container">
            <div class="filter-item">
              <el-input
                v-model="convSearchQuery"
                placeholder="搜索用户"
                class="search-input"
                clearable
                @clear="handleConvSearch"
                @input="handleConvSearch"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
            </div>
            <div class="filter-actions">
              <el-button type="primary" @click="handleConvSearch">
                <el-icon><Search /></el-icon>搜索
              </el-button>
              <el-button @click="resetConvFilters">
                <el-icon><RefreshLeft /></el-icon>重置
              </el-button>
            </div>
          </div>
        </el-card>
        
        <el-card shadow="hover" class="data-card">
          <div class="table-toolbar">
            <div class="toolbar-left">
              <el-tag type="success">共 {{ totalConversations }} 条记录</el-tag>
            </div>
            <div class="toolbar-right">
              <el-tooltip content="刷新" placement="top">
                <el-button circle @click="refreshConversations">
                  <el-icon><Refresh /></el-icon>
                </el-button>
              </el-tooltip>
            </div>
          </div>
          
          <el-table 
            :data="conversationsList" 
            style="width: 100%" 
            v-loading="convLoading"
            :header-cell-style="{ background: '#f5f7fa' }"
            border
          >
            <el-table-column label="用户1" width="200">
              <template #default="scope">
                <div class="user-info-cell">
                  <el-avatar :size="32" :src="scope.row.user1Avatar">
                    {{ scope.row.user1Name?.charAt(0) || '?' }}
                  </el-avatar>
                  <span>{{ scope.row.user1Name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="用户2" width="200">
              <template #default="scope">
                <div class="user-info-cell">
                  <el-avatar :size="32" :src="scope.row.user2Avatar">
                    {{ scope.row.user2Name?.charAt(0) || '?' }}
                  </el-avatar>
                  <span>{{ scope.row.user2Name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="lastMessage" label="最新消息" show-overflow-tooltip />
            <el-table-column label="最后时间" width="170">
              <template #default="scope">
                {{ formatDateTime(scope.row.lastTime) }}
              </template>
            </el-table-column>
            <el-table-column label="未读" width="100">
              <template #default="scope">
                <el-badge :value="scope.row.unreadCount || 0" :hidden="!scope.row.unreadCount" type="danger" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="scope">
                <el-button type="primary" size="small" @click="viewConversation(scope.row)">
                  <el-icon><View /></el-icon>查看对话
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          
          <div class="pagination-container">
            <el-pagination
              :current-page="convCurrentPage"
              :page-size="convPageSize"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="totalConversations"
              @size-change="handleConvSizeChange"
              @current-change="handleConvCurrentChange"
              background
            />
          </div>
        </el-card>
      </el-tab-pane>
    </el-tabs>
    
    <!-- 消息详情对话框 -->
    <el-dialog
      v-model="messageDetailVisible"
      title="消息详情"
      width="600px"
      destroy-on-close
    >
      <div v-if="currentMessage" class="message-detail">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="消息类型">
            <el-tag :type="currentMessage.type === 'system' ? 'danger' : 'primary'">
              {{ currentMessage.type === 'system' ? '系统消息' : '用户消息' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="发送者">{{ currentMessage.sender }}</el-descriptions-item>
          <el-descriptions-item label="接收者">{{ currentMessage.receiver }}</el-descriptions-item>
          <el-descriptions-item label="发送时间">{{ formatDateTime(currentMessage.createTime) }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="currentMessage.status === '已读' ? 'success' : 'info'">
              {{ currentMessage.status }}
            </el-tag>
            <span v-if="currentMessage.readTime">&nbsp;({{ formatDateTime(currentMessage.readTime) }})</span>
          </el-descriptions-item>
        </el-descriptions>
        
        <el-divider content-position="center">消息内容</el-divider>
        
        <div class="message-content">
          {{ currentMessage.content }}
        </div>
      </div>
    </el-dialog>
    
    <!-- 对话详情对话框 -->
    <el-dialog
      v-model="conversationDetailVisible"
      title="消息会话"
      width="800px"
      destroy-on-close
    >
      <div v-if="currentConversation" class="conversation-detail">
        <div class="conversation-header">
          <div class="user-info-cell">
            <el-avatar :size="40" :src="currentConversation.user1Avatar">
              {{ currentConversation.user1Name?.charAt(0) || '?' }}
            </el-avatar>
            <span>{{ currentConversation.user1Name }}</span>
          </div>
          <el-divider direction="vertical" />
          <div class="user-info-cell">
            <el-avatar :size="40" :src="currentConversation.user2Avatar">
              {{ currentConversation.user2Name?.charAt(0) || '?' }}
            </el-avatar>
            <span>{{ currentConversation.user2Name }}</span>
          </div>
        </div>
        
        <div class="conversation-body" v-loading="historyLoading">
          <div v-if="messageHistory.length === 0" class="empty-history">
            暂无消息记录
          </div>
          <div v-else class="message-list">
            <div 
              v-for="(msg, index) in messageHistory" 
              :key="index" 
              class="message-item"
              :class="{'self-message': msg.senderId === currentConversation.user1Id}"
            >
              <el-avatar :size="32" :src="msg.senderAvatar" class="message-avatar">
                {{ msg.sender?.charAt(0) || '?' }}
              </el-avatar>
              <div class="message-bubble">
                <div class="message-sender">{{ msg.sender }} <span class="message-time">{{ formatDateTime(msg.createTime) }}</span></div>
                <div class="message-text">{{ msg.content }}</div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="conversation-footer">
          <el-pagination
            :current-page="historyCurrentPage"
            :page-size="historyPageSize"
            layout="prev, pager, next"
            :total="totalHistory"
            @current-change="loadMoreHistory"
            small
            background
          />
        </div>
      </div>
    </el-dialog>
    
    <!-- 发送系统消息对话框 -->
    <el-dialog
      v-model="sendMessageVisible"
      title="发送系统消息"
      width="600px"
      destroy-on-close
    >
      <el-form :model="messageForm" :rules="messageRules" ref="messageFormRef" label-width="100px">
        <el-form-item label="接收者类型">
          <el-radio-group v-model="messageForm.receiverType">
            <el-radio label="all">所有用户</el-radio>
            <el-radio label="single">指定用户</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="接收者ID" v-if="messageForm.receiverType === 'single'" prop="receiverId">
          <el-input v-model.number="messageForm.receiverId" placeholder="请输入用户ID" type="number" />
        </el-form-item>
        <el-form-item label="消息标题" prop="title">
          <el-input v-model="messageForm.title" placeholder="请输入消息标题" />
        </el-form-item>
        <el-form-item label="消息内容" prop="content">
          <el-input 
            v-model="messageForm.content" 
            type="textarea" 
            :rows="5" 
            placeholder="请输入消息内容" 
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="sendMessageVisible = false">取消</el-button>
          <el-button type="primary" @click="submitSystemMessage" :loading="sending">发送</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { 
  getMessages, 
  getConversations, 
  getMessageHistory, 
  sendSystemMessageApi, 
  deleteMessage as deleteMessageApi 
} from '../../api/admin'

// 标签页控制
const activeTab = ref('messages')

// 消息列表数据
const searchQuery = ref('')
const typeFilter = ref('')
const dateRange = ref(null)
const loading = ref(false)
const messagesList = ref([])
const totalMessages = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 会话列表数据
const convSearchQuery = ref('')
const convLoading = ref(false)
const conversationsList = ref([])
const totalConversations = ref(0)
const convCurrentPage = ref(1)
const convPageSize = ref(10)

// 会话历史数据
const messageHistory = ref([])
const historyLoading = ref(false)
const historyCurrentPage = ref(1)
const historyPageSize = ref(20)
const totalHistory = ref(0)

// 对话框控制
const messageDetailVisible = ref(false)
const currentMessage = ref(null)
const conversationDetailVisible = ref(false)
const currentConversation = ref(null)
const sendMessageVisible = ref(false)
const sending = ref(false)

// 发送消息表单
const messageFormRef = ref(null)
const messageForm = reactive({
  receiverType: 'all',
  receiverId: '',
  title: '',
  content: ''
})
const messageRules = {
  receiverId: [
    { required: true, message: '请输入接收者ID', trigger: 'blur', type: 'number' }
  ],
  title: [
    { required: true, message: '请输入消息标题', trigger: 'blur' },
    { min: 2, max: 50, message: '标题长度应为2-50个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入消息内容', trigger: 'blur' },
    { min: 2, max: 500, message: '内容长度应为2-500个字符', trigger: 'blur' }
  ]
}

// 获取消息列表
const loadMessages = async () => {
  loading.value = true
  try {
    // 构建查询参数
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      search: searchQuery.value || undefined,
      type: typeFilter.value || undefined,
      startDate: dateRange.value && dateRange.value.length === 2 ? formatDate(dateRange.value[0]) : undefined,
      endDate: dateRange.value && dateRange.value.length === 2 ? formatDate(dateRange.value[1]) : undefined
    }
    
    // 添加日期范围参数
    if (dateRange.value && dateRange.value.length === 2) {
      params.startDate = formatDate(dateRange.value[0])
      params.endDate = formatDate(dateRange.value[1])
    }
    
    // 调用API
    const response = await getMessages(params)
    
    if (response && response.code === 200) {
      messagesList.value = response.data.list || []
      totalMessages.value = response.data.total || 0
    } else {
      // 如果API未完成，使用模拟数据
      useMessagesMockData()
    }
  } catch (error) {
    console.error('加载消息列表失败:', error)
    // 使用模拟数据
    useMessagesMockData()
  } finally {
    loading.value = false
  }
}

// 使用模拟的消息列表数据
const useMessagesMockData = () => {
  messagesList.value = [
    {
      id: 1,
      type: 'user',
      senderId: 1,
      sender: '张三',
      receiverId: 2,
      receiver: '李四',
      content: '你好，请问商品还在吗？',
      createTime: '2024-05-20 14:30:22',
      status: '已读',
      readTime: '2024-05-20 14:35:10'
    },
    {
      id: 2,
      type: 'system',
      senderId: 0,
      sender: '系统',
      receiverId: 1,
      receiver: '张三',
      content: '您的订单已完成',
      createTime: '2024-05-20 14:40:00',
      status: '未读',
      readTime: null
    },
    {
      id: 3,
      type: 'user',
      senderId: 2,
      sender: '李四',
      receiverId: 1,
      receiver: '张三',
      content: '在的，还可以购买',
      createTime: '2024-05-20 14:31:15',
      status: '已读',
      readTime: '2024-05-20 14:32:20'
    },
    {
      id: 4,
      type: 'user',
      senderId: 3,
      sender: '王五',
      receiverId: 4,
      receiver: '赵六',
      content: '这个价格能便宜点吗？',
      createTime: '2024-05-19 09:20:45',
      status: '已读',
      readTime: '2024-05-19 09:25:30'
    },
    {
      id: 5,
      type: 'system',
      senderId: 0,
      sender: '系统',
      receiverId: 0,
      receiver: '所有用户',
      content: '系统将于2024-05-25 02:00-04:00进行维护，期间将无法访问',
      createTime: '2024-05-18 10:00:00',
      status: '未读',
      readTime: null
    }
  ]
  totalMessages.value = 150
}

// 获取会话列表
const loadConversations = async () => {
  convLoading.value = true
  try {
    // 构建查询参数
    const params = {
      page: convCurrentPage.value,
      pageSize: convPageSize.value,
      search: convSearchQuery.value || undefined
    }
    
    // 调用API
    const response = await getConversations(params)
    
    if (response && response.code === 200) {
      conversationsList.value = response.data.list || []
      totalConversations.value = response.data.total || 0
    } else {
      // 如果API未完成，使用模拟数据
      useConversationsMockData()
    }
  } catch (error) {
    console.error('加载会话列表失败:', error)
    // 使用模拟数据
    useConversationsMockData()
  } finally {
    convLoading.value = false
  }
}

// 使用模拟的会话列表数据
const useConversationsMockData = () => {
  conversationsList.value = [
    {
      id: 1,
      user1Id: 1,
      user1Name: '张三',
      user1Avatar: 'https://via.placeholder.com/40',
      user2Id: 2,
      user2Name: '李四',
      user2Avatar: 'https://via.placeholder.com/40',
      lastMessage: '你好，请问商品还在吗？',
      lastTime: '2024-05-20 14:30:22',
      unreadCount: 0
    },
    {
      id: 2,
      user1Id: 3,
      user1Name: '王五',
      user1Avatar: 'https://via.placeholder.com/40',
      user2Id: 4,
      user2Name: '赵六',
      user2Avatar: 'https://via.placeholder.com/40',
      lastMessage: '这个价格能便宜点吗？',
      lastTime: '2024-05-19 09:20:45',
      unreadCount: 1
    },
    {
      id: 3,
      user1Id: 5,
      user1Name: '孙七',
      user1Avatar: 'https://via.placeholder.com/40',
      user2Id: 6,
      user2Name: '周八',
      user2Avatar: 'https://via.placeholder.com/40',
      lastMessage: '好的，我们约个时间地点交易吧',
      lastTime: '2024-05-18 16:42:35',
      unreadCount: 3
    }
  ]
  totalConversations.value = 45
}

// 查看消息详情
const viewMessage = (message) => {
  currentMessage.value = message
  messageDetailVisible.value = true
}

// 查看会话详情
const viewConversation = async (conversation) => {
  currentConversation.value = conversation
  conversationDetailVisible.value = true
  historyCurrentPage.value = 1
  await loadMessageHistory(conversation)
}

// 加载消息历史
const loadMessageHistory = async (conversation) => {
  if (!conversation) return
  
  historyLoading.value = true
  try {
    // 构建查询参数
    const params = {
      user1Id: conversation.user1Id,
      user2Id: conversation.user2Id,
      page: historyCurrentPage.value,
      pageSize: historyPageSize.value
    }
    
    // 调用API
    const response = await getMessageHistory(params)
    
    if (response && response.code === 200) {
      messageHistory.value = response.data.list || []
      totalHistory.value = response.data.total || 0
    } else {
      // 如果API未完成，使用模拟数据
      useHistoryMockData(conversation)
    }
  } catch (error) {
    console.error('加载消息历史失败:', error)
    // 使用模拟数据
    useHistoryMockData(conversation)
  } finally {
    historyLoading.value = false
  }
}

// 使用模拟的消息历史数据
const useHistoryMockData = (conversation) => {
  const mockHistory = [
    {
      id: 1,
      senderId: conversation.user1Id,
      sender: conversation.user1Name,
      senderAvatar: conversation.user1Avatar,
      receiverId: conversation.user2Id,
      content: '你好，请问商品还在吗？',
      createTime: '2024-05-20 14:30:22',
      status: '已读'
    },
    {
      id: 2,
      senderId: conversation.user2Id,
      sender: conversation.user2Name,
      senderAvatar: conversation.user2Avatar,
      receiverId: conversation.user1Id,
      content: '在的，还可以购买',
      createTime: '2024-05-20 14:31:15',
      status: '已读'
    },
    {
      id: 3,
      senderId: conversation.user1Id,
      sender: conversation.user1Name,
      senderAvatar: conversation.user1Avatar,
      receiverId: conversation.user2Id,
      content: '价格能便宜一点吗？',
      createTime: '2024-05-20 14:32:30',
      status: '已读'
    },
    {
      id: 4,
      senderId: conversation.user2Id,
      sender: conversation.user2Name,
      senderAvatar: conversation.user2Avatar,
      receiverId: conversation.user1Id,
      content: '可以便宜50元',
      createTime: '2024-05-20 14:33:40',
      status: '已读'
    },
    {
      id: 5,
      senderId: conversation.user1Id,
      sender: conversation.user1Name,
      senderAvatar: conversation.user1Avatar,
      receiverId: conversation.user2Id,
      content: '好的，我买了，我们约个时间地点交易吧',
      createTime: '2024-05-20 14:35:20',
      status: '已读'
    }
  ]
  
  messageHistory.value = mockHistory
  totalHistory.value = 28
}

// 加载更多历史消息
const loadMoreHistory = () => {
  if (currentConversation.value) {
    loadMessageHistory(currentConversation.value)
  }
}

// 发送系统消息对话框
const sendSystemMessage = () => {
  messageForm.receiverType = 'all'
  messageForm.receiverId = ''
  messageForm.title = ''
  messageForm.content = ''
  sendMessageVisible.value = true
}

// 提交系统消息
const submitSystemMessage = async () => {
  if (!messageFormRef.value) return
  
  await messageFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    sending.value = true
    try {
      // 构建请求体
      const data = {
        receiverId: messageForm.receiverType === 'single' ? messageForm.receiverId : 0,
        content: messageForm.content,
        title: messageForm.title
      }
      
      // 调用API
      const response = await sendSystemMessageApi(data)
      
      if (response && response.code === 200) {
        ElMessage.success('系统消息发送成功')
        sendMessageVisible.value = false
        // 重新加载消息列表
        loadMessages()
      } else {
        ElMessage.error(response?.message || '发送失败，请稍后再试')
      }
    } catch (error) {
      console.error('发送系统消息失败:', error)
      ElMessage.error('发送失败，请稍后再试')
    } finally {
      sending.value = false
    }
  })
}

// 删除消息
const deleteMessage = async (messageId) => {
  ElMessageBox.confirm(
    '确定要删除该消息吗？删除后无法恢复。',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      // 调用API
      const response = await deleteMessageApi(messageId)
      
      if (response && response.code === 200) {
        ElMessage.success('消息删除成功')
        // 重新加载消息列表
        loadMessages()
      } else {
        ElMessage.error(response?.message || '删除失败，请稍后再试')
      }
    } catch (error) {
      console.error('删除消息失败:', error)
      ElMessage.error('删除失败，请稍后再试')
    }
  }).catch(() => {
    // 取消删除操作
  })
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  loadMessages()
}

// 会话搜索处理
const handleConvSearch = () => {
  convCurrentPage.value = 1
  loadConversations()
}

// 重置筛选条件
const resetFilters = () => {
  searchQuery.value = ''
  typeFilter.value = ''
  dateRange.value = null
  currentPage.value = 1
  loadMessages()
}

// 重置会话筛选条件
const resetConvFilters = () => {
  convSearchQuery.value = ''
  convCurrentPage.value = 1
  loadConversations()
}

// 刷新表格
const refreshTable = () => {
  loadMessages()
}

// 刷新会话列表
const refreshConversations = () => {
  loadConversations()
}

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size
  loadMessages()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  loadMessages()
}

// 会话分页处理
const handleConvSizeChange = (size) => {
  convPageSize.value = size
  loadConversations()
}

const handleConvCurrentChange = (page) => {
  convCurrentPage.value = page
  loadConversations()
}

// 格式化日期
const formatDate = (date) => {
  if (!date) return ''
  
  // 确保date是Date对象
  const d = date instanceof Date ? date : new Date(date)
  
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  
  return `${year}-${month}-${day}`
}

// 格式化日期时间为YYYY-MM-DD HH:MM:SS
const formatDateTime = (dateTime) => {
  if (!dateTime) return ''
  
  // 确保dateTime是Date对象
  const d = dateTime instanceof Date ? dateTime : new Date(dateTime)
  
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  const seconds = String(d.getSeconds()).padStart(2, '0')
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 根据标签页加载相应数据
watch(activeTab, (newValue) => {
  if (newValue === 'messages') {
    loadMessages()
  } else if (newValue === 'conversations') {
    loadConversations()
  }
})

// 加载初始数据
onMounted(() => {
  loadMessages()
})
</script>

<style scoped>
.admin-messages-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  font-size: 24px;
  color: #303133;
  margin: 0;
}

.message-tabs {
  margin-bottom: 20px;
}

.filter-card {
  margin-bottom: 20px;
}

.filter-container {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  align-items: flex-start;
}

.filter-item {
  min-width: 200px;
  flex: 1;
}

.filter-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.search-input {
  width: 100%;
}

.data-card {
  margin-bottom: 20px;
}

.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.toolbar-left, .toolbar-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.action-group {
  margin-left: 15px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 5px;
}

.user-info-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.message-detail {
  padding: 10px;
}

.message-content {
  background-color: #f8f8f8;
  padding: 15px;
  border-radius: 4px;
  margin-top: 15px;
  white-space: pre-wrap;
  word-break: break-word;
}

.conversation-detail {
  display: flex;
  flex-direction: column;
  height: 500px;
}

.conversation-header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  padding: 15px;
  border-bottom: 1px solid #ebeef5;
}

.conversation-body {
  flex: 1;
  overflow-y: auto;
  padding: 15px;
  background-color: #f5f7fa;
}

.empty-history {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #909399;
  font-size: 14px;
}

.message-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.message-item {
  display: flex;
  gap: 10px;
  align-items: flex-start;
}

.self-message {
  flex-direction: row-reverse;
}

.self-message .message-bubble {
  background-color: #ecf5ff;
}

.message-avatar {
  flex-shrink: 0;
}

.message-bubble {
  background-color: #fff;
  padding: 10px;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  max-width: 70%;
}

.message-sender {
  font-size: 12px;
  color: #909399;
  margin-bottom: 5px;
}

.message-time {
  font-size: 12px;
  color: #c0c4cc;
  margin-left: 10px;
}

.message-text {
  word-break: break-word;
}

.conversation-footer {
  padding: 10px;
  display: flex;
  justify-content: center;
  border-top: 1px solid #ebeef5;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 