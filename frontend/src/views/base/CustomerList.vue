<template>
  <CrudList prefix="/customers" :columns="columns" show-status :default-form="{ status: 'active' }">
    <template #form="{ form }">
      <el-form-item label="客户编号" required>
        <el-input v-model="form.code" />
      </el-form-item>
      <el-form-item label="客户名称" required>
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item label="简称">
        <el-input v-model="form.shortName" />
      </el-form-item>
      <el-form-item label="联系人">
        <el-input v-model="form.contactName" />
      </el-form-item>
      <el-form-item label="联系电话">
        <el-input v-model="form.contactPhone" />
      </el-form-item>
      <el-form-item label="地址">
        <el-input v-model="form.address" />
      </el-form-item>
      <el-form-item label="备注">
        <el-input v-model="form.remark" type="textarea" />
      </el-form-item>
    </template>
  </CrudList>
</template>

<script setup lang="ts">
import CrudList from '@/components/table/CrudList.vue'

const columns = [
  { prop: 'code', label: '编号', width: 120 },
  { prop: 'name', label: '名称', minWidth: 160 },
  { prop: 'shortName', label: '简称', width: 120 },
  { prop: 'contactName', label: '联系人', width: 100 },
  { prop: 'contactPhone', label: '联系电话', width: 130 },
]
</script>
