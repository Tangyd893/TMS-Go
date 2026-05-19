<template>
  <main class="login">
    <el-card class="login__panel" shadow="never">
      <h1>TMS-Go</h1>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="handleLogin">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-button type="primary" class="login__button" :loading="loading" native-type="submit">
          {{ loading ? '登录中...' : '登录' }}
        </el-button>
      </el-form>
      <p v-if="errorMsg" class="login__error">{{ errorMsg }}</p>
    </el-card>
  </main>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { login } from '@/api/modules/auth'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'

const router = useRouter()
const authStore = useAuthStore()
const permissionStore = usePermissionStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const errorMsg = ref('')

const form = reactive({
  username: '',
  password: '',
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

async function handleLogin() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  errorMsg.value = ''

  try {
    const result = await login({ username: form.username, password: form.password })
    authStore.setToken(result.accessToken)
    authStore.setUser(result.user)
    permissionStore.setPermissions(result.user.permissions)
    permissionStore.setRoles(result.user.roles)
    router.push('/')
  } catch (err: any) {
    errorMsg.value = err?.message || '登录失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login {
  min-height: 100vh;
  display: grid;
  place-items: center;
  background: #f6f8fb;
}

.login__panel {
  width: min(420px, calc(100vw - 32px));
}

.login__panel h1 {
  margin: 0 0 24px;
  font-size: 24px;
}

.login__button {
  width: 100%;
}

.login__error {
  color: #f56c6c;
  margin: 12px 0 0;
  text-align: center;
}
</style>
