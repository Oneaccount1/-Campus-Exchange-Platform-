import { createRouter, createWebHistory } from 'vue-router'
import { useAdminStore, useUserStore } from '../stores'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue')
  },
  {
    path: '/products',
    name: 'Products',
    component: () => import('../views/Products.vue')
  },
  {
    path: '/product/:id',
    name: 'ProductDetail',
    component: () => import('../views/ProductDetail.vue')
  },

  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue')
  },
  {
    path: '/publish',
    name: 'Publish',
    component: () => import('../views/Publish.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/messages',
    name: 'Messages',
    component: () => import('../views/Messages.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/user',
    name: 'UserCenter',
    component: () => import('../views/UserCenter.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('../views/admin/AdminLogin.vue'),
    meta: { requiresAuth: false, isAdmin: true }
  },
  {
    path: '/admin',
    component: () => import('../views/admin/AdminLayout.vue'),
    redirect: '/admin/dashboard',
    meta: { requiresAuth: true, isAdmin: true },
    children: [
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: () => import('../views/admin/Dashboard.vue'),
        meta: { requiresAuth: true, isAdmin: true }
      },
      {
        path: 'products',
        name: 'AdminProducts',
        component: () => import('../views/admin/Products.vue'),
        meta: { requiresAuth: true, isAdmin: true }
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: () => import('../views/admin/Users.vue'),
        meta: { requiresAuth: true, isAdmin: true }
      },
      {
        path: 'orders',
        name: 'AdminOrders',
        component: () => import('../views/admin/Orders.vue'),
        meta: { requiresAuth: true, isAdmin: true }
      },
      {
        path: 'messages',
        name: 'AdminMessages',
        component: () => import('../views/admin/Messages.vue'),
        meta: { requiresAuth: true, isAdmin: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const adminStore = useAdminStore()
  const userStore = useUserStore()
  
  // Check if route requires auth
  if (to.meta.requiresAuth) {
    // Check if it's an admin route
    if (to.meta.isAdmin) {
      if (adminStore.isLoggedIn) {
        next()
      } else {
        next({ path: '/admin/login', query: { redirect: to.fullPath } })
      }
    } else {
      // Regular user auth check
      if (userStore.isLoggedIn) {
        next()
      } else {
        next({ path: '/login', query: { redirect: to.fullPath } })
      }
    }
  } else {
    // Public route - no auth needed
    next()
  }
})

export default router