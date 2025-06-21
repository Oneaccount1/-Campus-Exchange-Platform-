<template>
  <div>
    <div class="user-center-container" v-if="userStore.isLoggedIn">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-card class="user-info-card">
            <div class="user-avatar">
              <el-avatar :size="100" :src="userStore.userInfo?.avatar">
                {{ userStore.userInfo?.username?.charAt(0) }}
              </el-avatar>
            </div>
            <h2 class="username">{{ userStore.userInfo?.username }}</h2>
            <div class="user-stats">
              <div class="stat-item">
                <div class="stat-value">{{ userProductsTotal }}</div>
                <div class="stat-label">发布商品</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ favoritesTotal }}</div>
                <div class="stat-label">收藏商品</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ ordersTotal }}</div>
                <div class="stat-label">交易订单</div>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="18">
          <el-tabs v-model="activeTab" @tab-change="handleTabChange">
            <el-tab-pane label="个人资料" name="profile">
              <el-card>
                <el-form 
                  :model="profileForm" 
                  :rules="profileRules" 
                  ref="profileFormRef" 
                  label-width="100px"
                  v-loading="loading.profile"
                >
                  <el-form-item label="用户名" prop="username">
                    <el-input v-model="profileForm.username" placeholder="请输入用户名"></el-input>
                  </el-form-item>
                  <el-form-item label="邮箱" prop="email">
                    <el-input v-model="profileForm.email" placeholder="请输入邮箱"></el-input>
                  </el-form-item>
                  <el-form-item label="手机号码" prop="phone">
                    <el-input v-model="profileForm.phone" placeholder="请输入手机号码"></el-input>
                  </el-form-item>
                  <el-form-item label="头像">
                    <el-upload
                      class="avatar-uploader"
                      action="/api/v1/user/upload-avatar"
                      :headers="{ Authorization: `Bearer ${userStore.token}` }"
                      :show-file-list="false"
                      :on-success="handleAvatarSuccess"
                      :before-upload="beforeAvatarUpload"
                    >
                      <img v-if="profileForm.avatar" :src="profileForm.avatar" class="avatar" />
                      <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
                    </el-upload>
                  </el-form-item>
                  <el-form-item>
                    <el-button type="primary" @click="updateProfile">保存信息</el-button>
                    <el-button @click="showChangePasswordDialog">修改密码</el-button>
                  </el-form-item>
                </el-form>
              </el-card>
            </el-tab-pane>
            
            <el-tab-pane label="我的发布" name="products">
              <div v-loading="loading.products">
                <el-row :gutter="20" v-if="userProducts.length > 0">
                  <el-col :xs="24" :sm="12" :md="8" v-for="product in userProducts" :key="product.id">
                    <el-card shadow="hover" class="product-card">
                      <img :src="product.images && product.images.length > 0 ? product.images[0].image_url : '/default-product.png'" class="product-image">
                      <div class="product-info">
                        <h3>{{ product.title }}</h3>
                        <p class="price">¥{{ product.price }}</p>
                        <el-tag size="small" :type="getStatusType(product.status)">
                          {{ product.status }}
                        </el-tag>
                        <div class="product-actions">
                          <el-button type="primary" size="small" @click="editProduct(product.id)">编辑</el-button>
                          <el-button type="danger" size="small" @click="deleteProduct(product.id)">删除</el-button>
                        </div>
                      </div>
                    </el-card>
                  </el-col>
                </el-row>
                <el-empty v-else description="暂无发布的商品" />
                
                <el-pagination
                  v-if="userProductsTotal > 0"
                  layout="prev, pager, next"
                  :total="userProductsTotal"
                  :page-size="productPageSize"
                  :current-page="productPage"
                  @current-change="handleProductPageChange"
                  class="pagination"
                />
              </div>
            </el-tab-pane>
            
            <el-tab-pane label="我的收藏" name="favorites">
              <div v-loading="loading.favorites">
                <el-row :gutter="20" v-if="favorites.length > 0">
                  <el-col :xs="24" :sm="12" :md="8" v-for="favorite in favorites" :key="favorite.id">
                    <el-card shadow="hover" class="product-card">
                      <img :src="favorite.product.images && favorite.product.images.length > 0 ? favorite.product.images[0].image_url : '/default-product.png'" class="product-image">
                      <div class="product-info">
                        <h3>{{ favorite.product.title }}</h3>
                        <p class="price">¥{{ favorite.product.price }}</p>
                        <div class="product-actions">
                          <el-button type="primary" size="small" @click="viewProduct(favorite.product_id)">查看</el-button>
                          <el-button type="danger" size="small" @click="removeFavorite(favorite.product_id)">取消收藏</el-button>
                        </div>
                      </div>
                    </el-card>
                  </el-col>
                </el-row>
                <el-empty v-else description="暂无收藏的商品" />
                
                <el-pagination
                  v-if="favoritesTotal > 0"
                  layout="prev, pager, next"
                  :total="favoritesTotal"
                  :page-size="favoritePageSize"
                  :current-page="favoritePage"
                  @current-change="handleFavoritePageChange"
                  class="pagination"
                />
              </div>
            </el-tab-pane>
            
            <el-tab-pane label="我的订单" name="orders">
              <div v-loading="loading.orders">
                <el-tabs v-model="orderStatusTab" @tab-click="handleOrderStatusChange">
                  <el-tab-pane label="全部" name="all"></el-tab-pane>
                  <el-tab-pane label="待付款" name="pending"></el-tab-pane>
                  <el-tab-pane label="待发货" name="paid"></el-tab-pane>
                  <el-tab-pane label="待收货" name="shipped"></el-tab-pane>
                  <el-tab-pane label="已完成" name="completed"></el-tab-pane>
                  <el-tab-pane label="已取消" name="cancelled"></el-tab-pane>
                </el-tabs>
                
                <div v-if="orders.length > 0">
                  <el-card v-for="order in orders" :key="order.id" class="order-card">
                    <div class="order-header">
                      <div>
                        <span class="order-number">订单号: {{ order.id }}</span>
                        <span class="order-date">{{ formatDate(order.created_at) }}</span>
                      </div>
                      <el-tag :type="getOrderStatusType(order.status)">{{ getOrderStatusText(order.status) }}</el-tag>
                    </div>
                    
                    <div class="order-product">
                      <div class="order-product-info">
                        <h3>商品ID: {{ order.product_id }}</h3>
                        <p class="price">¥{{ order.price }}</p>
                        <p>买家ID: {{ order.buyer_id }} | 卖家ID: {{ order.seller_id }}</p>
                      </div>
                      <div class="order-actions">
                        <el-button 
                          v-if="order.status === '卖家未处理' || order.status === 'pending'" 
                          type="primary" 
                          size="small" 
                          @click="payOrder(order.id)"
                        >
                          立即付款
                        </el-button>
                        <el-button 
                          v-if="order.status === '已发货' || order.status === 'shipped'" 
                          type="success" 
                          size="small" 
                          @click="confirmReceive(order.id)"
                        >
                          确认收货
                        </el-button>
                        <el-button 
                          v-if="order.status === '卖家未处理' || order.status === 'pending'" 
                          type="danger" 
                          size="small" 
                          @click="cancelOrder(order.id)"
                        >
                          取消订单
                        </el-button>
                        <el-button 
                          v-if="order.status === '已完成' || order.status === 'completed'" 
                          type="warning" 
                          size="small" 
                          @click="reviewOrder(order.id)"
                        >
                          评价订单
                        </el-button>
                        <el-button 
                          size="small" 
                          @click="viewOrderDetail(order.id)"
                        >
                          订单详情
                        </el-button>
                      </div>
                    </div>
                  </el-card>
                </div>
                <el-empty v-else description="暂无相关订单" />
                
                <el-pagination
                  v-if="ordersTotal > 0"
                  layout="prev, pager, next"
                  :total="ordersTotal"
                  :page-size="orderPageSize"
                  :current-page="orderPage"
                  @current-change="handleOrderPageChange"
                  class="pagination"
                />
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-col>
      </el-row>
    </div>
    
    <div v-else class="not-logged-in">
      <el-result
        icon="warning"
        title="请先登录"
        sub-title="登录后才能访问个人中心"
      >
        <template #extra>
          <el-button type="primary" @click="$router.push('/login')">去登录</el-button>
        </template>
      </el-result>
    </div>
    
    <!-- 修改密码对话框 -->
    <el-dialog
      title="修改密码"
      v-model="passwordDialogVisible"
      width="400px"
    >
      <el-form 
        :model="passwordForm" 
        :rules="passwordRules" 
        ref="passwordFormRef" 
        label-width="100px"
      >
        <el-form-item label="原密码" prop="oldPassword">
          <el-input v-model="passwordForm.oldPassword" type="password" placeholder="请输入原密码" show-password></el-input>
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="passwordForm.newPassword" type="password" placeholder="请输入新密码" show-password></el-input>
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" placeholder="请再次输入新密码" show-password></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="passwordDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="changePassword">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import * as userApi from '../api/user'
import * as productApi from '../api/product'
import * as orderApi from '../api/order'

const router = useRouter()
const userStore = useUserStore()

const activeTab = ref('profile')
const orderStatusTab = ref('all')

// 个人资料表单
const profileFormRef = ref(null)
const profileForm = reactive({
  username: '',
  email: '',
  phone: '',
  avatar: ''
})

// 表单验证规则
const profileRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在2到20个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ]
}

// 密码修改表单
const passwordDialogVisible = ref(false)
const passwordFormRef = ref(null)
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 密码表单验证规则
const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入原密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { 
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      }, 
      trigger: 'blur' 
    }
  ]
}

// 用户发布的商品数据
const userProducts = ref([])
const userProductsTotal = ref(0)
const productPage = ref(1)
const productPageSize = ref(9)

// 用户收藏的商品数据
const favorites = ref([])
const favoritesTotal = ref(0)
const favoritePage = ref(1)
const favoritePageSize = ref(9)

// 用户订单数据
const orders = ref([])
const ordersTotal = ref(0)
const orderPage = ref(1)
const orderPageSize = ref(5)

// 加载状态
const loading = ref({
  profile: false,
  products: false,
  favorites: false,
  orders: false
})

// 初始化
onMounted(() => {
  if (userStore.isLoggedIn) {
    // 加载基本资料
    fetchUserProfile()
    
    // 直接加载所有统计数据，无需等待切换标签
    fetchUserProducts()
    fetchUserFavorites()
    fetchUserOrders()
  }
})

// 获取用户个人资料
const fetchUserProfile = async () => {
  if (!userStore.isLoggedIn) return
  
  loading.value.profile = true
  try {
    const response = await userApi.getUserInfo()
    // 填充表单
    profileForm.username = response.data.username
    profileForm.email = response.data.email
    profileForm.phone = response.data.phone
    profileForm.avatar = response.data.avatar
  } catch (error) {
    ElMessage.error('获取个人资料失败')
    console.error('获取个人资料失败:', error)
  } finally {
    loading.value.profile = false
  }
}

// 更新个人资料
const updateProfile = async () => {
  if (!profileFormRef.value) return
  
  await profileFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value.profile = true
    try {
      await userApi.updateUserInfo(profileForm)
      ElMessage.success('个人资料更新成功')
      // 更新本地存储的用户信息
      userStore.setUserInfo({
        ...userStore.userInfo,  // 保留原有信息
        username: profileForm.username,
        email: profileForm.email,
        phone: profileForm.phone,
        avatar: profileForm.avatar
      })
    } catch (error) {
      ElMessage.error('更新失败: ' + (error.response?.data?.message || '未知错误'))
    } finally {
      loading.value.profile = false
    }
  })
}

// 显示修改密码对话框
const showChangePasswordDialog = () => {
  passwordDialogVisible.value = true
  // 重置表单
  if (passwordFormRef.value) {
    passwordFormRef.value.resetFields()
  }
}

// 修改密码
const changePassword = async () => {
  if (!passwordFormRef.value) return
  
  await passwordFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    try {
      await userApi.changePassword({
        old_password: passwordForm.oldPassword,
        new_password: passwordForm.newPassword
      })
      ElMessage.success('密码修改成功')
      passwordDialogVisible.value = false
    } catch (error) {
      ElMessage.error('密码修改失败: ' + (error.response?.data?.message || '未知错误'))
    }
  })
}

// 头像上传前的处理
const beforeAvatarUpload = (file) => {
  const isJPG = file.type === 'image/jpeg'
  const isPNG = file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG && !isPNG) {
    ElMessage.error('上传头像图片只能是 JPG 或 PNG 格式!')
  }
  if (!isLt2M) {
    ElMessage.error('上传头像图片大小不能超过 2MB!')
  }
  return (isJPG || isPNG) && isLt2M
}

// 头像上传成功的处理
const handleAvatarSuccess = (response) => {
  profileForm.avatar = response.data.url
}

// 获取用户发布的商品
const fetchUserProducts = async () => {
  if (!userStore.isLoggedIn || !userStore.userInfo) return
  
  loading.value.products = true
  try {
    const response = await userApi.getUserProducts(userStore.userInfo.id, productPage.value, productPageSize.value)
    userProducts.value = response.data.products || []
    userProductsTotal.value = response.data.total || 0
  } catch (error) {
    ElMessage.error('获取发布商品失败')
    console.error('获取发布商品失败:', error)
  } finally {
    loading.value.products = false
  }
}

// 获取用户收藏的商品
const fetchUserFavorites = async () => {
  if (!userStore.isLoggedIn) return
  
  loading.value.favorites = true
  try {
    const response = await userApi.getUserFavorites(favoritePage.value, favoritePageSize.value)
    favorites.value = response.data.favorites || []
    favoritesTotal.value = response.data.total || 0
  } catch (error) {
    ElMessage.error('获取收藏商品失败')
    console.error('获取收藏商品失败:', error)
  } finally {
    loading.value.favorites = false
  }
}

// 获取用户订单
const fetchUserOrders = async () => {
  if (!userStore.isLoggedIn || !userStore.userInfo) return
  
  loading.value.orders = true
  try {
    const status = orderStatusTab.value !== 'all' ? orderStatusTab.value : ''
    const response = await orderApi.getUserOrders(userStore.userInfo.id, orderPage.value, orderPageSize.value, status)
    orders.value = response.data.orders || []
    ordersTotal.value = response.data.total || 0
  } catch (error) {
    ElMessage.error('获取订单失败')
    console.error('获取订单失败:', error)
  } finally {
    loading.value.orders = false
  }
}

// 处理标签切换
const handleTabChange = (tabName) => {
  if (tabName === 'profile' && !profileForm.username) {
    fetchUserProfile()
  } else if (tabName === 'products' && userProducts.value.length === 0) {
    fetchUserProducts()
  } else if (tabName === 'favorites' && favorites.value.length === 0) {
    fetchUserFavorites()
  } else if (tabName === 'orders' && orders.value.length === 0) {
    fetchUserOrders()
  }
}

// 处理订单状态切换
const handleOrderStatusChange = () => {
  orderPage.value = 1
  fetchUserOrders()
}

// 处理商品分页变化
const handleProductPageChange = (page) => {
  productPage.value = page
  fetchUserProducts()
}

// 处理收藏分页变化
const handleFavoritePageChange = (page) => {
  favoritePage.value = page
  fetchUserFavorites()
}

// 处理订单分页变化
const handleOrderPageChange = (page) => {
  orderPage.value = page
  fetchUserOrders()
}

// 编辑商品
const editProduct = (productId) => {
  router.push(`/publish?id=${productId}`)
}

// 删除商品
const deleteProduct = (productId) => {
  ElMessageBox.confirm(
    '确定要删除这个商品吗？',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await productApi.deleteProduct(productId)
      // 重新加载商品列表
      fetchUserProducts()
      ElMessage.success('删除成功')
    } catch (error) {
      ElMessage.error('删除失败: ' + (error.response?.data?.message || '未知错误'))
    }
  }).catch(() => {})
}

// 查看商品
const viewProduct = (productId) => {
  router.value.push(`/product/${productId}`)
}

// 取消收藏
const removeFavorite = (productId) => {
  ElMessageBox.confirm(
    '确定要取消收藏这个商品吗？',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await userApi.removeFavorite(productId)
      // 重新加载收藏列表
      fetchUserFavorites()
      ElMessage.success('已取消收藏')
    } catch (error) {
      ElMessage.error('取消收藏失败: ' + (error.response?.data?.message || '未知错误'))
    }
  }).catch(() => {})
}

// 支付订单
const payOrder = (orderId) => {
  ElMessageBox.confirm(
    '确认支付此订单吗？',
    '提示',
    {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await orderApi.payOrder(orderId)
      // 重新加载订单列表
      fetchUserOrders()
      ElMessage.success('支付成功')
    } catch (error) {
      ElMessage.error('支付失败: ' + (error.response?.data?.message || '未知错误'))
    }
  }).catch(() => {})
}

// 确认收货
const confirmReceive = (orderId) => {
  ElMessageBox.confirm(
    '确认已收到商品吗？',
    '提示',
    {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await orderApi.confirmReceive(orderId)
      // 重新加载订单列表
      fetchUserOrders()
      ElMessage.success('确认收货成功')
    } catch (error) {
      ElMessage.error('操作失败: ' + (error.response?.data?.message || '未知错误'))
    }
  }).catch(() => {})
}

// 取消订单
const cancelOrder = (orderId) => {
  ElMessageBox.confirm(
    '确定要取消这个订单吗？',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await orderApi.cancelOrder(orderId)
      // 重新加载订单列表
      fetchUserOrders()
      ElMessage.success('订单已取消')
    } catch (error) {
      ElMessage.error('取消失败: ' + (error.response?.data?.message || '未知错误'))
    }
  }).catch(() => {})
}

// 订单评价
const reviewOrder = (orderId) => {
  ElMessageBox.prompt('请输入您对此订单的评价', '订单评价', {
    confirmButtonText: '提交',
    cancelButtonText: '取消',
    inputType: 'textarea',
    inputPlaceholder: '请输入评价内容'
  }).then(({ value }) => {
    if (value.trim()) {
      orderApi.reviewOrder(orderId, { content: value.trim() })
        .then(() => {
          ElMessage.success('评价成功')
          fetchUserOrders()
        })
        .catch(error => {
          ElMessage.error('评价失败: ' + (error.response?.data?.message || '未知错误'))
        })
    } else {
      ElMessage.warning('评价内容不能为空')
    }
  }).catch(() => {})
}

// 查看订单详情
const viewOrderDetail = async (orderId) => {
  try {
    const response = await orderApi.getOrderDetail(orderId)
    ElMessageBox.alert(
      `
        <div>
          <p><strong>订单ID:</strong> ${response.data.id}</p>
          <p><strong>买家ID:</strong> ${response.data.buyer_id}</p>
          <p><strong>卖家ID:</strong> ${response.data.seller_id}</p>
          <p><strong>商品ID:</strong> ${response.data.product_id}</p>
          <p><strong>价格:</strong> ¥${response.data.price}</p>
          <p><strong>状态:</strong> ${getOrderStatusText(response.data.status)}</p>
          <p><strong>创建时间:</strong> ${formatDate(response.data.created_at)}</p>
          <p><strong>更新时间:</strong> ${formatDate(response.data.updated_at)}</p>
        </div>
      `,
      '订单详情',
      {
        dangerouslyUseHTMLString: true,
        confirmButtonText: '关闭'
      }
    )
  } catch (error) {
    ElMessage.error('获取订单详情失败: ' + (error.response?.data?.message || '未知错误'))
  }
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

// 获取商品状态类型
const getStatusType = (status) => {
  const statusMap = {
    '售卖中': 'success',
    '已售出': 'info',
    '已下架': 'danger'
  }
  return statusMap[status] || 'info'
}

// 获取订单状态类型
const getOrderStatusType = (status) => {
  const statusMap = {
    'pending': 'warning',
    'paid': 'primary',
    'shipped': 'success',
    'completed': 'success',
    'cancelled': 'info',
    '卖家未处理': 'warning',
    '卖家已同意': 'primary',
    '卖家已拒绝': 'danger',
    '已发货': 'success',
    '已完成': 'success',
    '已取消': 'info'
  }
  return statusMap[status] || 'info'
}

// 获取订单状态文本
const getOrderStatusText = (status) => {
  // 如果状态已经是中文，直接返回
  if (/[\u4e00-\u9fa5]/.test(status)) {
    return status
  }
  
  const statusMap = {
    'pending': '待付款',
    'paid': '待发货',
    'shipped': '待收货',
    'completed': '已完成',
    'cancelled': '已取消'
  }
  return statusMap[status] || '未知状态'
}

// 监听用户信息变化
watch(() => userStore.isLoggedIn, (newVal) => {
  if (newVal) {
    fetchUserProfile()
  }
})
</script>

<style scoped>
.user-center-container {
  padding: 20px;
}

.user-info-card {
  text-align: center;
  padding: 20px;
}

.user-avatar {
  margin-bottom: 15px;
}

.username {
  font-size: 1.2rem;
  margin-bottom: 15px;
}

.user-stats {
  display: flex;
  justify-content: space-around;
  margin-top: 20px;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: bold;
  color: #409eff;
}

.stat-label {
  font-size: 0.9rem;
  color: #606266;
}

.product-card {
  margin-bottom: 20px;
  overflow: hidden;
}

.product-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.product-info {
  padding: 10px;
}

.product-info h3 {
  margin-bottom: 10px;
  font-size: 1rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.price {
  color: #f56c6c;
  font-weight: bold;
  margin-bottom: 10px;
}

.product-actions {
  margin-top: 10px;
  display: flex;
  justify-content: space-between;
}

.pagination {
  margin-top: 20px;
  text-align: center;
}

.avatar-uploader {
  display: inline-block;
}

.avatar-uploader .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}

.avatar-uploader .el-upload:hover {
  border-color: #409eff;
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 100px;
  height: 100px;
  line-height: 100px;
  text-align: center;
}

.avatar {
  width: 100px;
  height: 100px;
  display: block;
  object-fit: cover;
}

.order-card {
  margin-bottom: 15px;
}

.order-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #ebeef5;
}

.order-number {
  font-weight: bold;
  margin-right: 20px;
}

.order-date {
  color: #909399;
}

.order-product {
  display: flex;
  align-items: center;
}

.order-product-image {
  width: 80px;
  height: 80px;
  object-fit: cover;
  margin-right: 15px;
}

.order-product-info {
  flex: 1;
}

.order-product-info h3 {
  margin-top: 0;
  margin-bottom: 5px;
}

.order-actions {
  margin-left: 20px;
}

.not-logged-in {
  margin-top: 40px;
}
</style>