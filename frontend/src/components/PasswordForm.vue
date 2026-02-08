<template>
  <div class="password-section">
    <n-form ref="pwdFormRef" :model="pwdForm" :rules="pwdRules">
      <n-form-item label="当前密码" path="oldPassword">
        <n-input
          v-model:value="pwdForm.oldPassword"
          type="password"
          show-password-on="click"
          placeholder="输入当前密码"
        />
      </n-form-item>

      <n-form-item label="新密码" path="newPassword">
        <n-input
          v-model:value="pwdForm.newPassword"
          type="password"
          show-password-on="click"
          placeholder="输入新密码（至少6位）"
        />
      </n-form-item>

      <n-form-item label="确认新密码" path="confirmPassword">
        <n-input
          v-model:value="pwdForm.confirmPassword"
          type="password"
          show-password-on="click"
          placeholder="再次输入新密码"
        />
      </n-form-item>

      <n-button
        type="warning"
        block
        :loading="changingPassword"
        @click="handleChangePassword"
      >
        修改密码
      </n-button>
    </n-form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { NForm, NFormItem, NInput, NButton, type FormInst, type FormRules } from 'naive-ui'
import { changePassword } from '@/api/user'

const message = useMessage()
const changingPassword = ref(false)
const pwdFormRef = ref<FormInst | null>(null)

const pwdForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const pwdRules: FormRules = {
  oldPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string) => {
        if (value !== pwdForm.value.newPassword) {
          return new Error('两次输入的密码不一致')
        }
        return true
      },
      trigger: 'blur'
    }
  ]
}

const handleChangePassword = async () => {
  try {
    await pwdFormRef.value?.validate()
    changingPassword.value = true

    await changePassword({
      oldPassword: pwdForm.value.oldPassword,
      newPassword: pwdForm.value.newPassword
    })

    message.success('密码修改成功')
    pwdForm.value = { oldPassword: '', newPassword: '', confirmPassword: '' }
  } catch (error: any) {
    if (error?.message) {
      message.error(error.message)
    }
  } finally {
    changingPassword.value = false
  }
}
</script>

<style scoped>
.password-section {
  padding: 8px 0;
}

:deep(.n-button--warning) {
  background: linear-gradient(135deg, #f59e0b, #fbbf24) !important;
  color: #000 !important;
  border: none !important;
}
</style>
