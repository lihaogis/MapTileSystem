<template>
  <div class="space-y-6">
    <!-- 操作栏 -->
    <el-card class="shadow-md">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-lg font-semibold text-gray-800">用户管理</h3>
          <p class="text-sm text-gray-500 mt-1">管理系统用户及其权限</p>
        </div>
        <el-button type="primary" :icon="Plus" @click="handleAdd" size="large">新增用户</el-button>
      </div>
    </el-card>

    <!-- 用户列表 -->
    <el-card class="shadow-md">
      <el-table :data="userList" stripe v-loading="loading" class="modern-table">
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="role" label="角色" width="120">
        <template #default="{ row }">
          <el-tag :type="row.role === 'admin' ? 'danger' : 'primary'">
            {{ roleLabel(row.role) }}
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
          <el-button link type="primary" size="small" @click="handleEdit(row)">
            <el-icon class="mr-1"><Edit /></el-icon>编辑
          </el-button>
          <el-button link type="warning" size="small" @click="handleResetPassword(row)">
            <el-icon class="mr-1"><RefreshRight /></el-icon>重置密码
          </el-button>
          <el-button
            link
            type="danger"
            size="small"
            @click="handleDelete(row)"
            :disabled="row.username === 'admin'"
          >
            <el-icon class="mr-1"><Delete /></el-icon>删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑用户' : '新增用户'"
      width="440px"
      @close="resetForm"
    >
      <el-form :model="form" :rules="formRules" ref="formRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" placeholder="请选择角色" class="w-full">
            <el-option
              v-for="r in roleOptions"
              :key="r.value"
              :label="r.label"
              :value="r.value"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 重置密码对话框 -->
    <el-dialog v-model="passwordDialogVisible" title="重置密码" width="440px" @close="resetPasswordForm">
      <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="80px">
        <el-form-item label="新密码" prop="password">
          <el-input v-model="passwordForm.password" type="password" placeholder="请输入新密码" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" placeholder="请再次输入新密码" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="passwordDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handlePasswordSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete, RefreshRight } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/api/request'

interface User {
  id: string
  username: string
  role: string
  createdAt: string
  updatedAt: string
}

const loading = ref(false)
const submitting = ref(false)
const userList = ref<User[]>([])

// 新增/编辑
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const form = ref({ id: '', username: '', password: '', role: 'operator' })

// 重置密码
const passwordDialogVisible = ref(false)
const passwordFormRef = ref<FormInstance>()
const passwordForm = ref({ id: '', password: '', confirmPassword: '' })

const roleOptions = [
  { label: '管理员', value: 'admin' },
  { label: '操作员', value: 'operator' },
  { label: '访客', value: 'viewer' }
]

const roleLabel = (role: string) => roleOptions.find(r => r.value === role)?.label ?? role

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const formRules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码不少于6位', trigger: 'blur' }
  ],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }]
}

const passwordRules: FormRules = {
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码不少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (_rule, value, callback) => {
        if (value !== passwordForm.value.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const fetchUsers = async () => {
  loading.value = true
  try {
    const res: any = await request.get('/api/users')
    if (res.code === 0) {
      userList.value = res.data
    }
  } catch (e) {
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  form.value = { id: '', username: '', password: '', role: 'operator' }
  dialogVisible.value = true
}

const handleEdit = (row: User) => {
  isEdit.value = true
  form.value = { id: row.id, username: row.username, password: '', role: row.role }
  dialogVisible.value = true
}

const handleDelete = async (row: User) => {
  await ElMessageBox.confirm(`确定删除用户 "${row.username}" 吗？`, '提示', {
    type: 'warning',
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  })
  try {
    const res: any = await request.delete(`/api/users/${row.id}`)
    if (res.code === 0) {
      ElMessage.success('删除成功')
      fetchUsers()
    } else {
      ElMessage.error(res.message || '删除失败')
    }
  } catch (e) {
    ElMessage.error('删除失败')
  }
}

const handleResetPassword = (row: User) => {
  passwordForm.value = { id: row.id, password: '', confirmPassword: '' }
  passwordDialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      const res: any = isEdit.value
        ? await request.put(`/api/users/${form.value.id}`, {
            username: form.value.username,
            role: form.value.role
          })
        : await request.post('/api/users', {
            username: form.value.username,
            password: form.value.password,
            role: form.value.role
          })
      if (res.code === 0) {
        ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
        dialogVisible.value = false
        fetchUsers()
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

const handlePasswordSubmit = async () => {
  if (!passwordFormRef.value) return
  await passwordFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      const res: any = await request.put(`/api/users/${passwordForm.value.id}/password`, {
        password: passwordForm.value.password
      })
      if (res.code === 0) {
        ElMessage.success('密码重置成功')
        passwordDialogVisible.value = false
      } else {
        ElMessage.error(res.message || '重置失败')
      }
    } catch (e) {
      ElMessage.error('重置失败')
    } finally {
      submitting.value = false
    }
  })
}

const resetForm = () => formRef.value?.resetFields()
const resetPasswordForm = () => passwordFormRef.value?.resetFields()

onMounted(fetchUsers)
</script>
