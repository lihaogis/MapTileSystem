<template>
  <div class="layout-root">
    <!-- 侧边栏 -->
    <aside class="sidebar">
      <!-- Logo -->
      <div class="sidebar-logo">
        <div class="logo-icon">
          <svg viewBox="0 0 24 24" fill="none" class="w-5 h-5">
            <rect x="1" y="1" width="9" height="9" rx="1" fill="currentColor" opacity="0.9"/>
            <rect x="14" y="1" width="9" height="9" rx="1" fill="currentColor" opacity="0.6"/>
            <rect x="1" y="14" width="9" height="9" rx="1" fill="currentColor" opacity="0.6"/>
            <rect x="14" y="14" width="9" height="9" rx="1" fill="currentColor" opacity="0.3"/>
          </svg>
        </div>
        <span class="logo-text">瓦片管理系统</span>
      </div>

      <!-- 导航菜单 -->
      <nav class="sidebar-nav">
        <router-link
          v-for="item in menuItems"
          :key="item.path"
          :to="item.path"
          class="nav-item"
          :class="{ active: activeMenu === item.path }"
        >
          <el-icon :size="16"><component :is="item.icon" /></el-icon>
          <span>{{ item.label }}</span>
        </router-link>
      </nav>

      <!-- 底部版本 -->
      <div class="sidebar-footer">
        <span class="version-tag">z/x/y · v1.0</span>
      </div>
    </aside>

    <!-- 主内容区 -->
    <div class="main-wrap">
      <!-- 顶部栏 -->
      <header class="topbar">
        <div class="page-title">
          <span class="title-dot"></span>
          {{ pageTitle }}
        </div>
        <el-dropdown @command="handleCommand">
          <div class="user-btn">
            <div class="user-avatar">
              <el-icon :size="14"><User /></el-icon>
            </div>
            <span>管理员</span>
            <el-icon :size="12" class="arrow-icon"><ArrowDown /></el-icon>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </header>

      <!-- 页面内容 -->
      <main class="page-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { DataLine, FolderOpened, Key, User, UserFilled, ArrowDown, Document } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

const menuItems = [
  { path: '/dashboard', label: '仪表盘', icon: DataLine },
  { path: '/datasource', label: '数据源管理', icon: FolderOpened },
  { path: '/auth', label: '授权管理', icon: Key },
  { path: '/access-log', label: '访问日志', icon: Document },
  { path: '/users', label: '用户管理', icon: UserFilled },
]

const activeMenu = computed(() => route.path)

const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    '/dashboard': '仪表盘',
    '/datasource': '数据源管理',
    '/auth': '授权管理',
    '/access-log': '访问日志',
    '/statistics': '统计报表',
    '/users': '用户管理'
  }
  return titles[route.path] || '地图瓦片管理系统'
})

const handleCommand = (command: string) => {
  if (command === 'logout') {
    localStorage.removeItem('token')
    router.push('/login')
  }
}
</script>

<style scoped>
/* ===== 布局根 ===== */
.layout-root {
  display: flex;
  height: 100vh;
  background: #f5f7fa;
  overflow: hidden;
}

/* ===== 侧边栏 ===== */
.sidebar {
  width: 220px;
  flex-shrink: 0;
  background: rgba(10, 14, 22, 0.95);
  border-right: 1px solid rgba(20, 184, 166, 0.12);
  display: flex;
  flex-direction: column;
  position: relative;
  z-index: 10;
}

.sidebar::before {
  content: '';
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 1px;
  background: linear-gradient(to bottom, transparent, rgba(20, 184, 166, 0.3), transparent);
}

/* Logo */
.sidebar-logo {
  height: 60px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 20px;
  border-bottom: 1px solid rgba(20, 184, 166, 0.1);
}

.logo-icon {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #0d9488, #0891b2);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.logo-text {
  font-size: 14px;
  font-weight: 700;
  color: #e2e8f0;
  letter-spacing: 0.02em;
  white-space: nowrap;
}

/* 导航 */
.sidebar-nav {
  flex: 1;
  padding: 12px 10px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  color: rgba(100, 116, 139, 0.9);
  font-size: 13.5px;
  text-decoration: none;
  transition: all 0.2s;
  cursor: pointer;
}

.nav-item:hover {
  background: rgba(20, 184, 166, 0.08);
  color: rgba(148, 163, 184, 0.9);
}

.nav-item.active {
  background: rgba(13, 148, 136, 0.15);
  color: #14b8a6;
  border-left: 2px solid #0d9488;
  padding-left: 10px;
}

.nav-item.active .el-icon {
  color: #0d9488;
}

/* 底部 */
.sidebar-footer {
  padding: 16px 20px;
  border-top: 1px solid rgba(20, 184, 166, 0.1);
}

.version-tag {
  font-size: 11px;
  font-family: 'SF Mono', monospace;
  color: rgba(20, 184, 166, 0.3);
  letter-spacing: 0.08em;
}

/* ===== 主内容区 ===== */
.main-wrap {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #f5f7fa;
}

/* 顶部栏 */
.topbar {
  height: 60px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: white;
  border-bottom: 1px solid #e5e7eb;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.page-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  color: #1f2937;
}

.title-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #0d9488;
  box-shadow: 0 0 6px rgba(13, 148, 136, 0.6);
}

.user-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
  color: #6b7280;
  font-size: 13px;
}

.user-btn:hover {
  background: #f3f4f6;
}

.user-avatar {
  width: 28px;
  height: 28px;
  background: linear-gradient(135deg, #0d9488, #0891b2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.arrow-icon {
  color: #9ca3af;
}

/* 页面内容 */
.page-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  background: #f5f7fa;
}
</style>
