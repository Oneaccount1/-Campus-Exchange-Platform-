import request from './request'

// 获取用户订单列表
export const getUserOrders = (userId, page = 1, size = 10, status = '') => {
  let url = `/order/user?user_id=${userId}&page=${page}&size=${size}`
  if (status) {
    url += `&status=${status}`
  }
  return request({
    url,
    method: 'get'
  })
}

// 获取订单详情
export const getOrderDetail = (orderId) => {
  return request({
    url: `/order/${orderId}`,
    method: 'get'
  })
}

// 创建订单
export const createOrder = (data) => {
  return request({
    url: '/order',
    method: 'post',
    data
  })
}

// 更新订单状态
export const updateOrderStatus = (orderId, status) => {
  return request({
    url: `/order/${orderId}/status`,
    method: 'put',
    data: { status }
  })
}

// 删除订单
export const deleteOrder = (orderId) => {
  return request({
    url: `/order/${orderId}`,
    method: 'delete'
  })
}

// 支付订单 - 实际上是更新订单状态
export const payOrder = (orderId) => {
  return updateOrderStatus(orderId, '卖家已同意')
}

// 确认收货 - 实际上是更新订单状态
export const confirmReceive = (orderId) => {
  return updateOrderStatus(orderId, '已完成')
}

// 取消订单 - 实际上是更新订单状态
export const cancelOrder = (orderId) => {
  return updateOrderStatus(orderId, '卖家已拒绝')
}

// 评价订单 - 这个功能可能需要单独实现
export const reviewOrder = (orderId, data) => {
  return request({
    url: `/order/${orderId}/review`,
    method: 'post',
    data
  })
}

// 获取订单统计信息
export const getOrderStats = () => {
  return request({
    url: '/order/stats',
    method: 'get'
  })
}

// 获取卖家订单列表
export const getSellerOrders = (page = 1, size = 10, status = '') => {
  let url = `/order/seller?page=${page}&size=${size}`
  if (status) {
    url += `&status=${status}`
  }
  return request({
    url,
    method: 'get'
  })
}

// 同意订单
export const approveOrder = (orderId) => {
  return updateOrderStatus(orderId, '已发货')
}

// 拒绝订单
export const rejectOrder = (orderId, reason) => {
  return request({
    url: `/order/${orderId}/reject`,
    method: 'put',
    data: { reason }
  })
}

// 获取订单日志
export const getOrderLogs = (orderId) => {
  return request({
    url: `/order/${orderId}/logs`,
    method: 'get'
  })
}

// 添加订单备注
export const addOrderRemark = (orderId, remark) => {
  return request({
    url: `/order/${orderId}/remark`,
    method: 'post',
    data: { remark }
  })
} 