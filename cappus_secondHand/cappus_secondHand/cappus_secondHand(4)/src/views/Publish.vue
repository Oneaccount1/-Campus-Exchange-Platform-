<template>
  <div class="publish-container">
    <el-card class="publish-card">
      <template #header>
        <div class="card-header">
          <h2>{{ isEdit ? '编辑商品' : '发布商品' }}</h2>
        </div>
      </template>
      
      <el-form :model="productForm" :rules="rules" ref="productFormRef" label-position="top">
        <el-form-item label="商品标题" prop="title">
          <el-input v-model="productForm.title" placeholder="请输入商品标题" maxlength="50" show-word-limit />
        </el-form-item>
        
        <el-form-item label="商品分类" prop="category">
          <el-select v-model="productForm.category" placeholder="请选择商品分类" style="width: 100%">
            <el-option 
              v-for="category in productStore.categories" 
              :key="category.id" 
              :label="category.name" 
              :value="category.id" 
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="商品价格" prop="price">
          <el-input-number v-model="productForm.price" :min="0" :precision="2" :step="10" style="width: 100%" />
        </el-form-item>
        
        <el-form-item label="商品图片" prop="image">
          <el-upload
            class="product-image-uploader"
            action="#"
            :auto-upload="false"
            :limit="1"
            :on-change="handleImageChange"
            :on-exceed="handleExceed"
            :file-list="fileList"
          >
            <el-button type="primary">
              <el-icon><Upload /></el-icon> 选择图片
            </el-button>
            <template #tip>
              <div class="el-upload__tip">请上传商品图片，支持JPG/PNG格式，大小不超过5MB</div>
            </template>
          </el-upload>
          
          <div class="image-preview" v-if="productForm.image">
            <img :src="productForm.image" class="preview-image">
            <div class="preview-actions">
              <el-button type="danger" size="small" @click="removeImage">
                <el-icon><Delete /></el-icon> 删除图片
              </el-button>
            </div>
          </div>
        </el-form-item>
        
        <el-form-item label="商品描述" prop="description">
          <el-input 
            v-model="productForm.description" 
            type="textarea" 
            :rows="5" 
            placeholder="请详细描述商品的成色、使用年限、出售原因等信息"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="商品成色" prop="condition">
          <el-select v-model="productForm.condition" placeholder="请选择商品成色">
            <el-option value="全新" label="全新" />
            <el-option value="九成新" label="九成新" />
            <el-option value="八成新" label="八成新" />
            <el-option value="七成新" label="七成新" />
            <el-option value="六成新" label="六成新" />
            <el-option value="五成新" label="五成新" />
            <el-option value="四成新" label="四成新" />
            <el-option value="三成新" label="三成新" />
            <el-option value="二成新" label="二成新" />
            <el-option value="一成新" label="一成新" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="商品状态" prop="status">
          <el-select v-model="productForm.status" placeholder="请选择商品状态">
            <el-option value="售卖中" label="售卖中" />
            <!-- 后端只支持"售卖中"和"已下架"两种状态 -->
            <el-option value="已下架" label="已下架" />
          </el-select>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" class="submit-button" @click="handleSubmit" :loading="loading">
            {{ isEdit ? '保存修改' : '发布商品' }}
          </el-button>
          <el-button @click="$router.go(-1)">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useProductStore, useUserStore } from '../stores'
import { ElMessage } from 'element-plus'
import { productApi } from '../api'
import request from '../api/request'
import { getUserInfo } from '../api/user'
import { uploadProductImage } from '../api/product'
import { Plus, Upload, Delete } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const productStore = useProductStore()
const userStore = useUserStore()

const productFormRef = ref(null)
const loading = ref(false)
const fileList = ref([])
const uploadLoading = ref(false)

// 判断是新增还是编辑
const isEdit = computed(() => {
  return route.query.id !== undefined
})

// 商品表单数据
const productForm = reactive({
  title: '',
  category: '',
  price: 0,
  images: [],
  image: '', // 保留用于前端显示
  description: '',
  condition: '全新',
  status: '售卖中'
})

// 表单验证规则
const rules = {
  title: [
    { required: true, message: '请输入商品标题', trigger: 'blur' },
    { min: 2, max: 50, message: '标题长度应为2-50个字符', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请选择商品分类', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入商品价格', trigger: 'blur' },
    { type: 'number', min: 0, message: '价格不能小于0', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入商品描述', trigger: 'blur' },
    { min: 10, max: 500, message: '描述长度应为10-500个字符', trigger: 'blur' }
  ],
  condition: [
    { required: true, message: '请选择商品成色', trigger: 'change' }
  ]
}

// 获取商品详情
const fetchProductDetail = async (productId) => {
  loading.value = true
  try {
    const response = await productApi.getProductById(productId)
    const productData = response.data
    
    // 填充表单
    Object.assign(productForm, {
      title: productData.title,
      category: productData.category,
      price: productData.price,
      description: productData.description,
      condition: productData.condition || '全新',
      status: productData.status || '售卖中'
    })
    
    // 处理图片
    if (productData.images && productData.images.length > 0) {
      productForm.images = productData.images
      productForm.image = productData.images[0] // 显示第一张图片
      
      // 添加到文件列表
      fileList.value = productData.images.map((url, index) => ({
        name: `商品图片${index + 1}`,
        url
      }))
    }
  } catch (error) {
    ElMessage.error('获取商品详情失败')
    console.error('获取商品详情失败:', error)
    router.push('/products')
  } finally {
    loading.value = false
  }
}

  // 初始化
onMounted(async () => {
  // 检查用户是否登录
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  // 确保有用户信息
  if (!userStore.userInfo) {
    try {
      // 尝试获取用户信息
      const response = await getUserInfo()
      if (response.code === 200) {
        userStore.setUserInfo(response.data)
      } else {
        ElMessage.warning('无法获取用户信息，请重新登录')
        router.push('/login')
        return
      }
    } catch (error) {
      console.error('获取用户信息失败:', error)
      ElMessage.warning('获取用户信息失败，请重新登录')
      router.push('/login')
      return
    }
  }
  
  // 如果是编辑模式，获取商品信息
  if (isEdit.value) {
    const productId = parseInt(route.query.id)
    fetchProductDetail(productId)
  }
})

// 上传图片到服务器
const uploadImage = async (file) => {
  uploadLoading.value = true
  try {
    // 创建FormData对象
    const formData = new FormData()
    formData.append('file', file)
    
    // 调用上传API - 使用专用的商品图片上传接口
    const response = await uploadProductImage(formData)
    
    console.log('图片上传响应:', response)
    if (response.code === 200 && response.data && response.data.url) {
      return response.data.url // 返回图片URL
    } else {
      ElMessage.error(response.message || '图片上传失败')
      return null
    }
  } catch (error) {
    ElMessage.error('图片上传失败')
    console.error('图片上传失败:', error)
    return null
  } finally {
    uploadLoading.value = false
  }
}

// 处理图片上传
const handleImageChange = async (file) => {
  // 检查文件类型
  const isImage = file.raw.type === 'image/jpeg' || file.raw.type === 'image/png'
  if (!isImage) {
    ElMessage.error('只能上传JPG或PNG格式的图片')
    return false
  }
  
  // 检查文件大小
  const isLt5M = file.raw.size / 1024 / 1024 < 5
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过5MB')
    return false
  }
  
  // 临时显示本地预览
  productForm.image = URL.createObjectURL(file.raw)
  
  // 在实际提交表单时再上传图片到服务器
}

// 处理超出上传数量限制
const handleExceed = () => {
  ElMessage.warning('最多只能上传1张图片')
}

// 移除图片
const removeImage = () => {
  productForm.image = ''
  fileList.value = []
}

// 提交表单
const handleSubmit = () => {
  productFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    if (!productForm.image) {
      ElMessage.warning('请上传商品图片')
      return
    }
    
    // 检查用户是否登录 - 只检查登录状态，不检查userInfo
    if (!userStore.isLoggedIn) {
      ElMessage.warning('请先登录')
      router.push('/login')
      return
    }
    
    loading.value = true
    
    try {
      // 如果是本地文件URL，需要先上传图片
      let imageUrl = productForm.image
      if (fileList.value.length > 0 && fileList.value[0].raw) {
        imageUrl = await uploadImage(fileList.value[0].raw)
        if (!imageUrl) {
          loading.value = false
          return
        }
      }
      
      // 获取分类名称
      const categoryName = productStore.categories.find(c => c.id === productForm.category)?.name || '其他'
      
      // 准备提交的数据
      const submitData = {
        title: productForm.title.trim(),  // 确保去掉首尾空格
        category: categoryName, // 使用分类名称而不是ID
        price: productForm.price,
        images: [imageUrl], // 将图片URL放入images数组
        description: productForm.description || '无描述',  // 确保描述字段不为空
        condition: productForm.condition || '全新',  // 确保成色字段不为空
        user_id: userStore.userInfo?.id,  // 使用当前登录用户的ID
        status: productForm.status === '已售出' ? '已下架' : productForm.status  // 确保状态符合后端要求
      }
      
      // 验证关键字段
      if (!submitData.title) {
        ElMessage.error('商品标题不能为空')
        loading.value = false
        return
      }
      
      if (!submitData.user_id) {
        ElMessage.error('用户未登录或用户ID无效')
        loading.value = false
        return
      }
      
      // 确保价格是数值类型且大于等于0
      submitData.price = Number(submitData.price)
      if (isNaN(submitData.price) || submitData.price < 0) {
        ElMessage.error('商品价格必须是大于等于0的数字')
        loading.value = false
        return
      }
      
      // 确保状态符合后端要求
      if (submitData.status !== '售卖中' && submitData.status !== '已下架') {
        submitData.status = '售卖中'  // 默认设为售卖中
      }
      
      let response
      
      if (isEdit.value) {
        // 更新商品
        const productId = parseInt(route.query.id)
        response = await productApi.updateProduct(productId, submitData)
        ElMessage.success('商品修改成功')
      } else {
        // 添加商品
        response = await productApi.addProduct(submitData)
        ElMessage.success('商品发布成功')
        
        // 添加到商品列表
        const newProduct = response.data
        productStore.addProduct(newProduct)
      }
      
      router.push('/products')
    } catch (error) {
      ElMessage.error('操作失败，请稍后再试')
      console.error('保存商品失败:', error)
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.publish-container {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.card-header {
  text-align: center;
}

.submit-button {
  width: 120px;
}

.product-image-uploader {
  margin-bottom: 15px;
}

.image-preview {
  margin-top: 15px;
  position: relative;
  width: 200px;
  border-radius: 4px;
  overflow: hidden;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.preview-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.preview-actions {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: rgba(0, 0, 0, 0.6);
  padding: 5px;
  text-align: center;
  opacity: 0;
  transition: opacity 0.3s;
}

.image-preview:hover .preview-actions {
  opacity: 1;
}
</style>