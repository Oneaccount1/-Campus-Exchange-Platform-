import request from './request'

// 获取联系人列表
export const getContactList = () => {
    return request({
        url: '/messages/contacts',
        method: 'get'
    })
}

// 别名，保持兼容性
export const getContacts = getContactList

// 发送消息
export const sendMessage = (data) => {
    // 创建数据的副本以避免修改原始对象
    const cleanData = { ...data }
    
    // 确保参数名称与后端一致
    if (cleanData.receiverId && !cleanData.receiver_id) {
        cleanData.receiver_id = cleanData.receiverId
        delete cleanData.receiverId
    }
    
    if (cleanData.productId && !cleanData.product_id) {
        cleanData.product_id = cleanData.productId
        delete cleanData.productId
    }
    
    // 移除无效的product_id（未定义、null或0值）
    if (!cleanData.product_id || cleanData.product_id <= 0) {
        delete cleanData.product_id
    }
    
    // 确保消息类型字段存在且有效
    if (!cleanData.type) {
        cleanData.type = 'text'
    }
    
    console.log('发送消息清理后的数据:', cleanData)
    
    return request({
        url: '/messages',
        method: 'post',
        data: cleanData
    })
}

// 获取与联系人的消息历史
export const getMessageHistory = (contactId, page = 1, size = 20) => {
    // 将page/size转换为offset/limit
    const offset = (page - 1) * size
    const limit = size
    
    return request({
        url: `/messages/${contactId}?limit=${limit}&offset=${offset}`,
        method: 'get'
    })
}

// 标记消息为已读
export const markAsRead = (contactId, messageIds = []) => {
    return request({
        url: `/messages/${contactId}/read`,
        method: 'put',
        data: {
            message_ids: messageIds
        }
    })
}

// 获取未读消息数量
export const getUnreadCount = () => {
    return request({
        url: '/messages/unread/count',
        method: 'get'
    })
}

// 获取最近一条消息
export const getLastMessage = (contactId) => {
    return request({
        url: `/messages/${contactId}/last`,
        method: 'get'
    })
}

// 创建新的联系人会话
export const createConversation = (data) => {
    // 标准化数据格式
    const cleanData = { ...data }
    
    // 确保使用user_id而不是userId
    if (cleanData.userId && !cleanData.user_id) {
        cleanData.user_id = cleanData.userId
        delete cleanData.userId
    }
    
    // 确保使用product_id而不是productId
    if (cleanData.productId && !cleanData.product_id) {
        cleanData.product_id = cleanData.productId
        delete cleanData.productId
    }
    
    console.log('API调用: 创建新会话', cleanData)
    
    return request({
        url: '/messages/conversation',
        method: 'post',
        data: cleanData
    }).then(response => {
        console.log('API响应: 创建新会话成功', response)
        return response
    }).catch(error => {
        console.error('API错误: 创建新会话失败', error)
        throw error
    })
}

// 删除消息
export const deleteMessage = (messageId) => {
    return request({
        url: `/messages/${messageId}`,
        method: 'delete'
    })
}

// 获取系统消息列表
export const getSystemMessages = (page = 1, size = 10) => {
    return request({
        url: `/messages/system?page=${page}&size=${size}`,
        method: 'get'
    })
}

// 标记系统消息为已读
export const markSystemMessageRead = (messageId) => {
    return request({
        url: `/messages/system/${messageId}/read`,
        method: 'put'
    })
}

// 标记所有系统消息为已读
export const markAllSystemMessagesRead = () => {
    return request({
        url: '/messages/system/read-all',
        method: 'put'
    })
}

// 批量删除消息
export const batchDeleteMessages = (messageIds) => {
    return request({
        url: '/messages/batch',
        method: 'delete',
        data: {
            message_ids: messageIds
        }
    })
}

// 获取会话消息历史
export const getConversationHistory = (user1Id, user2Id, page = 1, size = 20) => {
    return request({
        url: `/messages/history?user1Id=${user1Id}&user2Id=${user2Id}&page=${page}&size=${size}`,
        method: 'get'
    })
}

// 获取消息会话列表
export const getMessageConversations = (page = 1, size = 10, search = '') => {
    let url = `/messages/conversations?page=${page}&size=${size}`
    if (search) {
        url += `&search=${encodeURIComponent(search)}`
    }
    return request({
        url,
        method: 'get'
    })
}

// 上传消息图片
export const uploadMessageImage = (formData) => {
    return request({
        url: '/messages/upload',
        method: 'post',
        data: formData,
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    })
}

// 转发消息
export const forwardMessage = (messageId, receiverId) => {
    return request({
        url: '/messages/forward',
        method: 'post',
        data: {
            message_id: messageId,
            receiver_id: receiverId
        }
    })
}
