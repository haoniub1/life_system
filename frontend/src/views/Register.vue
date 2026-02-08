<template>
  <div class="register-container">
    <div class="stars-bg"></div>
    <n-card class="register-card" :segmented="true">
      <div class="register-header">
        <h1 class="title">开始冒险</h1>
        <p class="subtitle">创建你的英雄账号</p>
      </div>

      <n-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        @submit.prevent="handleRegister"
      >
        <n-form-item label="用户名" path="username">
          <n-input
            v-model:value="formData.username"
            placeholder="输入用户名"
            clearable
            size="large"
          />
        </n-form-item>

        <n-form-item label="密码" path="password">
          <n-input
            v-model:value="formData.password"
            type="password"
            placeholder="输入密码（最少6位）"
            clearable
            size="large"
            show-password-on="click"
          />
        </n-form-item>

        <n-form-item label="确认密码" path="confirmPassword">
          <n-input
            v-model:value="formData.confirmPassword"
            type="password"
            placeholder="再次输入密码"
            clearable
            size="large"
            show-password-on="click"
          />
        </n-form-item>

        <n-button type="primary" block size="large" :loading="loading" @click="handleRegister">
          注册
        </n-button>
      </n-form>

      <div class="register-footer">
        <span>已有账号？</span>
        <router-link to="/login" class="link">返回登录</router-link>
      </div>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { NCard, NForm, NFormItem, NInput, NButton } from 'naive-ui'
import { useUserStore } from '@/stores/user'
import type { FormInst, FormItemRule } from 'naive-ui'

const router = useRouter()
const message = useMessage()
const userStore = useUserStore()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)
const formData = ref({
  username: '',
  password: '',
  confirmPassword: ''
})

const validatePasswordLength = (_rule: FormItemRule, value: string) => {
  if (!value) {
    return new Error('请输入密码')
  } else if (value.length < 6) {
    return new Error('密码最少6位')
  }
  return true
}

const validatePasswordMatch = (_rule: FormItemRule, value: string) => {
  if (!value) {
    return new Error('请确认密码')
  } else if (value !== formData.value.password) {
    return new Error('两次输入的密码不一致')
  }
  return true
}

const rules = {
  username: [
    {
      required: true,
      message: '请输入用户名',
      trigger: 'blur'
    },
    {
      min: 3,
      message: '用户名最少3位',
      trigger: 'blur'
    }
  ],
  password: [
    {
      validator: validatePasswordLength,
      trigger: 'blur'
    }
  ],
  confirmPassword: [
    {
      validator: validatePasswordMatch,
      trigger: 'blur'
    }
  ]
}

const handleRegister = async () => {
  try {
    await formRef.value?.validate()
    loading.value = true

    await userStore.register(formData.value.username, formData.value.password)
    message.success('注册成功，欢迎来到人生RPG！')
    setTimeout(() => {
      router.push('/')
    }, 1000)
  } catch (error: any) {
    const errorMsg = error?.response?.data?.message || error?.message || '注册失败'
    message.error(errorMsg)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
  position: relative;
  width: 100%;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0f0f0f 0%, #1a1a2e 50%, #16213e 100%);
  overflow: hidden;
}

.stars-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image:
    radial-gradient(2px 2px at 20px 30px, #eee, rgba(0,0,0,0)),
    radial-gradient(2px 2px at 60px 70px, #fff, rgba(0,0,0,0)),
    radial-gradient(1px 1px at 50px 50px, #ddd, rgba(0,0,0,0)),
    radial-gradient(1px 1px at 130px 80px, #fff, rgba(0,0,0,0)),
    radial-gradient(2px 2px at 90px 10px, #fff, rgba(0,0,0,0));
  background-repeat: repeat;
  background-size: 200px 200px;
  animation: twinkle 5s infinite;
  opacity: 0.5;
  pointer-events: none;
}

@keyframes twinkle {
  0%, 100% { opacity: 0.3; }
  50% { opacity: 0.8; }
}

.register-card {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 400px;
  background: linear-gradient(135deg, rgba(30, 30, 50, 0.95) 0%, rgba(20, 20, 40, 0.95) 100%);
  border: 2px solid rgba(255, 215, 0, 0.3);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.6), inset 0 1px 0 rgba(255, 215, 0, 0.1);
  border-radius: 8px;
  padding: 40px;
}

.register-header {
  text-align: center;
  margin-bottom: 32px;
}

.title {
  font-size: 32px;
  font-weight: bold;
  background: linear-gradient(135deg, #ffd700, #ffed4e, #d4af37);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 14px;
  color: #a0a0b0;
}

:deep(.n-form-item) {
  margin-bottom: 20px;
}

:deep(.n-form-item__label) {
  color: #d0d0e0 !important;
  font-weight: 500;
}

:deep(.n-input__input) {
  background-color: rgba(255, 255, 255, 0.05) !important;
  color: #e0e0e0 !important;
  border-color: rgba(255, 215, 0, 0.2) !important;
}

:deep(.n-input__input::placeholder) {
  color: #707080 !important;
}

:deep(.n-input__input:hover) {
  border-color: rgba(255, 215, 0, 0.4) !important;
}

:deep(.n-input__input:focus) {
  border-color: #ffd700 !important;
  box-shadow: 0 0 8px rgba(255, 215, 0, 0.3) !important;
}

:deep(.n-button--primary) {
  background: linear-gradient(135deg, #ffd700, #ffed4e);
  color: #000 !important;
  font-weight: 600;
  border: none !important;
}

:deep(.n-button--primary:hover) {
  box-shadow: 0 4px 16px rgba(255, 215, 0, 0.4) !important;
  transform: translateY(-2px);
}

.register-footer {
  text-align: center;
  margin-top: 24px;
  font-size: 14px;
  color: #a0a0b0;
}

.link {
  color: #ffd700;
  text-decoration: none;
  font-weight: 600;
  margin-left: 4px;
  transition: all 0.3s;
}

.link:hover {
  color: #ffed4e;
  text-shadow: 0 0 8px rgba(255, 215, 0, 0.3);
}
</style>
