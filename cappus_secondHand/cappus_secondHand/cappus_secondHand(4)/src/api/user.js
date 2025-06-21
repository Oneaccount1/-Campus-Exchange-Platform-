import request from './request'
import { useUserStore } from '../stores'

// 用户登录
export const login = (data) => {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

// 用户注册
export const register = (data) => {
  return request({
    url: '/register',
    method: 'post',
    data
  })
}

// 获取用户个人信息
export function getUserInfo(userId) {
  if (userId) {
    // 获取指定用户信息
    return request({
      url: `/user/${userId}`,
      method: 'get'
    })
  } else {
    // 获取当前登录用户信息
    return request({
      url: '/user/profile',
      method: 'get'
    })
  }
}

// 更新用户信息
export const updateUserInfo = (data) => {
  return request({
    url: '/user/profile',
    method: 'put',
    data
  })
}

// 修改密码
export const changePassword = (data) => {
  return request({
    url: '/user/change-password',
    method: 'post',
    data
  })
}

// 获取用户发布的商品
export const getUserProducts = (userId, page = 1, size = 10) => {
  return request({
    url: `/product/user?user_id=${userId}&page=${page}&size=${size}`,
    method: 'get'
  })
}

// 获取用户收藏的商品
export const getUserFavorites = (page = 1, size = 10) => {
  return request({
    url: `/user/favorites?page=${page}&size=${size}`,
    method: 'get'
  })
}

// 添加收藏
export const addFavorite = (productId) => {
  return request({
    url: '/user/favorites',
    method: 'post',
    data: {
      product_id: productId
    }
  })
}

// 取消收藏
export const removeFavorite = (productId) => {
  return request({
    url: `/user/favorites/${productId}`,
    method: 'delete'
  })
}

// 检查商品是否已收藏
export const checkFavorite = (productId) => {
  return request({
    url: `/user/favorites/check/${productId}`,
    method: 'get'
  })
}

// 用户登出
export function logout() {
  return request({
    url: '/user/logout',
    method: 'post'
  })
}

// 上传头像
export const uploadAvatar = (formData) => {
  return request({
    url: '/user/avatar',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 获取用户统计信息
export const getUserStats = (userId) => {
  return request({
    url: `/user/stats${userId ? `?user_id=${userId}` : ''}`,
    method: 'get'
  })
}

// 获取用户的订单
export const getUserOrders = (page = 1, size = 10, status = '') => {
  let url = `/user/orders?page=${page}&size=${size}`
  if (status) {
    url += `&status=${status}`
  }
  return request({
    url,
    method: 'get'
  })
}

// 验证邮箱
export const verifyEmail = (token) => {
  return request({
    url: `/auth/verify-email?token=${token}`,
    method: 'get'
  })
}

// 找回密码请求
export const forgotPassword = (email) => {
  return request({
    url: '/auth/forgot-password',
    method: 'post',
    data: { email }
  })
}

// 重置密码
export const resetPassword = (token, password) => {
  return request({
    url: '/auth/reset-password',
    method: 'post',
    data: { token, password }
  })
}

// 获取用户评价
export const getUserReviews = (userId, page = 1, size = 10) => {
  return request({
    url: `/user/reviews?user_id=${userId}&page=${page}&size=${size}`,
    method: 'get'
  })
}

// 检查用户名是否可用
export const checkUsername = (username) => {
  return request({
    url: `/auth/check-username?username=${username}`,
    method: 'get'
  })
}

// 检查邮箱是否可用
export const checkEmail = (email) => {
  return request({
    url: `/auth/check-email?email=${email}`,
    method: 'get'
  })
}