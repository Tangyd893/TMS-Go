<template>
  <div class="order-detail" v-if="order">
    <el-page-header @back="router.back()" title="订单列表" :content="order.orderNo" />

    <el-card header="基础信息" shadow="never" class="order-detail__card">
      <el-descriptions :column="3" border>
        <el-descriptions-item label="订单号">{{ order.orderNo }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusTag(order.status)">{{ statusText(order.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="客户">{{ order.customerName }}</el-descriptions-item>
        <el-descriptions-item label="发件人">{{ order.shipperName }}</el-descriptions-item>
        <el-descriptions-item label="发件电话">{{ order.shipperPhone }}</el-descriptions-item>
        <el-descriptions-item label="发件地址">{{ order.shipperAddress }}</el-descriptions-item>
        <el-descriptions-item label="收件人">{{ order.receiverName }}</el-descriptions-item>
        <el-descriptions-item label="收件电话">{{ order.receiverPhone }}</el-descriptions-item>
        <el-descriptions-item label="收件地址">{{ order.receiverAddress }}</el-descriptions-item>
        <el-descriptions-item label="起点">{{ order.originName }}</el-descriptions-item>
        <el-descriptions-item label="终点">{{ order.destName }}</el-descriptions-item>
        <el-descriptions-item label="运输要求">{{ order.transportRequirement || '-' }}</el-descriptions-item>
        <el-descriptions-item label="计划提货">{{ order.planPickupTime || '-' }}</el-descriptions-item>
        <el-descriptions-item label="计划送达">{{ order.planDeliveryTime || '-' }}</el-descriptions-item>
        <el-descriptions-item label="备注">{{ order.remark || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card header="货物明细" shadow="never" class="order-detail__card">
      <el-table :data="order.cargoItems || []" border>
        <el-table-column prop="cargoName" label="货物名称" />
        <el-table-column prop="cargoType" label="货物类型" />
        <el-table-column prop="quantity" label="数量" />
        <el-table-column prop="weight" label="重量(kg)" />
        <el-table-column prop="volume" label="体积(方)" />
        <el-table-column prop="unit" label="单位" />
      </el-table>
    </el-card>

    <div class="order-detail__actions" v-if="['draft','submitted','pending_dispatch'].includes(order.status)">
      <el-button v-if="order.status==='draft'" type="success" @click="handleSubmit">提交订单</el-button>
      <el-button v-if="order.status==='draft'" type="primary" @click="router.push(`/order/${order.id}/edit`)">编辑订单</el-button>
      <el-button type="danger" @click="handleCancel">取消订单</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getOrder, submitOrder, cancelOrder, type Order } from '@/api/modules/order'

const route = useRoute()
const router = useRouter()
const order = ref<Order | null>(null)

onMounted(async () => {
  const id = route.params.id as string
  try { order.value = await getOrder(id) } catch { ElMessage.error('加载订单详情失败') }
})

async function handleSubmit() {
  try {
    await submitOrder(order.value!.id)
    ElMessage.success('提交成功')
    order.value = await getOrder(order.value!.id)
  } catch (e: any) { ElMessage.error(e?.message || '提交失败') }
}

async function handleCancel() {
  try {
    await ElMessageBox.confirm(`确认取消订单？`, '提示', { type: 'warning' })
    await cancelOrder(order.value!.id)
    ElMessage.success('已取消')
    order.value = await getOrder(order.value!.id)
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e?.message || '操作失败') }
}

function statusTag(s: string) {
  const m: Record<string,string>={draft:'info',submitted:'',pending_dispatch:'warning',dispatched:'',in_transit:'',signed:'success',cancelled:'info'}
  return m[s]||''
}
function statusText(s: string) {
  const m: Record<string,string>={draft:'草稿',submitted:'已提交',pending_dispatch:'待调度',dispatched:'已调度',in_transit:'运输中',signed:'已签收',cancelled:'已取消'}
  return m[s]||s
}
</script>

<style scoped>
.order-detail__card { margin-top: 16px; }
.order-detail__actions { margin-top: 16px; display: flex; gap: 8px; }
</style>
