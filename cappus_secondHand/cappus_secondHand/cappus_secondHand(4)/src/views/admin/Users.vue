<template>
  <div class="admin-users-container">
    <div class="page-header">
      <h2 class="page-title">用户管理</h2>
      <el-button type="primary" @click="exportUsers">
        <el-icon><Download /></el-icon>导出数据
      </el-button>
    </div>

    <el-card shadow="hover" class="filter-card">
      <div class="filter-container">
        <div class="filter-item">
          <el-input
            v-model="searchQuery"
            placeholder="搜索用户名"
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
          <el-select v-model="statusFilter" placeholder="状态筛选" clearable @change="handleSearch">
            <el-option label="正常" value="正常" />
            <el-option label="禁用" value="禁用" />
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
          <el-tag type="success">共 {{ totalUsers }} 条记录</el-tag>
          <el-button-group class="action-group">
            <el-button type="primary" plain size="small" @click="batchEnable" :disabled="selectedUsers.length === 0">
              <el-icon><Check /></el-icon>批量启用
            </el-button>
            <el-button type="danger" plain size="small" @click="batchDisable" :disabled="selectedUsers.length === 0">
              <el-icon><Close /></el-icon>批量禁用
            </el-button>
          </el-button-group>
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
        :data="filteredUsers" 
        style="width: 100%" 
        v-loading="loading"
        @selection-change="handleSelectionChange"
        :header-cell-style="{ background: '#f5f7fa' }"
        border
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" sortable />
        <el-table-column label="用户信息">
          <template #default="scope">
            <div class="user-info-cell">
              <el-avatar :size="40" :src="scope.row.avatar" class="user-avatar">
                {{ scope.row.username.charAt(0) }}
              </el-avatar>
              <div class="user-info-text">
                <div class="user-name">{{ scope.row.username }}</div>
                <div class="user-role">普通用户</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="注册时间" width="180" sortable>
          <template #default="scope">
            {{ formatDateTime(scope.row.created_at || scope.row.registerTime) }}
          </template>
        </el-table-column>
        <el-table-column label="发布商品数" width="120" sortable>
          <template #default="scope">
            <el-badge :value="scope.row.product_count || scope.row.productCount" :max="99" type="primary" />
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === '正常' ? 'success' : 'danger'">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="viewUserDetail(scope.row.id)">
              <el-icon><View /></el-icon>查看
            </el-button>
            <el-button 
              :type="scope.row.status === '正常' ? 'danger' : 'success'" 
              size="small" 
              @click="toggleUserStatus(scope.row)"
            >
              <el-icon v-if="scope.row.status === '正常'"><Lock /></el-icon>
              <el-icon v-else><Unlock /></el-icon>
              {{ scope.row.status === '正常' ? '禁用' : '启用' }}
            </el-button>
            <el-dropdown trigger="click" @command="handleCommand($event, scope.row)">
              <el-button size="small" plain>
                <el-icon><MoreFilled /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="resetPassword">
                    <el-icon><Key /></el-icon>重置密码
                  </el-dropdown-item>
                  <el-dropdown-item command="message">
                    <el-icon><Message /></el-icon>发送消息
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
      
      <div class="pagination-container">
        <el-pagination
          :current-page="currentPage"
          :page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="totalUsers"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          background
        />
      </div>
    </el-card>
    
    <!-- 用户详情对话框 -->
    <el-dialog
      v-model="userDetailVisible"
      title="用户详情"
      width="800px"
      destroy-on-close
    >
      <div v-if="currentUser" class="user-detail">
        <el-row :gutter="20">
          <el-col :span="8">
            <div class="user-profile-card">
              <div class="user-profile-header">
                <el-avatar :size="80" :src="currentUser.avatar" class="profile-avatar">
                  {{ currentUser.username.charAt(0) }}
                </el-avatar>
                <h2>{{ currentUser.username }}</h2>
                <el-tag :type="currentUser.status === '正常' ? 'success' : 'danger'" class="status-tag">
                  {{ currentUser.status }}
                </el-tag>
              </div>
              
              <el-divider />
              
              <div class="user-profile-stats">
                <div class="stat-item">
                  <div class="stat-value">{{ currentUser.productCount }}</div>
                  <div class="stat-label">发布商品</div>
                </div>
                <div class="stat-item">
                  <div class="stat-value">{{ currentUser.orderCount || 0 }}</div>
                  <div class="stat-label">交易订单</div>
                </div>
                <div class="stat-item">
                  <div class="stat-value">{{ currentUser.favoriteCount || 0 }}</div>
                  <div class="stat-label">收藏商品</div>
                </div>
              </div>
              
              <el-divider />
              
              <div class="user-profile-actions">
                <el-button type="primary" @click="sendMessage(currentUser)">
                  <el-icon><Message /></el-icon>发送消息
                </el-button>
                <el-button 
                  :type="currentUser.status === '正常' ? 'danger' : 'success'"
                  @click="toggleUserStatus(currentUser)"
                >
                  {{ currentUser.status === '正常' ? '禁用账号' : '启用账号' }}
                </el-button>
              </div>
            </div>
          </el-col>
          
          <el-col :span="16">
            <el-tabs>
              <el-tab-pane label="基本信息">
                <el-descriptions :column="2" border>
                  <el-descriptions-item label="用户ID">{{ currentUser.id }}</el-descriptions-item>
                  <el-descriptions-item label="注册时间">{{ formatDateTime(currentUser.created_at || currentUser.registerTime) }}</el-descriptions-item>
                  <el-descriptions-item label="电子邮箱">{{ currentUser.email || '未设置' }}</el-descriptions-item>
                  <el-descriptions-item label="手机号码">{{ currentUser.phone || '未设置' }}</el-descriptions-item>
                  <el-descriptions-item label="上次登录">{{ currentUser.lastLogin || '未记录' }}</el-descriptions-item>
                  <el-descriptions-item label="登录IP">{{ currentUser.lastIp || '未记录' }}</el-descriptions-item>
                </el-descriptions>
              </el-tab-pane>
              
              <el-tab-pane label="发布的商品" v-if="currentUser.products && currentUser.products.length > 0">
                <el-table :data="currentUser.products" style="width: 100%" :header-cell-style="{ background: '#f5f7fa' }" border>
                  <el-table-column prop="id" label="ID" width="60" />
                  <el-table-column prop="title" label="商品名称" show-overflow-tooltip />
                  <el-table-column prop="price" label="价格" width="100">
                    <template #default="scope">
                      <span class="price-tag">¥{{ scope.row.price }}</span>
                    </template>
                  </el-table-column>
                  <el-table-column prop="createTime" label="发布时间" width="180" />
                  <el-table-column label="操作" width="120">
                    <template #default="scope">
                      <el-button type="primary" size="small" plain @click="viewProductDetail(scope.row.id)">
                        查看商品
                      </el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </el-tab-pane>
              
              <el-tab-pane label="操作记录">
                <el-empty description="暂无操作记录" v-if="!currentUser.activities || currentUser.activities.length === 0" />
                <el-timeline v-else>
                  <el-timeline-item
                    v-for="(activity, index) in currentUser.activities"
                    :key="index"
                    :timestamp="activity.time"
                    :type="activity.type"
                  >
                    {{ activity.content }}
                  </el-timeline-item>
                </el-timeline>
              </el-tab-pane>
            </el-tabs>
          </el-col>
        </el-row>
      </div>
    </el-dialog>
    
    <!-- 发送消息对话框 -->
    <el-dialog
      v-model="messageDialogVisible"
      title="发送消息"
      width="500px"
    >
      <el-form :model="messageForm" :rules="messageRules" ref="messageFormRef" label-position="top">
        <el-form-item label="接收用户" prop="receiver">
          <el-input v-model="messageForm.receiver" disabled />
        </el-form-item>
        <el-form-item label="消息标题" prop="title">
          <el-input v-model="messageForm.title" placeholder="请输入消息标题" />
        </el-form-item>
        <el-form-item label="消息内容" prop="content">
          <el-input 
            v-model="messageForm.content" 
            type="textarea" 
            :rows="4" 
            placeholder="请输入消息内容" 
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="messageDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitMessage" :loading="messageSending">发送</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '../../stores'
import { ElMessageBox, ElMessage } from 'element-plus'
import { 
  getUsers, 
  getUserDetail, 
  updateUserStatus, 
  resetUserPassword, 
  exportUsersData as exportUsersAPI,
  sendSystemMessageApi
} from '../../api/admin'

const router = useRouter()
const adminStore = useAdminStore()

// 检查管理员是否登录
onMounted(() => {
  // 加载用户数据
  loadUsers()
})

const loading = ref(false)
const searchQuery = ref('')
const statusFilter = ref('')
const dateRange = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const totalUsers = ref(0)
const selectedUsers = ref([])
const users = ref([]) // 用户数据列表

// 用户详情相关
const userDetailVisible = ref(false)
const currentUser = ref(null)

// 消息对话框相关
const messageDialogVisible = ref(false)
const messageFormRef = ref(null)
const messageSending = ref(false)
const messageForm = reactive({
  receiver: '',
  receiverId: null,
  title: '',
  content: ''
})
const messageRules = {
  title: [
    { required: true, message: '请输入消息标题', trigger: 'blur' },
    { max: 50, message: '标题长度不能超过50个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入消息内容', trigger: 'blur' },
    { max: 500, message: '内容长度不能超过500个字符', trigger: 'blur' }
  ]
}

// 过滤后的用户列表（用于前端过滤）
const filteredUsers = computed(() => {
  // 如果是从服务器获取的已过滤数据，则直接返回
  return users.value
})

// 加载用户数据
const loadUsers = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      size: pageSize.value,
      search: searchQuery.value || undefined,
      status: statusFilter.value || undefined
    }
    
    // 添加日期范围参数
    if (dateRange.value && dateRange.value.length === 2) {
      params.startDate = formatDate(dateRange.value[0])
      params.endDate = formatDate(dateRange.value[1], true) // 设置为当天结束时间
    }
    
    const response = await getUsers(params)
    console.log('加载用户列表响应数据:', response)
    
          if (response.code === 200 && response.data) {
      // 兼容不同的返回数据结构
      let userData = []
      if (Array.isArray(response.data)) {
        // 如果直接返回数组
        userData = response.data
        totalUsers.value = response.data.length
      } else if (response.data.users) {
        // 如果返回的是users字段
        userData = response.data.users
        totalUsers.value = response.data.total || 0
      } else if (response.data.list) {
        // 如果返回带分页的数据结构
        userData = response.data.list
        totalUsers.value = response.data.total || 0
      } else {
        // 尝试从data对象中找到可能的数组
        const possibleList = Object.values(response.data).find(value => Array.isArray(value))
        if (possibleList) {
          userData = possibleList
          totalUsers.value = possibleList.length
        } else {
          users.value = []
          totalUsers.value = 0
          console.warn('无法解析返回的数据结构:', response.data)
        }
      }
      
      // 统一处理字段名不一致的问题
      users.value = userData.map(user => {
        // 统一字段名，优先使用后端返回的字段，兼容原有字段名
        return {
          ...user,
          registerTime: user.created_at || user.registerTime,
          productCount: user.product_count !== undefined ? user.product_count : user.productCount
        }
      })
      
      console.log('解析后的用户列表:', users.value)
      console.log('解析后的用户总数:', totalUsers.value)
    } else {
      ElMessage.error(response.message || '获取用户列表失败')
    }
  } catch (error) {
    console.error('加载用户数据失败:', error)
    ElMessage.error('加载用户数据失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 格式化日期时间函数（ISO格式转换为YYYY-MM-DD HH:MM:SS）
const formatDateTime = (dateTimeStr) => {
  if (!dateTimeStr) return '';
  
  try {
    // 处理ISO格式的时间
    const date = new Date(dateTimeStr);
    
    // 如果日期无效，则直接返回原字符串
    if (isNaN(date.getTime())) {
      return dateTimeStr;
    }
    
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    const seconds = String(date.getSeconds()).padStart(2, '0');
    
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
  } catch (error) {
    console.error('日期格式化错误:', error);
    return dateTimeStr; // 出错时返回原始字符串
  }
};

// 格式化日期函数
const formatDate = (date, isEndOfDay = false) => {
  const d = new Date(date)
  if (isEndOfDay) {
    d.setHours(23, 59, 59, 999)
  }
  return d.toISOString()
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  loadUsers()
}

// 重置过滤条件
const resetFilters = () => {
  searchQuery.value = ''
  statusFilter.value = ''
  dateRange.value = []
  handleSearch()
}

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size
  loadUsers()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  loadUsers()
}

// 查看用户详情
const viewUserDetail = async (userId) => {
  try {
    const response = await getUserDetail(userId)
    if (response.code === 200 && response.data) {
      currentUser.value = response.data
      userDetailVisible.value = true
    } else {
      ElMessage.error(response.message || '获取用户详情失败')
    }
  } catch (error) {
    console.error('获取用户详情失败:', error)
    ElMessage.error('获取用户详情失败，请稍后重试')
  }
}

// 查看商品详情
const viewProductDetail = (productId) => {
  router.push(`/admin/products/${productId}`)
}

// 切换用户状态（启用/禁用）
const toggleUserStatus = (user) => {
  const action = user.status === '正常' ? '禁用' : '启用'
  const newStatus = user.status === '正常' ? '禁用' : '正常'
  
  ElMessageBox.confirm(
    `确定要${action}用户 ${user.username} 吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const response = await updateUserStatus(user.id, newStatus)
      if (response.code === 200) {
        // 更新本地状态
        const targetUser = users.value.find(u => u.id === user.id)
        if (targetUser) {
          targetUser.status = newStatus
        }
        
        // 如果当前正在查看该用户详情，也更新详情中的状态
        if (currentUser.value && currentUser.value.id === user.id) {
          currentUser.value.status = newStatus
        }
        
        ElMessage.success(`已${action}用户 ${user.username}`)
      } else {
        ElMessage.error(response.message || `${action}用户失败`)
      }
    } catch (error) {
      console.error(`${action}用户失败:`, error)
      ElMessage.error(`${action}用户失败，请稍后重试`)
    }
  }).catch(() => {
    // 取消操作
  })
}

// 处理表格选择
const handleSelectionChange = (selection) => {
  selectedUsers.value = selection
}

// 批量启用
const batchEnable = async () => {
  if (selectedUsers.value.length === 0) return
  
  ElMessageBox.confirm(
    `确定要批量启用选中的 ${selectedUsers.value.length} 个用户吗？`,
    '批量操作',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      // 获取所有禁用状态的用户ID
      const userIds = selectedUsers.value
        .filter(user => user.status !== '正常')
        .map(user => user.id)
      
      if (userIds.length === 0) {
        ElMessage.info('所选用户均已启用')
        return
      }
      
      // 批量请求API
      let successCount = 0
      for (const userId of userIds) {
        try {
          const response = await updateUserStatus(userId, '正常')
          if (response.code === 200) {
            successCount++
            // 更新本地状态
            const targetUser = users.value.find(u => u.id === userId)
            if (targetUser) {
              targetUser.status = '正常'
            }
          }
        } catch (error) {
          console.error(`启用用户 ${userId} 失败:`, error)
        }
      }
      
      if (successCount > 0) {
        ElMessage.success(`已成功启用 ${successCount} 个用户`)
        // 如果当前正在查看详情的用户也在批量操作中，更新其状态
        if (currentUser.value && userIds.includes(currentUser.value.id)) {
          currentUser.value.status = '正常'
        }
      } else {
        ElMessage.warning('批量启用失败，请稍后重试')
      }
    } catch (error) {
      console.error('批量启用用户失败:', error)
      ElMessage.error('批量启用用户失败，请稍后重试')
    }
  }).catch(() => {
    // 取消操作
  })
}

// 批量禁用
const batchDisable = async () => {
  if (selectedUsers.value.length === 0) return
  
  ElMessageBox.confirm(
    `确定要批量禁用选中的 ${selectedUsers.value.length} 个用户吗？`,
    '批量操作',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      // 获取所有正常状态的用户ID
      const userIds = selectedUsers.value
        .filter(user => user.status !== '禁用')
        .map(user => user.id)
      
      if (userIds.length === 0) {
        ElMessage.info('所选用户均已禁用')
        return
      }
      
      // 批量请求API
      let successCount = 0
      for (const userId of userIds) {
        try {
          const response = await updateUserStatus(userId, '禁用')
          if (response.code === 200) {
            successCount++
            // 更新本地状态
            const targetUser = users.value.find(u => u.id === userId)
            if (targetUser) {
              targetUser.status = '禁用'
            }
          }
        } catch (error) {
          console.error(`禁用用户 ${userId} 失败:`, error)
        }
      }
      
      if (successCount > 0) {
        ElMessage.success(`已成功禁用 ${successCount} 个用户`)
        // 如果当前正在查看详情的用户也在批量操作中，更新其状态
        if (currentUser.value && userIds.includes(currentUser.value.id)) {
          currentUser.value.status = '禁用'
        }
      } else {
        ElMessage.warning('批量禁用失败，请稍后重试')
      }
    } catch (error) {
      console.error('批量禁用用户失败:', error)
      ElMessage.error('批量禁用用户失败，请稍后重试')
    }
  }).catch(() => {
    // 取消操作
  })
}

// 刷新表格
const refreshTable = () => {
  loadUsers()
}

// 导出数据
const exportUsers = async () => {
  try {
    const params = {
      search: searchQuery.value || undefined,
      status: statusFilter.value || undefined
    }
    
    // 添加日期范围参数
    if (dateRange.value && dateRange.value.length === 2) {
      params.startDate = formatDate(dateRange.value[0])
      params.endDate = formatDate(dateRange.value[1], true)
    }
    
    const response = await exportUsersAPI(params)
    // 处理blob响应
    const blob = new Blob([response], { type: 'application/vnd.ms-excel' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = `用户数据_${new Date().toLocaleDateString()}.xlsx`
    link.click()
    URL.revokeObjectURL(link.href)
    
    ElMessage.success('用户数据导出成功')
  } catch (error) {
    console.error('导出用户数据失败:', error)
    ElMessage.error('导出用户数据失败，请稍后重试')
  }
}

// 发送消息
const sendMessage = (user) => {
  messageForm.receiver = user.username
  messageForm.receiverId = user.id
  messageForm.title = ''
  messageForm.content = ''
  messageDialogVisible.value = true
}

// 提交消息
const submitMessage = async () => {
  messageFormRef.value.validate(async (valid) => {
    if (valid) {
      messageSending.value = true
      
      try {
        const response = await sendSystemMessageApi({
          receiverId: messageForm.receiverId,
          title: messageForm.title,
          content: messageForm.content
        })
        
        if (response.code === 200) {
          messageDialogVisible.value = false
          ElMessage.success(`消息已成功发送给 ${messageForm.receiver}`)
        } else {
          ElMessage.error(response.message || '发送消息失败')
        }
      } catch (error) {
        console.error('发送消息失败:', error)
        ElMessage.error('发送消息失败，请稍后重试')
      } finally {
        messageSending.value = false
      }
    }
  })
}

// 处理更多操作
const handleCommand = (command, user) => {
  switch (command) {
    case 'resetPassword':
      ElMessageBox.confirm(
        `确定要重置用户 ${user.username} 的密码吗？`,
        '重置密码',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).then(async () => {
        try {
          const response = await resetUserPassword(user.id)
          if (response.code === 200) {
            ElMessage.success(`已重置用户 ${user.username} 的密码`)
          } else {
            ElMessage.error(response.message || '重置密码失败')
          }
        } catch (error) {
          console.error('重置密码失败:', error)
          ElMessage.error('重置密码失败，请稍后重试')
        }
      }).catch(() => {
        // 取消操作
      })
      break
    case 'message':
      sendMessage(user)
      break
  }
}
</script>

<style scoped>
.admin-users-container {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  margin: 0;
  font-size: 24px;
  font-weight: 500;
  color: #303133;
}

.filter-card {
  margin-bottom: 20px;
  border-radius: 8px;
}

.filter-container {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  align-items: flex-start;
}

.filter-item {
  min-width: 200px;
}

.filter-actions {
  display: flex;
  gap: 10px;
  margin-left: auto;
}

.data-card {
  border-radius: 8px;
}

.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.action-group {
  margin-left: 10px;
}

.toolbar-right {
  display: flex;
  gap: 10px;
}

.search-input {
  width: 100%;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.user-info-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-avatar {
  background-color: #409EFF;
}

.user-info-text {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-weight: 500;
}

.user-role {
  font-size: 12px;
  color: #909399;
}

.user-detail {
  padding: 10px;
}

.user-profile-card {
  background-color: #f5f7fa;
  border-radius: 8px;
  padding: 20px;
  height: 100%;
}

.user-profile-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.profile-avatar {
  margin-bottom: 10px;
  background-color: #409EFF;
  font-size: 24px;
}

.user-profile-header h2 {
  margin: 10px 0;
  font-size: 18px;
}

.status-tag {
  margin-top: 5px;
}

.user-profile-stats {
  display: flex;
  justify-content: space-around;
  padding: 10px 0;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.stat-label {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

.user-profile-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.price-tag {
  color: #F56C6C;
  font-weight: bold;
}

@media (max-width: 768px) {
  .filter-container {
    flex-direction: column;
  }
  
  .filter-item {
    width: 100%;
  }
  
  .filter-actions {
    margin-left: 0;
    width: 100%;
    justify-content: flex-end;
  }
}
</style>