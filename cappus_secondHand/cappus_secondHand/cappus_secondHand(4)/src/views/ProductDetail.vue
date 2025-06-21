<template>
  <div class="product-detail-container" v-if="product">
    <el-row :gutter="30">
      <el-col :md="12">
        <div class="product-image-container">
          <el-carousel v-if="product.images && product.images.length > 0" height="400px">
            <el-carousel-item v-for="(image, index) in product.images" :key="index">
              <el-image 
                :src="image.image_url" 
                fit="cover"
                class="product-image"
              >
                <template #error>
                  <div class="image-error">
                    <el-icon><Picture /></el-icon>
                    <span>图片加载失败</span>
                  </div>
                </template>
              </el-image>
            </el-carousel-item>
          </el-carousel>
          <el-image 
            v-else
            src="/default-product.png" 
            fit="cover"
            class="product-image"
          >
            <template #error>
              <div class="image-error">
                <el-icon><Picture /></el-icon>
                <span>图片加载失败</span>
              </div>
            </template>
          </el-image>
        </div>
      </el-col>
      <el-col :md="12">
        <div class="product-info">
          <h1 class="product-title">{{ product.title }}</h1>
          <div class="product-price">¥{{ product.price }}</div>
          
          <div class="product-meta">
            <div class="meta-item">
              <span class="meta-label">分类：</span>
              <el-tag size="small">{{ product.category || '未分类' }}</el-tag>
            </div>
            <div class="meta-item">
              <span class="meta-label">状态：</span>
              <el-tag :type="product.status === '售卖中' ? 'success' : 'info'" size="small">{{ product.status }}</el-tag>
            </div>
            <div class="meta-item">
              <span class="meta-label">成色：</span>
              <span>{{ product.condition || '未知' }}</span>
            </div>
            <div class="meta-item">
              <span class="meta-label">发布时间：</span>
              <span>{{ formatTime(product.created_at) }}</span>
            </div>
            <div class="meta-item">
              <span class="meta-label">卖家ID：</span>
              <span>{{ product.user_id }}</span>
            </div>
          </div>
          
          <div class="product-description">
            <h3>商品描述</h3>
            <p>{{ product.description }}</p>
          </div>
          
          <div class="action-buttons">
            <el-button type="primary" @click="contactSeller" :disabled="!userStore.isLoggedIn">
              <el-icon><ChatDotRound /></el-icon> 联系卖家
            </el-button>
            <el-button :type="isFavorited ? 'danger' : 'default'" @click="addToFavorite" :disabled="!userStore.isLoggedIn">
              <el-icon><Star /></el-icon> {{ isFavorited ? '取消收藏' : '收藏商品' }}
            </el-button>
            <el-button type="success" @click="createOrder" :disabled="!userStore.isLoggedIn || product.user_id === userStore.userInfo?.id || product.status !== '售卖中'">
              <el-icon><ShoppingCart /></el-icon> 立即购买
            </el-button>
          </div>
          
          <el-alert
            v-if="!userStore.isLoggedIn"
            title="请先登录"
            type="warning"
            show-icon
            :closable="false"
          >
            登录后才能联系卖家或收藏商品
          </el-alert>
        </div>
      </el-col>
    </el-row>
    
    <el-divider />
    
    <el-row>
      <el-col :span="24">
        <h2 class="section-title">相似商品推荐</h2>
        <div class="similar-products">
          <el-row :gutter="20">
            <el-col :xs="24" :sm="12" :md="8" :lg="6" v-for="item in similarProducts" :key="item.id">
              <el-card shadow="hover" @click="viewProduct(item.id)" class="similar-product-card">
                <img :src="item.images && item.images.length > 0 ? item.images[0].image_url : '/default-product.png'" class="similar-product-image">
                <div class="similar-product-info">
                  <h3>{{ item.title }}</h3>
                  <p class="similar-product-price">¥{{ item.price }}</p>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </div>
      </el-col>
    </el-row>
  </div>
  
  <div v-else class="loading-container">
    <el-skeleton :rows="10" animated />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProductStore, useUserStore } from '../stores'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Picture, ChatDotRound, Star, ShoppingCart } from '@element-plus/icons-vue'
import * as productApi from '../api/product'
import * as userApi from '../api/user'
import * as orderApi from '../api/order'

const route = useRoute()
const router = useRouter()
const productStore = useProductStore()
const userStore = useUserStore()

const product = ref(null)
const similarProducts = ref([])
const loading = ref(true)
const isFavorited = ref(false)

onMounted(async () => {
  const productId = route.params.id
  
  try {
    loading.value = true
    // 使用API获取商品详情
    const response = await productApi.getProductById(productId)
    product.value = response.data
    
    // 设置当前商品到store
    productStore.setCurrentProduct(product.value)
    
    // 获取相似商品
    fetchSimilarProducts()
    
    // 检查是否已收藏
    if (userStore.isLoggedIn) {
      checkIsFavorite(productId)
    }
  } catch (error) {
    ElMessage.error('获取商品详情失败')
    console.error('获取商品详情失败:', error)
  } finally {
    loading.value = false
  }
})

// 获取相似商品
const fetchSimilarProducts = async () => {
  if (!product.value) return
  
  try {
    // 这里可以根据分类或其他条件获取相似商品
    // 目前简单实现为获取商品列表的前4个
    const response = await productApi.getProductList(1, 4)
    similarProducts.value = (response.data.products || [])
      .filter(p => p.id !== product.value.id)
      .slice(0, 4)
  } catch (error) {
    console.error('获取相似商品失败:', error)
  }
}

// 添加检查收藏状态的方法
const checkIsFavorite = async (productId) => {
  try {
    const response = await userApi.checkFavorite(productId)
    isFavorited.value = response.data.is_favorite
  } catch (error) {
    console.error('检查收藏状态失败:', error)
  }
}

// 修改添加收藏的方法
const addToFavorite = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }
  
  try {
    if (isFavorited.value) {
      await userApi.removeFavorite(product.value.id)
      isFavorited.value = false
      ElMessage.success('已取消收藏')
    } else {
      await userApi.addFavorite(product.value.id)
      isFavorited.value = true
      ElMessage.success('收藏成功')
    }
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
    console.error('收藏操作失败:', error)
  }
}

// 格式化时间
const formatTime = (timeString) => {
  const date = new Date(timeString)
  return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日 ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
}

// 联系卖家
const contactSeller = () => {
  console.log('联系卖家按钮被点击')
  
  if (!userStore.isLoggedIn) {
    console.log('用户未登录，无法联系卖家')
    ElMessage.warning('请先登录')
    return
  }
  
  if (!userStore.userInfo) {
    console.log('用户信息不完整，尝试重新获取用户信息')
    ElMessage.warning('用户信息不完整，请尝试重新登录')
    return
  }
  
  console.log('用户已登录，用户信息:', userStore.userInfo)
  console.log('商品信息:', product.value)
  
  // 确保不是自己发布的商品
  if (product.value.user_id === userStore.userInfo.id) {
    console.log('无法联系自己')
    ElMessage.warning('不能联系自己')
    return
  }
  
  console.log('准备跳转到消息页面，参数:', {
    userId: product.value.user_id,
    productId: product.value.id
  })
  
  // 跳转到消息页面，携带seller_id和product_id参数
  router.push({
    path: '/messages',
    query: { 
      userId: product.value.user_id,
      productId: product.value.id
    }
  })
  
  console.log('路由跳转完成')
}

// 查看其他商品
const viewProduct = (productId) => {
  router.push(`/product/${productId}`)
}

// 创建订单
const createOrder = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }
  
  if (!product.value) {
    ElMessage.error('商品信息不完整')
    return
  }
  
  // 检查是否是自己的商品
  if (product.value.user_id === userStore.userInfo.id) {
    ElMessage.warning('不能购买自己的商品')
    return
  }
  
  try {
    // 创建订单数据
    const orderData = {
      buyer_id: userStore.userInfo.id,
      seller_id: product.value.user_id,
      product_id: product.value.id
    }
    
    // 确认创建订单
    ElMessageBox.confirm(
      `确定要购买商品 "${product.value.title}" 吗？价格: ¥${product.value.price}`,
      '确认购买',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    ).then(async () => {
      try {
        const response = await orderApi.createOrder(orderData)
        ElMessage.success('订单创建成功')
        // 可以跳转到订单详情页或用户中心的订单列表
        router.push('/user?tab=orders')
      } catch (error) {
        ElMessage.error(error.response?.data?.message || '创建订单失败')
      }
    }).catch(() => {
      // 用户取消操作
    })
  } catch (error) {
    ElMessage.error('创建订单失败')
    console.error('创建订单失败:', error)
  }
}
</script>

<style scoped>
.product-detail-container {
  padding: 20px;
}

.product-image-container {
  width: 100%;
  height: 400px;
  overflow: hidden;
  border-radius: 8px;
  margin-bottom: 20px;
}

.product-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #909399;
  background-color: #f5f7fa;
}

.product-info {
  padding: 20px;
}

.product-title {
  margin: 0 0 20px;
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.product-price {
  font-size: 32px;
  color: #f56c6c;
  font-weight: bold;
  margin-bottom: 20px;
}

.product-meta {
  background-color: #f8f8f8;
  padding: 15px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.meta-item {
  margin-bottom: 10px;
  display: flex;
  align-items: center;
}

.meta-label {
  font-weight: bold;
  color: #606266;
  margin-right: 10px;
  width: 80px;
}

.product-description {
  margin-bottom: 20px;
}

.product-description h3 {
  margin: 0 0 10px;
  font-size: 18px;
  color: #303133;
}

.product-description p {
  color: #606266;
  line-height: 1.6;
  white-space: pre-wrap;
}

.action-buttons {
  margin-bottom: 15px;
}

.section-title {
  font-size: 20px;
  margin: 20px 0;
  font-weight: bold;
  color: #303133;
}

.similar-products {
  margin-bottom: 30px;
}

.similar-product-card {
  cursor: pointer;
  transition: all 0.3s;
  margin-bottom: 20px;
}

.similar-product-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 15px rgba(0,0,0,0.1);
}

.similar-product-image {
  width: 100%;
  height: 150px;
  object-fit: cover;
}

.similar-product-info {
  padding: 10px;
}

.similar-product-info h3 {
  font-size: 16px;
  margin: 0 0 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.similar-product-price {
  color: #f56c6c;
  font-weight: bold;
  margin: 0;
}

.loading-container {
  padding: 40px;
}
</style>