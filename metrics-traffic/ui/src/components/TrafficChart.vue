<template><div ref="element" class="traffic-chart" /></template>

<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import * as echarts from 'echarts'

const props = defineProps({ series: { type: Array, default: () => [] }, metric: String, format: Function })
const element = ref()
let chart

const metricFields = {
  requests: [['requests_total', '总请求'], ['requests_http1', 'HTTP/1'], ['requests_http2', 'HTTP/2'], ['requests_http3', 'HTTP/3'], ['requests_https', 'HTTPS']],
  traffic: [['traffic_bytes', '总流量']],
  bandwidth: [['bandwidth_bps', '平均带宽']],
  hitRate: [['hit_rate_2xx', '2xx 占比'], ['hit_rate_3xx', '3xx 占比'], ['hit_rate_4xx', '4xx 占比'], ['hit_rate_5xx', '5xx 占比']],
  hits: [['hits_total', '总请求'], ['hits_2xx', '2xx'], ['hits_3xx', '3xx'], ['hits_4xx', '4xx'], ['hits_5xx', '5xx']],
}
function render() {
  if (!chart) return
  const fields = metricFields[props.metric] || metricFields.requests
  const values = props.series || []
  chart.setOption({
    color: ['#165dff', '#00b42a', '#f7ba1e', '#722ed1', '#f53f3f', '#86909c'],
    tooltip: { trigger: 'axis', appendToBody: true, valueFormatter: value => props.format?.(value) ?? value },
    legend: { type: 'scroll', bottom: 0 }, grid: { left: 20, right: 20, top: 42, bottom: 52, containLabel: true },
    xAxis: { type: 'time', boundaryGap: false, axisLabel: { color: '#86909c' }, axisLine: { lineStyle: { color: '#e5e6eb' } } },
    yAxis: { type: 'value', axisLabel: { color: '#86909c', formatter: value => props.format?.(value) ?? value }, axisLine: { lineStyle: { color: '#e5e6eb' } }, splitLine: { show: true, lineStyle: { color: '#f2f3f5' } } },
    series: fields.map(([field, name]) => ({ name, type: 'line', smooth: true, showSymbol: false, symbol: 'circle', symbolSize: 6, areaStyle: { opacity: 0.05 }, data: values.map(item => [item._time, Number(item[field] || 0)]) })),
  }, true)
}
onMounted(() => { chart = echarts.init(element.value); render(); window.addEventListener('resize', chart.resize) })
onBeforeUnmount(() => { window.removeEventListener('resize', chart?.resize); chart?.dispose() })
watch(() => [props.series, props.metric], render, { deep: true })
</script>
