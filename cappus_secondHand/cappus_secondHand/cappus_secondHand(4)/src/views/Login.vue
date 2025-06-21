<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <h2>用户登录</h2>
        </div>
      </template>
      
      <el-form :model="loginForm" :rules="rules" ref="loginFormRef" label-position="top">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="loginForm.username" placeholder="请输入用户名">
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="请输入密码" show-password>
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        
        <div class="remember-forgot">
          <el-checkbox v-model="rememberMe">记住我</el-checkbox>
          <el-link type="primary" :underline="false">忘记密码?</el-link>
        </div>
        
        <el-form-item>
          <el-button type="primary" class="login-button" @click="handleLogin" :loading="loading">登录</el-button>
        </el-form-item>
        
        <div class="register-link">
          还没有账号? <el-link type="primary" @click="goToRegister">立即注册</el-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loginFormRef = ref(null)
const loading = ref(false)
const rememberMe = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度应为3-20个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度应为6-20个字符', trigger: 'blur' }
  ]
}

const handleLogin = () => {
  loginFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value = true
    
    try {
      // 调用登录API
      const { login } = await import('../api/user.js')
      const response = await login({
        user_name: loginForm.username,
        pass_word: loginForm.password
      })
      
      console.log('登录响应:', response)
      
      if (response.code === 200 && response.data) {
        // 检查token是否存在
        if (!response.data.token) {
          throw new Error('服务器未返回认证token')
        }
        
        // 从响应数据中构建用户信息对象
        const userInfo = {
          id: response.data.user_id,
          username: response.data.username,
          roles: response.data.roles || [],
          // 可以添加其他可能的字段
          avatar: response.data.avatar || ''
        }
        
        // 确保token不包含双引号
        const cleanToken = response.data.token.replace(/^"(.*)"$/, '$1')
        
        // 使用loginSuccess方法一次性处理登录和WebSocket连接
        await userStore.loginSuccess(userInfo, cleanToken)
        
        ElMessage.success('登录成功')
        
        // 如果有重定向，登录后跳转回原始请求页面
        const redirectPath = route.query.redirect || '/'
        router.push(redirectPath)
      } else {
        // 检查是否是账号被禁用的错误
        if (response.code === 403 && response.message && response.message.includes('禁用')) {
          ElMessage.error({
            message: '账号已被禁用，请联系管理员',
            duration: 5000,
            showClose: true
          })
        } else {
          ElMessage.error(response.message || '登录失败，服务器返回错误')
        }
      }
    } catch (error) {
      console.error('登录错误:', error)
      
      // 检查错误信息是否包含禁用关键词
      if (error.response && error.response.data) {
        const errorData = error.response.data
        if (errorData.code === 403 && errorData.message && errorData.message.includes('禁用')) {
          ElMessage.error({
            message: '账号已被禁用，请联系管理员',
            duration: 5000,
            showClose: true
          })
          return
        }
      }
      
      ElMessage.error('登录失败，请检查用户名和密码或网络连接')
    } finally {
      loading.value = false
    }
  })
}

const goToRegister = () => {
  router.push('/register')
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 180px);
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 400px;
}

.card-header {
  text-align: center;
}

.login-button {
  width: 100%;
  margin-top: 10px;
}

.remember-forgot {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.register-link {
  text-align: center;
  margin-top: 15px;
}
</style>