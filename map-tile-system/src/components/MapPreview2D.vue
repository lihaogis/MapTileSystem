<template>
  <div class="relative w-full h-full">
    <div ref="mapContainer" class="w-full h-full"></div>
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
let infoControl: L.Control | null = null

const makeErrorTileUrl = () => {
  const size = 256
  const canvas = document.createElement('canvas')
  canvas.width = size
  canvas.height = size
  const ctx = canvas.getContext('2d')!
  ctx.fillStyle = '#f0f0f0'
  ctx.fillRect(0, 0, size, size)
  ctx.strokeStyle = '#ddd'
  ctx.strokeRect(0, 0, size, size)
  ctx.fillStyle = '#999'
  ctx.font = '13px sans-serif'
  ctx.textAlign = 'center'
  ctx.textBaseline = 'middle'
  ctx.fillText('此区域暂无地图', size / 2, size / 2)
  return canvas.toDataURL()
}

const initMap = () => {
  if (!mapContainer.value) return

  const center: [number, number] = [
    props.centerLat ?? 39.9,
    props.centerLng ?? 116.4
  ]
  const zoom = props.defaultZoom ?? 1

  map = L.map(mapContainer.value, {
    center,
    zoom,
    minZoom: 0,
    maxZoom: 18,
    attributionControl: false,
    zoomControl: true,
  })

  currentZoom.value = zoom

  // 自定义信息控件，固定在左下角
  const InfoControl = L.Control.extend({
    options: { position: 'bottomleft' },
    onAdd() {
      const div = L.DomUtil.create('div', '')
      div.style.cssText = 'background:rgba(255,255,255,0.9);padding:4px 8px;border-radius:4px;font:12px monospace;line-height:1.6;pointer-events:none;box-shadow:0 1px 4px rgba(0,0,0,0.2)'
      div.id = 'map-info-control'
      return div
    }
  })
  infoControl = new InfoControl()
  infoControl.addTo(map)

  const updateInfo = () => {
    const el = document.getElementById('map-info-control')
    if (el) {
      el.innerHTML = `级别: ${currentZoom.value}<br>坐标: ${mouseCoords.value}`
    }
  }

  L.tileLayer(props.url, {
    maxZoom: 18,
    errorTileUrl: makeErrorTileUrl(),
    crossOrigin: true
  }).addTo(map)

  map.on('zoomend', () => {
    if (map) currentZoom.value = map.getZoom()
    updateInfo()
  })

  map.on('mousemove', (e: L.LeafletMouseEvent) => {
    mouseCoords.value = `${e.latlng.lat.toFixed(6)}, ${e.latlng.lng.toFixed(6)}`
    updateInfo()
  })

  map.on('mouseout', () => {
    mouseCoords.value = '---, ---'
    updateInfo()
  })

  updateInfo()

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
    map = null
  }
  initMap()
})
</script>
