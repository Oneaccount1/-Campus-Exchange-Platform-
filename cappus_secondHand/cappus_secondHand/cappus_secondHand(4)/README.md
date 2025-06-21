# 校园二手交易平台

## 项目简介

校园二手交易平台是一个面向大学校园的二手物品交易系统，旨在为大学生提供一个便捷、安全的二手物品交易环境。该平台采用前后端分离架构，前端使用Vue.js，后端使用Go语言。

## 技术栈

### 前端
- Vue 3
- Vue Router
- Pinia (状态管理)
- Element Plus (UI组件库)
- Axios (HTTP客户端)
- WebSocket (实时消息)

### 后端
- Go
- Gin (Web框架)
- GORM (ORM框架)
- MySQL (数据库)
- Redis (缓存)
- JWT (身份认证)
- Casbin (权限控制)
- RabbitMQ (消息队列)
- WebSocket (实时通信)

## 功能模块

### 用户模块
- 用户注册、登录、退出
- 个人信息管理
- 收藏商品管理

### 商品模块
- 商品发布、编辑、下架
- 商品分类、搜索、筛选
- 商品详情查看

### 消息模块
- 用户实时聊天
- 系统通知

### 订单模块
- 下单、支付、确认收货
- 订单状态跟踪

### 管理员后台模块
- 管理员登录认证
- 仪表盘数据统计
- 用户管理
- 商品管理
- 系统设置

## 项目结构

### 前端结构
```
cappus_secondHand/
  ├── public/              # 静态资源
  ├── src/
  │   ├── api/             # API接口
  │   ├── assets/          # 资源文件
  │   ├── components/      # 公共组件
  │   ├── router/          # 路由配置
  │   ├── stores/          # Pinia状态管理
  │   ├── utils/           # 工具函数
  │   ├── views/           # 页面视图
  │   │   ├── admin/       # 管理员后台页面
  │   │   └── ...          # 其他前台页面
  │   ├── App.vue          # 根组件
  │   └── main.js          # 入口文件
  ├── package.json         # 依赖配置
  └── vite.config.js       # Vite配置
```

### 后端结构
```
cmd/
  └── main.go              # 程序入口
configs/                   # 配置文件
internal/
  ├── auth/                # 认证相关
  ├── bootstrap/           # 启动初始化
  ├── config/              # 配置加载
  ├── database/            # 数据库连接
  ├── middleware/          # 中间件
  ├── models/              # 数据模型
  ├── modules/             # 业务模块
  │   ├── message/         # 消息模块
  │   ├── order/           # 订单模块
  │   ├── permission/      # 权限模块
  │   ├── product/         # 商品模块
  │   └── user/            # 用户模块
  ├── rabbitMQ/            # 消息队列
  ├── router/              # 路由注册
  ├── utils/               # 工具函数
  └── websocket/           # WebSocket服务
```

## 管理员后台API接口

管理员后台API接口规范详见 [API接口规范文档](./cappus_secondHand(3)/src/api/README.md)。

主要包括以下接口：
1. 管理员登录认证
2. 仪表盘数据统计
3. 用户管理相关接口
4. 商品管理相关接口
5. 系统设置相关接口

## 运行项目

### 前端
```bash
cd cappus_secondHand
npm install
npm run dev
```

### 后端
```bash
go run cmd/main.go
```

## 项目预览

### 前台页面
- 首页
- 商品列表页
- 商品详情页
- 用户中心
- 消息中心

### 管理员后台页面
- 登录页
- 仪表盘
- 用户管理
- 商品管理
- 系统设置

This template should help get you started developing with Vue 3 in Vite. The template uses Vue 3 `<script setup>` SFCs, check out the [script setup docs](https://v3.vuejs.org/api/sfc-script-setup.html#sfc-script-setup) to learn more.

Learn more about IDE Support for Vue in the [Vue Docs Scaling up Guide](https://vuejs.org/guide/scaling-up/tooling.html#ide-support).
