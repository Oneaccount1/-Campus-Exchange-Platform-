<template>
  <div class="admin-products-container">
    <div class="page-header">
      <h2 class="page-title">商品管理</h2>
      <el-button type="primary" @click="exportProducts">
        <el-icon><Download /></el-icon>导出数据
      </el-button>
    </div>

    <el-card shadow="hover" class="filter-card">
      <div class="filter-container">
        <div class="filter-item">
          <el-input
            v-model="searchQuery"
            placeholder="搜索商品名称"
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
          <el-select v-model="categoryFilter" placeholder="分类筛选" clearable @change="handleSearch">
            <el-option
              v-for="category in productStore.categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </div>
        <div class="filter-item">
          <el-select v-model="statusFilter" placeholder="状态筛选" clearable @change="handleSearch">
            <el-option label="已上架" value="已上架" />
            <el-option label="已下架" value="已下架" />
            <el-option label="待审核" value="待审核" />
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
          <el-tag type="success">共 {{ totalProducts }} 条记录</el-tag>
          <el-button-group class="action-group">
            <el-button type="primary" plain size="small" @click="batchApprove" :disabled="selectedProducts.length === 0">
              <el-icon><Check /></el-icon>批量上架
            </el-button>
            <el-button type="danger" plain size="small" @click="batchRemove" :disabled="selectedProducts.length === 0">
              <el-icon><Close /></el-icon>批量下架
            </el-button>
          </el-button-group>
        </div>
        <div class="toolbar-right">
          <el-tooltip content="刷新" placement="top">
            <el-button circle @click="refreshTable">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </el-tooltip>
          <el-tooltip content="列设置" placement="top">
            <el-button circle @click="columnSettingVisible = true">
              <el-icon><SetUp /></el-icon>
            </el-button>
          </el-tooltip>
        </div>
      </div>
      
      <el-table 
        :data="filteredProducts" 
        style="width: 100%" 
        v-loading="loading"
        @selection-change="handleSelectionChange"
        :header-cell-style="{ background: '#f5f7fa' }"
        border
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" sortable />
        <el-table-column label="商品图片" width="100" v-if="columnVisible.image">
          <template #default="scope">
            <el-image 
              style="width: 60px; height: 60px" 
              :src="scope.row.image" 
              fit="cover"
              :preview-src-list="[scope.row.image]"
              :initial-index="0"
            >
              <template #error>
                <div class="image-error">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="商品名称" show-overflow-tooltip sortable />
        <el-table-column prop="category" label="分类" width="120" v-if="columnVisible.category" />
        <el-table-column prop="price" label="价格" width="120" sortable>
          <template #default="scope">
            <span class="price-tag">¥{{ scope.row.price }}</span>
          </template>
        </el-table-column>
        <el-table-column label="卖家" width="120" v-if="columnVisible.seller">
          <template #default="scope">
            <span v-if="scope.row.user_id">{{ getUserName(scope.row.user_id) }}</span>
            <span v-else>{{ scope.row.seller || '未知' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="发布时间" width="180" sortable v-if="columnVisible.createTime">
          <template #default="scope">
            {{ formatDateTime(scope.row.created_at || scope.row.createTime) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="viewProductDetail(scope.row.id)">
              <el-icon><View /></el-icon>查看
            </el-button>
            <el-button 
              :type="scope.row.status === '已上架' ? 'danger' : 'success'" 
              size="small" 
              @click="toggleProductStatus(scope.row)"
            >
              <el-icon v-if="scope.row.status === '已上架'"><Close /></el-icon>
              <el-icon v-else><Check /></el-icon>
              {{ scope.row.status === '已上架' ? '下架' : '上架' }}
            </el-button>
            <el-dropdown trigger="click" @command="handleCommand($event, scope.row)">
              <el-button size="small" plain>
                <el-icon><MoreFilled /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="edit">
                    <el-icon><Edit /></el-icon>编辑
                  </el-dropdown-item>
                  <el-dropdown-item command="delete" divided>
                    <el-icon><Delete /></el-icon>删除
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
          :total="totalProducts"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          background
        />
      </div>
    </el-card>
    
    <!-- 商品详情对话框 -->
    <el-dialog
      v-model="productDetailVisible"
      title="商品详情"
      width="700px"
      destroy-on-close
    >
      <div v-if="currentProduct" class="product-detail">
        <div class="product-header">
          <el-carousel height="300px" indicator-position="outside" class="product-carousel" v-if="currentProduct.images && currentProduct.images.length > 0">
            <el-carousel-item v-for="(image, index) in currentProduct.images" :key="index">
              <el-image 
                class="carousel-image" 
                :src="image" 
                fit="cover"
                :preview-src-list="currentProduct.images"
              />
            </el-carousel-item>
          </el-carousel>
          <el-image 
            v-else
            class="product-image" 
            :src="currentProduct.image" 
            fit="cover"
            :preview-src-list="[currentProduct.image]"
          >
            <template #error>
              <div class="image-error large">
                <el-icon><Picture /></el-icon>
                <span>暂无图片</span>
              </div>
            </template>
          </el-image>
          <div class="product-info">
            <h2>{{ currentProduct.title }}</h2>
            <p class="product-price">¥{{ currentProduct.price }}</p>
            <el-descriptions :column="1" border>
              <el-descriptions-item label="分类">{{ currentProduct.category }}</el-descriptions-item>
              <el-descriptions-item label="卖家">
                <span v-if="currentProduct.user_id">{{ getUserName(currentProduct.user_id) }}</span>
                <span v-else>{{ currentProduct.seller || '未知' }}</span>
              </el-descriptions-item>
              <el-descriptions-item label="发布时间">{{ formatDateTime(currentProduct.created_at || currentProduct.createTime) }}</el-descriptions-item>
              <el-descriptions-item label="状态">
                <el-tag :type="getStatusType(currentProduct.status)">
                  {{ currentProduct.status }}
                </el-tag>
              </el-descriptions-item>
            </el-descriptions>
          </div>
        </div>
        
        <el-divider />
        
        <div class="product-description">
          <h3>商品描述</h3>
          <p>{{ currentProduct.description || '暂无描述' }}</p>
        </div>
        
        <div class="dialog-footer">
          <el-button @click="productDetailVisible = false">关闭</el-button>
          <el-button 
            :type="currentProduct.status === '已上架' ? 'danger' : 'success'" 
            @click="toggleProductStatus(currentProduct); productDetailVisible = false"
          >
            {{ currentProduct.status === '已上架' ? '下架商品' : '上架商品' }}
          </el-button>
        </div>
      </div>
    </el-dialog>
    
    <!-- 列设置对话框 -->
    <el-dialog
      v-model="columnSettingVisible"
      title="列设置"
      width="400px"
    >
      <el-checkbox-group v-model="visibleColumns">
        <el-checkbox label="image">商品图片</el-checkbox>
        <el-checkbox label="category">分类</el-checkbox>
        <el-checkbox label="seller">卖家</el-checkbox>
        <el-checkbox label="createTime">发布时间</el-checkbox>
      </el-checkbox-group>
      <template #footer>
        <el-button @click="columnSettingVisible = false">取消</el-button>
        <el-button type="primary" @click="saveColumnSettings">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore, useProductStore } from '../../stores'
import { ElMessageBox, ElMessage } from 'element-plus'
import { 
  getProducts, 
  getProductDetail, 
  updateProductStatus, 
  batchUpdateProductStatus,
  deleteProduct,
  exportProductsData
} from '../../api/admin'
import { getUserInfo } from '../../api/user' // 导入获取用户信息的API

const router = useRouter()
const adminStore = useAdminStore()
const productStore = useProductStore()
const products = ref([])
const userCache = ref({}) // 用户信息缓存

// 检查管理员是否登录
onMounted(() => {
  // 加载商品数据
  loadProducts()
})

const loading = ref(false)
const searchQuery = ref('')
const categoryFilter = ref('')
const statusFilter = ref('')
const dateRange = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const totalProducts = ref(0)
const selectedProducts = ref([])

// 列显示设置
const columnSettingVisible = ref(false)
const visibleColumns = ref(['image', 'category', 'seller', 'createTime'])
const columnVisible = computed(() => {
  const result = {}
  visibleColumns.value.forEach(col => {
    result[col] = true
  })
  return result
})

// 商品详情相关
const productDetailVisible = ref(false)
const currentProduct = ref(null)

// 过滤后的商品列表
const filteredProducts = computed(() => {
  return products.value
})

// 加载商品数据
const loadProducts = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      size: pageSize.value,
      search: searchQuery.value || undefined,
      category: categoryFilter.value || undefined,
      status: statusFilter.value || undefined
    }
    
    // 添加日期范围参数
    if (dateRange.value && dateRange.value.length === 2 && dateRange.value[0] && dateRange.value[1]) {
      params.start_date = formatDate(dateRange.value[0])
      params.end_date = formatDate(dateRange.value[1], true) // 设置为当天结束时间
    }
    
    const response = await getProducts(params)
    console.log('加载商品列表响应数据:', response)
    
    if (response.code === 200 && response.data) {
      // 兼容不同的返回数据结构
      let productData = []
      if (Array.isArray(response.data)) {
        // 如果直接返回数组
        productData = response.data
        totalProducts.value = response.data.length
      } else if (response.data.products) {
        // 如果返回带分页的数据结构
        productData = response.data.products
        totalProducts.value = response.data.total || 0
      } else {
        // 尝试从data对象中找到可能的数组
        const possibleList = Object.values(response.data).find(value => Array.isArray(value))
        if (possibleList) {
          productData = possibleList
          totalProducts.value = possibleList.length
        } else {
          products.value = []
          totalProducts.value = 0
          console.warn('无法解析返回的数据结构:', response.data)
          return // 早期返回，避免后续处理
        }
      }
      
      // 处理数据，适配字段
      products.value = productData.map(product => {
        return {
          ...product,
          createTime: formatDateTime(product.created_at || product.createTime),
          // 保留原始的user_id用于获取用户信息
          user_id: product.user_id !== undefined ? product.user_id : null,
          // 如果有description字段，保留它
          description: product.description || ''
        }
      })
      
      console.log('解析后的商品列表:', products.value)
      console.log('解析后的商品总数:', totalProducts.value)
    } else {
      ElMessage.error(response.message || '获取商品列表失败')
    }
  } catch (error) {
    console.error('加载商品数据失败:', error)
    ElMessage.error('加载商品数据失败，请稍后重试')
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
  if (!date) return ''
  
  const d = new Date(date)
  if (isEndOfDay) {
    d.setHours(23, 59, 59, 999)
  }
  // 使用toISOString()方法生成符合RFC3339标准的日期字符串
  return d.toISOString()
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  loadProducts()
}

// 重置过滤条件
const resetFilters = () => {
  searchQuery.value = ''
  categoryFilter.value = ''
  statusFilter.value = ''
  dateRange.value = []
  handleSearch()
}

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size
  loadProducts()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  loadProducts()
}

// 查看商品详情
const viewProductDetail = async (productId) => {
  try {
    const response = await getProductDetail(productId)
    if (response.code === 200 && response.data) {
      currentProduct.value = response.data
      productDetailVisible.value = true
    } else {
      ElMessage.error(response.message || '获取商品详情失败')
    }
  } catch (error) {
    console.error('获取商品详情失败:', error)
    ElMessage.error('获取商品详情失败，请稍后重试')
  }
}

// 切换商品状态（上架/下架）
const toggleProductStatus = async (product) => {
  const newStatus = product.status === '已上架' ? '已下架' : '已上架'
  const action = product.status === '已上架' ? '下架' : '上架'
  
  ElMessageBox.confirm(
    `确定要${action}商品 ${product.title} 吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const response = await updateProductStatus(product.id, newStatus)
      if (response.code === 200) {
        // 更新本地状态
        const targetProduct = products.value.find(p => p.id === product.id)
        if (targetProduct) {
          targetProduct.status = newStatus
        }
        
        // 如果当前正在查看该商品详情，也更新详情中的状态
        if (currentProduct.value && currentProduct.value.id === product.id) {
          currentProduct.value.status = newStatus
        }
        
        ElMessage.success(`已${action}商品 ${product.title}`)
      } else {
        ElMessage.error(response.message || `${action}商品失败`)
      }
    } catch (error) {
      console.error(`${action}商品失败:`, error)
      ElMessage.error(`${action}商品失败，请稍后重试`)
    }
  }).catch(() => {
    // 取消操作
  })
}

// 批量上架
const batchApprove = async () => {
  if (selectedProducts.value.length === 0) return
  
  ElMessageBox.confirm(
    `确定要批量上架选中的 ${selectedProducts.value.length} 个商品吗？`,
    '批量操作',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      // 获取所有非上架状态的商品ID
      const productIds = selectedProducts.value
        .filter(product => product.status !== '已上架')
        .map(product => product.id)
      
      if (productIds.length === 0) {
        ElMessage.info('所选商品均已上架')
        return
      }
      
      const response = await batchUpdateProductStatus(productIds, '已上架')
      if (response.code === 200) {
        // 更新本地状态
        productIds.forEach(id => {
          const product = products.value.find(p => p.id === id)
          if (product) {
            product.status = '已上架'
          }
        })
        
        ElMessage.success(`已成功上架 ${productIds.length} 个商品`)
        
        // 如果当前正在查看详情的商品也在批量操作中，更新其状态
        if (currentProduct.value && productIds.includes(currentProduct.value.id)) {
          currentProduct.value.status = '已上架'
        }
      } else {
        ElMessage.error(response.message || '批量上架失败')
      }
    } catch (error) {
      console.error('批量上架商品失败:', error)
      ElMessage.error('批量上架商品失败，请稍后重试')
    }
  }).catch(() => {
    // 取消操作
  })
}

// 批量下架
const batchRemove = async () => {
  if (selectedProducts.value.length === 0) return
  
  ElMessageBox.confirm(
    `确定要批量下架选中的 ${selectedProducts.value.length} 个商品吗？`,
    '批量操作',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      // 获取所有已上架状态的商品ID
      const productIds = selectedProducts.value
        .filter(product => product.status === '已上架')
        .map(product => product.id)
      
      if (productIds.length === 0) {
        ElMessage.info('所选商品均已下架')
        return
      }
      
      const response = await batchUpdateProductStatus(productIds, '已下架')
      if (response.code === 200) {
        // 更新本地状态
        productIds.forEach(id => {
          const product = products.value.find(p => p.id === id)
          if (product) {
            product.status = '已下架'
          }
        })
        
        ElMessage.success(`已成功下架 ${productIds.length} 个商品`)
        
        // 如果当前正在查看详情的商品也在批量操作中，更新其状态
        if (currentProduct.value && productIds.includes(currentProduct.value.id)) {
          currentProduct.value.status = '已下架'
        }
      } else {
        ElMessage.error(response.message || '批量下架失败')
      }
    } catch (error) {
      console.error('批量下架商品失败:', error)
      ElMessage.error('批量下架商品失败，请稍后重试')
    }
  }).catch(() => {
    // 取消操作
  })
}

// 删除商品
const deleteProductItem = async (product) => {
  ElMessageBox.confirm(
    `确定要删除商品 ${product.title} 吗？此操作不可恢复！`,
    '警告',
    {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const response = await deleteProduct(product.id)
      if (response.code === 200) {
        // 从列表中移除
        products.value = products.value.filter(p => p.id !== product.id)
        ElMessage.success(`已删除商品 ${product.title}`)
        
        // 如果正在查看详情，关闭详情对话框
        if (currentProduct.value && currentProduct.value.id === product.id) {
          productDetailVisible.value = false
        }
      } else {
        ElMessage.error(response.message || '删除商品失败')
      }
    } catch (error) {
      console.error('删除商品失败:', error)
      ElMessage.error('删除商品失败，请稍后重试')
    }
  }).catch(() => {
    // 取消操作
  })
}

// 处理更多操作
const handleCommand = (command, product) => {
  if (command === 'edit') {
    router.push(`/admin/products/edit/${product.id}`)
  } else if (command === 'delete') {
    deleteProductItem(product)
  }
}

// 刷新表格
const refreshTable = () => {
  loadProducts()
}

// 导出数据
const exportProducts = async () => {
  try {
    const params = {
      search: searchQuery.value || undefined,
      category: categoryFilter.value || undefined,
      status: statusFilter.value || undefined
    }
    
    // 添加日期范围参数
    if (dateRange.value && dateRange.value.length === 2 && dateRange.value[0] && dateRange.value[1]) {
      params.start_date = formatDate(dateRange.value[0])
      params.end_date = formatDate(dateRange.value[1], true)
    }
    
    const response = await exportProductsData(params)
    // 处理blob响应
    const blob = new Blob([response], { type: 'application/vnd.ms-excel' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = `商品数据_${new Date().toLocaleDateString()}.xlsx`
    link.click()
    URL.revokeObjectURL(link.href)
    
    ElMessage.success('商品数据导出成功')
  } catch (error) {
    console.error('导出商品数据失败:', error)
    ElMessage.error('导出商品数据失败，请稍后重试')
  }
}

// 处理表格选择
const handleSelectionChange = (selection) => {
  selectedProducts.value = selection
}

// 获取状态标签类型
const getStatusType = (status) => {
  switch (status) {
    case '已上架': return 'success'
    case '已下架': return 'danger'
    case '待审核': return 'warning'
    default: return 'info'
  }
}

// 根据userId获取用户名
const getUserName = (userId) => {
  // 如果缓存中存在，直接返回缓存的用户名
  if (userCache.value[userId]) {
    return userCache.value[userId];
  }

  // 初始显示 "加载中..."
  userCache.value[userId] = "加载中...";

  // 异步获取用户信息
  getUserInfo(userId).then(response => {
    if (response && response.code === 200 && response.data) {
      userCache.value[userId] = response.data.username || response.data.name || `用户${userId}`;
    } else {
      userCache.value[userId] = `用户${userId}`;
    }
  }).catch(error => {
    console.error('获取用户信息失败:', error);
    userCache.value[userId] = `用户${userId}`;
  });

  return userCache.value[userId];
};

// 保存列设置
const saveColumnSettings = () => {
  // 这里可以保存到本地存储或用户偏好设置
  columnSettingVisible.value = false
  ElMessage.success('列设置已保存')
}
</script>

<style scoped>
.admin-products-container {
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
  height: 300px;
  flex-direction: column;
  gap: 10px;
}

.product-detail {
  padding: 10px;
}

.product-header {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-bottom: 20px;
}

.product-image {
  width: 100%;
  height: 300px;
  border-radius: 8px;
  overflow: hidden;
}

.product-carousel {
  width: 100%;
  border-radius: 8px;
  overflow: hidden;
}

.carousel-image {
  width: 100%;
  height: 100%;
}

.product-info {
  flex: 1;
}

.product-price {
  font-size: 24px;
  color: #f56c6c;
  font-weight: bold;
  margin-bottom: 20px;
}

.product-description h3 {
  margin-bottom: 10px;
  font-size: 18px;
  font-weight: 500;
}

.product-description p {
  line-height: 1.6;
  color: #606266;
}

.dialog-footer {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.price-tag {
  color: #F56C6C;
  font-weight: bold;
}

@media (min-width: 768px) {
  .product-header {
    flex-direction: row;
  }
  
  .product-image,
  .product-carousel {
    width: 300px;
    flex-shrink: 0;
  }
}
</style>