<template>
  <div class="dispatch-create">
    <el-page-header @back="router.back()" title="调度计划列表" content="创建调度计划" />

    <el-card header="选择待调度订单" shadow="never" class="dispatch-create__card">
      <el-table :data="pendingOrders" v-loading="loadingOrders" stripe border ref="tableRef" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="orderNo" label="订单号" width="150" />
        <el-table-column prop="originName" label="起点" width="120" />
        <el-table-column prop="destName" label="终点" width="120" />
      </el-table>

      <div style="margin-top:16px" v-if="selectedOrders.length">
        <el-divider />
        <p>已选择 <strong>{{ selectedOrders.length }}</strong> 个订单, 确认创建调度计划?</p>
        <el-button type="primary" :loading="creating" @click="handleCreate">创建调度计划</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { type PendingOrder } from '@/api/modules/order'
import { listPendingOrders as getPendingOrders } from '@/api/modules/order'
import { createDispatchPlan } from '@/api/modules/order'

const router = useRouter()
const loadingOrders = ref(false)
const creating = ref(false)
const pendingOrders = ref<PendingOrder[]>([])
const selectedOrders = ref<PendingOrder[]>([])

onMounted(async () => {
  loadingOrders.value = true
  try { pendingOrders.value = await getPendingOrders() } catch { ElMessage.error('加载待调度订单失败') }
  finally { loadingOrders.value = false }
})

function handleSelectionChange(val: PendingOrder[]) {
  selectedOrders.value = val
}

async function handleCreate() {
  creating.value = true
  try {
    const orderIds = selectedOrders.value.map(o => o.id)
    await createDispatchPlan(orderIds)
    ElMessage.success('调度计划创建成功')
    router.push('/dispatch/plans')
  } catch (e: any) { ElMessage.error(e?.message || '创建失败') }
  finally { creating.value = false }
}
</script>

<style scoped>
.dispatch-create__card { margin-top: 16px; }
</style>
