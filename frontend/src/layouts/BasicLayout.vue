<template>
  <el-container class="layout">
    <el-aside class="layout__aside" width="240px">
      <div class="layout__brand">TMS-Go</div>
      <el-menu router :default-active="activeMenu" class="layout__menu">
        <el-menu-item index="/dashboard">
          <el-icon><Monitor /></el-icon>
          <span>工作台</span>
        </el-menu-item>

        <el-sub-menu index="base">
          <template #title><el-icon><FolderOpened /></el-icon><span>基础资料</span></template>
          <el-menu-item index="/base/customers">客户管理</el-menu-item>
          <el-menu-item index="/base/carriers">承运商管理</el-menu-item>
          <el-menu-item index="/base/vehicles">车辆管理</el-menu-item>
          <el-menu-item index="/base/drivers">司机管理</el-menu-item>
          <el-menu-item index="/base/routes">线路管理</el-menu-item>
          <el-menu-item index="/base/stations">站点管理</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="order">
          <template #title><el-icon><Document /></el-icon><span>订单管理</span></template>
          <el-menu-item index="/order/list">订单列表</el-menu-item>
          <el-menu-item index="/order/create">创建订单</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="dispatch">
          <template #title><el-icon><Tickets /></el-icon><span>调度管理</span></template>
          <el-menu-item index="/dispatch/plans">调度计划</el-menu-item>
          <el-menu-item index="/dispatch/create">创建调度</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="transport">
          <template #title><el-icon><Van /></el-icon><span>运输执行</span></template>
          <el-menu-item index="/transport/tasks">运输任务</el-menu-item>
        </el-sub-menu>

        <el-menu-item index="/exception/list">
          <el-icon><Warning /></el-icon>
          <span>异常管理</span>
        </el-menu-item>

        <el-sub-menu index="finance">
          <template #title><el-icon><Money /></el-icon><span>费用结算</span></template>
          <el-menu-item index="/finance/receivables">应收管理</el-menu-item>
          <el-menu-item index="/finance/payables">应付管理</el-menu-item>
          <el-menu-item index="/finance/statements">对账管理</el-menu-item>
          <el-menu-item index="/finance/settlements">结算管理</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="report">
          <template #title><el-icon><DataAnalysis /></el-icon><span>报表中心</span></template>
          <el-menu-item index="/report/order-stats">订单统计</el-menu-item>
          <el-menu-item index="/report/transport-efficiency">运输效率</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="system">
          <template #title><el-icon><Setting /></el-icon><span>系统管理</span></template>
          <el-menu-item index="/system/dict">字典管理</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="layout__header">
        <span>运输管理系统</span>
        <div class="layout__user">
          <el-dropdown @command="handleCommand">
            <span class="layout__username">
              {{ authStore.username || '管理员' }}
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main class="layout__main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Monitor, ArrowDown, FolderOpened, Setting, Document, Tickets, Van, Warning, Money, DataAnalysis } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const permissionStore = usePermissionStore()

const activeMenu = computed(() => route.path)

function handleCommand(command: string) {
  if (command === 'logout') {
    authStore.clearAll()
    permissionStore.clearAll()
    router.push('/login')
  }
}
</script>

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
  border-bottom: 1px solid #e5e7eb;
}

.layout__menu {
  border-right: none;
}

.layout__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #e5e7eb;
  background: #ffffff;
}

.layout__user {
  display: flex;
  align-items: center;
}

.layout__username {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  color: #606266;
}

.layout__main {
  background: #f6f8fb;
}
</style>
