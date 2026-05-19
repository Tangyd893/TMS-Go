<template>
  <div class="page-list">
    <div class="page-list__header">
      <h3>对账管理</h3>
      <el-button type="primary" @click="showCreateDialog">创建对账单</el-button>
    </div>
    <el-table :data="items" v-loading="loading" stripe border>
      <el-table-column prop="statementNo" label="对账单号" width="150" />
      <el-table-column prop="partnerName" label="合作方" width="150" />
      <el-table-column label="类型" width="80">
        <template #default="{ row }">{{ row.partnerType === 'customer' ? '客户' : '承运商' }}</template>
      </el-table-column>
      <el-table-column prop="statementPeriodStart" label="开始日期" width="120" />
      <el-table-column prop="statementPeriodEnd" label="结束日期" width="120" />
      <el-table-column label="应收金额" width="120">
        <template #default="{ row }">¥{{ row.totalReceivable?.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="应付金额" width="120">
        <template #default="{ row }">¥{{ row.totalPayable?.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }"><el-tag :type="statusTag(row.status)" size="small">{{ statusText(row.status) }}</el-tag></template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.status==='draft'" size="small" type="primary" @click="handleConfirm(row)">确认</el-button>
          <el-button v-if="row.status==='confirmed'" size="small" type="success" @click="handleSettle(row)">创建结算</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div class="page-list__pagination">
      <el-pagination v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[10,20,50]" :total="total" layout="total,sizes,prev,pager,next" @change="fetchData" />
    </div>

    <el-dialog v-model="dialogVisible" title="新建对账单" width="500px">
      <el-form label-position="top">
        <el-form-item label="合作方ID"><el-input v-model="form.partnerId" placeholder="客户或承运商ID" /></el-form-item>
        <el-form-item label="合作方名称"><el-input v-model="form.partnerName" /></el-form-item>
        <el-form-item label="合作方类型">
          <el-radio-group v-model="form.partnerType">
            <el-radio value="customer">客户</el-radio>
            <el-radio value="carrier">承运商</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="账期起始"><el-date-picker v-model="form.periodStart" type="date" value-format="YYYY-MM-DD" style="width:100%" /></el-form-item>
        <el-form-item label="账期结束"><el-date-picker v-model="form.periodEnd" type="date" value-format="YYYY-MM-DD" style="width:100%" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible=false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="submitting">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listStatements, createStatement, confirmStatement, createSettlement, type Statement } from '@/api/modules/finance'

const loading = ref(false); const submitting = ref(false)
const items = ref<Statement[]>([])
const page = ref(1); const pageSize = ref(20); const total = ref(0)
const dialogVisible = ref(false)

const form = reactive({ partnerId: '', partnerName: '', partnerType: 'customer', periodStart: '', periodEnd: '' })

onMounted(() => fetchData())

async function fetchData() {
  loading.value = true
  try {
    const result = await listStatements({ page: page.value, pageSize: pageSize.value })
    items.value = result.items; total.value = result.total
  } finally { loading.value = false }
}

function showCreateDialog() { dialogVisible.value = true }

async function handleCreate() {
  if (!form.partnerId || !form.partnerName || !form.periodStart || !form.periodEnd) {
    ElMessage.warning('请填写完整信息'); return
  }
  submitting.value = true
  try {
    await createStatement({ partnerId: form.partnerId, partnerName: form.partnerName, partnerType: form.partnerType, periodStart: form.periodStart, periodEnd: form.periodEnd })
    ElMessage.success('对账单创建成功'); dialogVisible.value = false; fetchData()
  } catch (e: any) { ElMessage.error(e?.message || '创建失败') }
  finally { submitting.value = false }
}

async function handleConfirm(row: Statement) {
  try {
    await ElMessageBox.confirm('确认此对账单？确认后不可修改。', '提示', { type: 'warning' })
    await confirmStatement(row.id); ElMessage.success('已确认'); fetchData()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e?.message || '操作失败') }
}

async function handleSettle(row: Statement) {
  try {
    await ElMessageBox.confirm('为该对账单创建结算单？', '提示', { type: 'info' })
    await createSettlement(row.id); ElMessage.success('结算单已创建'); fetchData()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e?.message || '操作失败') }
}

function statusTag(s: string) { const m: Record<string, string> = { draft: 'info', confirmed: '', settled: 'success' }; return m[s] || '' }
function statusText(s: string) { const m: Record<string, string> = { draft: '草稿', confirmed: '已确认', settled: '已结算' }; return m[s] || s }
</script>

<style scoped>
.page-list__header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.page-list__header h3 { margin: 0; }
.page-list__pagination { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
