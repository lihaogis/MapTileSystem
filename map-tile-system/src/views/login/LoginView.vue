<template>
  <div class="login-container">
    <!-- 地图网格背景 -->
    <div class="map-grid">
      <div class="grid-lines"></div>
    </div>

    <!-- 瓦片坐标装饰 -->
    <div class="tile-coords">
      <span v-for="coord in tileCoords" :key="coord.id" :style="coord.style" class="coord-tag">
        {{ coord.text }}
      </span>
    </div>

    <!-- 扫描线动画 -->
    <div class="scan-line"></div>

    <!-- 登录卡片 -->
    <div class="login-card">
      <!-- 顶部标识 -->
      <div class="card-header">
        <div class="logo-wrap">
          <div class="logo-inner">
            <svg viewBox="0 0 32 32" class="logo-icon" fill="none">
              <!-- 瓦片网格图标 -->
              <rect x="2" y="2" width="12" height="12" rx="1.5" fill="currentColor" opacity="0.9"/>
              <rect x="18" y="2" width="12" height="12" rx="1.5" fill="currentColor" opacity="0.6"/>
              <rect x="2" y="18" width="12" height="12" rx="1.5" fill="currentColor" opacity="0.6"/>
              <rect x="18" y="18" width="12" height="12" rx="1.5" fill="currentColor" opacity="0.3"/>
              <!-- 中心十字 -->
              <line x1="16" y1="6" x2="16" y2="26" stroke="white" stroke-width="1" opacity="0.4"/>
              <line x1="6" y1="16" x2="26" y2="16" stroke="white" stroke-width="1" opacity="0.4"/>
            </svg>
          </div>
          <div class="logo-pulse"></div>
        </div>
        <h1 class="system-title">地图瓦片管理系统</h1>
        <p class="system-subtitle">Map Tile Management System</p>
      </div>

      <!-- 登录表单 -->
      <div class="form-wrap">
        <el-form :model="loginForm" :rules="rules" ref="formRef">
          <el-form-item prop="username">
            <div class="input-wrap">
              <span class="input-label">用户名</span>
              <el-input
                v-model="loginForm.username"
                placeholder="admin"
                class="map-input"
              >
                <template #prefix>
                  <el-icon class="input-icon"><User /></el-icon>
                </template>
              </el-input>
            </div>
          </el-form-item>

          <el-form-item prop="password">
            <div class="input-wrap">
              <span class="input-label">密码</span>
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="••••••••"
                show-password
                class="map-input"
                @keyup.enter="handleLogin"
              >
                <template #prefix>
                  <el-icon class="input-icon"><Lock /></el-icon>
                </template>
              </el-input>
            </div>
          </el-form-item>

          <div class="form-options">
            <el-checkbox v-model="loginForm.remember" class="remember-check">记住登录</el-checkbox>
          </div>

          <el-button
            class="login-btn"
            :loading="loading"
            @click="handleLogin"
          >
            <span v-if="!loading">进入系统</span>
            <span v-else>验证中...</span>
          </el-button>
        </el-form>
      </div>

      <!-- 底部装饰信息 -->
      <div class="card-footer">
        <div class="footer-stat">
          <span class="stat-dot"></span>
          <span>XYZ 瓦片服务</span>
        </div>
        <div class="footer-divider"></div>
        <div class="footer-stat">
          <span class="stat-dot stat-dot-3d"></span>
          <span>3D Tiles 服务</span>
        </div>
      </div>
    </div>

    <!-- 右下角坐标信息 -->
    <div class="corner-info">
      <span>z/x/y · Tile Server v1.0</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/api/request'
import CryptoJS from 'crypto-js'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

// 随机生成瓦片坐标装饰标签
const tileCoords = Array.from({ length: 12 }, (_, i) => ({
  id: i,
  text: `${Math.floor(Math.random() * 18 + 3)}/${Math.floor(Math.random() * 256)}/${Math.floor(Math.random() * 256)}`,
  style: {
    left: `${Math.random() * 90}%`,
    top: `${Math.random() * 90}%`,
    animationDelay: `${Math.random() * 4}s`,
    animationDuration: `${Math.random() * 3 + 3}s`,
  }
}))

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

<style scoped>
/* ===== 基础容器 ===== */
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #0d1117;
  position: relative;
  overflow: hidden;
  font-family: 'SF Pro Display', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

/* ===== 地图网格背景 ===== */
.map-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(20, 184, 166, 0.06) 1px, transparent 1px),
    linear-gradient(90deg, rgba(20, 184, 166, 0.06) 1px, transparent 1px);
  background-size: 48px 48px;
}

.grid-lines {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(20, 184, 166, 0.12) 1px, transparent 1px),
    linear-gradient(90deg, rgba(20, 184, 166, 0.12) 1px, transparent 1px);
  background-size: 192px 192px;
}

/* ===== 扫描线 ===== */
.scan-line {
  position: absolute;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, rgba(20, 184, 166, 0.4), transparent);
  animation: scan 6s linear infinite;
  pointer-events: none;
}

@keyframes scan {
  0% { top: -2px; opacity: 0; }
  5% { opacity: 1; }
  95% { opacity: 1; }
  100% { top: 100%; opacity: 0; }
}

/* ===== 瓦片坐标装饰 ===== */
.tile-coords {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.coord-tag {
  position: absolute;
  font-size: 11px;
  font-family: 'SF Mono', 'Fira Code', monospace;
  color: rgba(20, 184, 166, 0.25);
  animation: fadeInOut linear infinite;
  white-space: nowrap;
}

@keyframes fadeInOut {
  0%, 100% { opacity: 0; transform: translateY(4px); }
  20%, 80% { opacity: 1; transform: translateY(0); }
}

/* ===== 登录卡片 ===== */
.login-card {
  position: relative;
  z-index: 10;
  width: 400px;
  background: rgba(15, 20, 30, 0.85);
  border: 1px solid rgba(20, 184, 166, 0.2);
  border-radius: 16px;
  backdrop-filter: blur(20px);
  box-shadow:
    0 0 0 1px rgba(20, 184, 166, 0.1),
    0 24px 64px rgba(0, 0, 0, 0.6),
    0 0 80px rgba(20, 184, 166, 0.05);
  overflow: hidden;
}

.login-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(20, 184, 166, 0.5), transparent);
}

/* ===== 卡片头部 ===== */
.card-header {
  padding: 36px 32px 24px;
  text-align: center;
}

.logo-wrap {
  position: relative;
  width: 64px;
  height: 64px;
  margin: 0 auto 20px;
}

.logo-inner {
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #0d9488, #0891b2);
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  z-index: 1;
}

.logo-icon {
  width: 36px;
  height: 36px;
  color: white;
}

.logo-pulse {
  position: absolute;
  inset: -4px;
  border-radius: 18px;
  border: 1px solid rgba(20, 184, 166, 0.4);
  animation: pulse 2.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.4; transform: scale(1); }
  50% { opacity: 0.8; transform: scale(1.05); }
}

.system-title {
  font-size: 20px;
  font-weight: 700;
  color: #e2e8f0;
  letter-spacing: 0.02em;
  margin: 0 0 6px;
}

.system-subtitle {
  font-size: 12px;
  color: rgba(20, 184, 166, 0.7);
  letter-spacing: 0.12em;
  text-transform: uppercase;
  margin: 0;
  font-family: 'SF Mono', monospace;
}

/* ===== 表单区域 ===== */
.form-wrap {
  padding: 0 32px;
}

.input-wrap {
  width: 100%;
}

.input-label {
  display: block;
  font-size: 12px;
  color: rgba(148, 163, 184, 0.8);
  margin-bottom: 6px;
  font-weight: 500;
  letter-spacing: 0.05em;
}

/* Element Plus 输入框深色主题覆盖 */
:deep(.map-input .el-input__wrapper) {
  background: rgba(30, 41, 59, 0.8) !important;
  border: 1px solid rgba(51, 65, 85, 0.8) !important;
  border-radius: 8px !important;
  box-shadow: none !important;
  transition: border-color 0.2s, box-shadow 0.2s;
  padding: 0 12px;
  height: 44px;
}

:deep(.map-input .el-input__wrapper:hover) {
  border-color: rgba(20, 184, 166, 0.5) !important;
}

:deep(.map-input .el-input__wrapper.is-focus) {
  border-color: rgba(20, 184, 166, 0.8) !important;
  box-shadow: 0 0 0 3px rgba(20, 184, 166, 0.12) !important;
}

:deep(.map-input .el-input__inner) {
  color: #e2e8f0 !important;
  font-size: 14px;
  background: transparent !important;
}

:deep(.map-input .el-input__inner::placeholder) {
  color: rgba(100, 116, 139, 0.7) !important;
}

.input-icon {
  color: rgba(20, 184, 166, 0.7);
}

:deep(.el-form-item) {
  margin-bottom: 16px;
}

:deep(.el-form-item__error) {
  color: #f87171;
  font-size: 12px;
}

/* ===== 表单选项 ===== */
.form-options {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

:deep(.remember-check .el-checkbox__label) {
  color: rgba(148, 163, 184, 0.8);
  font-size: 13px;
}

:deep(.remember-check .el-checkbox__inner) {
  background: rgba(30, 41, 59, 0.8);
  border-color: rgba(51, 65, 85, 0.8);
}

:deep(.remember-check.is-checked .el-checkbox__inner) {
  background: #0d9488;
  border-color: #0d9488;
}

/* ===== 登录按钮 ===== */
.login-btn {
  width: 100%;
  height: 46px;
  background: linear-gradient(135deg, #0d9488, #0891b2) !important;
  border: none !important;
  border-radius: 8px !important;
  color: white !important;
  font-size: 15px;
  font-weight: 600;
  letter-spacing: 0.05em;
  cursor: pointer;
  transition: opacity 0.2s, transform 0.15s, box-shadow 0.2s;
  box-shadow: 0 4px 16px rgba(13, 148, 136, 0.35);
  margin-bottom: 4px;
}

.login-btn:hover {
  opacity: 0.92;
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(13, 148, 136, 0.45) !important;
}

.login-btn:active {
  transform: translateY(0);
}

/* ===== 卡片底部 ===== */
.card-footer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 16px 32px 24px;
  margin-top: 8px;
}

.footer-stat {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: rgba(100, 116, 139, 0.7);
}

.stat-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #10b981;
  animation: blink 2s ease-in-out infinite;
}

.stat-dot-3d {
  background: #0891b2;
  animation-delay: 1s;
}

@keyframes blink {
  0%, 100% { opacity: 0.5; }
  50% { opacity: 1; }
}

.footer-divider {
  width: 1px;
  height: 14px;
  background: rgba(51, 65, 85, 0.8);
}

/* ===== 右下角版本信息 ===== */
.corner-info {
  position: absolute;
  bottom: 20px;
  right: 24px;
  font-size: 11px;
  font-family: 'SF Mono', monospace;
  color: rgba(20, 184, 166, 0.3);
  letter-spacing: 0.08em;
  pointer-events: none;
}
</style>
