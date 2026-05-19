<template>
  <div class="crud-list">
    <div class="crud-list__header">
      <el-input v-model="keyword" placeholder="请输入关键字搜索" clearable class="crud-list__search" @clear="fetchData" @keyup.enter="fetchData">
        <template #prefix><el-icon><Search /></el-icon></template>
      </el-input>
      <el-button type="primary" @click="openCreate" v-if="canCreate">
        <el-icon><Plus /></el-icon> 新增
      </el-button>
    </div>

    <el-table :data="items" v-loading="loading" stripe border>
      <el-table-column v-for="col in columns" :key="col.prop" :prop="col.prop" :label="col.label" :width="col.width" :min-width="col.minWidth">
        <template #default="{ row }">
          <slot :name="'col-' + col.prop" :row="row">
            {{ row[col.prop] }}
          </slot>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="80" v-if="showStatus">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
            {{ row.status === 'active' ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="openEdit(row)" v-if="canUpdate">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)" v-if="canDelete">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="crud-list__pagination">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next"
        @change="fetchData"
      />
    </div>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑' : '新增'" width="560px" @close="resetForm">
      <el-form ref="formRef" :model="form" label-position="top">
        <slot name="form" :form="form" :editing="!!editingId" />
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Search, Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listEntities, getEntity, createEntity, updateEntity, deleteEntity, type PageResult } from '@/api/modules/base'

interface Column {
  prop: string
  label: string
  width?: string | number
  minWidth?: string | number
}

const props = defineProps<{
  prefix: string
  columns: Column[]
  showStatus?: boolean
  canCreate?: boolean
  canUpdate?: boolean
  canDelete?: boolean
  defaultForm?: Record<string, any>
}>()

const loading = ref(false)
const items = ref<any[]>([])
  const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const keyword = ref('')

const dialogVisible = ref(false)
const editingId = ref('')
const submitting = ref(false)
const form = ref<Record<string, any>>({ ...props.defaultForm || {} })

onMounted(() => fetchData())

async function fetchData() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: page.value, pageSize: pageSize.value }
    if (keyword.value) params.keyword = keyword.value
    const result = await listEntities<any>(props.prefix, params)
    items.value = result.items
    total.value = result.total
  } catch (e: any) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = ''
  form.value = { ...(props.defaultForm || {}), status: 'active' }
  dialogVisible.value = true
}

async function openEdit(row: any) {
  try {
    const entity = await getEntity<any>(props.prefix, row.id)
    editingId.value = entity.id
    form.value = { ...entity }
    dialogVisible.value = true
  } catch (e) {
    ElMessage.error('获取数据失败')
  }
}

async function handleDelete(row: any) {
  try {
    await ElMessageBox.confirm('确认删除该记录？', '提示', { type: 'warning' })
    await deleteEntity(props.prefix, row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error(e?.message || '删除失败')
  }
}

async function handleSubmit() {
  submitting.value = true
  try {
    if (editingId.value) {
      await updateEntity(props.prefix, editingId.value, form.value)
      ElMessage.success('更新成功')
    } else {
      await createEntity(props.prefix, form.value)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (e: any) {
    ElMessage.error(e?.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

function resetForm() {
  form.value = { ...(props.defaultForm || {}) }
  editingId.value = ''
}
</script>

<style scoped>
.crud-list__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.crud-list__search {
  width: 280px;
}
.crud-list__pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
