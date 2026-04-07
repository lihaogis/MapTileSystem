<template>
  <div class="space-y-6">
    <el-card class="shadow-md">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-lg font-semibold text-gray-800">访问日志</h3>
          <p class="text-sm text-gray-500 mt-1">查看瓦片服务调用日志</p>
        </div>
        <el-button :icon="Refresh" @click="fetchLogs" :loading="loading">刷新</el-button>
      </div>
    </el-card>

    <el-card class="shadow-md">
      <el-form :inline="true" :model="queryForm" class="mb-4">
        <el-form-item label="API Key">
          <el-select v-model="queryForm.apiKeyId" placeholder="全部" clearable style="width: 180px">
            <el-option v-for="key in apiKeys" :key="key.id" :label="key.name" :value="key.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="数据源">
          <el-select v-model="queryForm.dataSourceId" placeholder="全部" clearable style="width: 180px">
            <el-option v-for="ds in dataSources" :key="ds.id" :label="ds.name" :value="ds.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="IP地址">
          <el-input v-model="queryForm.ipAddress" placeholder="请输入IP" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="请求时间">
          <el-date-picker
            v-model="queryForm.dateRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            style="width: 360px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleQuery">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
      <el-table :data="logList" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="API Key" width="150">
          <template #default="{ row }">
            <el-text class="text-sm" truncated>{{ row.apiKeyName || '-' }}</el-text>
          </template>
        </el-table-column>
        <el-table-column label="数据源" width="150">
          <template #default="{ row }">
            <el-text class="text-sm" truncated>{{ row.dataSourceName || '-' }}</el-text>
          </template>
        </el-table-column>
        <el-table-column label="请求类型" width="100">
          <template #default="{ row }">
            <el-tag v-if="getDataSourceType(row.dataSourceId) === '3dtiles'" type="warning" size="small">3D Tiles</el-tag>
            <el-tag v-else-if="getDataSourceType(row.dataSourceId) === 'vector'" type="success" size="small">Vector</el-tag>
            <el-tag v-else type="info" size="small">XYZ</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="瓦片坐标" width="150">
          <template #default="{ row }">
            <span class="font-mono text-xs" v-if="getDataSourceType(row.dataSourceId) === '3dtiles'">
              tileset.json
            </span>
            <span class="font-mono text-xs" v-else>
              {{ row.tileZ }}/{{ row.tileX }}/{{ row.tileY }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="状态码" width="100">
          <template #default="{ row }">
            <el-tag :type="row.statusCode === 200 ? 'success' : 'danger'" size="small">
              {{ row.statusCode }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="responseTime" label="响应时间(ms)" width="120" />
        <el-table-column prop="ipAddress" label="IP地址" width="140" />
        <el-table-column label="User-Agent" min-width="200">
          <template #default="{ row }">
            <el-text class="text-xs" truncated>{{ row.userAgent }}</el-text>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="请求时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
      </el-table>

      <div class="mt-4 flex justify-end">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import request from '@/api/request'
import type { CallLog } from '@/types'

const loading = ref(false)
const logList = ref<CallLog[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const apiKeys = ref<any[]>([])
const dataSources = ref<any[]>([])

const queryForm = reactive({
  apiKeyId: '',
  dataSourceId: '',
  ipAddress: '',
  dateRange: null as any
})

const fetchLogs = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      pageSize: pageSize.value
    }
    if (queryForm.apiKeyId) params.apiKeyId = queryForm.apiKeyId
    if (queryForm.dataSourceId) params.dataSourceId = queryForm.dataSourceId
    if (queryForm.ipAddress) params.ipAddress = queryForm.ipAddress
    if (queryForm.dateRange && queryForm.dateRange.length === 2) {
      params.startDate = queryForm.dateRange[0].toISOString()
      params.endDate = queryForm.dateRange[1].toISOString()
    }

    const res: any = await request.get('/api/statistics/details', { params })
    if (res.code === 0) {
      logList.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (e) {
    ElMessage.error('获取访问记录失败')
  } finally {
    loading.value = false
  }
}

const fetchApiKeys = async () => {
  try {
    const res: any = await request.get('/api/apikeys')
    if (res.code === 0) {
      apiKeys.value = res.data || []
    }
  } catch (e) {
    console.error('获取API Key列表失败', e)
  }
}

const fetchDataSources = async () => {
  try {
    const res: any = await request.get('/api/datasources')
    if (res.code === 0) {
      dataSources.value = res.data || []
    }
  } catch (e) {
    console.error('获取数据源列表失败', e)
  }
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
  fetchLogs()
}

const handlePageChange = (val: number) => {
  currentPage.value = val
  fetchLogs()
}

const handleQuery = () => {
  currentPage.value = 1
  fetchLogs()
}

const handleReset = () => {
  queryForm.apiKeyId = ''
  queryForm.dataSourceId = ''
  queryForm.ipAddress = ''
  queryForm.dateRange = null
  currentPage.value = 1
  fetchLogs()
}

const getDataSourceType = (id: string) => {
  return dataSources.value.find((ds: any) => ds.id === id)?.type ?? 'xyz'
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchApiKeys()
  fetchDataSources()
  fetchLogs()
})
</script>
