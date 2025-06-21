import request from './request'

// 获取商品列表（带分页）
export const getProductList = (page = 1, size = 12) => {
    return request({
        url: `/product?page=${page}&size=${size}`,
        method: 'get'
    })
}

// 搜索商品
export const searchProducts = (keyword, page = 1, size = 12) => {
    return request({
        url: `/product/search?keyword=${keyword}&page=${page}&size=${size}`,
        method: 'get'
    })
}

// 通过id获取商品详情
export const getProductById = (id) => {
    return request({
        url: `/product/${id}`,
        method: 'get'
    })
}

// 添加商品
export const addProduct = (data) => {
    return request({
        url: '/product',
        method: 'post',
        data
    })
}

// 更新商品
export const updateProduct = (id, data) => {
    return request({
        url: `/product/${id}`,
        method: 'put',
        data
    })
}

// 删除商品
export const deleteProduct = (id) => {
    return request({
        url: `/product/${id}`,
        method: 'delete'
    })
}

// 获取商品分类列表
export const getCategoryList = () => {
    return request({
        url: '/product/categories',
        method: 'get'
    })
}

// 按分类获取商品
export const getProductsByCategory = (categoryId, page = 1, size = 12) => {
    return request({
        url: `/product/category/${categoryId}?page=${page}&size=${size}`,
        method: 'get'
    })
}

// 获取推荐商品
export const getRecommendProducts = (limit = 6) => {
    return request({
        url: `/product/recommend?limit=${limit}`,
        method: 'get'
    })
}

// 获取最新商品
export const getLatestProducts = (limit = 8) => {
    return request({
        url: '/product/latest',
        method: 'get',
        params: { limit }
    })
}

// 上传商品图片
export const uploadProductImage = (formData) => {
    return request({
        url: '/product/upload',
        method: 'post',
        data: formData,
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    })
}

// 更新商品状态
export const updateProductStatus = (id, status) => {
    return request({
        url: `/product/${id}/status`,
        method: 'put',
        data: { status }
    })
}




