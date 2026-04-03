<template>
  <div class="space-y-6">
    <el-card class="shadow-md">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-lg font-semibold text-gray-800">授权管理</h3>
          <p class="text-sm text-gray-500 mt-1">管理 API Key 及其访问权限</p>
        </div>
        <el-button type="primary" :icon="Plus" @click="handleAdd" size="large">新增 API Key</el-button>
      </div>
    </el-card>

    <el-card class="shadow-md">
      <el-table :data="keyList" stripe v-loading="loading">
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="key" label="Key" width="300">
          <template #default="{ row }">
            <div class="flex items-center gap-2">
              <span class="font-mono text-sm">{{ maskKey(row.key) }}</span>
              <el-button link type="primary" size="small" @click="copyKey(row.key)">
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="dataSources" label="关联数据源" width="150">
          <template #default="{ row }">
            <el-tag v-if="row.dataSources && row.dataSources.length > 0" size="small">
              {{ row.dataSources.length }} 个
            </el-tag>
            <span v-else class="text-gray-400">未关联</span>
          </template>
        </el-table-column>
        <el-table-column prop="url" label="瓦片URL" min-width="250">
          <template #default="{ row }">
            <div v-if="row.dataSources && row.dataSources.length > 0">
              <el-popover placement="bottom" :width="400" trigger="click">
                <template #reference>
                  <el-button link type="primary" size="small">
                    查看URL ({{ row.dataSources.length }}个)
                  </el-button>
                </template>
                <div class="space-y-2">
                  <div v-for="dsId in row.dataSources" :key="dsId" class="border-b pb-2">
                    <div class="text-xs text-gray-500 mb-1">{{ getDataSourceName(dsId) }}</div>
                    <div class="flex items-center gap-2">
                      <code class="flex-1 text-xs bg-gray-50 p-1 rounded break-all">
                        {{ getTileUrl(row.key, dsId) }}
                      </code>
                      <el-button link type="primary" size="small" @click="copyUrl(getTileUrl(row.key, dsId))">
                        <el-icon><CopyDocument /></el-icon>
                      </el-button>
                    </div>
                  </div>
                </div>
              </el-popover>
            </div>
            <span v-else class="text-gray-400 text-sm">未关联数据源</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'enabled' ? 'success' : 'info'">
              {{ row.status === 'enabled' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleEdit(row)">
              <el-icon class="mr-1"><Edit /></el-icon>编辑
            </el-button>
            <el-button link type="danger" size="small" @click="handleDelete(row)">
              <el-icon class="mr-1"><Delete /></el-icon>删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑 API Key' : '新增 API Key'"
      width="500px"
      @close="resetForm"
    >
      <el-form :model="form" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入 API Key 名称" />
        </el-form-item>
        <el-form-item label="关联数据源" prop="dataSources">
          <el-select v-model="form.dataSources" multiple placeholder="请选择数据源" class="w-full">
            <el-option
              v-for="ds in dataSourceList"
              :key="ds.id"
              :label="ds.name"
              :value="ds.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio value="enabled">启用</el-radio>
            <el-radio value="disabled">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-alert
          v-if="!isEdit"
          title="API Key 将在创建后显示，请妥善保管"
          type="warning"
          :closable="false"
          class="mb-4"
        />
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="keyDisplayVisible" title="API Key 创建成功" width="500px">
      <el-alert
        title="请复制并保存 API Key，关闭后将无法再次查看完整密钥"
        type="success"
        :closable="false"
        class="mb-4"
      />
      <div class="bg-gray-50 p-4 rounded">
        <div class="text-sm text-gray-600 mb-2">API Key:</div>
        <div class="flex items-center gap-2">
          <code class="flex-1 bg-white p-2 rounded border text-sm break-all">{{ newApiKey }}</code>
          <el-button type="primary" @click="copyKey(newApiKey)">
            <el-icon><CopyDocument /></el-icon>
          </el-button>
        </div>
      </div>
      <template #footer>
        <el-button type="primary" @click="keyDisplayVisible = false">我已保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, CopyDocument } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/api/request'
import type { ApiKey, DataSource } from '@/types'

const loading = ref(false)
const submitting = ref(false)
const keyList = ref<ApiKey[]>([])
const dataSourceList = ref<DataSource[]>([])
const dialogVisible = ref(false)
const keyDisplayVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const newApiKey = ref('')
const form = ref({
  id: '',
  name: '',
  dataSources: [] as string[],
  status: 'enabled'
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入 API Key 名称', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

const maskKey = (key: string) => {
  if (!key || key.length < 10) return key
  return `${key.slice(0, 8)}...${key.slice(-4)}`
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}

const copyKey = (key: string) => {
  navigator.clipboard.writeText(key)
  ElMessage.success('已复制到剪贴板')
}

const copyUrl = (url: string) => {
  navigator.clipboard.writeText(url)
  ElMessage.success('已复制到剪贴板')
}

const getTileUrl = (key: string, dataSourceId: string) => {
  const ds = dataSourceList.value.find(d => d.id === dataSourceId)
  if (ds?.type === '3dtiles') {
    return `http://localhost:8080/tiles/${dataSourceId}/tileset.json?key=${key}`
  }
  return `http://localhost:8080/tiles/${dataSourceId}/{z}/{x}/{y}?key=${key}`
}

const getDataSourceName = (id: string) => {
  const ds = dataSourceList.value.find(d => d.id === id)
  return ds ? ds.name : id
}

const fetchApiKeys = async () => {
  loading.value = true
  try {
    const res: any = await request.get('/api/apikeys')
    if (res.code === 0) {
      keyList.value = res.data
    }
  } catch (e) {
    ElMessage.error('获取 API Key 列表失败')
  } finally {
    loading.value = false
  }
}

const fetchDataSources = async () => {
  try {
    const res: any = await request.get('/api/datasources')
    if (res.code === 0) {
      dataSourceList.value = res.data
    }
  } catch (e) {
    console.error('获取数据源列表失败', e)
  }
}

const handleAdd = () => {
  isEdit.value = false
  form.value = { id: '', name: '', dataSources: [], status: 'enabled' }
  dialogVisible.value = true
}

const handleEdit = (row: ApiKey) => {
  isEdit.value = true
  form.value = {
    id: row.id,
    name: row.name,
    dataSources: row.dataSources || [],
    status: row.status
  }
  dialogVisible.value = true
}

const handleDelete = async (row: ApiKey) => {
  await ElMessageBox.confirm(`确定删除 API Key "${row.name}" 吗？`, '提示', {
    type: 'warning',
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  })
  try {
    const res: any = await request.delete(`/api/apikeys/${row.id}`)
    if (res.code === 0) {
      ElMessage.success('删除成功')
      fetchApiKeys()
    } else {
      ElMessage.error(res.message || '删除失败')
    }
  } catch (e) {
    ElMessage.error('删除失败')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      const res: any = isEdit.value
        ? await request.put(`/api/apikeys/${form.value.id}`, form.value)
        : await request.post('/api/apikeys', form.value)
      if (res.code === 0) {
        ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
        dialogVisible.value = false
        if (!isEdit.value && res.data.key) {
          newApiKey.value = res.data.key
          keyDisplayVisible.value = true
        }
        fetchApiKeys()
      } else {
        ElMessage.error(res.message || '操作失败')
      }
    } catch (e) {
      ElMessage.error('操作失败')
    } finally {
      submitting.value = false
    }
  })
}

const resetForm = () => formRef.value?.resetFields()

onMounted(() => {
  fetchApiKeys()
  fetchDataSources()
})
</script>
