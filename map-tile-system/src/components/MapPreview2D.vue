<template>
  <div class="relative w-full h-full">
    <div ref="mapContainer" class="w-full h-full"></div>
    <div class="absolute bottom-6 left-4 bg-white/90 px-3 py-2 rounded shadow text-xs font-mono z-[1000] pointer-events-none">
      <div>级别: {{ currentZoom }}</div>
      <div>坐标: {{ mouseCoords }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

const props = defineProps<{
  url: string
  format: string
  centerLat?: number
  centerLng?: number
  defaultZoom?: number
}>()

const mapContainer = ref<HTMLElement>()
const currentZoom = ref(1)
const mouseCoords = ref('---, ---')
let map: L.Map | null = null

const initMap = () => {
  if (!mapContainer.value) return

  // 使用数据源配置或默认值（北京，zoom 1）
  const center: [number, number] = [
    props.centerLat ?? 39.9,
    props.centerLng ?? 116.4
  ]
  const zoom = props.defaultZoom ?? 1

  map = L.map(mapContainer.value, {
    center,
    zoom,
    minZoom: 0,
    maxZoom: 18
  })

  currentZoom.value = zoom

  L.tileLayer(props.url, {
    maxZoom: 18,
    attribution: '© Map Tile System'
  }).addTo(map)

  // 监听缩放变化
  map.on('zoomend', () => {
    if (map) currentZoom.value = map.getZoom()
  })

  // 监听鼠标移动
  map.on('mousemove', (e: L.LeafletMouseEvent) => {
    const lat = e.latlng.lat.toFixed(6)
    const lng = e.latlng.lng.toFixed(6)
    mouseCoords.value = `${lat}, ${lng}`
  })

  // 鼠标离开地图
  map.on('mouseout', () => {
    mouseCoords.value = '---, ---'
  })

  // 延迟修正地图尺寸，确保容器已完全渲染
  nextTick(() => {
    setTimeout(() => {
      if (map) map.invalidateSize()
    }, 100)
  })
}

onMounted(() => {
  initMap()
})

onUnmounted(() => {
  if (map) {
    map.remove()
    map = null
  }
})

watch(() => props.url, () => {
  if (map) {
    map.remove()
  }
  initMap()
})
</script>
