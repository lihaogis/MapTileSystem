<template>
  <el-container class="h-screen">
    <el-aside width="240px" class="bg-gradient-to-b from-gray-900 to-gray-800 shadow-xl">
      <div class="h-16 flex items-center justify-center text-white border-b border-gray-700">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-lg flex items-center justify-center">
            <el-icon :size="20"><MapLocation /></el-icon>
          </div>
          <span class="text-lg font-bold">地图瓦片管理系统</span>
        </div>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="border-none mt-4 sidebar-menu"
        background-color="transparent"
        text-color="#9ca3af"
        active-text-color="#3b82f6"
        @select="handleMenuSelect"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataLine /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/datasource">
          <el-icon><FolderOpened /></el-icon>
          <span>数据源管理</span>
        </el-menu-item>
        <el-menu-item index="/auth">
          <el-icon><Key /></el-icon>
          <span>授权管理</span>
        </el-menu-item>
        <el-menu-item index="/access-log">
          <el-icon><Document /></el-icon>
          <span>访问日志</span>
        </el-menu-item>
        <el-menu-item index="/users">
          <el-icon><UserFilled /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="bg-white shadow-md flex items-center justify-between px-8 border-b border-gray-100">
        <div class="text-xl font-semibold text-gray-800">{{ pageTitle }}</div>
        <el-dropdown @command="handleCommand">
          <span class="cursor-pointer flex items-center space-x-2 px-3 py-2 rounded-lg hover:bg-gray-50 transition-colors">
            <div class="w-8 h-8 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-full flex items-center justify-center">
              <el-icon class="text-white" :size="16"><User /></el-icon>
            </div>
            <span class="font-medium text-gray-700">管理员</span>
            <el-icon class="text-gray-400"><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>
      <el-main class="bg-gray-50">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { DataLine, FolderOpened, Key, TrendCharts, User, UserFilled, MapLocation, ArrowDown, Document } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

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

const handleMenuSelect = (index: string) => {
  router.push(index)
}

const handleCommand = (command: string) => {
  if (command === 'logout') {
    localStorage.removeItem('token')
    router.push('/login')
  }
}
</script>

<style scoped>
.sidebar-menu :deep(.el-menu-item.is-active) {
  background-color: rgba(59, 130, 246, 0.15) !important;
  border-left: 3px solid #3b82f6;
  color: #3b82f6 !important;
  font-weight: 600;
}

.sidebar-menu :deep(.el-menu-item:hover) {
  background-color: rgba(255, 255, 255, 0.05) !important;
}

.sidebar-menu :deep(.el-menu-item) {
  margin: 4px 12px;
  border-radius: 8px;
  transition: all 0.3s;
}
</style>
