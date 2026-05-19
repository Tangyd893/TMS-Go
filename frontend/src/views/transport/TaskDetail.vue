<template>
  <div class="task-detail" v-if="task">
    <el-page-header @back="router.back()" title="运输任务列表" :content="task.taskNo" />

    <el-card header="任务信息" shadow="never" class="task-detail__card">
      <el-descriptions :column="3" border>
        <el-descriptions-item label="任务号">{{ task.taskNo }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="taskStatusTag(task.status)">{{ taskStatusText(task.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="车牌号">{{ task.plateNo || '-' }}</el-descriptions-item>
        <el-descriptions-item label="司机">{{ task.driverName || '-' }}</el-descriptions-item>
        <el-descriptions-item label="起点">{{ task.originName }}</el-descriptions-item>
        <el-descriptions-item label="终点">{{ task.destName }}</el-descriptions-item>
        <el-descriptions-item label="实际发车">{{ task.actualDepartTime || '-' }}</el-descriptions-item>
        <el-descriptions-item label="实际到达">{{ task.actualArriveTime || '-' }}</el-descriptions-item>
        <el-descriptions-item label="备注">{{ task.remark || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <div class="task-detail__actions" style="margin-top:16px">
      <el-button v-if="task.status==='pending'" type="primary" @click="handleDepart">发车确认</el-button>
      <el-button v-if="['departed','in_transit'].includes(task.status)" type="warning" @click="handleArrive">到达确认</el-button>
      <el-button v-if="task.status==='arrived'" type="success" @click="handleSign">签收确认</el-button>
      <el-button v-if="!['pending'].includes(task.status)" type="danger" @click="showException = true">上报异常</el-button>
    </div>

    <el-card header="运输节点" shadow="never" class="task-detail__card">
      <div class="task-detail__actions">
        <el-button size="small" type="primary" @click="showNodeForm = true">新增节点</el-button>
      </div>
      <el-timeline v-if="nodes.length" style="margin-top:16px">
        <el-timeline-item
          v-for="node in nodes" :key="node.id"
          :timestamp="node.arrivedAt || node.createdAt"
          :type="nodeTypeColor(node.nodeType)"
        >
          <strong>{{ node.nodeName || nodeTypeText(node.nodeType) }}</strong>
          <p v-if="node.locationName">{{ node.locationName }}</p>
          <p v-if="node.remark" style="color:#909399">{{ node.remark }}</p>
        </el-timeline-item>
      </el-timeline>
      <el-empty v-else description="暂无节点记录" />
    </el-card>

    <el-card header="签收回单" shadow="never" class="task-detail__card" v-if="receipt">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="回单号">{{ receipt.receiptNo }}</el-descriptions-item>
        <el-descriptions-item label="签收人">{{ receipt.signBy }}</el-descriptions-item>
        <el-descriptions-item label="签收时间">{{ receipt.signAt }}</el-descriptions-item>
        <el-descriptions-item label="备注">{{ receipt.remark || '-' }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-dialog v-model="showNodeForm" title="新增运输节点" width="500px">
      <el-form :model="nodeForm">
        <el-form-item label="节点类型" required>
          <el-select v-model="nodeForm.nodeType">
            <el-option label="提货" value="pickup" />
            <el-option label="发车" value="depart" />
            <el-option label="中转到达" value="transit_arrive" />
            <el-option label="中转离开" value="transit_depart" />
            <el-option label="到达目的地" value="arrive" />
          </el-select>
        </el-form-item>
        <el-form-item label="节点名称">
          <el-input v-model="nodeForm.nodeName" />
        </el-form-item>
        <el-form-item label="位置">
          <el-input v-model="nodeForm.locationName" />
        </el-form-item>
        <el-form-item label="到达时间">
          <el-date-picker v-model="nodeForm.arrivedAt" type="datetime" />
        </el-form-item>
        <el-form-item label="离开时间">
          <el-date-picker v-model="nodeForm.departedAt" type="datetime" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="nodeForm.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showNodeForm = false">取消</el-button>
        <el-button type="primary" @click="submitNode">确认</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showException" title="上报异常" width="500px">
      <el-form :model="exceptionForm">
        <el-form-item label="异常类型" required>
          <el-select v-model="exceptionForm.exceptionType">
            <el-option label="车辆故障" value="车辆故障" />
            <el-option label="交通延误" value="交通延误" />
            <el-option label="货损货差" value="货损货差" />
            <el-option label="客户拒收" value="客户拒收" />
            <el-option label="地址异常" value="地址异常" />
            <el-option label="天气原因" value="天气原因" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item label="严重程度" required>
          <el-radio-group v-model="exceptionForm.severity">
            <el-radio value="normal">一般</el-radio>
            <el-radio value="serious">严重</el-radio>
            <el-radio value="critical">紧急</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="描述" required>
          <el-input v-model="exceptionForm.description" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showException = false">取消</el-button>
        <el-button type="primary" @click="submitException" :loading="submitting">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getTransportTask, departTask, arriveTask, signTask, addNode, getNodes, getReceipt, type TransportTask, type TransportNode, type Receipt } from '@/api/modules/transport'
import { createException } from '@/api/modules/exception'

const route = useRoute()
const router = useRouter()
const task = ref<TransportTask | null>(null)
const nodes = ref<TransportNode[]>([])
const receipt = ref<Receipt | null>(null)
const submitting = ref(false)

const showNodeForm = ref(false)
const showException = ref(false)

const nodeForm = ref({ nodeType: '', nodeName: '', locationName: '', arrivedAt: '', departedAt: '', remark: '' })
const exceptionForm = ref({ exceptionType: '', severity: 'normal', description: '' })

onMounted(async () => {
  const id = route.params.id as string
  try {
    task.value = await getTransportTask(id)
    nodes.value = await getNodes(id)
    try { receipt.value = await getReceipt(id) } catch { receipt.value = null }
  } catch { ElMessage.error('加载运输任务失败') }
})

async function handleDepart() {
  try { await departTask(task.value!.id); ElMessage.success('已发车'); reloadTask() } catch (e: any) { ElMessage.error(e?.message || '操作失败') }
}

async function handleArrive() {
  try { await arriveTask(task.value!.id); ElMessage.success('已到达'); reloadTask() } catch (e: any) { ElMessage.error(e?.message || '操作失败') }
}

async function handleSign() {
  try {
    await ElMessageBox.confirm('确认签收？', '提示', { type: 'warning' })
    await signTask(task.value!.id); ElMessage.success('已签收'); reloadTask()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e?.message || '操作失败') }
}

async function submitNode() {
  try {
    await addNode(task.value!.id, nodeForm.value)
    ElMessage.success('节点已添加')
    showNodeForm.value = false
    nodeForm.value = { nodeType: '', nodeName: '', locationName: '', arrivedAt: '', departedAt: '', remark: '' }
    nodes.value = await getNodes(task.value!.id)
  } catch (e: any) { ElMessage.error(e?.message || '操作失败') }
}

async function submitException() {
  submitting.value = true
  try {
    await createException({
      taskId: task.value!.id,
      orderId: task.value!.orderId,
      exceptionType: exceptionForm.value.exceptionType,
      description: exceptionForm.value.description,
      severity: exceptionForm.value.severity,
    })
    ElMessage.success('异常已上报')
    showException.value = false
    exceptionForm.value = { exceptionType: '', severity: 'normal', description: '' }
  } catch (e: any) { ElMessage.error(e?.message || '操作失败') }
  finally { submitting.value = false }
}

async function reloadTask() {
  task.value = await getTransportTask(task.value!.id)
}

function taskStatusTag(s: string) {
  const m: Record<string,string>={pending:'info',departed:'warning',in_transit:'',arrived:'primary',signed:'success',closed:'info'}
  return m[s]||''
}
function taskStatusText(s: string) {
  const m: Record<string,string>={pending:'待执行',departed:'已发车',in_transit:'运输中',arrived:'已到达',signed:'已签收',closed:'已关闭'}
  return m[s]||s
}
function nodeTypeColor(t: string) {
  const m: Record<string,string>={pickup:'primary',depart:'warning',transit_arrive:'',transit_depart:'info',arrive:'success',sign:'success'}
  return m[t]||''
}
function nodeTypeText(t: string) {
  const m: Record<string,string>={pickup:'提货',depart:'发车',transit_arrive:'中转到达',transit_depart:'中转离开',arrive:'到达目的地',sign:'签收'}
  return m[t]||t
}
</script>

<style scoped>
.task-detail__card { margin-top: 16px; }
.task-detail__actions { display: flex; gap: 8px; }
</style>
