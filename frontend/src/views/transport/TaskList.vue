<template>
  <div class="task-list">
    <div class="task-list__header">
      <el-input v-model="keyword" placeholder="搜索任务号/司机/车牌" clearable class="task-list__search" @keyup.enter="fetchData"><template #prefix><el-icon><Search /></el-icon></template></el-input>
      <el-select v-model="statusFilter" placeholder="运输状态" clearable @change="fetchData" style="width:140px; margin-left:8px">
        <el-option label="待执行" value="pending" />
        <el-option label="已发车" value="departed" />
        <el-option label="运输中" value="in_transit" />
        <el-option label="已到达" value="arrived" />
        <el-option label="已签收" value="signed" />
        <el-option label="已关闭" value="closed" />
      </el-select>
    </div>

    <el-table :data="items" v-loading="loading" stripe border @row-click="(row: any) => router.push(`/transport/task/${row.id}`)" style="cursor:pointer">
      <el-table-column prop="taskNo" label="任务号" width="150" />
      <el-table-column prop="plateNo" label="车牌号" width="120" />
      <el-table-column prop="driverName" label="司机" width="100" />
      <el-table-column prop="originName" label="起点" width="120" />
      <el-table-column prop="destName" label="终点" width="120" />
      <el-table-column label="状态" width="90">
        <template #default="{ row }"><el-tag :type="taskStatusTag(row.status)" size="small">{{ taskStatusText(row.status) }}</el-tag></template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right" @click.stop>
        <template #default="{ row }">
          <el-button v-if="row.status==='pending'" size="small" type="primary" @click="handleDepart(row)">发车</el-button>
          <el-button v-if="['departed','in_transit'].includes(row.status)" size="small" type="warning" @click="handleArrive(row)">到达</el-button>
          <el-button v-if="row.status==='arrived'" size="small" type="success" @click="handleSign(row)">签收</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="task-list__pagination">
      <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[10,20,50]" :total="total" layout="total,sizes,prev,pager,next" @change="fetchData" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listTransportTasks, departTask, arriveTask, signTask, type TransportTask } from '@/api/modules/transport'

const router = useRouter()
const loading = ref(false)
const items = ref<TransportTask[]>([])
const page = ref(1); const pageSize = ref(20); const total = ref(0)
const keyword = ref(''); const statusFilter = ref('')

onMounted(() => fetchData())

async function fetchData() {
  loading.value = true
  try {
    const params: any = { page: page.value, pageSize: pageSize.value }
    if (keyword.value) params.keyword = keyword.value
    if (statusFilter.value) params.status = statusFilter.value
    const result = await listTransportTasks(params)
    items.value = result.items; total.value = result.total
  } finally { loading.value = false }
}

async function handleDepart(row: TransportTask) {
  try { await departTask(row.id); ElMessage.success('已发车'); fetchData() } catch (e: any) { ElMessage.error(e?.message || '操作失败') }
}

async function handleArrive(row: TransportTask) {
  try { await arriveTask(row.id); ElMessage.success('已到达'); fetchData() } catch (e: any) { ElMessage.error(e?.message || '操作失败') }
}

async function handleSign(row: TransportTask) {
  try {
    await ElMessageBox.confirm('确认签收？', '提示', { type: 'warning' })
    await signTask(row.id); ElMessage.success('已签收'); fetchData()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e?.message || '操作失败') }
}

function taskStatusTag(s: string) {
  const m: Record<string,string>={pending:'info',departed:'warning',in_transit:'',arrived:'primary',signed:'success',closed:'info'}
  return m[s]||''
}
function taskStatusText(s: string) {
  const m: Record<string,string>={pending:'待执行',departed:'已发车',in_transit:'运输中',arrived:'已到达',signed:'已签收',closed:'已关闭'}
  return m[s]||s
}
</script>

<style scoped>
.task-list__header { display: flex; align-items: center; margin-bottom: 16px; }
.task-list__search { width: 300px; }
.task-list__pagination { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
