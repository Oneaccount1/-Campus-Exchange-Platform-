<template>
  <div class="admin-login-container">
    <div class="login-content">
      <div class="login-header">
     
        <h1>校园二手交易平台</h1>
        <p class="login-subtitle">管理员后台系统</p>
      </div>
      
      <el-card class="admin-login-card" shadow="hover">
        <div class="card-header">
          <h2>管理员登录</h2>
          <p class="header-tip">请输入您的管理员账号和密码</p>
        </div>
        
        <el-form :model="loginForm" :rules="rules" ref="loginFormRef" label-position="top">
          <el-form-item label="用户名" prop="username">
            <el-input 
              v-model="loginForm.username" 
              placeholder="请输入管理员用户名"
              :prefix-icon="User"
              size="large"
            />
          </el-form-item>
          
          <el-form-item label="密码" prop="password">
            <el-input 
              v-model="loginForm.password" 
              type="password" 
              placeholder="请输入管理员密码" 
              show-password
              :prefix-icon="Lock"
              size="large"
            />
          </el-form-item>
          
          <div class="remember-forgot">
            <el-checkbox v-model="rememberMe">记住我</el-checkbox>
            <a href="#" class="forgot-link">忘记密码?</a>
          </div>
          
          <el-form-item>
            <el-button type="primary" class="login-button" @click="handleLogin" :loading="loading" size="large">
              登录系统
            </el-button>
          </el-form-item>
        </el-form>
        
        <div class="login-footer">
          <p>© 2025 校园二手交易平台 | 管理员系统</p>
        </div>
      </el-card>
      
      <div class="login-help">
        <el-button type="text" @click="goToFrontend">
          <el-icon><Back /></el-icon>
          返回前台首页
        </el-button>
      </div>
    </div>
    
    <div class="login-background">
      <div class="bg-pattern"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '../../stores'
import { ElMessage } from 'element-plus'
import { User, Lock, Back } from '@element-plus/icons-vue'
import { adminLogin } from '../../api/admin'

const router = useRouter()
const adminStore = useAdminStore()

const loginFormRef = ref(null)
const loading = ref(false)
const rememberMe = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入管理员用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度应为3-20个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入管理员密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度应为6-20个字符', trigger: 'blur' }
  ]
}

const handleLogin = () => {
  loginFormRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value = true
    
    try {
      // 调用管理员登录API
      const response = await adminLogin({
        user_name: loginForm.username,
        pass_word: loginForm.password
      })
      
      if (response.code === 200) {
        // 登录成功，保存管理员信息和token
        const adminInfo = response.data.adminInfo || {
          id: response.data.user_id || 0,
          username: response.data.username || loginForm.username,
          role: 'admin'
        }
        
        const token = response.data.token
        
        adminStore.setAdminInfo(adminInfo)
        adminStore.setToken(token)
        
        ElMessage.success('登录成功')
        
        // 如果有重定向参数，则跳转到对应页面
        const redirect = router.currentRoute.value.query.redirect || '/admin/dashboard'
        router.push(redirect)
      } else {
        ElMessage.error(response.message || '用户名或密码错误')
      }
    } catch (error) {
      console.error('登录错误:', error)
      ElMessage.error('登录失败，请检查用户名和密码')
    } finally {
      loading.value = false
    }
  })
}

// 前往前台首页
const goToFrontend = () => {
  router.push('/')
}
</script>

<style scoped>
.admin-login-container {
  display: flex;
  min-height: 100vh;
  position: relative;
  overflow: hidden;
}

.login-content {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100%;
  max-width: 500px;
  padding: 40px 20px;
  margin: 0 auto;
  position: relative;
  z-index: 2;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
  margin: 0;
  margin-bottom: 5px;
}

.login-subtitle {
  font-size: 16px;
  color: #606266;
  margin: 0;
}

.admin-login-card {
  width: 100%;
  border-radius: 10px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.card-header {
  text-align: center;
  margin-bottom: 20px;
}

.card-header h2 {
  font-size: 24px;
  font-weight: 500;
  color: #303133;
  margin: 0;
  margin-bottom: 10px;
}

.header-tip {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

.remember-forgot {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.forgot-link {
  font-size: 14px;
  color: #409EFF;
  text-decoration: none;
}

.forgot-link:hover {
  text-decoration: underline;
}

.login-button {
  width: 100%;
  height: 50px;
  font-size: 16px;
  border-radius: 4px;
}

.login-footer {
  text-align: center;
  margin-top: 30px;
  font-size: 12px;
  color: #909399;
}

.login-help {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: #606266;
}

.login-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e7ed 100%);
  z-index: 1;
}

.bg-pattern {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: 
    radial-gradient(circle at 25px 25px, #dcdfe6 2px, transparent 0),
    radial-gradient(circle at 75px 75px, #dcdfe6 2px, transparent 0);
  background-size: 100px 100px;
  opacity: 0.6;
}

@media (min-width: 768px) {
  .admin-login-container {
    background: transparent;
  }
  
  .login-content {
    margin-left: 50%;
    transform: translateX(-50%);
  }
  
  .login-background {
    width: 100%;
    background: linear-gradient(135deg, #f5f7fa 0%, #e4e7ed 100%);
  }
}
</style>