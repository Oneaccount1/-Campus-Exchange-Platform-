<template>
  <div class="products-container">
    <el-row :gutter="20">
      <el-col :span="24">
        <div class="search-bar">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索商品"
            class="search-input"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="filter-section">
      <el-col :span="24">
        <el-card shadow="never">
          <div class="filter-item">
            <span class="filter-label">分类：</span>
            <el-radio-group v-model="selectedCategory" @change="handleCategoryChange">
              <el-radio :label="0">全部</el-radio>
              <el-radio v-for="category in productStore.categories" :key="category.id" :label="category.id">
                {{ category.name }}
              </el-radio>
            </el-radio-group>
          </div>
          <div class="filter-item">
            <span class="filter-label">价格：</span>
            <el-select v-model="priceRange" placeholder="价格区间" @change="handlePriceChange">
              <el-option label="全部" value="all" />
              <el-option label="0-100元" value="0-100" />
              <el-option label="100-500元" value="100-500" />
              <el-option label="500-1000元" value="500-1000" />
              <el-option label="1000-5000元" value="1000-5000" />
              <el-option label="5000元以上" value="5000+" />
            </el-select>
          </div>
          <div class="filter-item">
            <span class="filter-label">排序：</span>
            <el-select v-model="sortBy" placeholder="排序方式" @change="handleSortChange">
              <el-option label="最新发布" value="newest" />
              <el-option label="价格从低到高" value="price-asc" />
              <el-option label="价格从高到低" value="price-desc" />
            </el-select>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="product-list">
      <template v-if="products.length > 0">
        <el-col :xs="24" :sm="12" :md="8" :lg="6" v-for="product in displayProducts" :key="product.id" class="product-item">
          <el-card shadow="hover" @click="viewProduct(product.id)">
            <img :src="product.images && product.images.length > 0 ? product.images[0].image_url : '/default-product.png'" class="product-image">
            <div class="product-info">
              <h3 class="product-title">{{ product.title }}</h3>
              <p class="product-price">¥{{ product.price }}</p>
              <p class="product-description">{{ product.description }}</p>
              <div class="product-meta">
                <span class="publish-time">{{ formatTime(product.created_at) }}</span>
                <el-tag size="small" :type="product.status === '售卖中' ? 'success' : 'info'">
                  {{ product.status }}
                </el-tag>
              </div>
            </div>
          </el-card>
        </el-col>
      </template>
      <el-col :span="24" v-else>
        <el-empty v-if="!loading" description="暂无商品" />
        <div v-else class="loading-container">
          <el-skeleton :rows="3" animated />
        </div>
      </el-col>
    </el-row>

    <el-pagination
      v-if="totalProducts > 0"
      layout="prev, pager, next"
      :total="totalProducts"
      :page-size="pageSize"
      :current-page="currentPage"
      @current-change="handlePageChange"
      class="pagination"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useProductStore } from '../stores'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import * as productApi from '../api/product'

const router = useRouter()
const route = useRoute()
const productStore = useProductStore()

// 分页相关
const currentPage = ref(1)
const pageSize = ref(12)
const totalProducts = ref(0)
const loading = ref(false)

// 筛选相关
const searchKeyword = ref('')
const selectedCategory = ref(parseInt(route.query.category) || 0)
const priceRange = ref('all')
const sortBy = ref('newest')

// 商品列表数据
const products = ref([])
const isSearchMode = ref(false)

// 获取商品列表
const fetchProducts = async () => {
  loading.value = true
  try {
    const response = await productApi.getProductList(currentPage.value, pageSize.value)
    products.value = response.data.products || []
    totalProducts.value = response.data.total || 0
    productStore.setProducts(products.value)
    isSearchMode.value = false
  } catch (error) {
    ElMessage.error('获取商品列表失败')
    console.error('获取商品列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索商品
const searchProducts = async () => {
  if (!searchKeyword.value.trim()) {
    fetchProducts()
    return
  }
  
  loading.value = true
  try {
    const response = await productApi.searchProducts(searchKeyword.value, currentPage.value, pageSize.value)
    products.value = response.data.products || []
    totalProducts.value = response.data.total || 0
    isSearchMode.value = true
  } catch (error) {
    ElMessage.error('搜索商品失败')
    console.error('搜索商品失败:', error)
  } finally {
    loading.value = false
  }
}

// 根据筛选和排序显示的商品
const displayProducts = computed(() => {
  let result = [...products.value]
  
  // 分类筛选
  if (selectedCategory.value !== 0) {
    result = result.filter(item => {
      // 根据后端分类数据结构调整这里的匹配逻辑
      return item.category === productStore.categories.find(c => c.id === selectedCategory.value)?.name
    })
  }
  
  // 价格筛选
  if (priceRange.value !== 'all') {
    if (priceRange.value === '5000+') {
      result = result.filter(item => item.price >= 5000)
    } else {
      const [min, max] = priceRange.value.split('-')
      result = result.filter(item => item.price >= Number(min) && item.price <= Number(max))
    }
  }
  
  // 排序
  if (sortBy.value === 'price-asc') {
    result.sort((a, b) => a.price - b.price)
  } else if (sortBy.value === 'price-desc') {
    result.sort((a, b) => b.price - a.price)
  } else {
    // 默认按最新发布排序
    result.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
  }
  
  return result
})

// 监听筛选条件变化，如果在搜索模式下，不自动刷新数据
watch([selectedCategory, priceRange, sortBy], () => {
  if (!isSearchMode.value) {
    fetchProducts()
  }
})

// 初始化
onMounted(() => {
  fetchProducts()
  
  // 从路由参数获取分类
  if (route.query.category) {
    selectedCategory.value = parseInt(route.query.category)
  }
  
  // 从路由参数获取搜索关键词
  if (route.query.keyword) {
    searchKeyword.value = route.query.keyword
    searchProducts()
  }
})

// 事件处理函数
const handleSearch = () => {
  currentPage.value = 1
  searchProducts()
}

const handleCategoryChange = () => {
  currentPage.value = 1
  if (isSearchMode.value) {
    // 如果在搜索模式下点击分类，则退出搜索模式，回到常规浏览
    isSearchMode.value = false
    fetchProducts()
  }
}

const handlePriceChange = () => {
  currentPage.value = 1
}

const handleSortChange = () => {
  currentPage.value = 1
}

const handlePageChange = (page) => {
  currentPage.value = page
  if (isSearchMode.value) {
    searchProducts()
  } else {
    fetchProducts()
  }
}

const viewProduct = (productId) => {
  router.push(`/product/${productId}`)
}

// 格式化时间
const formatTime = (timeString) => {
  const date = new Date(timeString)
  return `${date.getMonth() + 1}月${date.getDate()}日`
}
</script>

<style scoped>
.products-container {
  padding: 20px;
}

.search-bar {
  display: flex;
  margin-bottom: 20px;
}

.search-input {
  margin-right: 10px;
  width: 300px;
}

.filter-section {
  margin-bottom: 20px;
}

.filter-item {
  margin-bottom: 15px;
  display: flex;
  align-items: center;
}

.filter-label {
  margin-right: 10px;
  font-weight: bold;
  width: 60px;
}

.product-list {
  margin-bottom: 20px;
}

.product-item {
  margin-bottom: 20px;
}

.product-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.product-info {
  padding: 14px;
}

.product-title {
  margin: 0 0 10px;
  font-size: 16px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.product-price {
  color: #f56c6c;
  font-size: 20px;
  font-weight: bold;
  margin: 5px 0;
}

.product-description {
  color: #909399;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 10px;
}

.product-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #909399;
}

.pagination {
  margin-top: 20px;
  text-align: center;
}

.loading-container {
  padding: 20px;
}
</style>