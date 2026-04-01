<template>
  <div class="space-y-6">
    <div class="grid grid-cols-5 gap-6">
      <el-card class="stat-card shadow-md hover:shadow-lg transition-shadow">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-gray-500 text-sm mb-2">今日总请求</div>
            <div class="text-3xl font-bold text-blue-600">{{ stats.todayTotal.toLocaleString() }}</div>
          </div>
          <div class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
            <el-icon :size="24" class="text-blue-600"><TrendCharts /></el-icon>
          </div>
        </div>
      </el-card>
      <el-card class="stat-card shadow-md hover:shadow-lg transition-shadow">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-gray-500 text-sm mb-2">昨日总请求</div>
            <div class="text-3xl font-bold text-green-600">{{ stats.yesterdayTotal.toLocaleString() }}</div>
          </div>
          <div class="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center">
            <el-icon :size="24" class="text-green-600"><DataLine /></el-icon>
          </div>
        </div>
      </el-card>
      <el-card class="stat-card shadow-md hover:shadow-lg transition-shadow">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-gray-500 text-sm mb-2">本月累计</div>
            <div class="text-3xl font-bold text-purple-600">{{ stats.monthTotal.toLocaleString() }}</div>
          </div>
          <div class="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center">
            <el-icon :size="24" class="text-purple-600"><DataAnalysis /></el-icon>
          </div>
        </div>
      </el-card>
      <el-card class="stat-card shadow-md hover:shadow-lg transition-shadow">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-gray-500 text-sm mb-2">全部累计</div>
            <div class="text-3xl font-bold text-rose-600">{{ stats.totalAllTime.toLocaleString() }}</div>
          </div>
          <div class="w-12 h-12 bg-rose-100 rounded-lg flex items-center justify-center">
            <el-icon :size="24" class="text-rose-600"><Histogram /></el-icon>
          </div>
        </div>
      </el-card>
      <el-card class="stat-card shadow-md hover:shadow-lg transition-shadow">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-gray-500 text-sm mb-2">活跃 Key 数</div>
            <div class="text-3xl font-bold text-orange-600">{{ stats.activeKeys }}</div>
          </div>
          <div class="w-12 h-12 bg-orange-100 rounded-lg flex items-center justify-center">
            <el-icon :size="24" class="text-orange-600"><Key /></el-icon>
          </div>
        </div>
      </el-card>
    </div>

    <el-card class="shadow-md">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-semibold text-lg">近7天请求趋势</span>
          <el-icon class="text-gray-400"><TrendCharts /></el-icon>
        </div>
      </template>
      <div ref="chartRef" class="h-80"></div>
    </el-card>

    <el-card class="shadow-md">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-semibold text-lg">Top 5 调用 Key 排行</span>
          <el-icon class="text-gray-400"><Rank /></el-icon>
        </div>
      </template>
      <el-table :data="topKeys" stripe class="modern-table">
        <el-table-column prop="name" label="Key 名称" />
        <el-table-column prop="key" label="Key" width="200">
          <template #default="{ row }">
            <span class="font-mono text-sm text-gray-600">{{ maskKey(row.key) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="callCount" label="调用次数" width="120">
          <template #default="{ row }">
            <el-tag type="success">{{ row.callCount.toLocaleString() }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { TrendCharts, DataLine, DataAnalysis, Key, Rank, Histogram } from '@element-plus/icons-vue'
import request from '@/api/request'
import type { Statistics, ApiKey } from '@/types'

const chartRef = ref<HTMLElement>()
const loading = ref(false)
let chartInstance: echarts.ECharts | null = null

const stats = ref<Statistics>({
  todayTotal: 0,
  yesterdayTotal: 0,
  monthTotal: 0,
  activeKeys: 0,
  totalAllTime: 0
})

const topKeys = ref<ApiKey[]>([])

const maskKey = (key: string) => {
  return `${key.slice(0, 6)}...${key.slice(-4)}`
}

const fetchStatistics = async () => {
  loading.value = true
  try {
    const res: any = await request.get('/api/statistics/overview')
    if (res.code === 0) {
      stats.value = res.data
    }
  } catch (e) {
    ElMessage.error('获取统计数据失败')
  } finally {
    loading.value = false
  }
}

const fetchTrend = async () => {
  try {
    const res: any = await request.get('/api/statistics/trend')
    if (res.code === 0) {
      const dates = res.data.map((item: any) => item.date)
      const counts = res.data.map((item: any) => item.count)

      await nextTick()
      if (!chartRef.value) return

      if (!chartInstance) {
        chartInstance = echarts.init(chartRef.value)
        window.addEventListener('resize', () => chartInstance?.resize())
      }

      chartInstance.setOption({
        tooltip: { trigger: 'axis' },
        grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
        xAxis: { type: 'category', data: dates, boundaryGap: false },
        yAxis: { type: 'value', minInterval: 1 },
        series: [{
          name: '请求次数',
          data: counts,
          type: 'line',
          smooth: true,
          areaStyle: {
            color: {
              type: 'linear',
              x: 0, y: 0, x2: 0, y2: 1,
              colorStops: [
                { offset: 0, color: 'rgba(59, 130, 246, 0.3)' },
                { offset: 1, color: 'rgba(59, 130, 246, 0.05)' }
              ]
            }
          },
          lineStyle: { color: '#3b82f6', width: 3 },
          itemStyle: { color: '#3b82f6' }
        }]
      })
    }
  } catch (e) {
    ElMessage.error('获取趋势数据失败')
  }
}

const fetchTopKeys = async () => {
  try {
    const res: any = await request.get('/api/statistics/top-keys')
    if (res.code === 0) {
      topKeys.value = res.data
    }
  } catch (e) {
    ElMessage.error('获取 Top Keys 失败')
  }
}

onMounted(() => {
  fetchStatistics()
  fetchTrend()
  fetchTopKeys()
})

onUnmounted(() => {
  chartInstance?.dispose()
  chartInstance = null
})
</script>
