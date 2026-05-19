<template>
  <div class="dispatch-list">
    <div class="dispatch-list__header">
      <el-button type="primary" @click="router.push('/dispatch/create')"><el-icon><Plus /></el-icon> 新增调度计划</el-button>
    </div>

    <el-table :data="items" v-loading="loading" stripe border @row-click="(row: any) => router.push(`/dispatch/plan/${row.id}`)" style="cursor:pointer">
      <el-table-column prop="planNo" label="计划编号" width="150" />
      <el-table-column prop="status" label="状态" width="90">
        <template #default="{ row }"><el-tag :type="statusTag(row.status)" size="small">{{ statusText(row.status) }}</el-tag></template>
      </el-table-column>
      <el-table-column label="关联订单" min-width="200">
        <template #default="{ row }">
          <el-tag v-for="d in (row.details || [])" :key="d.id" size="small" style="margin-right:4px">订单{{ d.orderId?.substring(0,8) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="createdAt" label="创建时间" width="180" />
      <el-table-column label="操作" width="160" fixed="right" @click.stop>
        <template #default="{ row }">
          <el-button v-if="row.status==='pending'" size="small" type="danger" @click="handleCancel(row)">取消</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="dispatch-list__pagination">
      <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[10,20,50]" :total="total" layout="total,sizes,prev,pager,next" @change="fetchData" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listDispatchPlans, cancelDispatch, type DispatchPlan } from '@/api/modules/order'

const router = useRouter()
const loading = ref(false)
const items = ref<DispatchPlan[]>([])
const page = ref(1); const pageSize = ref(20); const total = ref(0)

onMounted(() => fetchData())

async function fetchData() {
  loading.value = true
  try {
    const result = await listDispatchPlans({ page: page.value, pageSize: pageSize.value })
    items.value = result.items; total.value = result.total
  } finally { loading.value = false }
}

async function handleCancel(row: DispatchPlan) {
  try {
    await ElMessageBox.confirm(`确认取消调度计划 ${row.planNo}？`, '提示', { type: 'warning' })
    await cancelDispatch(row.id); ElMessage.success('已取消'); fetchData()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e?.message || '操作失败') }
}

function statusTag(s: string) { const m: Record<string,string>={pending:'info',assigned:'',cancelled:'danger'}; return m[s]||'' }
function statusText(s: string) { const m: Record<string,string>={pending:'待调度',assigned:'已分配',cancelled:'已取消'}; return m[s]||s }
</script>

<style scoped>
.dispatch-list__header { margin-bottom: 16px; }
.dispatch-list__pagination { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
