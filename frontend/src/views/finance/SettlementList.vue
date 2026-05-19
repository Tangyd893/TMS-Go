<template>
  <div class="page-list">
    <div class="page-list__header">
      <h3>结算管理</h3>
    </div>
    <el-table :data="items" v-loading="loading" stripe border>
      <el-table-column prop="settlementNo" label="结算单号" width="150" />
      <el-table-column prop="partnerName" label="合作方" width="150" />
      <el-table-column label="类型" width="80">
        <template #default="{ row }">{{ row.partnerType === 'customer' ? '客户' : '承运商' }}</template>
      </el-table-column>
      <el-table-column label="结算金额" width="120">
        <template #default="{ row }">¥{{ row.settlementAmount?.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="结算方式" width="100">
        <template #default="{ row }">{{ row.settlementMethod || '--' }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }"><el-tag :type="statusTag(row.status)" size="small">{{ statusText(row.status) }}</el-tag></template>
      </el-table-column>
      <el-table-column prop="settledAt" label="结算时间" width="180">
        <template #default="{ row }">{{ row.settledAt ? new Date(row.settledAt).toLocaleString() : '--' }}</template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" min-width="140" show-overflow-tooltip />
      <el-table-column label="操作" width="100" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.status==='pending'" size="small" type="success" @click="handleComplete(row)">完成结算</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="page-list__pagination">
      <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[10,20,50]" :total="total" layout="total,sizes,prev,pager,next" @change="fetchData" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listSettlements, completeSettlement, type Settlement } from '@/api/modules/finance'

const loading = ref(false)
const items = ref<Settlement[]>([])
const page = ref(1); const pageSize = ref(20); const total = ref(0)

onMounted(() => fetchData())

async function fetchData() {
  loading.value = true
  try {
    const result = await listSettlements({ page: page.value, pageSize: pageSize.value })
    items.value = result.items; total.value = result.total
  } finally { loading.value = false }
}

async function handleComplete(row: Settlement) {
  try {
    await ElMessageBox.confirm('确认完成此结算？', '提示', { type: 'warning' })
    await completeSettlement(row.id); ElMessage.success('结算已完成'); fetchData()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e?.message || '操作失败') }
}

function statusTag(s: string) { const m: Record<string, string> = { pending: 'warning', completed: 'success' }; return m[s] || '' }
function statusText(s: string) { const m: Record<string, string> = { pending: '待结算', completed: '已结算' }; return m[s] || s }
</script>

<style scoped>
.page-list__header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.page-list__header h3 { margin: 0; }
.page-list__pagination { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
