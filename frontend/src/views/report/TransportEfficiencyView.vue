<template>
  <div class="report-page">
    <h3>运输效率</h3>
    <el-row :gutter="16" v-loading="loading">
      <el-col :span="8">
        <el-card shadow="never">
          <el-statistic title="总任务数" :value="eff.totalTasks" />
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="never">
          <el-statistic title="已完成任务" :value="eff.completedTasks" />
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="never">
          <el-statistic title="平均运输时长 (小时)" :value="eff.avgTransitHours" :precision="1" />
        </el-card>
      </el-col>
    </el-row>
    <el-card shadow="never" style="margin-top:16px">
      <template #header>完成率</template>
      <el-progress type="circle" :percentage="Math.round(eff.totalTasks > 0 ? (eff.completedTasks / eff.totalTasks) * 100 : 0)" :width="160" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { getTransportEfficiency, type TransportEfficiency } from '@/api/modules/report'

const loading = ref(false)
const eff = reactive<TransportEfficiency>({ totalTasks: 0, completedTasks: 0, avgTransitHours: 0 })

onMounted(() => fetchData())

async function fetchData() {
  loading.value = true
  try {
    const result = await getTransportEfficiency()
    Object.assign(eff, result)
  } finally { loading.value = false }
}
</script>
