<template>
  <div class="dispatch-page">
    <div class="dispatch-page__header">
      <el-input v-model="keyword" placeholder="搜索计划号" style="width:240px" @keyup.enter="fetchPlans"><template #prefix><el-icon><Search /></el-icon></template></el-input>
      <el-button type="primary" @click="showCreateDialog"><el-icon><Plus /></el-icon> 创建调度计划</el-button>
    </div>

    <el-table :data="plans" v-loading="loading" stripe border>
      <el-table-column prop="planNo" label="计划号" width="150" />
      <el-table-column label="状态" width="100"><template #default="{row}"><el-tag size="small">{{row.status}}</el-tag></template></el-table-column>
      <el-table-column label="订单数" width="80"><template #default="{row}">{{row.details?.length||0}}</template></el-table-column>
      <el-table-column prop="createdAt" label="创建时间" width="180">
        <template #default="{row}">{{ new Date(row.createdAt).toLocaleString() }}</template>
      </el-table-column>
    </el-table>

    <div style="display:flex;justify-content:flex-end;margin-top:16px">
      <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :total="total" layout="total,prev,pager,next" @change="fetchPlans" />
    </div>

    <el-dialog v-model="dialogVisible" title="新建调度计划" width="500px">
      <el-card header="待调度订单" shadow="never" style="max-height:350px;overflow:auto">
        <el-checkbox-group v-model="selectedOrders">
          <div v-for="o in pendingOrders" :key="o.id" style="padding:8px 0;border-bottom:1px solid #eee">
            <el-checkbox :value="o.id">{{ o.orderNo }} - {{ o.originName }} → {{ o.destName }}</el-checkbox>
          </div>
        </el-checkbox-group>
        <p v-if="pendingOrders.length===0" style="color:#999">暂无待调度订单</p>
      </el-card>
      <template #footer>
        <el-button @click="dialogVisible=false">取消</el-button>
        <el-button type="primary" :disabled="selectedOrders.length===0" @click="handleCreate">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Search, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { listDispatchPlans, createDispatchPlan, listPendingOrders, type DispatchPlan, type PendingOrder } from '@/api/modules/order'

const loading = ref(false)
const plans = ref<DispatchPlan[]>([])
const page = ref(1); const pageSize = ref(20); const total = ref(0)
const keyword = ref('')
const dialogVisible = ref(false)
const pendingOrders = ref<PendingOrder[]>([])
const selectedOrders = ref<string[]>([])

onMounted(() => fetchPlans())

async function fetchPlans() {
  loading.value = true
  try { const r = await listDispatchPlans({ page: page.value, pageSize: pageSize.value }); plans.value = r.items; total.value = r.total } finally { loading.value = false }
}

async function showCreateDialog() {
  pendingOrders.value = await listPendingOrders()
  selectedOrders.value = []
  dialogVisible.value = true
}

async function handleCreate() {
  try { await createDispatchPlan(selectedOrders.value); ElMessage.success('创建成功'); dialogVisible.value = false; fetchPlans() } catch (e: any) { ElMessage.error(e?.message || '创建失败') }
}
</script>

<style scoped>
.dispatch-page__header { display:flex; justify-content:space-between; margin-bottom:16px; }
</style>
