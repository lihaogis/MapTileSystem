<template>
  <div class="space-y-4">
    <el-card>
      <div class="flex items-center gap-4">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
        />
        <el-button type="primary">查询</el-button>
        <el-button>导出</el-button>
      </div>
    </el-card>

    <el-card>
      <template #header>
        <div class="font-semibold">调用趋势</div>
      </template>
      <div ref="chartRef" class="h-96"></div>
    </el-card>

    <el-card>
      <template #header>
        <div class="font-semibold">调用明细</div>
      </template>
      <el-table :data="callList" stripe>
        <el-table-column prop="date" label="日期" width="120" />
        <el-table-column prop="keyName" label="Key 名称" />
        <el-table-column prop="dataSource" label="数据源" />
        <el-table-column prop="count" label="调用次数" width="120" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import * as echarts from 'echarts'

const chartRef = ref<HTMLElement>()
const dateRange = ref<[Date, Date]>()
const callList = ref([
  { date: '2024-03-27', keyName: '生产环境', dataSource: '卫星影像', count: 12580 },
  { date: '2024-03-26', keyName: '测试环境', dataSource: '城市模型', count: 8234 }
])

onMounted(() => {
  if (chartRef.value) {
    const chart = echarts.init(chartRef.value)
    chart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: {
        type: 'category',
        data: ['3/21', '3/22', '3/23', '3/24', '3/25', '3/26', '3/27']
      },
      yAxis: { type: 'value' },
      series: [{
        name: '调用次数',
        data: [8200, 9320, 10100, 11340, 12900, 13000, 12580],
        type: 'bar'
      }]
    })
  }
})
</script>
