<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
    <el-card class="w-[420px] shadow-2xl">
      <template #header>
        <div class="text-center py-2">
          <div class="w-16 h-16 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-2xl mx-auto mb-4 flex items-center justify-center">
            <el-icon :size="32" class="text-white"><MapLocation /></el-icon>
          </div>
          <h2 class="text-2xl font-bold text-gray-800">地图瓦片管理系统</h2>
          <p class="text-gray-500 text-sm mt-2">欢迎回来，请登录您的账户</p>
        </div>
      </template>
      <el-form :model="loginForm" :rules="rules" ref="formRef" class="px-2">
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            size="large"
            class="login-input"
          >
            <template #prefix>
              <el-icon class="text-gray-400"><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            show-password
            class="login-input"
            @keyup.enter="handleLogin"
          >
            <template #prefix>
              <el-icon class="text-gray-400"><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item>
          <div class="flex justify-between w-full items-center">
            <el-checkbox v-model="loginForm.remember">记住我</el-checkbox>
          </div>
        </el-form-item>
        <el-form-item class="mb-2">
          <el-button
            type="primary"
            size="large"
            class="w-full login-btn"
            :loading="loading"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, MapLocation } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/api/request'
import CryptoJS from 'crypto-js'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: '',
  remember: false
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const encryptedPassword = CryptoJS.SHA256(loginForm.password).toString()
        const res: any = await request.post('/api/auth/login', {
          username: loginForm.username,
          password: encryptedPassword
        })

        if (res.code === 0) {
          localStorage.setItem('token', res.data.token)
          ElMessage.success('登录成功')
          router.push('/dashboard')
        } else {
          ElMessage.error(res.message || '登录失败')
        }
      } catch (error: any) {
        ElMessage.error(error.message || '登录失败，请检查网络连接')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>
