<template>
  <div class="dict-page">
    <div class="dict-page__header">
      <h3>字典管理</h3>
      <el-button type="primary" @click="openCreateType"><el-icon><Plus /></el-icon> 新增字典</el-button>
    </div>

    <el-table :data="types" v-loading="loading" stripe border>
      <el-table-column prop="code" label="字典编码" width="160" />
      <el-table-column prop="name" label="字典名称" min-width="160" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
            {{ row.status === 'active' ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="240">
        <template #default="{ row }">
          <el-button size="small" @click="openItems(row)">字典项</el-button>
          <el-button size="small" @click="openEditType(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDeleteType(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 字典类型弹窗 -->
    <el-dialog v-model="typeDialogVisible" :title="editingTypeId ? '编辑字典' : '新增字典'" width="480px">
      <el-form label-position="top">
        <el-form-item label="字典编码" required>
          <el-input v-model="typeForm.code" :disabled="!!editingTypeId" />
        </el-form-item>
        <el-form-item label="字典名称" required>
          <el-input v-model="typeForm.name" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="typeForm.status">
            <el-option label="启用" value="active" />
            <el-option label="禁用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="typeForm.remark" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="typeDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="typeSubmitting" @click="handleTypeSubmit">确认</el-button>
      </template>
    </el-dialog>

    <!-- 字典项弹窗 -->
    <el-dialog v-model="itemsDialogVisible" :title="`字典项 - ${currentType?.name || ''}`" width="600px">
      <el-button type="primary" size="small" @click="openCreateItem" style="margin-bottom:12px">
        <el-icon><Plus /></el-icon> 新增字典项
      </el-button>
      <el-table :data="currentItems" stripe border size="small">
        <el-table-column prop="itemCode" label="编码" width="120" />
        <el-table-column prop="itemName" label="名称" min-width="140" />
        <el-table-column prop="sortNo" label="排序" width="80" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140">
          <template #default="{ row }">
            <el-button size="small" text @click="openEditItem(row)">编辑</el-button>
            <el-button size="small" text type="danger" @click="handleDeleteItem(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <!-- 字典项编辑弹窗 -->
    <el-dialog v-model="itemDialogVisible" :title="editingItemId ? '编辑字典项' : '新增字典项'" width="420px">
      <el-form label-position="top">
        <el-form-item label="编码" required>
          <el-input v-model="itemForm.itemCode" :disabled="!!editingItemId" />
        </el-form-item>
        <el-form-item label="名称" required>
          <el-input v-model="itemForm.itemName" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="itemForm.sortNo" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="itemForm.status">
            <el-option label="启用" value="active" />
            <el-option label="禁用" value="disabled" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="itemDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="itemSubmitting" @click="handleItemSubmit">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listDictTypes, createDictType, updateDictType, deleteDictType, listDictItems, createDictItem, updateDictItem, deleteDictItem, type DictType, type DictItem } from '@/api/modules/dict'

const loading = ref(false)
const types = ref<DictType[]>([])

const typeDialogVisible = ref(false)
const typeSubmitting = ref(false)
const editingTypeId = ref('')
const typeForm = ref({ code: '', name: '', status: 'active', remark: '' })

const itemsDialogVisible = ref(false)
const currentType = ref<DictType | null>(null)
const currentItems = ref<DictItem[]>([])

const itemDialogVisible = ref(false)
const itemSubmitting = ref(false)
const editingItemId = ref('')
const itemForm = ref({ typeCode: '', itemCode: '', itemName: '', sortNo: 0, status: 'active' })

onMounted(() => fetchTypes())

async function fetchTypes() {
  loading.value = true
  try { types.value = await listDictTypes() } catch (e) { console.error(e) }
  finally { loading.value = false }
}

function resetTypeForm() { typeForm.value = { code: '', name: '', status: 'active', remark: '' }; editingTypeId.value = '' }

function openCreateType() { resetTypeForm(); typeDialogVisible.value = true }
function openEditType(row: DictType) { editingTypeId.value = row.id; typeForm.value = { code: row.code, name: row.name, status: row.status, remark: row.remark || '' }; typeDialogVisible.value = true }

async function handleTypeSubmit() {
  typeSubmitting.value = true
  try {
    if (editingTypeId.value) {
      await updateDictType(editingTypeId.value, typeForm.value)
    } else {
      await createDictType(typeForm.value)
    }
    ElMessage.success(editingTypeId.value ? '更新成功' : '创建成功')
    typeDialogVisible.value = false
    fetchTypes()
  } catch (e: any) { ElMessage.error(e?.message || '操作失败') }
  finally { typeSubmitting.value = false }
}

async function handleDeleteType(row: DictType) {
  try {
    await ElMessageBox.confirm(`确认删除字典「${row.name}」？`, '提示', { type: 'warning' })
    await deleteDictType(row.id)
    ElMessage.success('删除成功')
    fetchTypes()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e?.message || '删除失败') }
}

async function openItems(row: DictType) {
  currentType.value = row
  currentItems.value = await listDictItems(row.code)
  itemsDialogVisible.value = true
}

function resetItemForm() { itemForm.value = { typeCode: currentType.value?.code || '', itemCode: '', itemName: '', sortNo: 0, status: 'active' }; editingItemId.value = '' }

function openCreateItem() { resetItemForm(); itemDialogVisible.value = true }
function openEditItem(row: DictItem) { editingItemId.value = row.id; itemForm.value = { typeCode: row.typeCode, itemCode: row.itemCode, itemName: row.itemName, sortNo: row.sortNo, status: row.status }; itemDialogVisible.value = true }

async function handleItemSubmit() {
  itemSubmitting.value = true
  try {
    if (editingItemId.value) {
      await updateDictItem(editingItemId.value, itemForm.value)
    } else {
      await createDictItem(itemForm.value)
    }
    ElMessage.success(editingItemId.value ? '更新成功' : '创建成功')
    itemDialogVisible.value = false
    if (currentType.value) currentItems.value = await listDictItems(currentType.value.code)
  } catch (e: any) { ElMessage.error(e?.message || '操作失败') }
  finally { itemSubmitting.value = false }
}

async function handleDeleteItem(row: DictItem) {
  try {
    await ElMessageBox.confirm(`确认删除字典项「${row.itemName}」？`, '提示', { type: 'warning' })
    await deleteDictItem(row.id)
    ElMessage.success('删除成功')
    if (currentType.value) currentItems.value = await listDictItems(currentType.value.code)
  } catch (e: any) { if (e !== 'cancel') ElMessage.error(e?.message || '删除失败') }
}
</script>

<style scoped>
.dict-page__header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.dict-page__header h3 { margin: 0; }
</style>
