import { adminRequest } from './request'

/**
 * 管理员登录
 * @param {Object} data - 登录信息
 * @param {string} data.user_name - 用户名
 * @param {string} data.pass_word - 密码
 * @returns {Promise}
 */
export function adminLogin(data) {
  return adminRequest({
    url: '/admin/login',
    method: 'post',
    data
  })
}

/**
 * 获取仪表盘统计数据
 * @returns {Promise}
 */
export function getDashboardStats() {
  return adminRequest({
    url: '/admin/dashboard/stats',
    method: 'get'
  })
}

/**
 * 获取商品发布趋势数据
 * @param {string} timeRange - 时间范围，可选值：week, month
 * @returns {Promise}
 */
export function getProductTrend(timeRange = 'week') {
  return adminRequest({
    url: '/admin/dashboard/product-trend',
    method: 'get',
    params: { timeRange }
  })
}

/**
 * 获取商品分类统计数据
 * @returns {Promise}
 */
export function getCategoryStats() {
  return adminRequest({
    url: '/admin/dashboard/category-stats',
    method: 'get'
  })
}

/**
 * 获取最新商品列表
 * @param {Object} params - 查询参数
 * @param {number} params.limit - 限制数量
 * @returns {Promise}
 */
export function getLatestProducts(params) {
  return adminRequest({
    url: '/admin/dashboard/latest-products',
    method: 'get',
    params
  })
}

/**
 * 获取系统最近活动
 * @param {number} limit - 限制数量
 * @returns {Promise}
 */
export function getRecentActivities(limit = 5) {
  return adminRequest({
    url: '/admin/dashboard/activities',
    method: 'get',
    params: { limit }
  })
}

/**
 * 获取用户列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.size - 每页数量
 * @param {string} params.search - 搜索关键词
 * @param {string} params.status - 状态筛选
 * @param {Array} params.dateRange - 日期范围
 * @returns {Promise}
 */
export function getUsers(params) {
  return adminRequest({
    url: '/admin/users',
    method: 'get',
    params
  })
}

/**
 * 获取用户详情
 * @param {number} userId - 用户ID
 * @returns {Promise}
 */
export function getUserDetail(userId) {
  return adminRequest({
    url: `/admin/users/${userId}`,
    method: 'get'
  })
}

/**
 * 更新用户状态
 * @param {number} userId - 用户ID
 * @param {string} status - 状态，可选值：正常, 禁用
 * @returns {Promise}
 */
export function updateUserStatus(userId, status) {
  return adminRequest({
    url: `/admin/users/${userId}/status`,
    method: 'put',
    data: { status }
  })
}

/**
 * 重置用户密码
 * @param {number} userId - 用户ID
 * @returns {Promise}
 */
export function resetUserPassword(userId) {
  return adminRequest({
    url: `/admin/users/${userId}/reset-password`,
    method: 'post'
  })
}

/**
 * 获取商品列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.size - 每页数量
 * @param {string} params.search - 搜索关键词
 * @param {string} params.category - 分类筛选
 * @param {string} params.status - 状态筛选
 * @param {Array} params.dateRange - 日期范围
 * @returns {Promise}
 */
export function getProducts(params) {
  return adminRequest({
    url: '/admin/products',
    method: 'get',
    params
  })
}

/**
 * 获取商品详情
 * @param {number} productId - 商品ID
 * @returns {Promise}
 */
export function getProductDetail(productId) {
  return adminRequest({
    url: `/admin/products/${productId}`,
    method: 'get'
  })
}

/**
 * 更新商品状态
 * @param {number} productId - 商品ID
 * @param {string} status - 状态，可选值：已上架, 已下架, 待审核
 * @returns {Promise}
 */
export function updateProductStatus(productId, status) {
  return adminRequest({
    url: `/admin/products/${productId}/status`,
    method: 'put',
    data: { status }
  })
}

/**
 * 批量更新商品状态
 * @param {Array} productIds - 商品ID数组
 * @param {string} status - 状态，可选值：已上架, 已下架, 待审核
 * @returns {Promise}
 */
export function batchUpdateProductStatus(productIds, status) {
  return adminRequest({
    url: '/admin/products/batch-status',
    method: 'put',
    data: { productIds, status }
  })
}

/**
 * 删除商品
 * @param {number} productId - 商品ID
 * @returns {Promise}
 */
export function deleteProduct(productId) {
  return adminRequest({
    url: `/admin/products/${productId}`,
    method: 'delete'
  })
}

/**
 * 导出用户数据
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export function exportUsersData(params) {
  return adminRequest({
    url: '/admin/users/export',
    method: 'get',
    params,
    responseType: 'blob'
  })
}

/**
 * 导出商品数据
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export function exportProductsData(params) {
  return adminRequest({
    url: '/admin/products/export',
    method: 'get',
    params,
    responseType: 'blob'
  })
}

/**
 * 获取订单列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.size - 每页数量
 * @param {string} params.search - 搜索关键词
 * @param {string} params.status - 状态筛选
 * @param {string} params.startDate - 开始日期
 * @param {string} params.endDate - 结束日期
 * @returns {Promise}
 */
export function getOrders(params) {
  return adminRequest({
    url: '/admin/orders',
    method: 'get',
    params
  })
}

/**
 * 获取订单详情
 * @param {string} orderId - 订单ID
 * @returns {Promise}
 */
export function getOrderDetail(orderId) {
  return adminRequest({
    url: `/admin/orders/${orderId}`,
    method: 'get'
  })
}

/**
 * 更新订单状态
 * @param {string} orderId - 订单ID
 * @param {Object} data - 状态数据
 * @param {string} data.status - 状态
 * @param {string} data.remark - 备注
 * @returns {Promise}
 */
export function updateOrderStatus(orderId, data) {
  return adminRequest({
    url: `/admin/orders/${orderId}/status`,
    method: 'put',
    data
  })
}

/**
 * 导出订单数据
 * @param {Object} params - 查询参数
 * @returns {Promise}
 */
export function exportOrdersData(params) {
  return adminRequest({
    url: '/admin/orders/export',
    method: 'get',
    params,
    responseType: 'blob'
  })
}

/**
 * 获取消息列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.size - 每页数量
 * @param {string} params.search - 搜索关键词
 * @param {string} params.type - 消息类型，可选值：user, system
 * @param {string} params.startDate - 开始日期
 * @param {string} params.endDate - 结束日期
 * @returns {Promise}
 */
export function getMessages(params) {
  return adminRequest({
    url: '/admin/messages',
    method: 'get',
    params
  })
}

/**
 * 获取用户消息会话列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码
 * @param {number} params.size - 每页数量
 * @param {string} params.search - 搜索关键词
 * @returns {Promise}
 */
export function getConversations(params) {
  return adminRequest({
    url: '/admin/messages/conversations',
    method: 'get',
    params
  })
}

/**
 * 获取会话消息历史
 * @param {Object} params - 查询参数
 * @param {number} params.user1Id - 用户1 ID
 * @param {number} params.user2Id - 用户2 ID
 * @param {number} params.page - 页码
 * @param {number} params.size - 每页数量
 * @returns {Promise}
 */
export function getMessageHistory(params) {
  return adminRequest({
    url: '/admin/messages/history',
    method: 'get',
    params
  })
}

/**
 * 发送系统消息
 * @param {Object} data - 消息数据
 * @param {number} data.receiverId - 接收者ID，为0表示发送给所有用户
 * @param {string} data.content - 消息内容
 * @param {string} data.title - 消息标题
 * @returns {Promise}
 */
export function sendSystemMessageApi(data) {
  return adminRequest({
    url: '/admin/messages/system',
    method: 'post',
    data
  })
}

/**
 * 删除消息
 * @param {number} messageId - 消息ID
 * @returns {Promise}
 */
export function deleteMessage(messageId) {
  return adminRequest({
    url: `/admin/messages/${messageId}`,
    method: 'delete'
  })
} 