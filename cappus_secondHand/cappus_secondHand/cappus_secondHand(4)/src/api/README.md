# 校园二手交易平台 API 接口规范

## 基础URL

- 前台API: `http://localhost:8080/api/v1`
- 后台API: `http://localhost:8080/api/v1/admin`

## 通用响应格式

```json
{
  "code": 200,       // 状态码，200表示成功，非200表示失败
  "message": "操作成功", // 提示信息
  "data": {}         // 响应数据，可能是对象、数组或null
}
```

## 错误码

- 200: 成功
- 400: 请求参数错误
- 401: 未授权或token过期
- 403: 权限不足
- 404: 资源不存在
- 500: 服务器内部错误

## 管理员API接口

### 1. 管理员登录

- **URL**: `/api/v1/admin/auth/login`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "username": "admin",
    "password": "admin123"
  }
  ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "登录成功",
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "adminInfo": {
        "id": 1,
        "username": "admin",
        "role": "admin"
      }
    }
  }
  ```

### 2. 获取仪表盘统计数据

- **URL**: `/api/v1/admin/dashboard/stats`
- **方法**: `GET`
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "productCount": 125,         // 商品总数
      "productTrend": 15,          // 商品增长趋势，百分比
      "userCount": 368,            // 用户总数
      "userTrend": 8,              // 用户增长趋势，百分比
      "newUserCount": 12,          // 本周新增用户数
      "todayProductCount": 12,     // 今日新增商品数
      "todayProductTrend": -5,     // 今日商品增长趋势，百分比
      "yesterdayProductCount": 15, // 昨日新增商品数
      "totalAmount": 8526,         // 交易总额
      "amountTrend": 12,           // 交易额增长趋势，百分比
      "monthAmount": 3245          // 本月交易额
    }
  }
  ```

### 3. 获取商品发布趋势数据

- **URL**: `/api/v1/admin/dashboard/product-trend`
- **方法**: `GET`
- **参数**:
  - `timeRange`: 时间范围，可选值：`week`, `month`，默认为`week`
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": [
      { "label": "周一", "value": 32 },
      { "label": "周二", "value": 45 },
      { "label": "周三", "value": 38 },
      { "label": "周四", "value": 25 },
      { "label": "周五", "value": 60 },
      { "label": "周六", "value": 48 },
      { "label": "周日", "value": 35 }
    ]
  }
  ```

### 4. 获取商品分类统计数据

- **URL**: `/api/v1/admin/dashboard/category-stats`
- **方法**: `GET`
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": [
      { "name": "电子产品", "value": 42 },
      { "name": "图书教材", "value": 28 },
      { "name": "生活用品", "value": 18 },
      { "name": "服装鞋包", "value": 12 }
    ]
  }
  ```

### 5. 获取最新商品列表

- **URL**: `/api/v1/admin/products/latest`
- **方法**: `GET`
- **参数**:
  - `limit`: 限制数量，默认为5
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": [
      {
        "id": 1,
        "title": "全新iPad Pro 2023款",
        "price": 6999,
        "category": "电子产品",
        "seller": "张三",
        "createTime": "2024-05-20 14:30:22",
        "image": "https://example.com/image1.jpg",
        "status": "已上架"
      },
      // 更多商品...
    ]
  }
  ```

### 6. 获取系统最近活动

- **URL**: `/api/v1/admin/dashboard/activities`
- **方法**: `GET`
- **参数**:
  - `limit`: 限制数量，默认为5
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": [
      {
        "content": "管理员审核通过了商品 \"全新iPad Pro 2023款\"",
        "time": "10分钟前",
        "type": "primary",
        "color": "#409EFF"
      },
      // 更多活动...
    ]
  }
  ```

### 7. 获取用户列表

- **URL**: `/api/v1/admin/users`
- **方法**: `GET`
- **参数**:
  - `page`: 页码，默认为1
  - `pageSize`: 每页数量，默认为10
  - `search`: 搜索关键词，可选
  - `status`: 状态筛选，可选值：`正常`, `禁用`
  - `startDate`: 开始日期，格式：`YYYY-MM-DD`，可选
  - `endDate`: 结束日期，格式：`YYYY-MM-DD`，可选
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "total": 368,
      "list": [
        {
          "id": 1,
          "username": "张三",
          "avatar": "https://example.com/avatar1.jpg",
          "registerTime": "2024-01-15 10:30:00",
          "productCount": 5,
          "status": "正常",
          "email": "zhangsan@example.com",
          "phone": "13800138000"
        },
        // 更多用户...
      ]
    }
  }
  ```

### 8. 获取用户详情

- **URL**: `/api/v1/admin/users/:userId`
- **方法**: `GET`
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "id": 1,
      "username": "张三",
      "avatar": "https://example.com/avatar1.jpg",
      "registerTime": "2024-01-15 10:30:00",
      "email": "zhangsan@example.com",
      "phone": "13800138000",
      "lastLogin": "2024-05-20 15:30:00",
      "lastIp": "192.168.1.1",
      "status": "正常",
      "productCount": 5,
      "orderCount": 10,
      "favoriteCount": 8,
      "products": [
        {
          "id": 1,
          "title": "全新iPad Pro 2023款",
          "price": 6999,
          "createTime": "2024-05-20 14:30:22"
        },
        // 更多商品...
      ],
      "activities": [
        {
          "content": "发布了商品 \"全新iPad Pro 2023款\"",
          "time": "2024-05-20 14:30:22"
        },
        // 更多活动...
      ]
    }
  }
  ```

### 9. 更新用户状态

- **URL**: `/api/v1/admin/users/:userId/status`
- **方法**: `PUT`
- **请求体**:
  ```json
  {
    "status": "正常" // 或 "禁用"
  }
  ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "更新成功",
    "data": null
  }
  ```

### 10. 重置用户密码

- **URL**: `/api/v1/admin/users/:userId/reset-password`
- **方法**: `POST`
- **响应**:
  ```json
  {
    "code": 200,
    "message": "密码已重置，新密码为：123456",
    "data": {
      "newPassword": "123456"
    }
  }
  ```

### 11. 获取商品列表

- **URL**: `/api/v1/admin/products`
- **方法**: `GET`
- **参数**:
  - `page`: 页码，默认为1
  - `pageSize`: 每页数量，默认为10
  - `search`: 搜索关键词，可选
  - `category`: 分类筛选，可选
  - `status`: 状态筛选，可选值：`已上架`, `已下架`, `待审核`
  - `startDate`: 开始日期，格式：`YYYY-MM-DD`，可选
  - `endDate`: 结束日期，格式：`YYYY-MM-DD`，可选
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "total": 125,
      "list": [
        {
          "id": 1,
          "title": "全新iPad Pro 2023款",
          "price": 6999,
          "category": "电子产品",
          "seller": "张三",
          "createTime": "2024-05-20 14:30:22",
          "image": "https://example.com/image1.jpg",
          "status": "已上架"
        },
        // 更多商品...
      ]
    }
  }
  ```

### 12. 获取商品详情

- **URL**: `/api/v1/admin/products/:productId`
- **方法**: `GET`
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "id": 1,
      "title": "全新iPad Pro 2023款",
      "price": 6999,
      "category": "电子产品",
      "seller": "张三",
      "sellerId": 1,
      "createTime": "2024-05-20 14:30:22",
      "updateTime": "2024-05-20 14:30:22",
      "description": "全新未拆封，支持面容ID，M2芯片",
      "status": "已上架",
      "images": [
        "https://example.com/image1.jpg",
        "https://example.com/image2.jpg"
      ]
    }
  }
  ```

### 13. 更新商品状态

- **URL**: `/api/v1/admin/products/:productId/status`
- **方法**: `PUT`
- **请求体**:
  ```json
  {
    "status": "已上架" // 或 "已下架", "待审核"
  }
  ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "更新成功",
    "data": null
  }
  ```

### 14. 批量更新商品状态

- **URL**: `/api/v1/admin/products/batch-status`
- **方法**: `PUT`
- **请求体**:
  ```json
  {
    "productIds": [1, 2, 3],
    "status": "已上架" // 或 "已下架", "待审核"
  }
  ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "更新成功",
    "data": null
  }
  ```

### 15. 删除商品

- **URL**: `/api/v1/admin/products/:productId`
- **方法**: `DELETE`
- **响应**:
  ```json
  {
    "code": 200,
    "message": "删除成功",
    "data": null
  }
  ```

### 16. 导出用户数据

- **URL**: `/api/v1/admin/users/export`
- **方法**: `GET`
- **参数**: 同获取用户列表
- **响应**: Excel文件流

### 17. 导出商品数据

- **URL**: `/api/v1/admin/products/export`
- **方法**: `GET`
- **参数**: 同获取商品列表
- **响应**: Excel文件流

### 18. 获取订单列表

- **URL**: `/api/v1/admin/orders`
- **方法**: `GET`
- **参数**:
  - `page`: 页码，默认为1
  - `pageSize`: 每页数量，默认为10
  - `search`: 搜索关键词，可选
  - `status`: 状态筛选，可选值：`待付款`, `待发货`, `待收货`, `已完成`, `已取消`
  - `startDate`: 开始日期，格式：`YYYY-MM-DD`，可选
  - `endDate`: 结束日期，格式：`YYYY-MM-DD`，可选
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "total": 85,
      "list": [
        {
          "id": "O20240520001",
          "productTitle": "全新iPad Pro 2023款",
          "productImage": "https://example.com/image1.jpg",
          "price": 6999,
          "buyer": "李四",
          "seller": "张三",
          "status": "已完成",
          "createTime": "2024-05-20 14:30:22",
          "payTime": "2024-05-20 14:35:10",
          "completeTime": "2024-05-22 10:15:30"
        },
        // 更多订单...
      ]
    }
  }
  ```

### 19. 获取订单详情

- **URL**: `/api/v1/admin/orders/:orderId`
- **方法**: `GET`
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "id": "O20240520001",
      "productId": 1,
      "productTitle": "全新iPad Pro 2023款",
      "productImage": "https://example.com/image1.jpg",
      "price": 6999,
      "buyerId": 2,
      "buyer": "李四",
      "buyerPhone": "13900139000",
      "buyerAddress": "北京市海淀区清华大学学生公寓",
      "sellerId": 1,
      "seller": "张三",
      "sellerPhone": "13800138000",
      "status": "已完成",
      "createTime": "2024-05-20 14:30:22",
      "payTime": "2024-05-20 14:35:10",
      "deliveryTime": "2024-05-21 09:22:15",
      "completeTime": "2024-05-22 10:15:30",
      "remark": "请包装好再发货",
      "logs": [
        {
          "action": "创建订单",
          "time": "2024-05-20 14:30:22",
          "operator": "李四",
          "remark": null
        },
        {
          "action": "支付订单",
          "time": "2024-05-20 14:35:10",
          "operator": "李四",
          "remark": "微信支付"
        },
        {
          "action": "确认发货",
          "time": "2024-05-21 09:22:15",
          "operator": "张三",
          "remark": "已发顺丰快递，单号SF123456789"
        },
        {
          "action": "确认收货",
          "time": "2024-05-22 10:15:30",
          "operator": "李四",
          "remark": "商品完好"
        }
      ]
    }
  }
  ```

### 20. 更新订单状态

- **URL**: `/api/v1/admin/orders/:orderId/status`
- **方法**: `PUT`
- **请求体**:
  ```json
  {
    "status": "已完成", // 或其他状态
    "remark": "管理员手动更新状态"
  }
  ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "更新成功",
    "data": null
  }
  ```

### 21. 导出订单数据

- **URL**: `/api/v1/admin/orders/export`
- **方法**: `GET`
- **参数**: 同获取订单列表
- **响应**: Excel文件流

### 22. 获取消息列表

- **URL**: `/api/v1/admin/messages`
- **方法**: `GET`
- **参数**:
  - `page`: 页码，默认为1
  - `pageSize`: 每页数量，默认为10
  - `search`: 搜索关键词，可选
  - `type`: 消息类型，可选值：`user`, `system`，默认所有类型
  - `startDate`: 开始日期，格式：`YYYY-MM-DD`，可选
  - `endDate`: 结束日期，格式：`YYYY-MM-DD`，可选
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "total": 150,
      "list": [
        {
          "id": 1,
          "type": "user",
          "senderId": 1,
          "sender": "张三",
          "receiverId": 2,
          "receiver": "李四",
          "content": "你好，请问商品还在吗？",
          "createTime": "2024-05-20 14:30:22",
          "status": "已读",
          "readTime": "2024-05-20 14:35:10"
        },
        {
          "id": 2,
          "type": "system",
          "senderId": 0,
          "sender": "系统",
          "receiverId": 1,
          "receiver": "张三",
          "content": "您的订单已完成",
          "createTime": "2024-05-20 14:40:00",
          "status": "未读",
          "readTime": null
        },
        // 更多消息...
      ]
    }
  }
  ```

### 23. 获取用户消息会话列表

- **URL**: `/api/v1/admin/messages/conversations`
- **方法**: `GET`
- **参数**:
  - `page`: 页码，默认为1
  - `pageSize`: 每页数量，默认为10
  - `search`: 搜索关键词，可选
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "total": 45,
      "list": [
        {
          "id": 1,
          "user1Id": 1,
          "user1Name": "张三",
          "user1Avatar": "https://example.com/avatar1.jpg",
          "user2Id": 2,
          "user2Name": "李四",
          "user2Avatar": "https://example.com/avatar2.jpg",
          "lastMessage": "你好，请问商品还在吗？",
          "lastTime": "2024-05-20 14:30:22",
          "unreadCount": 0
        },
        // 更多会话...
      ]
    }
  }
  ```

### 24. 获取会话消息历史

- **URL**: `/api/v1/admin/messages/history`
- **方法**: `GET`
- **参数**:
  - `user1Id`: 用户1 ID
  - `user2Id`: 用户2 ID
  - `page`: 页码，默认为1
  - `pageSize`: 每页数量，默认为20
- **响应**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "total": 28,
      "list": [
        {
          "id": 1,
          "senderId": 1,
          "sender": "张三",
          "senderAvatar": "https://example.com/avatar1.jpg",
          "receiverId": 2,
          "content": "你好，请问商品还在吗？",
          "createTime": "2024-05-20 14:30:22",
          "status": "已读"
        },
        {
          "id": 2,
          "senderId": 2,
          "sender": "李四",
          "senderAvatar": "https://example.com/avatar2.jpg",
          "receiverId": 1,
          "content": "在的，还可以购买",
          "createTime": "2024-05-20 14:31:15",
          "status": "已读"
        },
        // 更多消息...
      ]
    }
  }
  ```

### 25. 发送系统消息

- **URL**: `/api/v1/admin/messages/system`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "receiverId": 1, // 接收者ID，为0表示发送给所有用户
    "content": "系统维护公告：系统将于2024-05-25 02:00-04:00进行维护，期间将无法访问",
    "title": "系统维护公告" // 可选
  }
  ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "发送成功",
    "data": null
  }
  ```

### 26. 删除消息

- **URL**: `/api/v1/admin/messages/:messageId`
- **方法**: `DELETE`
- **响应**:
  ```json
  {
    "code": 200,
    "message": "删除成功",
    "data": null
  }
  ```

## 前台API接口

(前台API接口已存在，此处省略)

## WebSocket通信规范

### WebSocket连接

WebSocket连接使用标准WebSocket协议，连接URL为：
```
ws(s)://{host}/api/v1/messages/ws?token={jwt_token}
```

认证通过查询参数`token`传递JWT令牌。

### 消息格式

所有WebSocket消息应当遵循以下标准JSON格式：

```json
{
  "type": "message_type",         // 消息类型，如 "message", "typing", "online_status" 等
  "data": {                       // 消息数据，所有业务数据应放在此对象中
    "senderId": 123,              // 发送者ID (如适用)
    "recipientId": 456,           // 接收者ID (如适用)
    "content": "消息内容",        // 消息内容 (如适用)
    // 其他特定消息类型的字段
  },
  "timestamp": 1629893718000      // 消息发送时间戳 (毫秒)
}
```

### 心跳机制

系统使用WebSocket标准的Ping/Pong控制帧进行心跳检测，不使用应用层Ping/Pong消息。
- 服务器每30秒发送一次Ping控制帧
- 客户端自动响应Pong控制帧
- 客户端60秒无响应则视为连接断开

### 标准消息类型

#### 1. 聊天消息

```json
{
  "type": "message",
  "data": {
    "id": 1234,                   // 消息ID
    "senderId": 123,              // 发送者ID
    "recipientId": 456,           // 接收者ID
    "content": "你好，世界！",    // 消息内容
    "isRead": false,              // 是否已读
    "productId": 789,             // 相关商品ID (可选)
    "createdAt": "2023-04-01T12:34:56Z" // 创建时间
  },
  "timestamp": 1629893718000
}
```

#### 2. 在线状态

```json
{
  "type": "online_status",
  "data": {
    "userId": 123,                // 用户ID
    "online": true,               // 是否在线
    "lastSeen": "2023-04-01T12:34:56Z" // 最后在线时间 (离线时)
  },
  "timestamp": 1629893718000
}
```

#### 3. 正在输入状态

```json
{
  "type": "typing",
  "data": {
    "senderId": 123,              // 发送者ID
    "recipientId": 456,           // 接收者ID
    "status": true                // true = 正在输入, false = 停止输入
  },
  "timestamp": 1629893718000
}
```

#### 4. 系统通知

```json
{
  "type": "notification",
  "data": {
    "message": "新版本已发布",    // 通知内容
    "type": "info",               // 通知类型: "info", "warning", "error", "success"
    "duration": 5000              // 显示时长 (毫秒)
  },
  "timestamp": 1629893718000
}
```

#### 5. 订阅请求

客户端发送:
```json
{
  "type": "subscribe",
  "data": {
    "recipientId": 456            // 想要订阅消息的用户ID
  },
  "timestamp": 1629893718000
}
```

服务器响应:
```json
{
  "type": "subscribe_success",
  "data": {
    "recipientId": 456            // 订阅成功的用户ID
  },
  "timestamp": 1629893718000
}
```

#### 6. 在线状态请求

```json
{
  "type": "online_status_request",
  "data": {
    "userIds": [123, 456, 789]    // 请求查询的用户ID列表
  },
  "timestamp": 1629893718000
}
```

### 错误处理

当发生错误时，服务器将发送错误消息：

```json
{
  "type": "error",
  "data": {
    "code": 403,                  // 错误代码
    "message": "权限不足"         // 错误消息
  },
  "timestamp": 1629893718000
}
```

### 连接管理

1. 客户端应实现自动重连机制，在连接断开时尝试重新连接
2. 重连应使用指数退避策略，避免频繁重连请求
3. 消息应进行排队，在连接恢复后自动发送 