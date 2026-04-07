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
  centerLat?: number
  centerLng?: number
  defaultZoom?: number
}>()

const mapContainer = ref<HTMLElement>()
const currentZoom = ref(5)
const mouseCoords = ref('---, ---')
let map: L.Map | null = null
let infoControl: L.Control | null = null
let geoJsonLayer: L.GeoJSON | null = null

const initMap = () => {
  if (!mapContainer.value) return

  const center: [number, number] = [
    props.centerLat ?? 35,
    props.centerLng ?? 105
  ]
  const zoom = props.defaultZoom ?? 5

  map = L.map(mapContainer.value, {
    center,
    zoom,
    minZoom: 0,
    maxZoom: 20,
    attributionControl: false,
    zoomControl: true,
  })

  currentZoom.value = zoom

  // 底图
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 20,
    crossOrigin: true,
  }).addTo(map)

  // 信息控件
  const InfoControl = L.Control.extend({
    options: { position: 'bottomleft' },
    onAdd() {
      const div = L.DomUtil.create('div', '')
      div.style.cssText = 'background:rgba(255,255,255,0.9);padding:4px 8px;border-radius:4px;font:12px monospace;line-height:1.6;pointer-events:none;box-shadow:0 1px 4px rgba(0,0,0,0.2)'
      div.id = 'geojson-info-control'
      return div
    }
  })
  infoControl = new InfoControl()
  infoControl.addTo(map)

  const updateInfo = () => {
    const el = document.getElementById('geojson-info-control')
    if (el) {
      el.innerHTML = `级别: ${currentZoom.value}<br>坐标: ${mouseCoords.value}`
    }
  }

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

  if (props.url) {
    loadGeoJSON()
  }

  nextTick(() => {
    setTimeout(() => {
      if (map) map.invalidateSize()
    }, 100)
  })
}

const loadGeoJSON = async () => {
  if (!map || !props.url) return

  if (geoJsonLayer) {
    map.removeLayer(geoJsonLayer)
    geoJsonLayer = null
  }

  try {
    const res = await fetch(props.url)
    if (!res.ok) return
    const data = await res.json()

    geoJsonLayer = L.geoJSON(data, {
      style: {
        color: '#3388ff',
        weight: 2,
        opacity: 0.8,
        fillOpacity: 0.3,
        fillColor: '#3388ff',
      },
      pointToLayer: (_feature, latlng) => {
        return L.circleMarker(latlng, {
          radius: 6,
          fillColor: '#e74c3c',
          color: '#c0392b',
          weight: 1.5,
          opacity: 1,
          fillOpacity: 0.9,
        })
      },
      onEachFeature: (feature, layer) => {
        if (feature.properties && Object.keys(feature.properties).length > 0) {
          const props = feature.properties
          const rows = Object.entries(props)
            .map(([k, v]) => `<tr><td style="padding:2px 8px;font-weight:600;color:#555">${k}</td><td style="padding:2px 8px">${v ?? ''}</td></tr>`)
            .join('')
          layer.bindPopup(`<table style="font-size:12px;min-width:160px">${rows}</table>`)
        }
      }
    }).addTo(map)

    // 自动缩放到数据范围（仅当未设置自定义中心时）
    if (!props.centerLat && !props.centerLng) {
      const bounds = geoJsonLayer.getBounds()
      if (bounds.isValid()) {
        map.fitBounds(bounds, { padding: [20, 20] })
      }
    }
  } catch (e) {
    console.error('加载 GeoJSON 失败', e)
  }
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
    loadGeoJSON()
  }
})
</script>
