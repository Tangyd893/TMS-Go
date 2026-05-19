<template>
  <CrudList prefix="/stations" :columns="columns" show-status :default-form="{ status: 'active' }">
    <template #form="{ form }">
      <el-form-item label="站点编号" required>
        <el-input v-model="form.code" />
      </el-form-item>
      <el-form-item label="站点名称" required>
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item label="省份">
        <el-input v-model="form.province" />
      </el-form-item>
      <el-form-item label="城市">
        <el-input v-model="form.city" />
      </el-form-item>
      <el-form-item label="区县">
        <el-input v-model="form.district" />
      </el-form-item>
      <el-form-item label="详细地址">
        <el-input v-model="form.address" />
      </el-form-item>
      <el-form-item label="联系人">
        <el-input v-model="form.contactName" />
      </el-form-item>
      <el-form-item label="联系电话">
        <el-input v-model="form.contactPhone" />
      </el-form-item>
      <el-form-item label="纬度">
        <el-input-number v-model="form.latitude" :min="-90" :max="90" :precision="7" />
      </el-form-item>
      <el-form-item label="经度">
        <el-input-number v-model="form.longitude" :min="-180" :max="180" :precision="7" />
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
  { prop: 'code', label: '编号', width: 100 },
  { prop: 'name', label: '名称', minWidth: 160 },
  { prop: 'province', label: '省份', width: 100 },
  { prop: 'city', label: '城市', width: 100 },
  { prop: 'district', label: '区县', width: 100 },
  { prop: 'contactName', label: '联系人', width: 100 },
  { prop: 'contactPhone', label: '电话', width: 130 },
]
</script>
