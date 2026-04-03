<template>
  <div ref="cesiumContainer" class="w-full h-full"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as Cesium from 'cesium'
import 'cesium/Build/Cesium/Widgets/widgets.css'

const props = defineProps<{
  url: string
}>()

const cesiumContainer = ref<HTMLElement>()
let viewer: Cesium.Viewer | null = null
let tileset: Cesium.Cesium3DTileset | null = null

const loadTileset = async () => {
  if (!viewer || !props.url) return

  try {
    // 移除旧的 tileset
    if (tileset) {
      viewer.scene.primitives.remove(tileset)
      tileset = null
    }

    // 创建新的 tileset
    tileset = await Cesium.Cesium3DTileset.fromUrl(props.url, {
      skipLevelOfDetail: true,
      baseScreenSpaceError: 1024,
      skipScreenSpaceErrorFactor: 16,
      skipLevels: 1,
      immediatelyLoadDesiredLevelOfDetail: false,
      loadSiblings: false,
      cullWithChildrenBounds: true
    })

    viewer.scene.primitives.add(tileset)

    // 等待 tileset 加载完成后再缩放
    await tileset.readyPromise
    await viewer.zoomTo(tileset)

  } catch (error) {
    console.error('加载 3D Tiles 失败:', error)
  }
}

onMounted(async () => {
  if (!cesiumContainer.value) return

  // 设置 Cesium 默认访问令牌（如果需要）
  // Cesium.Ion.defaultAccessToken = 'your_token_here'

  viewer = new Cesium.Viewer(cesiumContainer.value, {
    terrainProvider: undefined,
    baseLayerPicker: false,
    geocoder: false,
    homeButton: false,
    sceneModePicker: false,
    navigationHelpButton: false,
    animation: false,
    timeline: false,
    fullscreenButton: false,
    vrButton: false,
    infoBox: false,
    selectionIndicator: false
  })

  // 加载 tileset
  await loadTileset()
})

// 监听 URL 变化，重新加载 tileset
watch(() => props.url, () => {
  loadTileset()
})

onUnmounted(() => {
  if (tileset && viewer) {
    viewer.scene.primitives.remove(tileset)
    tileset = null
  }
  if (viewer) {
    viewer.destroy()
    viewer = null
  }
})
</script>
