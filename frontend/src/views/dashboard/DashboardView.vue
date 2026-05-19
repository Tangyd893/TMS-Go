<template>
  <section class="dashboard">
    <div class="dashboard__header">
      <h1>TMS 工作台</h1>
      <p>运输管理系统运营总览</p>
    </div>

    <el-row :gutter="16" v-loading="loading">
      <el-col :span="6" v-for="item in metrics" :key="item.label">
        <el-card shadow="never" :class="`metric-card metric-card--${item.color}`" @click="item.link ? router.push(item.link) : null" :style="item.link ? 'cursor:pointer' : ''">
          <div class="metric">
            <span class="metric__label">{{ item.label }}</span>
            <strong class="metric__value" :style="{ color: `var(--el-color-${item.color})` }">{{ item.value }}</strong>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { getDashboard } from '@/api/modules/report'

const router = useRouter()
const loading = ref(false)

const metrics = reactive([
  { label: '待调度订单', value: 0, color: 'primary', link: '/dispatch/list' },
  { label: '运输中任务', value: 0, color: 'warning', link: '/transport/tasks' },
  { label: '待处理异常', value: 0, color: 'danger', link: '/exception/list' },
  { label: '待结算费用', value: 0, color: 'success', link: '/finance/receivables' },
])

onMounted(() => fetchData())

async function fetchData() {
  loading.value = true
  try {
    const data = await getDashboard()
    metrics[0].value = data.pendingDispatch
    metrics[1].value = data.inTransit
    metrics[2].value = data.pendingException
    metrics[3].value = data.pendingSettlement
  } finally { loading.value = false }
}
</script>

<style scoped>
.dashboard { display: flex; flex-direction: column; gap: 16px; }
.dashboard__header { display: flex; flex-direction: column; gap: 6px; }
.dashboard__header h1 { margin: 0; font-size: 24px; }
.dashboard__header p { margin: 0; color: #64748b; }
.metric { display: flex; flex-direction: column; gap: 8px; }
.metric__label { color: #64748b; font-size: 14px; }
.metric__value { font-size: 32px; font-weight: 700; }
.metric-card { border-top: 3px solid transparent; transition: box-shadow 0.2s; }
.metric-card:hover { box-shadow: 0 2px 12px rgba(0,0,0,0.1); }
.metric-card--primary { border-top-color: var(--el-color-primary); }
.metric-card--warning { border-top-color: var(--el-color-warning); }
.metric-card--danger { border-top-color: var(--el-color-danger); }
.metric-card--success { border-top-color: var(--el-color-success); }
</style>
