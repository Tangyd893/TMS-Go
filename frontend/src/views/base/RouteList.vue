<template>
  <CrudList prefix="/routes" :columns="columns" show-status :default-form="{ status: 'active' }">
    <template #form="{ form }">
      <el-form-item label="线路编号" required>
        <el-input v-model="form.code" />
      </el-form-item>
      <el-form-item label="线路名称" required>
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item label="起点">
        <el-input v-model="form.originName" />
      </el-form-item>
      <el-form-item label="终点">
        <el-input v-model="form.destName" />
      </el-form-item>
      <el-form-item label="里程(km)">
        <el-input-number v-model="form.distanceKm" :min="0" :precision="2" />
      </el-form-item>
      <el-form-item label="预计耗时(h)">
        <el-input-number v-model="form.estimatedHours" :min="0" :precision="1" />
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
  { prop: 'originName', label: '起点', width: 140 },
  { prop: 'destName', label: '终点', width: 140 },
  { prop: 'distanceKm', label: '里程(km)', width: 100 },
  { prop: 'estimatedHours', label: '预计耗时(h)', width: 110 },
]
</script>
