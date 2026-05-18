<template>
  <el-container class="layout">
    <el-aside class="layout__aside" width="240px">
      <div class="layout__brand">TMS-Go</div>
      <el-menu router default-active="/">
        <el-menu-item index="/">
          <span>工作台</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="layout__header">
        <span>运输管理系统</span>
      </el-header>
      <el-main class="layout__main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.layout {
  min-height: 100vh;
}

.layout__aside {
  border-right: 1px solid #e5e7eb;
  background: #ffffff;
}

.layout__brand {
  height: 56px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  font-size: 18px;
  font-weight: 700;
}

.layout__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #e5e7eb;
  background: #ffffff;
}

.layout__main {
  background: #f6f8fb;
}
</style>
