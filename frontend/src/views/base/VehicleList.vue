<template>
  <CrudList prefix="/vehicles" :columns="columns" show-status :default-form="{ status: 'active' }">
    <template #form="{ form }">
      <el-form-item label="车牌号" required>
        <el-input v-model="form.plateNo" />
      </el-form-item>
      <el-form-item label="车辆类型">
        <el-input v-model="form.type" />
      </el-form-item>
      <el-form-item label="品牌型号">
        <el-input v-model="form.brandModel" />
      </el-form-item>
      <el-form-item label="颜色">
        <el-input v-model="form.color" />
      </el-form-item>
      <el-form-item label="最大载重(kg)">
        <el-input-number v-model="form.maxLoad" :min="0" :precision="2" />
      </el-form-item>
      <el-form-item label="最大容积(m³)">
        <el-input-number v-model="form.maxVolume" :min="0" :precision="3" />
      </el-form-item>
      <el-form-item label="归属类型">
        <el-select v-model="form.ownerType" placeholder="请选择">
          <el-option label="自有" value="self" />
          <el-option label="承运商" value="carrier" />
        </el-select>
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
  { prop: 'plateNo', label: '车牌号', width: 130 },
  { prop: 'type', label: '类型', width: 80 },
  { prop: 'brandModel', label: '品牌型号', width: 140 },
  { prop: 'color', label: '颜色', width: 80 },
  { prop: 'maxLoad', label: '载重(kg)', width: 100 },
  { prop: 'ownerType', label: '归属', width: 80 },
]
</script>
