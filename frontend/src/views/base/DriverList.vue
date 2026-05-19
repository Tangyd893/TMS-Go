<template>
  <CrudList prefix="/drivers" :columns="columns" show-status :default-form="{ status: 'active' }">
    <template #form="{ form }">
      <el-form-item label="司机编号" required>
        <el-input v-model="form.code" />
      </el-form-item>
      <el-form-item label="姓名" required>
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item label="手机号">
        <el-input v-model="form.phone" />
      </el-form-item>
      <el-form-item label="身份证号">
        <el-input v-model="form.idCard" />
      </el-form-item>
      <el-form-item label="驾驶证类型">
        <el-input v-model="form.licenseType" />
      </el-form-item>
      <el-form-item label="驾驶证号">
        <el-input v-model="form.licenseNo" />
      </el-form-item>
      <el-form-item label="驾驶证到期日">
        <el-date-picker v-model="form.licenseExpireDate" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" />
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
  { prop: 'name', label: '姓名', width: 100 },
  { prop: 'phone', label: '手机号', width: 130 },
  { prop: 'licenseType', label: '驾照类型', width: 100 },
  { prop: 'licenseNo', label: '驾照号', width: 140 },
]
</script>
