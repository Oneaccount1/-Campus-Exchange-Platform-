<template>
  <div class="admin-dashboard-container">
    <div v-if="isLoading" class="dashboard-loading">
      <el-card shadow="hover" class="loading-card">
        <template #header>
          <div class="card-header loading-header">
            <h3>加载中</h3>
          </div>
        </template>
        <div class="loading-content">
          <el-skeleton :rows="6" animated />
          <div class="loading-message">正在加载仪表盘数据，请稍候...</div>
        </div>
      </el-card>
    </div>
    
    <template v-else>
      <div class="dashboard-header">
        <h2 class="page-title">仪表盘</h2>
        <div class="date-info">
          <el-icon><Calendar /></el-icon>
          <span>{{ currentDate }}</span>
        </div>
      </div>

      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="6" v-for="(card, index) in statCards" :key="index">
          <el-card shadow="hover" class="stat-card" :body-style="{ padding: '0px' }">
            <div class="stat-card-content">
              <div class="stat-info">
                <div class="stat-title">{{ card.title }}</div>
                <div class="stat-value">{{ card.value }}</div>
                <div class="stat-trend" :class="{'up': card.trend > 0, 'down': card.trend < 0}" v-if="card.trend !== undefined">
                  {{ card.trend > 0 ? '+' : '' }}{{ card.trend }}% <el-icon><ArrowUp v-if="card.trend > 0" /><ArrowDown v-else /></el-icon>
                </div>
              </div>
              <div class="stat-icon" :style="{ backgroundColor: card.color }">
                <el-icon><component :is="card.icon" /></el-icon>
              </div>
            </div>
            <div class="stat-footer" v-if="card.footer">
              {{ card.footer }}
            </div>
          </el-card>
        </el-col>
      </el-row>
      
      <el-row :gutter="20" class="chart-row">
        <el-col :xs="24" :md="16">
          <el-card shadow="hover" class="chart-card">
            <template #header>
              <div class="card-header">
                <div class="header-title">
                  <span>商品发布趋势</span>
                  <el-tag size="small" effect="plain" type="info">最近7天</el-tag>
                </div>
                <div class="header-actions">
                  <el-radio-group v-model="trendTimeRange" size="small">
                    <el-radio-button label="week">周</el-radio-button>
                    <el-radio-button label="month">月</el-radio-button>
                  </el-radio-group>
                </div>
              </div>
            </template>
            <div class="trend-chart">
              <!-- 这里可以集成图表库如ECharts -->
              <div class="chart-placeholder">
                <div class="mock-chart">
                  <template v-if="trendData.length > 0">
                    <div v-for="(bar, index) in trendData" :key="index" class="mock-bar-item">
                      <div class="mock-bar" :style="{ height: bar.value + '%', backgroundColor: '#409EFF' }"></div>
                      <div class="mock-label">{{ bar.label }}</div>
                    </div>
                  </template>
                  <div v-else class="loading-placeholder">
                    <el-skeleton :rows="5" animated />
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :xs="24" :md="8">
          <el-card shadow="hover" class="chart-card">
            <template #header>
              <div class="card-header">
                <span>商品分类统计</span>
              </div>
            </template>
            <div class="category-chart">
              <div class="chart-placeholder">
                <div class="mock-pie-chart">
                  <div class="mock-pie-center"></div>
                  <template v-if="categoryData.length > 0">
                    <div v-for="(slice, index) in categoryData" :key="index" 
                         class="mock-pie-slice" 
                         :style="{ 
                           backgroundColor: slice.color,
                           transform: `rotate(${slice.startAngle}deg)`,
                           clipPath: `polygon(50% 50%, 50% 0%, ${50 + 50 * Math.cos(slice.endAngle * Math.PI / 180)}% ${50 - 50 * Math.sin(slice.endAngle * Math.PI / 180)}%, 50% 50%)`
                         }">
                    </div>
                  </template>
                </div>
                <div class="mock-pie-legend">
                  <template v-if="categoryData.length > 0">
                    <div v-for="(item, index) in categoryData" :key="index" class="legend-item">
                      <div class="legend-color" :style="{ backgroundColor: item.color }"></div>
                      <div class="legend-label">{{ item.name }}</div>
                      <div class="legend-value">{{ item.value }}</div>
                    </div>
                  </template>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
      
      <el-row :gutter="20">
        <el-col :span="24">
          <el-card shadow="hover" class="table-card">
            <template #header>
              <div class="card-header">
                <div class="header-title">
                  <span>最新商品</span>
                  <el-tag size="small" effect="plain" type="success">今日新增</el-tag>
                </div>
                <el-button type="primary" size="small" @click="$router.push('/admin/products')">
                  <el-icon><View /></el-icon>查看全部
                </el-button>
              </div>
            </template>
            
            <el-table :data="latestProducts" style="width: 100%" :header-cell-style="{ background: '#f5f7fa' }" stripe>
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column label="商品信息">
                <template #default="scope">
                  <div class="product-info-cell">
                    <el-image 
                      class="product-thumbnail" 
                      :src="scope.row.image || 'https://via.placeholder.com/40'" 
                      fit="cover"
                    >
                      <template #error>
                        <div class="image-error">
                          <el-icon><Picture /></el-icon>
                        </div>
                      </template>
                    </el-image>
                    <div class="product-info-text">
                      <div class="product-title">{{ scope.row.title }}</div>
                      <div class="product-category">{{ scope.row.category }}</div>
                    </div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="price" label="价格" width="120">
                <template #default="scope">
                  <span class="price-tag">¥{{ scope.row.price }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="seller" label="卖家" width="120" />
              <el-table-column prop="createTime" label="发布时间" width="180" />
              <el-table-column label="操作" width="150">
                <template #default="scope">
                  <el-button type="primary" size="small" plain @click="viewProduct(scope.row.id)">
                    <el-icon><View /></el-icon>查看
                  </el-button>
                  <el-button type="danger" size="small" plain @click="removeProduct(scope.row.id)">
                    <el-icon><Delete /></el-icon>下架
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
      </el-row>
      
      <el-row :gutter="20" class="chart-row">
        <el-col :xs="24" :md="12">
          <el-card shadow="hover" class="activity-card">
            <template #header>
              <div class="card-header">
                <span>系统动态</span>
              </div>
            </template>
            <div class="activity-timeline">
              <el-timeline>
                <template v-if="recentActivities.length > 0">
                  <el-timeline-item
                    v-for="(activity, index) in recentActivities"
                    :key="index"
                    :type="activity.type"
                    :color="activity.color"
                    :timestamp="activity.time"
                    :hollow="activity.hollow"
                  >
                    {{ activity.content }}
                  </el-timeline-item>
                </template>
                <el-timeline-item v-else>
                  <span class="no-data">暂无系统活动</span>
                </el-timeline-item>
              </el-timeline>
            </div>
          </el-card>
        </el-col>
        
        <el-col :xs="24" :md="12">
          <el-card shadow="hover" class="quick-actions-card">
            <template #header>
              <div class="card-header">
                <span>快捷操作</span>
              </div>
            </template>
            <div class="quick-actions">
              <el-row :gutter="20">
                <el-col :span="8" v-for="(action, index) in quickActions" :key="index">
                  <div class="quick-action-item" @click="handleQuickAction(action)">
                    <el-icon><component :is="action.icon" /></el-icon>
                    <span>{{ action.label }}</span>
                  </div>
                </el-col>
              </el-row>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '../../stores'
import { ElMessageBox, ElMessage } from 'element-plus'
import { 
  getDashboardStats, 
  getProductTrend, 
  getCategoryStats, 
  getLatestProducts, 
  getRecentActivities 
} from '../../api/admin'

const router = useRouter()
const adminStore = useAdminStore()

// 检查管理员是否登录
onMounted(() => {
  // 初始化数据结构，避免未加载时的渲染错误
  initDefaultData()
  
  // 加载仪表盘数据
  loadDashboardData()
})

// 当前日期
const currentDate = computed(() => {
  const now = new Date()
  const options = { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' }
  return now.toLocaleDateString('zh-CN', options)
})

// 统计卡片数据
const statCards = ref([
  {
    title: '商品总数',
    value: 0,
    icon: 'ShoppingBag',
    color: '#409EFF',
    trend: 0,
    footer: '加载中...'
  },
  {
    title: '用户总数',
    value: 0,
    icon: 'User',
    color: '#67C23A',
    trend: 0,
    footer: '加载中...'
  },
  {
    title: '今日新增商品',
    value: 0,
    icon: 'Plus',
    color: '#E6A23C',
    trend: 0,
    footer: '加载中...'
  },
  {
    title: '交易总额',
    value: 0,
    icon: 'Money',
    color: '#F56C6C',
    trend: 0,
    footer: '加载中...'
  }
])

// 商品趋势时间范围
const trendTimeRange = ref('week')
// 商品趋势数据
const trendData = ref([])
// 商品分类数据
const categoryData = ref([])
// 最新商品
const latestProducts = ref([])
// 系统动态
const recentActivities = ref([])
// 快捷操作
const quickActions = ref([
  { icon: 'Plus', label: '发布公告', action: 'publishNotice' },
  { icon: 'UserFilled', label: '添加用户', action: 'addUser' },
  { icon: 'ShoppingBag', label: '商品管理', action: 'manageProducts' },
  { icon: 'Setting', label: '系统设置', action: 'systemSettings' },
  { icon: 'Document', label: '查看日志', action: 'viewLogs' },
  { icon: 'Refresh', label: '刷新缓存', action: 'refreshCache' }
])

// 添加一个简单的全局加载状态
const isLoading = ref(false)

// 初始化默认数据
const initDefaultData = () => {
  // 初始化默认值以防止空数据渲染错误
  trendData.value = []
  categoryData.value = []
  latestProducts.value = []
  recentActivities.value = []
}

// 加载仪表盘数据
const loadDashboardData = async () => {
  isLoading.value = true
  try {
    // 加载统计数据
    await loadStatsData()
    
    // 加载商品趋势数据
    await loadTrendData()
    
    // 加载分类统计数据
    await loadCategoryData()
    
    // 加载最新商品
    await loadLatestProducts()
    
    // 加载系统动态
    await loadRecentActivities()
  } catch (error) {
    console.error('加载仪表盘数据失败:', error)
    ElMessage.error('加载仪表盘数据失败，请稍后再试')
  } finally {
    isLoading.value = false
  }
}

// 加载统计数据
const loadStatsData = async () => {
  try {
    const response = await getDashboardStats()
    
    if (response.code === 200 && response.data) {
      const data = response.data
      
      // 更新统计卡片数据
      statCards.value[0].value = data.productCount || 0
      statCards.value[0].trend = data.productTrend || 0
      statCards.value[0].footer = `总计${data.productCount || 0}件商品`
      
      statCards.value[1].value = data.userCount || 0
      statCards.value[1].trend = data.userTrend || 0
      statCards.value[1].footer = `本周新增${data.newUserCount || 0}名用户`
      
      statCards.value[2].value = data.todayProductCount || 0
      statCards.value[2].trend = data.todayProductTrend || 0
      statCards.value[2].footer = `昨日新增${data.yesterdayProductCount || 0}件商品`
      
      statCards.value[3].value = data.totalAmount || 0
      statCards.value[3].trend = data.amountTrend || 0
      statCards.value[3].footer = `本月交易额¥${data.monthAmount || 0}`
    } else {
      ElMessage.error(response.message || '获取统计数据失败')
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
    ElMessage.error('加载统计数据失败，请稍后再试')
  }
}

// 加载商品趋势数据
const loadTrendData = async () => {
  try {
    const response = await getProductTrend(trendTimeRange.value)
    
    if (response.code === 200 && response.data) {
      trendData.value = response.data || []
    } else {
      ElMessage.error(response.message || '获取商品趋势数据失败')
      trendData.value = []
    }
  } catch (error) {
    console.error('加载趋势数据失败:', error)
    ElMessage.error('加载趋势数据失败，请稍后再试')
    trendData.value = []
  }
}

// 监听趋势时间范围变化
watch(trendTimeRange, () => {
  loadTrendData()
})

// 加载分类统计数据
const loadCategoryData = async () => {
  try {
    const response = await getCategoryStats()
    
    if (response.code === 200 && response.data) {
      // 处理分类数据，计算饼图角度
      const data = response.data || []
      let startAngle = 0
      
      if (data.length > 0) {
        categoryData.value = data.map((item, index) => {
          const percentage = item.percentage || (item.value / data.reduce((sum, i) => sum + i.value, 0)) * 100
          const angle = (percentage / 100) * 360
          const result = {
            ...item,
            startAngle,
            endAngle: startAngle + angle,
            color: getCategoryColor(index)
          }
          startAngle += angle
          return result
        })
      } else {
        categoryData.value = []
      }
    } else {
      ElMessage.error(response.message || '获取分类统计数据失败')
      categoryData.value = []
    }
  } catch (error) {
    console.error('加载分类数据失败:', error)
    ElMessage.error('加载分类数据失败，请稍后再试')
    categoryData.value = []
  }
}

// 获取分类颜色
const getCategoryColor = (index) => {
  const colors = ['#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399', '#9B59B6', '#3498DB', '#1ABC9C']
  return colors[index % colors.length]
}

// 加载最新商品
const loadLatestProducts = async () => {
  try {
    const response = await getLatestProducts({ limit: 5 })
    
    if (response.code === 200 && response.data) {
      latestProducts.value = response.data || []
    } else {
      ElMessage.error(response.message || '获取最新商品失败')
      latestProducts.value = [] // 确保在失败时有一个空数组
    }
  } catch (error) {
    console.error('加载最新商品失败:', error)
    ElMessage.error('加载最新商品失败，请稍后再试')
    latestProducts.value = [] // 确保在出错时有一个空数组
  }
}

// 加载系统动态
const loadRecentActivities = async () => {
  try {
    const response = await getRecentActivities(5)
    
    if (response.code === 200 && response.data) {
      recentActivities.value = response.data || []
    } else {
      ElMessage.error(response.message || '获取系统动态失败')
      recentActivities.value = []
    }
  } catch (error) {
    console.error('加载系统动态失败:', error)
    ElMessage.error('加载系统动态失败，请稍后再试')
    recentActivities.value = []
  }
}

// 查看商品详情
const viewProduct = (productId) => {
  router.push(`/admin/products?id=${productId}`)
}

// 下架商品
const removeProduct = (productId) => {
  ElMessageBox.confirm(
    '确定要下架该商品吗？',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    // 这里应该调用API下架商品
    ElMessage.success('商品已下架')
    // 重新加载最新商品数据
    loadLatestProducts()
  }).catch(() => {
    // 取消操作
  })
}

// 处理快捷操作
const handleQuickAction = (action) => {
  switch (action.action) {
    case 'publishNotice':
      ElMessage.info('发布公告功能开发中')
      break
    case 'addUser':
      router.push('/admin/users?action=add')
      break
    case 'manageProducts':
      router.push('/admin/products')
      break
    case 'systemSettings':
      ElMessage.info('系统设置功能开发中')
      break
    case 'viewLogs':
      ElMessage.info('日志查看功能开发中')
      break
    case 'refreshCache':
      ElMessage.success('缓存刷新成功')
      loadDashboardData()
      break
  }
}
</script>

<style scoped>
.admin-dashboard-container {
  padding: 0;
}

.dashboard-header {
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

.date-info {
  display: flex;
  align-items: center;
  color: #606266;
  font-size: 14px;
}

.date-info .el-icon {
  margin-right: 8px;
}

.chart-row {
  margin-top: 20px;
  margin-bottom: 20px;
}

.stat-card {
  height: 100%;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0,0,0,0.1);
}

.stat-card-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
}

.stat-info {
  flex: 1;
}

.stat-title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 8px;
}

.stat-trend {
  font-size: 12px;
  display: flex;
  align-items: center;
}

.stat-trend.up {
  color: #67C23A;
}

.stat-trend.down {
  color: #F56C6C;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.stat-icon .el-icon {
  font-size: 24px;
  color: white;
}

.stat-footer {
  padding: 10px 20px;
  background-color: #f5f7fa;
  color: #606266;
  font-size: 12px;
  border-top: 1px solid #ebeef5;
}

.chart-card {
  height: 100%;
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.chart-placeholder {
  height: 300px;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
}

.mock-chart {
  width: 100%;
  height: 250px;
  display: flex;
  align-items: flex-end;
  justify-content: space-around;
  padding: 0 20px;
}

.mock-bar-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 40px;
}

.mock-bar {
  width: 30px;
  border-radius: 4px 4px 0 0;
  transition: all 0.3s;
}

.mock-label {
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
}

.mock-pie-chart {
  position: relative;
  width: 180px;
  height: 180px;
  border-radius: 50%;
  margin-bottom: 20px;
}

.mock-pie-center {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background-color: white;
  z-index: 2;
}

.mock-pie-slice {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  transform-origin: center;
}

.mock-pie-legend {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 2px;
}

.legend-label {
  font-size: 12px;
  color: #606266;
}

.legend-value {
  margin-left: auto;
  font-size: 12px;
  font-weight: bold;
  color: #303133;
}

.table-card {
  border-radius: 8px;
}

.product-info-cell {
  display: flex;
  align-items: center;
}

.product-thumbnail {
  width: 40px;
  height: 40px;
  border-radius: 4px;
  margin-right: 10px;
}

.product-info-text {
  display: flex;
  flex-direction: column;
}

.product-title {
  font-weight: 500;
}

.product-category {
  font-size: 12px;
  color: #909399;
}

.price-tag {
  color: #F56C6C;
  font-weight: bold;
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

.activity-timeline {
  padding: 10px;
}

.quick-actions {
  padding: 10px;
}

.quick-action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 80px;
  background-color: #f5f7fa;
  border-radius: 8px;
  margin-bottom: 20px;
  cursor: pointer;
  transition: all 0.3s;
}

.quick-action-item:hover {
  background-color: #ecf5ff;
  color: #409EFF;
}

.quick-action-item .el-icon {
  font-size: 24px;
  margin-bottom: 8px;
}

@media (max-width: 768px) {
  .stat-card {
    margin-bottom: 20px;
  }
  
  .chart-row {
    margin-top: 0;
  }
  
  .chart-card {
    margin-bottom: 20px;
  }
}

.loading-placeholder {
  padding: 20px 0;
  text-align: center;
}

.no-data {
  color: #909399;
  font-size: 14px;
}
</style>