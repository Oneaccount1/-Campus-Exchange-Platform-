<template>
  <div class="home-container">
    <el-row :gutter="20">
      <el-col :span="24">
        <el-carousel height="400px" class="banner">
          <el-carousel-item v-for="item in bannerItems" :key="item.id">
            <div class="banner-content" :style="{ backgroundImage: `url(${item.image})` }">
              <h2>{{ item.title }}</h2>
              <p>{{ item.description }}</p>
            </div>
          </el-carousel-item>
        </el-carousel>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="category-section">
      <el-col :span="24">
        <h2 class="section-title">商品分类</h2>
      </el-col>
      <el-col :span="4" v-for="category in productStore.categories" :key="category.id">
        <el-card class="category-card" shadow="hover" @click="handleCategoryClick(category)">
          <div class="category-name">{{ category.name }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="latest-products">
      <el-col :span="24">
        <div class="section-header">
          <h2 class="section-title">最新上架</h2>
          <el-link type="primary" @click="viewMoreProducts">查看更多</el-link>
        </div>
      </el-col>
      
      <!-- 空数据状态 -->
      <el-col :span="24" v-if="!loading && latestProducts.length === 0">
        <el-empty 
          description="暂无商品数据" 
          :image-size="200"
        >
          <template #extra>
            <el-button type="primary" @click="goToPublish">立即发布</el-button>
          </template>
        </el-empty>
      </el-col>
      
      <el-col :span="6" v-for="product in latestProducts" :key="product.id" v-loading="loading">
        <el-card class="product-card" shadow="hover" @click="viewProduct(product.id)">
          <div class="product-image-container">
            <el-image 
              :src="product.images && product.images.length > 0 ? product.images[0].image_url : 'https://via.placeholder.com/300x200'" 
              class="product-image"
              fit="cover"
            >
              <template #error>
                <div class="image-error">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
          </div>
          <div class="product-info">
            <h3>{{ product.title }}</h3>
            <p class="price">¥{{ product.price }}</p>
            <p class="description">{{ product.description }}</p>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useProductStore } from '../stores'
import { getLatestProducts } from '../api/product'
import { ElMessage } from 'element-plus'
import { Picture } from '@element-plus/icons-vue'

const router = useRouter()
const productStore = useProductStore()
const loading = ref(false)

const bannerItems = ref([
  {
    id: 1,
    title: '校园二手交易平台',
    description: '让闲置物品流转起来',
    image: 'https://img.freepik.com/free-photo/students-campus_23-2147778155.jpg'
  },
  {
    id: 2,
    title: '安全可靠的交易环境',
    description: '校园认证，安全保障',
    image: 'https://img.freepik.com/free-photo/group-students-university_23-2148888772.jpg'
  }
])

const latestProducts = ref([])

// 获取最新商品
const fetchLatestProducts = async () => {
  loading.value = true
  try {
    const response = await getLatestProducts(8) // 获取8条最新商品
    console.log('最新商品响应:', response)
    
    if (response && response.code === 200) {
      if (response.data && Array.isArray(response.data.products) && response.data.products.length > 0) {
        // 如果是数组格式的商品列表
        latestProducts.value = response.data.products
      } else if (response.data && !Array.isArray(response.data) && response.data.id) {
        // 如果是单个商品对象（有时API可能直接返回一个对象而非数组）
        latestProducts.value = [response.data]
      } else if (response.data && Array.isArray(response.data)) {
        // 如果直接返回的是数组
        latestProducts.value = response.data
      } else {
        console.warn('未找到商品数据或格式不正确:', response.data)
        // 清空数据，避免显示旧数据
        latestProducts.value = []
      }
    } else {
      console.warn('获取最新商品失败:', response)
      ElMessage.error('获取最新商品数据失败')
      latestProducts.value = []
    }
  } catch (error) {
    console.error('获取最新商品错误:', error)
    ElMessage.error('获取最新商品数据出错')
    latestProducts.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchLatestProducts()
})

const handleCategoryClick = (category) => {
  router.push({
    path: '/products',
    query: { category: category.id }
  })
}

const viewProduct = (productId) => {
  router.push(`/product/${productId}`)
}

const viewMoreProducts = () => {
  router.push('/products')
}

const goToPublish = () => {
  router.push('/publish')
}
</script>

<style scoped>
.home-container {
  padding: 20px;
}

.banner {
  margin-bottom: 40px;
  border-radius: 8px;
  overflow: hidden;
}

.banner-content {
  height: 100%;
  background-size: cover;
  background-position: center;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: white;
  text-align: center;
  background-color: rgba(0, 0, 0, 0.4);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-title {
  margin-bottom: 20px;
  font-size: 24px;
  color: #303133;
}

.category-section {
  margin-bottom: 40px;
}

.category-card {
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: transform 0.3s;
}

.category-card:hover {
  transform: translateY(-5px);
}

.category-name {
  font-size: 16px;
  color: #606266;
}

.product-card {
  margin-bottom: 20px;
  cursor: pointer;
  transition: transform 0.3s;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.product-card:hover {
  transform: translateY(-5px);
}

.product-image-container {
  width: 100%;
  height: 200px;
  overflow: hidden;
}

.product-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.product-card:hover .product-image {
  transform: scale(1.05);
}

.image-error {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
  color: #909399;
}

.product-info {
  padding: 14px;
  flex-grow: 1;
  display: flex;
  flex-direction: column;
}

.product-info h3 {
  margin: 0 0 10px;
  font-size: 16px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.price {
  color: #f56c6c;
  font-size: 20px;
  font-weight: bold;
  margin: 5px 0;
}

.description {
  color: #909399;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-top: auto;
}
</style>