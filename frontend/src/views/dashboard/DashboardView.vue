<template>
  <section class="dashboard">
    <div class="dashboard__header">
      <h1>TMS 工作台</h1>
      <p>项目骨架已就绪，后续可在这里接入订单、调度、运输和结算数据。</p>
    </div>

    <el-row :gutter="16">
      <el-col :span="6" v-for="item in metrics" :key="item.label">
        <el-card shadow="never">
          <div class="metric">
            <span>{{ item.label }}</span>
            <strong>{{ item.value }}</strong>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </section>
</template>

<script setup lang="ts">
const metrics = [
  { label: '待调度订单', value: 0 },
  { label: '运输中任务', value: 0 },
  { label: '待处理异常', value: 0 },
  { label: '待结算费用', value: 0 },
]
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.dashboard__header {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.dashboard__header h1 {
  margin: 0;
  font-size: 24px;
}

.dashboard__header p {
  margin: 0;
  color: #64748b;
}

.metric {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.metric span {
  color: #64748b;
}

.metric strong {
  font-size: 28px;
}
</style>
