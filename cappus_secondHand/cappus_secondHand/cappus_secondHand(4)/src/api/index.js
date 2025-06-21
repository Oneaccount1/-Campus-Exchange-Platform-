// API服务统一出口
import * as userApi from './user'
import * as productApi from './product'
import request from './request'
import * as messageApi from './message'
import * as orderApi from './order'
import * as adminApi from './admin'

export {
  userApi,
  productApi,
  messageApi,
  orderApi,
  adminApi
}

// 导出请求实例，供直接使用
export default request