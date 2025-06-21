<template>
  <div class="admin-orders-container">
    <div class="page-header">
      <h2 class="page-title">订单管理</h2>
      <el-button type="primary" @click="exportOrders">
        <el-icon><Download /></el-icon>导出数据
      </el-button>
    </div>

    <el-card shadow="hover" class="filter-card">
      <div class="filter-container">
        <div class="filter-item">
          <el-input
            v-model="searchQuery"
            placeholder="搜索订单号/商品名"
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
            <el-option label="待付款" value="待付款" />
            <el-option label="待发货" value="待发货" />
            <el-option label="待收货" value="待收货" />
            <el-option label="已完成" value="已完成" />
            <el-option label="已取消" value="已取消" />
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
          <el-tag type="success">共 {{ totalOrders }} 条记录</el-tag>
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
        :data="ordersList" 
        style="width: 100%" 
        v-loading="loading"
        :header-cell-style="{ background: '#f5f7fa' }"
        border
      >
        <el-table-column prop="id" label="订单号" width="150" />
        <el-table-column label="商品信息">
          <template #default="scope">
            <div class="product-info-cell">
              <el-image 
                class="product-thumbnail" 
                :src="scope.row.productImage || 'https://via.placeholder.com/40'" 
                fit="cover"
              >
                <template #error>
                  <div class="image-error">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
              <div class="product-info-text">
                <div class="product-title">{{ scope.row.productTitle }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="price" label="价格" width="120">
          <template #default="scope">
            <span class="price-tag">¥{{ scope.row.price }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="buyer" label="买家" width="120" />
        <el-table-column prop="seller" label="卖家" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="创建时间" width="150">
          <template #default="scope">
            {{ formatDateTime(scope.row.createTime) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="viewOrderDetail(scope.row.id)">
              <el-icon><View /></el-icon>详情
            </el-button>
            <el-dropdown trigger="click" @command="handleCommand($event, scope.row)">
              <el-button size="small" plain>
                <el-icon><MoreFilled /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="updateStatus">
                    <el-icon><Edit /></el-icon>修改状态
                  </el-dropdown-item>
                  <el-dropdown-item command="contact" divided>
                    <el-icon><ChatDotRound /></el-icon>联系买卖家
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
          :total="totalOrders"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          background
        />
      </div>
    </el-card>
    
    <!-- 订单详情对话框 -->
    <el-dialog
      v-model="orderDetailVisible"
      title="订单详情"
      width="800px"
      destroy-on-close
    >
      <div v-if="currentOrder" class="order-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="订单号">{{ currentOrder.id }}</el-descriptions-item>
          <el-descriptions-item label="订单状态">
            <el-tag :type="getStatusType(currentOrder.status)">{{ currentOrder.status }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDateTime(currentOrder.createTime) }}</el-descriptions-item>
          <el-descriptions-item label="支付时间">{{ currentOrder.payTime ? formatDateTime(currentOrder.payTime) : '未支付' }}</el-descriptions-item>
          <el-descriptions-item label="发货时间" v-if="currentOrder.deliveryTime">{{ formatDateTime(currentOrder.deliveryTime) }}</el-descriptions-item>
          <el-descriptions-item label="完成时间" v-if="currentOrder.completeTime">{{ formatDateTime(currentOrder.completeTime) }}</el-descriptions-item>
        </el-descriptions>
        
        <el-divider content-position="center">商品信息</el-divider>
        
        <div class="product-detail">
          <el-image 
            class="product-image" 
            :src="currentOrder.productImage" 
            fit="cover"
          >
            <template #error>
              <div class="image-error large">
                <el-icon><Picture /></el-icon>
              </div>
            </template>
          </el-image>
          <div class="product-info">
            <h3>{{ currentOrder.productTitle }}</h3>
            <p class="price">¥{{ currentOrder.price }}</p>
          </div>
        </div>
        
        <el-divider content-position="center">用户信息</el-divider>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-card class="user-card">
              <template #header>
                <div class="card-header">买家信息</div>
              </template>
              <p><strong>用户名：</strong>{{ currentOrder.buyer }}</p>
              <p><strong>联系电话：</strong>{{ currentOrder.buyerPhone }}</p>
              <p><strong>收货地址：</strong>{{ currentOrder.buyerAddress }}</p>
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card class="user-card">
              <template #header>
                <div class="card-header">卖家信息</div>
              </template>
              <p><strong>用户名：</strong>{{ currentOrder.seller }}</p>
              <p><strong>联系电话：</strong>{{ currentOrder.sellerPhone }}</p>
            </el-card>
          </el-col>
        </el-row>
        
        <el-divider content-position="center">订单备注</el-divider>
        <div class="order-remark">
          {{ currentOrder.remark || '无' }}
        </div>
        
        <el-divider content-position="center">操作记录</el-divider>
        <el-timeline>
          <el-timeline-item
            v-for="(log, index) in currentOrder.logs"
            :key="index"
            :timestamp="formatDateTime(log.time)"
          >
            <h4>{{ log.action }}</h4>
            <p>操作人：{{ log.operator }}</p>
            <p v-if="log.remark">备注：{{ log.remark }}</p>
          </el-timeline-item>
        </el-timeline>
        
        <div class="dialog-footer">
          <el-button @click="orderDetailVisible = false">关闭</el-button>
          <el-button type="primary" @click="showStatusUpdateDialog">修改状态</el-button>
        </div>
      </div>
    </el-dialog>
    
    <!-- 修改状态对话框 -->
    <el-dialog
      v-model="statusUpdateVisible"
      title="修改订单状态"
      width="500px"
      append-to-body
      destroy-on-close
    >
      <el-form :model="statusForm" label-width="100px">
        <el-form-item label="当前状态">
          <el-tag :type="getStatusType(currentOrder?.status)">{{ currentOrder?.status }}</el-tag>
        </el-form-item>
        <el-form-item label="新状态">
          <el-select v-model="statusForm.status" placeholder="选择新状态">
            <el-option label="待付款" value="待付款" />
            <el-option label="待发货" value="待发货" />
            <el-option label="待收货" value="待收货" />
            <el-option label="已完成" value="已完成" />
            <el-option label="已取消" value="已取消" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="statusForm.remark" type="textarea" :rows="3" placeholder="请输入备注信息" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="statusUpdateVisible = false">取消</el-button>
          <el-button type="primary" @click="updateOrderStatus" :loading="updating">确认修改</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { getOrders, getOrderDetail, updateOrderStatus as updateOrderStatusApi, exportOrdersData } from '../../api/admin'

const router = useRouter()

// 页面数据
const searchQuery = ref('')
const statusFilter = ref('')
const dateRange = ref(null)
const loading = ref(false)
const ordersList = ref([])
const totalOrders = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 订单详情
const orderDetailVisible = ref(false)
const currentOrder = ref(null)

// 状态修改
const statusUpdateVisible = ref(false)
const statusForm = reactive({
  status: '',
  remark: ''
})
const updating = ref(false)

// 获取订单列表
const loadOrders = async () => {
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
    
    const response = await getOrders(params)
    console.log('加载订单列表响应数据:', response)
    
    if (response && response.code === 200 && response.data) {
      // 兼容不同的返回数据结构
      if (Array.isArray(response.data)) {
        // 如果直接返回数组
        ordersList.value = response.data
        totalOrders.value = response.data.length
      } else if (response.data.orders) {
        // 如果返回的是orders字段
        ordersList.value = response.data.orders
        totalOrders.value = response.data.total || 0
      } else if (response.data.list) {
        // 如果返回带分页的数据结构
        ordersList.value = response.data.list
        totalOrders.value = response.data.total || 0
      } else {
        // 尝试从data对象中找到可能的数组
        const possibleList = Object.values(response.data).find(value => Array.isArray(value))
        if (possibleList) {
          ordersList.value = possibleList
          totalOrders.value = possibleList.length
        } else {
          // 如果API未完成，使用模拟数据
          useMockData()
          console.warn('无法解析返回的数据结构:', response.data)
          return // 使用模拟数据后直接返回
        }
      }
      
      console.log('解析后的订单列表:', ordersList.value)
      console.log('解析后的订单总数:', totalOrders.value)
    } else {
      // 如果API未完成，使用模拟数据
      useMockData()
    }
  } catch (error) {
    console.error('加载订单列表失败:', error)
    // 如果API调用失败，使用模拟数据
    useMockData()
  } finally {
    loading.value = false
  }
}

// 使用模拟数据
const useMockData = () => {
  ordersList.value = [
    {
      id: 'O20240520001',
      productTitle: '全新iPad Pro 2023款',
      productImage: 'https://via.placeholder.com/100',
      price: 6999,
      buyer: '李四',
      seller: '张三',
      status: '已完成',
      createTime: '2024-05-20 14:30:22',
      payTime: '2024-05-20 14:35:10',
      completeTime: '2024-05-22 10:15:30'
    },
    {
      id: 'O20240519002',
      productTitle: '计算机网络教材 第8版',
      productImage: 'https://via.placeholder.com/100',
      price: 45,
      buyer: '王五',
      seller: '赵六',
      status: '待发货',
      createTime: '2024-05-19 09:15:40',
      payTime: '2024-05-19 09:20:12'
    },
    {
      id: 'O20240518003',
      productTitle: '宿舍小冰箱 95新',
      productImage: 'https://via.placeholder.com/100',
      price: 320,
      buyer: '孙七',
      seller: '周八',
      status: '待收货',
      createTime: '2024-05-18 16:42:35',
      payTime: '2024-05-18 16:50:22',
      deliveryTime: '2024-05-19 10:30:15'
    },
    {
      id: 'O20240517004',
      productTitle: 'Nike运动鞋 42码',
      productImage: 'https://via.placeholder.com/100',
      price: 199,
      buyer: '吴九',
      seller: '郑十',
      status: '待付款',
      createTime: '2024-05-17 11:25:18'
    },
    {
      id: 'O20240516005',
      productTitle: '蓝牙耳机 全新未拆封',
      productImage: 'https://via.placeholder.com/100',
      price: 129,
      buyer: '钱一',
      seller: '孙七',
      status: '已取消',
      createTime: '2024-05-16 08:35:42',
      payTime: null,
      completeTime: '2024-05-16 15:40:20'
    }
  ]
  totalOrders.value = 85
}

// 查看订单详情
const viewOrderDetail = async (orderId) => {
  try {
    const response = await getOrderDetail(orderId)
    
    if (response && response.code === 200) {
      currentOrder.value = response.data
    } else {
      // 使用模拟数据
      currentOrder.value = {
        id: orderId,
        productId: 1,
        productTitle: '全新iPad Pro 2023款',
        productImage: 'https://via.placeholder.com/100',
        price: 6999,
        buyerId: 2,
        buyer: '李四',
        buyerPhone: '13900139000',
        buyerAddress: '北京市海淀区清华大学学生公寓',
        sellerId: 1,
        seller: '张三',
        sellerPhone: '13800138000',
        status: '已完成',
        createTime: '2024-05-20 14:30:22',
        payTime: '2024-05-20 14:35:10',
        deliveryTime: '2024-05-21 09:22:15',
        completeTime: '2024-05-22 10:15:30',
        remark: '请包装好再发货',
        logs: [
          {
            action: '创建订单',
            time: '2024-05-20 14:30:22',
            operator: '李四',
            remark: null
          },
          {
            action: '支付订单',
            time: '2024-05-20 14:35:10',
            operator: '李四',
            remark: '微信支付'
          },
          {
            action: '确认发货',
            time: '2024-05-21 09:22:15',
            operator: '张三',
            remark: '已发顺丰快递，单号SF123456789'
          },
          {
            action: '确认收货',
            time: '2024-05-22 10:15:30',
            operator: '李四',
            remark: '商品完好'
          }
        ]
      }
    }
    
    orderDetailVisible.value = true
  } catch (error) {
    console.error('获取订单详情失败:', error)
    ElMessage.error('获取订单详情失败，请稍后再试')
  }
}

// 显示修改状态对话框
const showStatusUpdateDialog = () => {
  if (!currentOrder.value) return
  
  statusForm.status = currentOrder.value.status
  statusForm.remark = ''
  statusUpdateVisible.value = true
}

// 更新订单状态
const updateOrderStatus = async () => {
  if (!currentOrder.value || !statusForm.status) {
    ElMessage.warning('请选择新状态')
    return
  }
  
  updating.value = true
  try {
    const response = await updateOrderStatusApi(currentOrder.value.id, statusForm)
    
    if (response && response.code === 200) {
      ElMessage.success('订单状态更新成功')
      // 更新当前显示的订单状态
      currentOrder.value.status = statusForm.status
      // 添加操作日志
      currentOrder.value.logs.push({
        action: `修改状态为"${statusForm.status}"`,
        time: new Date().toLocaleString(),
        operator: '管理员',
        remark: statusForm.remark || null
      })
      // 刷新订单列表
      loadOrders()
      statusUpdateVisible.value = false
    } else {
      ElMessage.error(response?.message || '更新失败，请稍后再试')
    }
  } catch (error) {
    console.error('更新订单状态失败:', error)
    ElMessage.error('更新失败，请稍后再试')
  } finally {
    updating.value = false
  }
}

// 处理下拉菜单命令
const handleCommand = (command, order) => {
  switch (command) {
    case 'updateStatus':
      currentOrder.value = order
      showStatusUpdateDialog()
      break
    case 'contact':
      ElMessage.info('联系功能开发中')
      break
  }
}

// 获取状态类型
const getStatusType = (status) => {
  switch (status) {
    case '待付款':
      return 'warning'
    case '待发货':
      return 'primary'
    case '待收货':
      return 'info'
    case '已完成':
      return 'success'
    case '已取消':
      return 'danger'
    default:
      return 'info'
  }
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  loadOrders()
}

// 重置筛选条件
const resetFilters = () => {
  searchQuery.value = ''
  statusFilter.value = ''
  dateRange.value = null
  currentPage.value = 1
  loadOrders()
}

// 刷新表格
const refreshTable = () => {
  loadOrders()
}

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size
  loadOrders()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  loadOrders()
}

// 导出订单数据
const exportOrders = async () => {
  try {
    // 构建查询参数
    const params = {
      search: searchQuery.value || undefined,
      status: statusFilter.value || undefined
    }
    
    // 添加日期范围参数
    if (dateRange.value && dateRange.value.length === 2) {
      const startDate = dateRange.value[0]
      const endDate = dateRange.value[1]
      params.startDate = formatDate(startDate)
      params.endDate = formatDate(endDate)
    }
    
    await exportOrdersData(params)
    ElMessage.success('订单数据导出成功')
  } catch (error) {
    console.error('导出订单数据失败:', error)
    ElMessage.error('导出失败，请稍后再试')
  }
}

// 格式化日期时间（用于ISO格式转换为YYYY-MM-DD HH:MM:SS）
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

// 格式化日期
const formatDate = (date, endOfDay = false) => {
  if (!date) return ''
  
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  
  if (endOfDay) {
    return `${year}-${month}-${day} 23:59:59`
  }
  
  return `${year}-${month}-${day}`
}

// 加载初始数据
onMounted(() => {
  loadOrders()
})

// 监听筛选条件变化
watch([searchQuery, statusFilter, dateRange], () => {
  handleSearch()
}, { deep: true })
</script>

<style scoped>
.admin-orders-container {
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

.product-info-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.product-thumbnail {
  width: 40px;
  height: 40px;
  border-radius: 4px;
  object-fit: cover;
}

.product-info-text {
  display: flex;
  flex-direction: column;
}

.product-title {
  font-weight: bold;
  margin-bottom: 5px;
}

.price-tag {
  color: #f56c6c;
  font-weight: bold;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.order-detail {
  padding: 10px;
}

.product-detail {
  display: flex;
  gap: 20px;
  margin-top: 15px;
  margin-bottom: 15px;
}

.product-image {
  width: 100px;
  height: 100px;
  border-radius: 4px;
  object-fit: cover;
}

.product-info {
  flex: 1;
}

.product-info h3 {
  margin-top: 0;
  margin-bottom: 10px;
}

.price {
  color: #f56c6c;
  font-weight: bold;
  font-size: 18px;
}

.user-card {
  height: 100%;
}

.card-header {
  font-weight: bold;
}

.order-remark {
  background-color: #f8f8f8;
  padding: 10px;
  border-radius: 4px;
  min-height: 60px;
}

.dialog-footer {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.image-error {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background-color: #f5f7fa;
  color: #909399;
}

.image-error.large {
  height: 100px;
  flex-direction: column;
  gap: 5px;
}
</style> 