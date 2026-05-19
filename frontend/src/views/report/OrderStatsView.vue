<template>
  <div class="report-page">
    <h3>订单统计</h3>
    <el-card v-loading="loading" class="report-page__card" shadow="never">
      <el-statistic title="订单总数" :value="stats.totalOrders" />
    </el-card>
    <el-card class="report-page__card" shadow="never" style="margin-top:16px">
      <template #header>按状态分布</template>
      <el-table :data="stats.byStatus" stripe border v-if="stats.byStatus.length">
        <el-table-column prop="status" label="状态" width="200">
          <template #default="{ row }"><el-tag :type="statusTag(row.status)" size="small">{{ statusText(row.status) }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="count" label="数量" width="150" />
        <el-table-column label="占比" min-width="200">
          <template #default="{ row }">
            <el-progress :percentage="Math.round(row.count / Math.max(stats.totalOrders, 1) * 100)" :color="statusColor(row.status)" />
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { getOrderStats, type OrderStats } from '@/api/modules/report'

const loading = ref(false)
const stats = reactive<OrderStats>({ totalOrders: 0, byStatus: [] })

onMounted(() => fetchData())

async function fetchData() {
  loading.value = true
  try {
    const result = await getOrderStats()
    Object.assign(stats, result)
  } finally { loading.value = false }
}

function statusTag(s: string) { const m: Record<string,string> = { draft:'info', submitted:'', pending_dispatch:'warning', dispatched:'', in_transit:'', signed:'success', settled:'success', cancelled:'danger' }; return m[s]||'info' }
function statusText(s: string) { const m: Record<string,string> = { draft:'草稿', submitted:'已提交', pending_dispatch:'待调度', dispatched:'已调度', in_transit:'运输中', signed:'已签收', settled:'已结算', cancelled:'已取消' }; return m[s]||s }
function statusColor(s: string) { const m: Record<string,string> = { draft:'#909399', submitted:'#409EFF', pending_dispatch:'#E6A23C', dispatched:'#67C23A', in_transit:'#409EFF', signed:'#67C23A', settled:'#67C23A', cancelled:'#F56C6C' }; return m[s]||'#409EFF' }
</script>

<style scoped>
.report-page__card { margin-bottom: 16px; }
</style>
