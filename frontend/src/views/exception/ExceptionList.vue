<template>
  <div class="exception-list">
    <div class="exception-list__header">
      <el-select v-model="excTypeFilter" placeholder="异常类型" clearable @change="fetchData" style="width:140px">
        <el-option label="车辆故障" value="车辆故障" />
        <el-option label="交通延误" value="交通延误" />
        <el-option label="货损货差" value="货损货差" />
        <el-option label="客户拒收" value="客户拒收" />
        <el-option label="地址异常" value="地址异常" />
        <el-option label="天气原因" value="天气原因" />
        <el-option label="其他" value="其他" />
      </el-select>
      <el-select v-model="statusFilter" placeholder="处理状态" clearable @change="fetchData" style="width:120px; margin-left:8px">
        <el-option label="待处理" value="pending" />
        <el-option label="处理中" value="processing" />
        <el-option label="已解决" value="resolved" />
        <el-option label="已关闭" value="closed" />
      </el-select>
    </div>

    <el-table :data="items" v-loading="loading" stripe border @row-click="(row: any) => showDetail(row)" style="cursor:pointer">
      <el-table-column prop="exceptionNo" label="异常编号" width="150" />
      <el-table-column prop="exceptionType" label="异常类型" width="100" />
      <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
      <el-table-column label="严重程度" width="90">
        <template #default="{ row }"><el-tag :type="severityTag(row.severity)" size="small">{{ severityText(row.severity) }}</el-tag></template>
      </el-table-column>
      <el-table-column label="状态" width="90">
        <template #default="{ row }"><el-tag :type="excStatusTag(row.status)" size="small">{{ excStatusText(row.status) }}</el-tag></template>
      </el-table-column>
      <el-table-column prop="reportByName" label="上报人" width="100" />
      <el-table-column prop="handlerName" label="处理人" width="100" />
      <el-table-column label="操作" width="160" fixed="right" @click.stop>
        <template #default="{ row }">
          <el-button v-if="row.status==='pending'" size="small" type="primary" @click="openHandle(row)">处理</el-button>
          <el-button v-if="row.status==='resolved'" size="small" type="success" @click="handleClose(row)">关闭</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="exception-list__pagination">
      <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[10,20,50]" :total="total" layout="total,sizes,prev,pager,next" @change="fetchData" />
    </div>

    <el-dialog v-model="showDetailDialog" title="异常详情" width="560px">
      <el-descriptions v-if="currentDetail" :column="2" border>
        <el-descriptions-item label="异常编号">{{ currentDetail.exceptionNo }}</el-descriptions-item>
        <el-descriptions-item label="异常类型">{{ currentDetail.exceptionType }}</el-descriptions-item>
        <el-descriptions-item label="严重程度">{{ severityText(currentDetail.severity) }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ excStatusText(currentDetail.status) }}</el-descriptions-item>
        <el-descriptions-item label="上报人">{{ currentDetail.reportByName }}</el-descriptions-item>
        <el-descriptions-item label="处理人">{{ currentDetail.handlerName || '--' }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentDetail.description }}</el-descriptions-item>
        <el-descriptions-item label="处理结果" :span="2">{{ currentDetail.handleResult || '--' }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
    <el-dialog v-model="showHandle" title="处理异常" width="500px">
      <el-form-item label="处理结果" required>
        <el-input v-model="handleResult" type="textarea" :rows="4" />
      </el-form-item>
      <template #footer>
        <el-button @click="showHandle = false">取消</el-button>
        <el-button type="primary" @click="submitHandle" :loading="submitting">确认处理</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listExceptions, handleException, closeException, type Exception } from '@/api/modules/exception'

const loading = ref(false)
const items = ref<Exception[]>([])
const page = ref(1); const pageSize = ref(20); const total = ref(0)
const excTypeFilter = ref(''); const statusFilter = ref('')
const submitting = ref(false)

const showHandle = ref(false)
const handleResult = ref('')
const currentException = ref<Exception | null>(null)

onMounted(() => fetchData())

async function fetchData() {
  loading.value = true
  try {
    const params: any = { page: page.value, pageSize: pageSize.value }
    if (excTypeFilter.value) params.exceptionType = excTypeFilter.value
    if (statusFilter.value) params.status = statusFilter.value
    const result = await listExceptions(params)
    items.value = result.items; total.value = result.total
  } finally { loading.value = false }
}

const showDetailDialog = ref(false)
const currentDetail = ref<Exception | null>(null)

function showDetail(row: Exception) {
  currentDetail.value = row
  showDetailDialog.value = true
}

function openHandle(row: Exception) {
  currentException.value = row
  handleResult.value = ''
  showHandle.value = true
}

async function submitHandle() {
  if (!currentException.value) return
  submitting.value = true
  try {
    await handleException(currentException.value.id, { handleResult: handleResult.value })
    ElMessage.success('处理成功')
    showHandle.value = false
    fetchData()
  } catch (e: any) { ElMessage.error(e?.message || '处理失败') }
  finally { submitting.value = false }
}

async function handleClose(row: Exception) {
  try {
    await ElMessageBox.confirm('确认关闭此异常？', '提示', { type: 'warning' })
    await closeException(row.id); ElMessage.success('已关闭'); fetchData()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e?.message || '操作失败') }
}

function severityTag(s: string) { const m: Record<string,string>={normal:'',serious:'warning',critical:'danger'}; return m[s]||'' }
function severityText(s: string) { const m: Record<string,string>={normal:'一般',serious:'严重',critical:'紧急'}; return m[s]||s }
function excStatusTag(s: string) { const m: Record<string,string>={pending:'warning',processing:'',resolved:'success',closed:'info'}; return m[s]||'' }
function excStatusText(s: string) { const m: Record<string,string>={pending:'待处理',processing:'处理中',resolved:'已解决',closed:'已关闭'}; return m[s]||s }
</script>

<style scoped>
.exception-list__header { display: flex; align-items: center; margin-bottom: 16px; }
.exception-list__pagination { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
