<template>
  <main class="login">
    <el-card class="login__panel" shadow="never">
      <h1>TMS-Go</h1>
      <el-form label-position="top">
        <el-form-item label="用户名">
          <el-input model-value="admin" disabled />
        </el-form-item>
        <el-form-item label="密码">
          <el-input model-value="admin123" type="password" disabled />
        </el-form-item>
        <el-button type="primary" class="login__button" @click="$router.push('/')">
          进入工作台
        </el-button>
      </el-form>
    </el-card>
  </main>
</template>

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
</style>
