<template>
  <div class="order-create">
    <h3>新增运输订单</h3>
    <el-form ref="formRef" :model="form" label-position="top" style="max-width:800px">
      <el-row :gutter="16">
        <el-col :span="12"><el-form-item label="客户"><el-input v-model="form.customerId" placeholder="客户UUID" /></el-form-item></el-col>
        <el-col :span="12"><el-form-item label="备注"><el-input v-model="form.remark" /></el-form-item></el-col>
      </el-row>
      <el-card header="发货方信息" style="margin-bottom:16px">
        <el-row :gutter="16">
          <el-col :span="8"><el-form-item label="发货方"><el-input v-model="form.shipperName" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="电话"><el-input v-model="form.shipperPhone" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="地址"><el-input v-model="form.shipperAddress" /></el-form-item></el-col>
        </el-row>
      </el-card>
      <el-card header="收货方信息" style="margin-bottom:16px">
        <el-row :gutter="16">
          <el-col :span="8"><el-form-item label="收货方"><el-input v-model="form.receiverName" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="电话"><el-input v-model="form.receiverPhone" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="地址"><el-input v-model="form.receiverAddress" /></el-form-item></el-col>
        </el-row>
      </el-card>
      <el-row :gutter="16">
        <el-col :span="12"><el-form-item label="起点"><el-input v-model="form.originName" /></el-form-item></el-col>
        <el-col :span="12"><el-form-item label="终点"><el-input v-model="form.destName" /></el-form-item></el-col>
        <el-col :span="12"><el-form-item label="计划提货时间"><el-date-picker v-model="form.planPickupTime" type="datetime" value-format="YYYY-MM-DDTHH:mm:ss+08:00" style="width:100%" /></el-form-item></el-col>
        <el-col :span="12"><el-form-item label="计划送达时间"><el-date-picker v-model="form.planDeliveryTime" type="datetime" value-format="YYYY-MM-DDTHH:mm:ss+08:00" style="width:100%" /></el-form-item></el-col>
      </el-row>

      <el-card header="货物信息" style="margin-bottom:16px">
        <div v-for="(c, idx) in form.cargoItems" :key="idx" style="margin-bottom:12px">
          <el-row :gutter="12">
            <el-col :span="6"><el-input v-model="c.cargoName" placeholder="货物名称" /></el-col>
            <el-col :span="4"><el-input v-model="c.cargoType" placeholder="类型" /></el-col>
            <el-col :span="3"><el-input-number v-model="c.quantity" :min="1" placeholder="数量" /></el-col>
            <el-col :span="3"><el-input-number v-model="c.weight" :min="0" :precision="2" placeholder="重量" /></el-col>
            <el-col :span="3"><el-input-number v-model="c.volume" :min="0" :precision="2" placeholder="体积" /></el-col>
            <el-col :span="3"><el-input v-model="c.unit" placeholder="单位" /></el-col>
            <el-col :span="2"><el-button type="danger" @click="form.cargoItems.splice(idx,1)" :disabled="form.cargoItems.length<=1">-</el-button></el-col>
          </el-row>
        </div>
        <el-button @click="form.cargoItems.push({cargoName:'',cargoType:'',quantity:1,weight:0,volume:0,unit:''})">+ 添加货物</el-button>
      </el-card>

      <el-button type="primary" :loading="submitting" @click="handleCreate">创建订单</el-button>
      <el-button @click="router.back()">取消</el-button>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { createOrder } from '@/api/modules/order'

const router = useRouter()
const submitting = ref(false)
const form = reactive({
  customerId: '', remark: '',
  shipperName: '', shipperPhone: '', shipperAddress: '',
  receiverName: '', receiverPhone: '', receiverAddress: '',
  originName: '', destName: '',
  planPickupTime: '', planDeliveryTime: '',
  cargoItems: [{ cargoName: '', cargoType: '', quantity: 1, weight: 0, volume: 0, unit: '' }],
})

async function handleCreate() {
  submitting.value = true
  try {
    await createOrder(form)
    ElMessage.success('订单创建成功')
    router.push('/order/list')
  } catch (e: any) { ElMessage.error(e?.message || '创建失败') }
  finally { submitting.value = false }
}
</script>
