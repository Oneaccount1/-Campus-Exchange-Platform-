<template>
  <div class="messages-container" v-if="userStore.isLoggedIn">
    <el-row :gutter="20" class="messages-content">
      <el-col :span="6" class="contact-list-container">
        <div class="contact-header">
          <h3>消息列表</h3>
          <el-badge :value="messageStore.totalUnread" :hidden="messageStore.totalUnread === 0" type="danger">
            <el-icon><Bell /></el-icon>
          </el-badge>
        </div>
        <div class="contact-search">
          <el-input
            v-model="searchContact"
            placeholder="搜索联系人"
            clearable
            prefix-icon="Search"
          />
        </div>
        <div class="contact-list" v-loading="messageStore.loadingContacts">
          <div 
            v-for="contact in filteredContacts" 
            :key="contact.id"
            class="contact-item"
            :class="{ 'active': messageStore.currentContactId === contact.id }"
            @click="selectContact(contact.id)"
          >
            <el-avatar :size="40" :src="contact.avatar">
              {{ contact.username?.charAt(0) }}
            </el-avatar>
            <div class="contact-info">
              <div class="contact-name">{{ contact.username }}</div>
              <div class="last-message">{{ contact.lastMessage }}</div>
            </div>
            <div class="message-meta">
              <div class="message-time">{{ formatTime(contact.lastTime) }}</div>
              <el-badge v-if="contact.unread" :value="contact.unread" class="unread-badge" type="danger" />
            </div>
          </div>
          <el-empty v-if="filteredContacts.length === 0 && !messageStore.loadingContacts" description="暂无消息" />
        </div>
      </el-col>
      
      <el-col :span="18" class="chat-container">
        <template v-if="currentContact">
          <div class="chat-header">
            <div class="contact-info">
              <h3>{{ currentContact.username }}</h3>
              <div class="online-status" v-if="currentContact.online">
                <el-tag size="small" type="success">在线</el-tag>
              </div>
            </div>
            <div class="chat-actions">
              <el-dropdown v-if="associatedProducts.length > 0" @command="handleProductSelect">
                <el-button type="primary" size="small">
                  <template v-if="currentProduct">
                    关联商品: {{ currentProduct.title?.substring(0, 10) }}{{ currentProduct.title?.length > 10 ? '...' : '' }}
                  </template>
                  <template v-else>
                    选择关联商品 <el-icon class="el-icon--right"><arrow-down /></el-icon>
                  </template>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item 
                      v-for="product in associatedProducts" 
                      :key="product.id" 
                      :command="product.id"
                      :class="{ 'active-product': currentProduct?.id === product.id }"
                    >
                      {{ product.title }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="currentProduct" command="clear" divided>
                      <span style="color: #f56c6c;">取消关联商品</span>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              
              <el-button type="primary" size="small" plain @click="viewProduct" v-if="currentProduct">
                <el-icon><InfoFilled /></el-icon> 查看商品
              </el-button>
            </div>
          </div>
          
          <div class="chat-messages" ref="messagesContainer" v-loading="messageStore.loadingMessages">
            <div v-if="loadingHistory" class="loading-history">
              <el-button type="primary" plain size="small" @click="loadMoreHistory">加载更多历史消息</el-button>
            </div>
            
            <div 
              v-for="(message, index) in sortedMessages" 
              :key="index"
              class="message-item"
              :class="{ 'self': message.sender_id === userStore.userInfo.id }"
            >
              <el-avatar 
                v-if="message.sender_id !== userStore.userInfo.id" 
                :size="36" 
                :src="currentContact.avatar"
              >
                {{ currentContact.username?.charAt(0) }}
              </el-avatar>
              <div class="message-content">
                <div class="message-text">{{ message.content }}</div>
                <div class="message-time">{{ formatTime(message.created_at) }}</div>
              </div>
              <el-avatar 
                v-if="message.sender_id === userStore.userInfo.id" 
                :size="36" 
                :src="userStore.userInfo.avatar"
              >
                {{ userStore.userInfo?.username?.charAt(0) }}
              </el-avatar>
            </div>
            
            <div class="typing-indicator" v-if="isTyping">
              <div class="typing-bubble"></div>
              <div class="typing-bubble"></div>
              <div class="typing-bubble"></div>
              <span>{{ currentContact.username }} 正在输入...</span>
            </div>
            
            <div v-if="currentMessages.length === 0 && !messageStore.loadingMessages" class="no-messages">
              <el-empty description="暂无消息记录，开始聊天吧" />
            </div>
          </div>
          
          <div class="chat-input">
            <el-input
              v-model="messageText"
              type="textarea"
              :rows="3"
              placeholder="输入消息..."
              @keyup.enter.exact="sendMessage"
              @input="handleTyping"
            />
            <div class="input-actions">
              <el-button type="primary" @click="sendMessage" :disabled="!messageText.trim() || sending">
                <el-icon><Position /></el-icon> {{ sending ? '发送中...' : '发送' }}
              </el-button>
            </div>
          </div>
        </template>
        
        <div v-else class="no-contact-selected">
          <el-empty description="选择一个联系人开始聊天" />
        </div>
      </el-col>
    </el-row>
  </div>
  
  <div v-else class="not-logged-in">
    <el-result
      icon="warning"
      title="请先登录"
      sub-title="登录后才能使用消息功能"
    >
      <template #extra>
        <el-button type="primary" @click="$router.push('/login')">去登录</el-button>
      </template>
    </el-result>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore, useProductStore, useMessageStore } from '../stores'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Bell, Search, Position, InfoFilled, ArrowDown, Check } from '@element-plus/icons-vue'
import * as messageApi from '../api/message'
import * as productApi from '../api/product'
import * as userApi from '../api/user'
import webSocketService from '../utils/websocket'
import notificationService from '../utils/notification'
import { h } from 'vue'
import { ElDropdown, ElDropdownMenu, ElDropdownItem } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const productStore = useProductStore()
const messageStore = useMessageStore()

// 默认头像URL
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'

const searchContact = ref('')
const messageText = ref('')
const messagesContainer = ref(null)
const currentProduct = ref(null)
const loadingHistory = ref(false)
const isTyping = ref(false)
const sending = ref(false)
const typingTimer = ref(null)
const currentPage = ref(1)
const pageSize = ref(20)
const hasMoreMessages = ref(true)
const associatedProducts = ref([])
const hasShownProductPrompt = ref(false)
const loading = ref(false)
const isConnected = ref(false)

// 当前联系人和消息，使用messageStore中的数据
const currentContact = computed(() => messageStore.currentContact)

// 按时间排序后的消息列表
const sortedMessages = computed(() => {
  // 后端返回的是最新消息在前，我们需要按时间递增排序
  return [...messageStore.currentMessages].sort((a, b) => {
    // 将时间字符串转为时间戳进行比较
    const timeA = new Date(a.created_at).getTime();
    const timeB = new Date(b.created_at).getTime();
    return timeA - timeB; // 升序排列，最早的消息在前
  });
});

const currentMessages = computed(() => messageStore.currentMessages)

// 过滤联系人
const filteredContacts = computed(() => {
  return messageStore.filterContacts(searchContact.value);
})

// 在script setup顶部添加userId状态变量
const userId = computed(() => userStore.userInfo?.id)

// 获取联系人列表
const fetchContacts = async () => {
  if (loading.value) return
  
  loading.value = true
  messageStore.setLoading('contacts', true)
  
  try {
    console.log('开始获取联系人列表')
    
    // 检查API函数是否存在
    if (typeof messageApi.getContacts !== 'function') {
      console.error('getContacts函数不存在，尝试使用getContactList')
      throw new Error('联系人API函数未定义')
    }
    
    const response = await messageApi.getContactList()
    console.log('获取联系人原始API响应:', response)
    
    // 检查响应格式
    if (!response) {
      throw new Error('获取联系人列表返回空响应')
    }
    
    // 适配不同的API响应格式
    let contactsData = null
    
    if (response.data && Array.isArray(response.data)) {
      // 直接是数组格式
      contactsData = response.data
      console.log('获取联系人使用数组格式')
    } else if (response.data && response.data.contacts && Array.isArray(response.data.contacts)) {
      // 嵌套在contacts字段中
      contactsData = response.data.contacts
      console.log('获取联系人使用嵌套格式contacts')
    } else if (response.data && response.data.list && Array.isArray(response.data.list)) {
      // 嵌套在list字段中
      contactsData = response.data.list
      console.log('获取联系人使用嵌套格式list')
    } else {
      console.error('无法识别的联系人列表格式:', response)
      throw new Error('无法解析联系人列表数据')
    }
    
    console.log('解析后的联系人数据:', contactsData)
    
    // 进一步处理联系人格式
    const formattedContacts = contactsData.map(contact => {
      // 打印原始联系人数据，帮助调试
      console.log('处理联系人原始数据:', contact)
      
      // 尝试从不同可能的字段名获取数据
      const contactId = contact.user_id || contact.id || contact.contact_id
      const username = contact.username || contact.name || contact.nickname || '未知用户'
      
      return {
        id: parseInt(contactId),
        username: username,
        avatar: contact.avatar || defaultAvatar,
        lastMessage: contact.last_message || contact.lastMessage || contact.content || '',
        lastTime: contact.last_time || contact.lastTime || contact.created_at || contact.time || new Date().toISOString(),
        unread: contact.unread_count || contact.unread || 0
      }
    })
    
    console.log('格式化后的联系人列表:', formattedContacts)
    
    // 设置联系人列表
    messageStore.setContacts(formattedContacts)
    
    // 如果有选中的联系人，获取消息历史
    if (messageStore.currentContactId) {
      await fetchMessages(messageStore.currentContactId)
    }
  } catch (error) {
    console.error('获取联系人列表失败:', error)
    ElMessage.error(`获取联系人列表失败: ${error.message}`)
  } finally {
    loading.value = false
    messageStore.setLoading('contacts', false)
  }
}

// 获取与联系人的消息历史
const fetchMessages = async (contactId) => {
  if (!contactId || messageStore.loadingMessages) return
  
  messageStore.setLoading('messages', true)
  
  try {
    console.log('开始获取联系人消息历史，联系人ID:', contactId)
    const response = await messageApi.getMessageHistory(contactId, currentPage.value, pageSize.value)
    console.log('获取消息历史API响应:', response)
    
    if (!response || !response.data) {
      throw new Error('获取消息历史失败: 无效响应')
    }
    
    // 解析消息数据，支持不同格式
    let messagesData = []
    if (Array.isArray(response.data)) {
      // 直接是消息数组
      messagesData = response.data
      console.log('消息历史使用数组格式')
    } else if (response.data.messages && Array.isArray(response.data.messages)) {
      // 嵌套在messages字段
      messagesData = response.data.messages
      console.log('消息历史使用嵌套格式messages')
    } else if (response.data.list && Array.isArray(response.data.list)) {
      // 嵌套在list字段
      messagesData = response.data.list
      console.log('消息历史使用嵌套格式list')
    } else {
      console.error('无法识别的消息历史格式:', response)
    }
    
    console.log(`获取到${messagesData.length}条消息历史`)
    
    // 收集消息中关联的所有商品ID
    const productIds = new Set()
    
    // 格式化消息数据
    const formattedMessages = messagesData.map(msg => {
      // 打印原始消息数据，帮助调试
      console.log('处理消息原始数据:', msg)
      
      // 如果消息有商品ID，添加到集合
      if (msg.product_id && msg.product_id > 0) {
        productIds.add(msg.product_id)
      }
      
      return {
        id: msg.id,
        sender_id: msg.sender_id,
        receiver_id: msg.receiver_id,
        content: msg.content,
        created_at: msg.created_at || msg.time || new Date().toISOString(),
        is_read: msg.is_read || false,
        product_id: msg.product_id || null
      }
    })
    
    // 添加到store
    if (formattedMessages.length > 0) {
      // 如果是第一页，直接设置消息列表
      if (currentPage.value === 1) {
        messageStore.setMessages(contactId, formattedMessages)
      } else {
        // 如果是加载更多，合并消息列表
        const currentMessages = messageStore.messages[contactId] || []
        messageStore.setMessages(contactId, [...formattedMessages, ...currentMessages])
      }
    } else {
      // 没有消息，设置空数组
      if (currentPage.value === 1) {
        messageStore.setMessages(contactId, [])
      }
    }
    
    // 检查是否还有更多消息
    hasMoreMessages.value = formattedMessages.length >= pageSize.value
    
    // 标记消息为已读
    markMessagesAsRead(contactId)
    
    // 加载关联的商品信息
    if (productIds.size > 0) {
      loadProductsFromIds(Array.from(productIds))
    }
    
    // 如果是第一页，滚动到底部
    if (currentPage.value === 1) {
      nextTick(() => {
        scrollToBottom()
      })
    }
  } catch (error) {
    console.error('获取消息历史失败:', error)
    ElMessage.error(`获取消息历史失败: ${error.message || '未知错误'}`)
  } finally {
    messageStore.setLoading('messages', false)
  }
}

// 从ID列表加载商品信息
const loadProductsFromIds = async (productIds) => {
  if (!productIds || productIds.length === 0) return
  
  try {
    console.log('加载商品信息，商品ID列表:', productIds)
    
    // 加载所有商品信息
    const productPromises = productIds.map(id => productApi.getProductById(id))
    const responses = await Promise.all(productPromises)
    
    // 过滤有效响应
    const validProducts = responses
      .filter(res => res && res.data)
      .map(res => res.data)
    
    console.log('加载到的商品信息:', validProducts)
    
    // 更新associatedProducts
    associatedProducts.value = validProducts
    
    // 如果当前没有选择商品，自动选择第一个
    if (!currentProduct.value && validProducts.length > 0) {
      // 查找最新消息中关联的商品
      const messages = messageStore.messages[messageStore.currentContactId] || []
      if (messages.length > 0) {
        const latestMsg = [...messages].reverse().find(m => m.product_id > 0)
        if (latestMsg && latestMsg.product_id) {
          // 选择与最新消息相同的商品
          const matchingProduct = validProducts.find(p => p.id === latestMsg.product_id)
          if (matchingProduct) {
            currentProduct.value = matchingProduct
            console.log('自动选择与最新消息关联的商品:', matchingProduct.title)
          }
        }
      }
      
      // 如果仍未选择商品，选择第一个
      if (!currentProduct.value) {
        currentProduct.value = validProducts[0]
        console.log('自动选择第一个商品:', validProducts[0].title)
      }
    }
    
    // 保存商品关联到store
    validProducts.forEach(product => {
      messageStore.addContactProduct(messageStore.currentContactId, product.id)
    })
  } catch (error) {
    console.error('加载商品信息失败:', error)
  }
}

// 加载更多历史消息
const loadMoreHistory = async () => {
  if (!messageStore.currentContactId || !hasMoreMessages.value) return
  
  loadingHistory.value = true
  try {
    // 记录当前滚动位置和第一条消息的引用
    const messagesEl = messagesContainer.value
    const scrollPosition = messagesEl.scrollTop
    const firstMessageHeight = messagesEl.firstElementChild?.nextElementSibling?.offsetHeight || 0
    
    currentPage.value++
    await fetchMessages(messageStore.currentContactId)
    
    // 让滚动条保持在原来查看的消息位置
    nextTick(() => {
      const newHeight = messagesEl.firstElementChild?.nextElementSibling?.offsetHeight || 0
      messagesEl.scrollTop = scrollPosition + (newHeight - firstMessageHeight)
    })
  } catch (error) {
    ElMessage.error('加载更多消息失败')
    console.error('加载更多消息失败:', error)
  } finally {
    loadingHistory.value = false
  }
}

// 滚动到底部
const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

// 发送消息后滚动到底部
watch(sortedMessages, () => {
  nextTick(() => {
    scrollToBottom()
  })
})

// 标记消息为已读
const markMessagesAsRead = async (contactId) => {
  try {
    await messageApi.markAsRead(contactId)
    messageStore.clearUnread(contactId)
  } catch (error) {
    console.error('标记消息为已读失败:', error)
  }
}

// 选择联系人
const selectContact = async (contactId) => {
  // 如果已经选中该联系人，无需重复操作
  if (messageStore.currentContactId === contactId) return
  
  // 设置当前联系人
  messageStore.setCurrentContact(contactId)
  
  // 重置分页
  currentPage.value = 1
  hasMoreMessages.value = true
  
  // 获取与该联系人的消息历史
  await fetchMessages(contactId)
  
  // 查找联系人关联的商品
  const contact = messageStore.getContactById(contactId)
  if (contact && contact.product_id) {
    try {
      const response = await productApi.getProductById(contact.product_id)
      currentProduct.value = response.data
    } catch (error) {
      console.error('获取商品信息失败:', error)
      currentProduct.value = null
    }
  } else {
    currentProduct.value = null
  }
}

// 修改发送消息函数
const sendMessage = async () => {
  if (!messageText.value.trim() || !messageStore.currentContactId || sending.value) return
  
  const content = messageText.value.trim()
  sending.value = true
  
  try {
    // 创建基本消息数据 - 使用后端期望的receiver_id参数名
    const messageData = {
      receiver_id: messageStore.currentContactId,
      content: content,
      type: 'text' // 明确指定消息类型
    }
    
    // 确定要关联的商品ID
    // 优先级: 1.当前选择的商品 > 2.对话历史中的商品 > 3.关联商品列表中的第一个
    let productId = null
    
    // 1. 检查当前选择的商品
    if (currentProduct.value && currentProduct.value.id > 0) {
      productId = currentProduct.value.id
      console.log('使用当前选择的商品ID:', productId)
    } 
    // 2. 检查历史消息中是否有商品ID
    else {
      const messages = messageStore.messages[messageStore.currentContactId] || []
      
      // 查找最近的带有product_id的消息，按时间逆序
      const recentMsgWithProduct = [...messages]
        .reverse()
        .find(msg => msg.product_id && msg.product_id > 0)
      
      if (recentMsgWithProduct) {
        productId = recentMsgWithProduct.product_id
        console.log('使用历史消息中的商品ID:', productId)
      }
      // 3. 如果没有历史商品ID，尝试使用关联商品列表中的第一个
      else if (associatedProducts.value.length > 0) {
        productId = associatedProducts.value[0].id
        console.log('使用关联商品列表中的第一个商品ID:', productId)
        
        // 如果用户尚未选择过商品，显示提示
        if (!hasShownProductPrompt.value) {
          ElMessageBox.confirm(
            `系统将自动关联商品"${associatedProducts.value[0].title}"，是否要选择其他商品？`,
            '商品关联提示',
            {
              confirmButtonText: '选择其他商品',
              cancelButtonText: '使用此商品',
              type: 'info'
            }
          ).then(() => {
            // 用户选择其他商品，显示商品选择对话框
            showProductSelector()
            sending.value = false
          }).catch(() => {
            // 用户确认使用自动选择的商品
            if (productId) {
              messageData.product_id = productId
              // 更新当前商品
              selectProduct(productId).then(() => {
                sendMessageWithData(messageData)
              })
            } else {
              sendMessageWithData(messageData)
            }
          })
          
          // 标记已显示提示，避免重复提示
          hasShownProductPrompt.value = true
          return
        }
      }
    }
    
    // 如果找到了有效的商品ID，添加到消息数据
    if (productId && productId > 0) {
      messageData.product_id = productId
    }
    
    console.log('发送消息数据:', messageData)
    
    // 发送消息
    sendMessageWithData(messageData)
  } catch (error) {
    ElMessage.error('发送消息失败')
    console.error('发送消息失败:', error)
    sending.value = false
  }
}

// 实际发送消息的函数
const sendMessageWithData = async (messageData) => {
  try {
    console.log('开始发送消息:', messageData)
    
    const response = await messageApi.sendMessage(messageData)
    console.log('消息发送API响应:', response)
    
    if (response.code !== 200 || !response.data) {
      throw new Error(`消息发送失败: ${response.message || '未知错误'}`)
    }
    
    // 添加消息到本地
    const newMessage = {
      id: response.data.id || Date.now(),
      sender_id: userStore.userInfo.id,
      receiver_id: parseInt(messageStore.currentContactId),
      content: messageData.content,
      created_at: response.data.created_at || new Date().toISOString(),
      is_read: false
    }
    
    // 只有当产品ID存在且有效时才添加product_id字段
    if (messageData.product_id && messageData.product_id > 0) {
      newMessage.product_id = messageData.product_id
    }
    
    console.log('添加消息到store:', newMessage)
    
    // 添加到store
    try {
      messageStore.addMessage(messageStore.currentContactId, newMessage)
    } catch (storeError) {
      console.error('添加消息到store失败:', storeError)
      // 继续执行，不中断流程
    }
    
    // 也通过WebSocket发送，实现双重保障
    if (webSocketService.isConnected) {
      webSocketService.sendMessage({
        type: 'message',
        data: newMessage
      })
    }
    
    // 清空输入框
    messageText.value = ''
    
    // 确保滚动到底部
    nextTick(() => {
      scrollToBottom()
    })
    
    return true
  } catch (error) {
    console.error('发送消息过程中出错:', error)
    ElMessage.error(`发送消息失败: ${error.message || '未知错误'}`)
    return false
  } finally {
    sending.value = false
  }
}

// 显示商品选择器
const showProductSelector = () => {
  if (associatedProducts.value.length === 0) {
    ElMessage.info('没有可关联的商品')
    return
  }
  
  ElMessageBox.confirm(
    h('div', { class: 'product-selector' }, [
      h('h3', '选择要关联的商品'),
      h('ul', { class: 'product-list' }, 
        associatedProducts.value.map(product => 
          h('li', {
            class: ['product-item', { active: currentProduct.value?.id === product.id }],
            onClick: () => selectProduct(product.id)
          }, [
            h('span', product.title),
            currentProduct.value?.id === product.id ? h('el-icon', { class: 'check-icon' }, [h(Check)]) : null
          ])
        )
      )
    ]),
    {
      title: '选择商品',
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      customClass: 'product-selector-dialog'
    }
  ).catch(() => {
    // 用户取消，不做任何操作
  })
}

// 处理商品选择
const handleProductSelect = async (command) => {
  if (command === 'clear') {
    // 清除当前商品关联
    currentProduct.value = null
    // 更新store
    messageStore.setContactProduct(messageStore.currentContactId, null)
    return
  }
  
  // 根据ID选择商品
  const productId = parseInt(command)
  await selectProduct(productId)
}

// 选择商品
const selectProduct = async (productId) => {
  try {
    const response = await productApi.getProductById(productId)
    currentProduct.value = response.data
    
    // 更新store
    if (currentProduct.value && currentProduct.value.id) {
      messageStore.setContactProduct(messageStore.currentContactId, currentProduct.value.id)
    }
    
    ElMessage.success(`已关联商品: ${currentProduct.value.title}`)
  } catch (error) {
    console.error('获取商品信息失败:', error)
    ElMessage.error('获取商品信息失败')
  }
}

// 处理用户输入，发送正在输入状态
const handleTyping = () => {
  if (!messageStore.currentContactId) return
  
  // 发送正在输入状态 - 使用标准格式
  if (!typingTimer.value) {
    // 发送输入状态
    webSocketService.sendMessage({
      type: 'typing',
      data: {
        recipientId: messageStore.currentContactId,
        status: true
      }
    })
  }
  
  // 清除之前的定时器
  if (typingTimer.value) {
    clearTimeout(typingTimer.value)
  }
  
  // 设置新的定时器，2秒后发送停止输入状态
  typingTimer.value = setTimeout(() => {
    webSocketService.sendMessage({
      type: 'typing',
      data: {
        recipientId: messageStore.currentContactId,
        status: false
      }
    })
    typingTimer.value = null
  }, 2000)
}

// 查看商品详情
const viewProduct = () => {
  if (currentProduct.value) {
    router.push(`/product/${currentProduct.value.id}`)
  }
}

// 初始化WebSocket连接
const initWebSocket = () => {
  if (!userStore.isLoggedIn) return
  
  // 检查WebSocket连接状态
  if (!webSocketService.isConnected) {
    console.log('Messages页面: WebSocket未连接，尝试重新连接')
    webSocketService.connect()
  } else {
    console.log('Messages页面: WebSocket已连接')
  }
  
  // 删除之前可能存在的监听器，避免重复
  webSocketService.offMessage(handleWebSocketMessage)
  
  // 添加消息监听器
  console.log('Messages页面: 添加WebSocket消息监听器')
  webSocketService.onMessage(handleWebSocketMessage)
  
  // 发送一个请求在线状态更新
  setTimeout(() => {
    if (webSocketService.isConnected) {
      // 请求联系人在线状态更新
      if (messageStore.currentContactId) {
        console.log('Messages页面: 请求联系人在线状态')
        webSocketService.sendMessage({
          type: 'online_status_request',
          data: {
            userIds: [messageStore.currentContactId]
          }
        })
      }
    } else {
      console.log('Messages页面: WebSocket未连接，无法请求状态')
      // 尝试重新连接
      webSocketService.connect()
    }
  }, 500)
  
  // 添加连接状态监听器
  webSocketService.onConnection((connected) => {
    console.log('Messages页面: WebSocket连接状态变更:', connected ? '已连接' : '已断开')
    isConnected.value = connected
    
    // 连接成功后重新订阅
    if (connected) {
      // 发送订阅请求，使用标准格式
      setTimeout(() => {
        console.log('Messages页面: 重新订阅消息')
        if (messageStore.currentContactId) {
          webSocketService.sendMessage({
            type: 'subscribe',
            data: {
              recipientId: messageStore.currentContactId
            }
          })
        }
      }, 1000)
    }
  })
}

// 处理收到的WebSocket消息
const handleWebSocketMessage = (message) => {
  try {
    console.log('Messages页面接收到WebSocket消息:', message)
    if (!message) return
    
    // 处理标准格式的消息
    const type = message.type || 'unknown'
    const data = message.data || {}
    
    switch (type) {
      case 'message':
        // 收到新消息
        handleNewMessage(data)
        break
      case 'typing':
        // 收到正在输入状态
        handleTypingStatus(data)
        break
      case 'online':
      case 'online_status':
        // 收到在线状态变更
        handleOnlineStatus(data)
        break
      case 'subscribe_success':
        // 订阅成功
        console.log('消息订阅成功:', data)
        break
      case 'raw':
        // 未知格式的原始消息
        console.log('收到原始消息:', data)
        break
      default:
        console.log('收到未知类型的WebSocket消息:', message)
        // 尝试兼容旧消息格式
        if (message.sender_id && message.content) {
          handleNewMessage(message)
        }
    }
  } catch (error) {
    console.error('处理WebSocket消息时出错:', error)
  }
}

// 处理新消息
const handleNewMessage = (messageData) => {
  if (!messageData || !messageData.sender_id) return
  
  // 如果是自己发送的消息，已在sendMessage处理，这里不再处理
  if (messageData.sender_id === userStore.userInfo.id) return
  
  const senderId = messageData.sender_id
  console.log('收到新消息，发送者ID:', senderId)
  
  // 添加消息到对应联系人
  messageStore.addMessage(senderId, messageData)
  
  // 查找联系人
  const contact = messageStore.getContactById(senderId)
  
  // 如果消息不是当前选中的联系人发送的，显示通知
  if (messageStore.currentContactId !== senderId && contact?.username) {
    showMessageNotification(contact.username, messageData.content)
  } else if (messageStore.currentContactId !== senderId) {
    // 如果没有联系人名称，使用ID显示通知
    showMessageNotification(`联系人 ${senderId}`, messageData.content)
    
    // 尝试获取联系人信息
    fetchContactInfo(senderId)
  }
}

// 获取联系人信息
const fetchContactInfo = async (userId) => {
  try {
    const response = await userApi.getUserInfo(userId)
    if (response.code === 200 && response.data) {
      // 创建或更新联系人
      const contact = {
        id: userId,
        username: response.data.username,
        avatar: response.data.avatar || defaultAvatar
      }
      
      // 添加到联系人列表
      messageStore.addOrUpdateContact(contact)
    }
  } catch (error) {
    console.error('获取联系人信息失败:', error)
  }
}

// 显示消息通知
const showMessageNotification = (sender, content) => {
  notificationService.showDesktopNotification({
    title: `来自 ${sender} 的新消息`,
    message: content,
    onClick: () => {
      // 点击通知时跳转到消息页面
      router.push('/messages')
    }
  })
}

// 处理输入状态
const handleTypingStatus = (data) => {
  if (!data || data.receiver_id !== messageStore.currentContactId) return
  
  isTyping.value = data.status
}

// 处理在线状态
const handleOnlineStatus = (data) => {
  if (!data || !data.user_id) return
  
  const contact = messageStore.getContactById(data.user_id)
  if (contact) {
    contact.online = data.status
  }
}

// 从商品详情页创建新联系人
const createNewContact = async () => {
  console.log('Messages页面：准备创建新联系人，路由参数:', route.query)
  
  // 检查路由参数
  if (!route.query.userId) {
    console.log('Messages页面：缺少必要的用户ID参数，无法创建新联系人')
    return
  }
  
  const userId = parseInt(route.query.userId)
  if (!userId || isNaN(userId)) {
    console.log('Messages页面：用户ID无效', route.query.userId)
    return
  }
  
  // 尝试解析商品ID
  let productId = null
  if (route.query.productId) {
    const parsedId = parseInt(route.query.productId)
    if (parsedId && !isNaN(parsedId) && parsedId > 0) {
      productId = parsedId
    }
  }
  
  console.log('Messages页面：解析后的参数:', {userId, productId})
  
  // 检查联系人是否已存在
  const existingContact = messageStore.getContactById(userId)
  if (existingContact) {
    console.log('Messages页面：联系人已存在，直接选择:', existingContact)
    // 如果已存在，直接选择
    selectContact(userId)
    
    // 如果有新的商品ID且有效，更新商品关联
    if (productId) {
      try {
        const productResponse = await productApi.getProductById(productId)
        currentProduct.value = productResponse.data
        console.log('Messages页面：更新现有联系人关联的商品:', currentProduct.value)
        
        // 同时更新store中的商品关联
        messageStore.setContactProduct(userId, productId)
        messageStore.addContactProduct(userId, productId)
        
        // 重新加载商品列表
        loadAssociatedProducts()
      } catch (error) {
        console.error('获取商品信息失败:', error)
      }
    }
    
    return
  }
  
  try {
    // 创建新会话请求数据
    const conversationData = {
      user_id: userId
    }
    
    // 如果有有效的商品ID，则添加
    if (productId) {
      // 获取商品信息
      console.log('Messages页面：开始获取商品信息，productId:', productId)
      const productResponse = await productApi.getProductById(productId)
      currentProduct.value = productResponse.data
      console.log('Messages页面：获取到商品信息:', currentProduct.value)
      
      // 添加商品ID到会话请求
      conversationData.product_id = productId
    }
    
    console.log('Messages页面：创建新会话请求数据:', conversationData)
    
    // 使用后端API创建会话
    const response = await messageApi.createConversation(conversationData)
    console.log('Messages页面：创建会话响应:', response)
    
    // 添加新联系人
    const newContact = {
      id: userId,
      username: response.data.username || '卖家',
      avatar: response.data.avatar || '',
      lastMessage: '',
      lastTime: new Date().toISOString(),
      unread: 0
    }
    
    // 如果有有效的商品ID，则添加到联系人
    if (productId) {
      newContact.product_id = productId
    }
    
    console.log('Messages页面：新建联系人信息:', newContact)
    
    // 添加到联系人列表
    messageStore.addOrUpdateContact(newContact)
    
    // 如果有商品ID，更新商品关联 
    if (productId) {
      messageStore.setContactProduct(userId, productId)
      messageStore.addContactProduct(userId, productId)
    }
    
    // 选择新联系人
    selectContact(userId)
    
    // 提示用户
    ElMessage.success('已为您创建新会话，可以开始聊天了')
  } catch (error) {
    console.error('Messages页面：创建会话失败:', error)
    ElMessage.error('创建会话失败')
  }
}

// 初始化
onMounted(async () => {
  // 检查用户是否登录
  if (!userStore.isLoggedIn) return
  
  // 获取联系人列表
  await fetchContacts()
  
  // 初始化WebSocket
  initWebSocket()
  
  // 处理从其他页面跳转过来的情况
  if (route.query.userId) {
    // 如果有联系人ID，创建或选择联系人
    createNewContact()
  }
  
  // 恢复商品映射
  messageStore.restoreProductMappings()
  
  // 如果已选择联系人，加载关联商品
  if (messageStore.currentContactId) {
    loadAssociatedProducts()
  }
})

// 监听当前联系人变化
watch(() => messageStore.currentContactId, (newVal) => {
  if (newVal) {
    // 重置分页
    currentPage.value = 1
    hasMoreMessages.value = true
    loadAssociatedProducts()
    // 重置标记
    hasShownProductPrompt.value = false
  } else {
    currentProduct.value = null
    associatedProducts.value = []
  }
})

// 清理
onBeforeUnmount(() => {
  // 清除定时器
  if (typingTimer.value) {
    clearTimeout(typingTimer.value)
  }
  
  // 移除WebSocket监听器
  webSocketService.offMessage(handleWebSocketMessage)
})

// 格式化时间
const formatTime = (timeString) => {
  if (!timeString) return ''
  
  const date = new Date(timeString)
  const now = new Date()
  
  // 如果是今天的消息，只显示时间
  if (date.toDateString() === now.toDateString()) {
    return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
  }
  
  // 如果是昨天的消息
  const yesterday = new Date(now)
  yesterday.setDate(now.getDate() - 1)
  if (date.toDateString() === yesterday.toDateString()) {
    return `昨天 ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
  }
  
  // 如果是今年的消息
  if (date.getFullYear() === now.getFullYear()) {
    return `${date.getMonth() + 1}月${date.getDate()}日 ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
  }
  
  // 其他情况显示完整日期
  return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日 ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

// 在初始化时加载关联商品
const loadAssociatedProducts = async () => {
  if (!messageStore.currentContactId) return
  
  // 获取当前联系人关联的商品列表
  const productIds = messageStore.getContactProducts(messageStore.currentContactId)
  
  // 如果只有一个商品ID，直接设置为当前商品
  if (productIds.length === 1) {
    try {
      const response = await productApi.getProductById(productIds[0])
      currentProduct.value = response.data
    } catch (error) {
      console.error('获取商品信息失败:', error)
    }
  }
  
  // 加载所有关联过的商品信息
  const productPromises = productIds.map(id => productApi.getProductById(id))
  try {
    const responses = await Promise.all(productPromises.map(p => p.catch(err => null)))
    associatedProducts.value = responses
      .filter(res => res && res.data)
      .map(res => res.data)
  } catch (error) {
    console.error('获取关联商品列表失败:', error)
  }
}
</script>

<style scoped>
.messages-container {
  padding: 20px;
  height: calc(100vh - 180px);
}

.messages-content {
  height: 100%;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.contact-list-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #e4e7ed;
}

.contact-header {
  padding: 15px;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.contact-header h3 {
  margin: 0;
  color: #303133;
}

.contact-search {
  padding: 10px;
  border-bottom: 1px solid #e4e7ed;
}

.contact-list {
  flex: 1;
  overflow-y: auto;
}

.contact-item {
  display: flex;
  align-items: center;
  padding: 10px 15px;
  cursor: pointer;
  transition: background-color 0.3s;
  border-bottom: 1px solid #f0f2f5;
}

.contact-item:hover {
  background-color: #f5f7fa;
}

.contact-item.active {
  background-color: #ecf5ff;
}

.contact-info {
  flex: 1;
  margin: 0 10px;
  overflow: hidden;
}

.contact-name {
  font-size: 14px;
  color: #303133;
  margin-bottom: 5px;
  font-weight: bold;
}

.last-message {
  font-size: 12px;
  color: #909399;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.message-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  min-width: 60px;
}

.message-time {
  font-size: 12px;
  color: #909399;
  margin-bottom: 5px;
}

.unread-badge {
  margin-top: 2px;
}

.chat-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.chat-header {
  padding: 15px;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chat-header .contact-info {
  display: flex;
  align-items: center;
}

.chat-header h3 {
  margin: 0;
  color: #303133;
  margin-right: 10px;
}

.chat-header .online-status {
  margin-left: 10px;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 15px;
  background-color: #f5f7fa;
}

.message-item {
  display: flex;
  margin-bottom: 15px;
  align-items: flex-start;
}

.message-item.self {
  flex-direction: row-reverse;
}

.message-content {
  margin: 0 10px;
  max-width: 70%;
}

.message-text {
  padding: 10px 15px;
  border-radius: 4px;
  background-color: white;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  word-break: break-word;
}

.message-item.self .message-text {
  background-color: #ecf5ff;
  color: #409EFF;
}

.message-time {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
  text-align: right;
}

.message-item.self .message-time {
  text-align: left;
}

.chat-input {
  padding: 15px;
  border-top: 1px solid #e4e7ed;
}

.input-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 10px;
}

.no-contact-selected,
.not-logged-in {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.loading-history {
  text-align: center;
  margin-bottom: 15px;
}

.typing-indicator {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
  color: #909399;
  font-size: 12px;
}

.typing-bubble {
  width: 8px;
  height: 8px;
  margin: 0 2px;
  background-color: #909399;
  border-radius: 50%;
  animation: typing-animation 1s infinite;
  display: inline-block;
}

.typing-bubble:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-bubble:nth-child(3) {
  animation-delay: 0.4s;
}

.typing-indicator span {
  margin-left: 8px;
}

@keyframes typing-animation {
  0% {
    opacity: 0.3;
    transform: translateY(0);
  }
  50% {
    opacity: 1;
    transform: translateY(-5px);
  }
  100% {
    opacity: 0.3;
    transform: translateY(0);
  }
}

/* 添加商品选择器样式 */
.product-selector {
  max-height: 300px;
  overflow-y: auto;
}

.product-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.product-item {
  padding: 10px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.product-item:hover {
  background-color: #f5f7fa;
}

.product-item.active {
  background-color: #ecf5ff;
  color: #409EFF;
}

.check-icon {
  color: #67c23a;
}

.active-product {
  color: #409EFF;
  font-weight: bold;
}
</style>