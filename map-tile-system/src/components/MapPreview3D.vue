<template>
  <div ref="cesiumContainer" class="w-full h-full"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import * as Cesium from 'cesium'
import 'cesium/Build/Cesium/Widgets/widgets.css'

const props = defineProps<{
  url: string
}>()

const cesiumContainer = ref<HTMLElement>()
let viewer: Cesium.Viewer | null = null

onMounted(() => {
  if (!cesiumContainer.value) return

  viewer = new Cesium.Viewer(cesiumContainer.value, {
    terrainProvider: undefined,
    baseLayerPicker: false,
    geocoder: false,
    homeButton: false,
    sceneModePicker: false,
    navigationHelpButton: false,
    animation: false,
    timeline: false
  })

  const tileset = viewer.scene.primitives.add(
    new Cesium.Cesium3DTileset({
      url: props.url
    })
  )

  viewer.zoomTo(tileset)
})

onUnmounted(() => {
  if (viewer) {
    viewer.destroy()
    viewer = null
  }
})
</script>
