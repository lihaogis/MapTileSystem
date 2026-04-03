<template>
  <div class="space-y-6">
    <el-card class="shadow-md">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-lg font-semibold text-gray-800">数据源管理</h3>
          <p class="text-sm text-gray-500 mt-1">管理地图瓦片数据源</p>
        </div>
        <el-button type="primary" :icon="Plus" @click="handleAdd" size="large">新增数据源</el-button>
      </div>
    </el-card>

    <el-card class="shadow-md">
      <el-table :data="dataList" stripe v-loading="loading">
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag :type="row.type === 'xyz' ? 'success' : 'primary'">
              {{ row.type.toUpperCase() }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="format" label="格式" width="100" />
        <el-table-column prop="path" label="路径" />
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
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button link type="success" size="small" @click="handlePreview(row)">
              <el-icon class="mr-1"><View /></el-icon>预览
            </el-button>
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
      :title="isEdit ? '编辑数据源' : '新增数据源'"
      width="500px"
      @close="resetForm"
    >
      <el-form :model="form" :rules="formRules" ref="formRef" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入数据源名称" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择类型" class="w-full" @change="handleTypeChange">
            <el-option label="XYZ 栅格瓦片" value="xyz" />
            <el-option label="3D Tiles" value="3dtiles" />
          </el-select>
        </el-form-item>
        <el-form-item label="格式" prop="format" v-if="form.type === 'xyz'">
          <el-select v-model="form.format" placeholder="请选择格式" class="w-full">
            <el-option label="PNG" value="png" />
            <el-option label="JPG" value="jpg" />
            <el-option label="WebP" value="webp" />
          </el-select>
        </el-form-item>
        <el-form-item label="路径" prop="path">
          <div class="flex gap-2">
            <el-input
              v-model="form.path"
              :placeholder="form.type === 'xyz' ? '选择瓦片目录（包含 {z}/{x}/{y} 结构）' : '选择包含 tileset.json 的目录'"
            />
            <el-button @click="showPathDialog = true">浏览</el-button>
          </div>
          <div class="text-xs text-gray-500 mt-1">
            <span v-if="form.type === 'xyz'">XYZ 瓦片：选择包含 z/x/y 目录结构的根目录</span>
            <span v-else>3D Tiles：选择包含 tileset.json 文件的目录</span>
          </div>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio value="enabled">启用</el-radio>
            <el-radio value="disabled">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-divider content-position="left">地图预览配置（可选）</el-divider>
        <el-form-item label="中心纬度">
          <el-input-number v-model="form.centerLat" :precision="6" :step="0.1" placeholder="如：39.9" class="w-full" />
        </el-form-item>
        <el-form-item label="中心经度">
          <el-input-number v-model="form.centerLng" :precision="6" :step="0.1" placeholder="如：116.4" class="w-full" />
        </el-form-item>
        <el-form-item label="缩放级别">
          <el-input-number v-model="form.defaultZoom" :min="0" :max="18" placeholder="如：10" class="w-full" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="previewVisible" :title="`预览 - ${previewData?.name}`" width="80%" top="5vh">
      <div class="h-[70vh]">
        <MapPreview2D
          v-if="previewData?.type === 'xyz'"
          :url="previewUrl"
          :format="previewData?.format || 'png'"
          :center-lat="previewData?.centerLat"
          :center-lng="previewData?.centerLng"
          :default-zoom="previewData?.defaultZoom"
        />
        <MapPreview3D v-if="previewData?.type === '3dtiles'" :url="previewUrl" />
      </div>
    </el-dialog>

    <el-dialog v-model="showPathDialog" title="选择目录" width="600px">
      <div class="mb-4">
        <el-button @click="loadDrives" size="small">盘符</el-button>
        <el-button @click="goParent" size="small" :disabled="!currentPath">上级</el-button>
        <span class="ml-4 text-sm text-gray-600">{{ currentPath || '请选择盘符' }}</span>
      </div>
      <el-table :data="fileList" v-loading="fileLoading" height="400" @row-click="handleFileClick">
        <el-table-column label="名称">
          <template #default="{ row }">
            <el-icon v-if="row.isDir" class="mr-2"><Folder /></el-icon>
            <el-icon v-else class="mr-2"><Document /></el-icon>
            {{ row.name }}
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="showPathDialog = false">取消</el-button>
        <el-button type="primary" @click="selectPath">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, Folder, Document, View } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/api/request'
import type { DataSource } from '@/types'
import MapPreview2D from '@/components/MapPreview2D.vue'
import MapPreview3D from '@/components/MapPreview3D.vue'

const loading = ref(false)
const submitting = ref(false)
const dataList = ref<DataSource[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const form = ref({
  id: '',
  name: '',
  type: 'xyz',
  format: 'png',
  path: '',
  status: 'enabled',
  centerLat: undefined as number | undefined,
  centerLng: undefined as number | undefined,
  defaultZoom: undefined as number | undefined
})

const showPathDialog = ref(false)
const fileLoading = ref(false)
const currentPath = ref('')
const fileList = ref<any[]>([])

const previewVisible = ref(false)
const previewData = ref<DataSource | null>(null)
const previewUrl = computed(() => {
  if (!previewData.value) return ''
  const token = localStorage.getItem('token') || ''
  if (previewData.value.type === 'xyz') {
    // XYZ 瓦片使用 /api/preview/xyz 路径
    return `/api/preview/xyz/${previewData.value.id}/{z}/{x}/{y}?token=${token}`
  } else {
    // 3D Tiles 使用 /api/preview/3dtiles 路径
    return `/api/preview/3dtiles/${previewData.value.id}/tileset.json?token=${token}`
  }
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入数据源名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  format: [
    {
      required: true,
      message: '请选择格式',
      trigger: 'change',
      validator: (_rule, _value, callback) => {
        // 只有 XYZ 类型才需要验证格式
        if (form.value.type === 'xyz' && !form.value.format) {
          callback(new Error('请选择格式'))
        } else {
          callback()
        }
      }
    }
  ],
  path: [{ required: true, message: '请输入路径', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

// 处理类型切换
const handleTypeChange = (type: string) => {
  if (type === '3dtiles') {
    // 3D Tiles 不需要格式字段，清空格式值
    form.value.format = ''
  } else if (type === 'xyz' && !form.value.format) {
    // XYZ 类型默认 PNG 格式
    form.value.format = 'png'
  }
}

const fetchDataSources = async () => {
  loading.value = true
  try {
    const res: any = await request.get('/api/datasources')
    if (res.code === 0) {
      dataList.value = res.data.sort((a: DataSource, b: DataSource) =>
        new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
      )
    }
  } catch (e) {
    ElMessage.error('获取数据源列表失败')
  } finally {
    loading.value = false
  }
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleString('zh-CN')
}

const handleAdd = () => {
  isEdit.value = false
  form.value = {
    id: '',
    name: '',
    type: 'xyz',
    format: 'png',
    path: '',
    status: 'enabled',
    centerLat: undefined,
    centerLng: undefined,
    defaultZoom: undefined
  }
  dialogVisible.value = true
}

const handleEdit = (row: DataSource) => {
  isEdit.value = true
  form.value = { ...row }
  dialogVisible.value = true
}

const handlePreview = (row: DataSource) => {
  previewData.value = row
  previewVisible.value = true
}

const handleDelete = async (row: DataSource) => {
  await ElMessageBox.confirm(`确定删除数据源 "${row.name}" 吗？`, '提示', {
    type: 'warning',
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  })
  try {
    const res: any = await request.delete(`/api/datasources/${row.id}`)
    if (res.code === 0) {
      ElMessage.success('删除成功')
      fetchDataSources()
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
        ? await request.put(`/api/datasources/${form.value.id}`, form.value)
        : await request.post('/api/datasources', form.value)
      if (res.code === 0) {
        ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
        dialogVisible.value = false
        fetchDataSources()
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

const loadDrives = async () => {
  fileLoading.value = true
  try {
    const res: any = await request.get('/api/files/drives')
    if (res.code === 0) {
      fileList.value = res.data
      currentPath.value = ''
    }
  } catch (e) {
    ElMessage.error('获取盘符失败')
  } finally {
    fileLoading.value = false
  }
}

const loadDirectory = async (path: string) => {
  fileLoading.value = true
  try {
    const res: any = await request.get('/api/files/directories', { params: { path } })
    if (res.code === 0) {
      fileList.value = res.data.filter((f: any) => f.isDir)
      currentPath.value = path
    } else {
      ElMessage.error(res.message || '读取目录失败')
    }
  } catch (e) {
    ElMessage.error('读取目录失败')
  } finally {
    fileLoading.value = false
  }
}

const handleFileClick = (row: any) => {
  if (row.isDir) {
    loadDirectory(row.path)
  }
}

const goParent = () => {
  if (!currentPath.value) return
  const parent = currentPath.value.split('\\').slice(0, -1).join('\\')
  if (parent) {
    loadDirectory(parent)
  } else {
    loadDrives()
  }
}

const selectPath = () => {
  if (currentPath.value) {
    form.value.path = currentPath.value
    showPathDialog.value = false
  } else {
    ElMessage.warning('请选择一个目录')
  }
}

onMounted(fetchDataSources)
</script>
