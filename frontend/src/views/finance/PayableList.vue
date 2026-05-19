<template>
  <div class="page-list">
    <div class="page-list__header">
      <h3>应付管理</h3>
      <el-select v-model="statusFilter" placeholder="状态筛选" clearable @change="fetchData" style="width:140px">
        <el-option label="待确认" value="pending" />
        <el-option label="已确认" value="confirmed" />
        <el-option label="已结算" value="settled" />
      </el-select>
    </div>
    <el-table :data="items" v-loading="loading" stripe border>
      <el-table-column prop="payableNo" label="应付单号" width="150" />
      <el-table-column prop="carrierName" label="承运商" width="120" />
      <el-table-column prop="feeItemName" label="费用项目" width="120" />
      <el-table-column label="金额" width="120">
        <template #default="{ row }">¥{{ row.totalAmount?.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }"><el-tag :type="statusTag(row.status)" size="small">{{ statusText(row.status) }}</el-tag></template>
      </el-table-column>
      <el-table-column label="含税" width="70">
        <template #default="{ row }">{{ row.hasTax ? '是' : '否' }}</template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" min-width="160" show-overflow-tooltip />
      <el-table-column prop="createdAt" label="创建时间" width="180">
        <template #default="{ row }">{{ new Date(row.createdAt).toLocaleString() }}</template>
      </el-table-column>
    </el-table>
    <div class="page-list__pagination">
      <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[10,20,50]" :total="total" layout="total,sizes,prev,pager,next" @change="fetchData" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { listPayables, type Payable } from '@/api/modules/finance'

const loading = ref(false)
const items = ref<Payable[]>([])
const page = ref(1); const pageSize = ref(20); const total = ref(0)
const statusFilter = ref('')

onMounted(() => fetchData())

async function fetchData() {
  loading.value = true
  try {
    const params: any = { page: page.value, pageSize: pageSize.value }
    if (statusFilter.value) params.status = statusFilter.value
    const result = await listPayables(params)
    items.value = result.items; total.value = result.total
  } finally { loading.value = false }
}

function statusTag(s: string) { const m: Record<string, string> = { pending: 'warning', confirmed: '', settled: 'success' }; return m[s] || '' }
function statusText(s: string) { const m: Record<string, string> = { pending: '待确认', confirmed: '已确认', settled: '已结算' }; return m[s] || s }
</script>

<style scoped>
.page-list__header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.page-list__header h3 { margin: 0; }
.page-list__pagination { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
