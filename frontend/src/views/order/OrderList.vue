<template>
  <div class="order-list">
    <div class="order-list__header">
      <el-input v-model="keyword" placeholder="搜索订单号/客户" clearable class="order-list__search" @keyup.enter="fetchData"><template #prefix><el-icon><Search /></el-icon></template></el-input>
      <el-select v-model="statusFilter" placeholder="订单状态" clearable @change="fetchData" style="width:140px; margin-left:8px">
        <el-option label="草稿" value="draft" />
        <el-option label="已提交" value="submitted" />
        <el-option label="待调度" value="pending_dispatch" />
        <el-option label="已调度" value="dispatched" />
        <el-option label="运输中" value="in_transit" />
        <el-option label="已签收" value="signed" />
        <el-option label="已取消" value="cancelled" />
      </el-select>
      <el-button type="primary" @click="router.push('/order/create')"><el-icon><Plus /></el-icon> 新增订单</el-button>
    </div>

    <el-table :data="items" v-loading="loading" stripe border @row-click="(row: any) => router.push(`/order/${row.id}`)" style="cursor:pointer">
      <el-table-column prop="orderNo" label="订单号" width="150" />
      <el-table-column prop="customerName" label="客户" min-width="120" />
      <el-table-column prop="originName" label="起点" width="120" />
      <el-table-column prop="destName" label="终点" width="120" />
      <el-table-column prop="cargoName" label="货物" width="100" />
      <el-table-column label="状态" width="90">
        <template #default="{ row }"><el-tag :type="statusTag(row.status)" size="small">{{ statusText(row.status) }}</el-tag></template>
      </el-table-column>
      <el-table-column label="操作" width="160" fixed="right" @click.stop>
        <template #default="{ row }">
          <el-button v-if="row.status==='draft'" size="small" type="success" @click="handleSubmit(row)">提交</el-button>
          <el-button v-if="['draft','submitted','pending_dispatch'].includes(row.status)" size="small" type="danger" @click="handleCancel(row)">取消</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="order-list__pagination">
      <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[10,20,50]" :total="total" layout="total,sizes,prev,pager,next" @change="fetchData" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search, Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listOrders, submitOrder, cancelOrder, type Order } from '@/api/modules/order'

const router = useRouter()
const loading = ref(false)
const items = ref<Order[]>([])
const page = ref(1); const pageSize = ref(20); const total = ref(0)
const keyword = ref(''); const statusFilter = ref('')

onMounted(() => fetchData())

async function fetchData() {
  loading.value = true
  try {
    const params: any = { page: page.value, pageSize: pageSize.value }
    if (keyword.value) params.keyword = keyword.value
    if (statusFilter.value) params.status = statusFilter.value
    const result = await listOrders(params)
    items.value = result.items; total.value = result.total
  } finally { loading.value = false }
}

async function handleSubmit(row: Order) {
  try { await submitOrder(row.id); ElMessage.success('提交成功'); fetchData() } catch (e: any) { ElMessage.error(e?.message || '提交失败') }
}

async function handleCancel(row: Order) {
  try {
    await ElMessageBox.confirm(`确认取消订单 ${row.orderNo}？`, '提示', { type: 'warning' })
    await cancelOrder(row.id); ElMessage.success('已取消'); fetchData()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e?.message || '操作失败') }
}

function statusTag(s: string) { const m: Record<string,string>={draft:'info',submitted:'',pending_dispatch:'warning',dispatched:'',in_transit:'',signed:'success',cancelled:'info'}; return m[s]||'' }
function statusText(s: string) { const m: Record<string,string>={draft:'草稿',submitted:'已提交',pending_dispatch:'待调度',dispatched:'已调度',in_transit:'运输中',signed:'已签收',settled:'已结算',cancelled:'已取消'}; return m[s]||s }
</script>

<style scoped>
.order-list__header { display:flex; align-items:center; margin-bottom:16px; flex-wrap:wrap; gap:8px; }
.order-list__search { width:260px; }
.order-list__pagination { display:flex; justify-content:flex-end; margin-top:16px; }
</style>
