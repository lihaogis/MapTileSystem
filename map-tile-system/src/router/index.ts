import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '@/components/layout/MainLayout.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/login/LoginView.vue')
    },
    {
      path: '/',
      component: MainLayout,
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/dashboard/DashboardView.vue')
        },
        {
          path: 'datasource',
          name: 'datasource',
          component: () => import('@/views/datasource/DataSourceView.vue')
        },
        {
          path: 'auth',
          name: 'auth',
          component: () => import('@/views/auth/AuthView.vue')
        },
        {
          path: 'access-log',
          name: 'access-log',
          component: () => import('@/views/access-log/AccessLogView.vue')
        },
        {
          path: 'statistics',
          name: 'statistics',
          component: () => import('@/views/statistics/StatisticsView.vue')
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('@/views/user/UserView.vue')
        }
      ]
    }
  ]
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
