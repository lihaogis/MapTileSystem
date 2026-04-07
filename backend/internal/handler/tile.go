package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"map-tile-system/internal/model"
)

// ServePreviewTile 内部预览专用，使用 JWT 认证，不记录调用统计
func (h *Handler) ServePreviewTile(c *gin.Context) {
	dataset := c.Param("dataset")
	z := c.Param("z")
	x := c.Param("x")
	y := c.Param("y")

	var dataSource model.DataSource
	if err := h.db.First(&dataSource, "id = ?", dataset).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "数据源不存在"})
		return
	}

	if dataSource.Status != "enabled" {
		c.JSON(http.StatusForbidden, gin.H{"error": "数据源已禁用"})
		return
	}

	tilePath := filepath.Join(dataSource.Path, z, x, y+"."+dataSource.Format)

	if _, err := os.Stat(tilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "瓦片不存在"})
		return
	}

	c.File(tilePath)
}

// ServePreviewGeoJSON 内部预览 GeoJSON 文件（JWT 认证，不记录统计）
func (h *Handler) ServePreviewGeoJSON(c *gin.Context) {
	dataset := c.Param("dataset")

	var dataSource model.DataSource
	if err := h.db.First(&dataSource, "id = ?", dataset).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "数据源不存在"})
		return
	}

	if dataSource.Status != "enabled" {
		c.JSON(http.StatusForbidden, gin.H{"error": "数据源已禁用"})
		return
	}

	if _, err := os.Stat(dataSource.Path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	c.File(dataSource.Path)
}

// ServeTile 提供 XYZ 瓦片服务
func (h *Handler) ServeTile(c *gin.Context) {
	startTime := time.Now()
	dataset := c.Param("dataset")
	z := c.Param("z")
	x := c.Param("x")
	y := c.Param("y")

	// 从中间件获取 API Key
	apiKeyValue, _ := c.Get("apiKey")
	apiKey := apiKeyValue.(model.ApiKey)

	// 查询数据源配置
	var dataSource model.DataSource
	if err := h.db.First(&dataSource, "id = ?", dataset).Error; err != nil {
		h.logCall(c, apiKey.ID, dataset, z, x, y, http.StatusNotFound, startTime)
		c.JSON(http.StatusNotFound, gin.H{"error": "数据源不存在"})
		return
	}

	// 检查数据源状态
	if dataSource.Status != "enabled" {
		h.logCall(c, apiKey.ID, dataset, z, x, y, http.StatusForbidden, startTime)
		c.JSON(http.StatusForbidden, gin.H{"error": "数据源已禁用"})
		return
	}

	// 构建瓦片文件路径
	tilePath := filepath.Join(dataSource.Path, z, x, y+"."+dataSource.Format)

	// 检查文件是否存在
	if _, err := os.Stat(tilePath); os.IsNotExist(err) {
		h.logCall(c, apiKey.ID, dataset, z, x, y, http.StatusNotFound, startTime)
		c.JSON(http.StatusNotFound, gin.H{"error": "瓦片不存在"})
		return
	}

	h.logCall(c, apiKey.ID, dataset, z, x, y, http.StatusOK, startTime)
	c.File(tilePath)
}

// ServeGeoJSON 提供 GeoJSON 文件服务（需 API Key 鉴权，记录日志）
func (h *Handler) ServeGeoJSON(c *gin.Context) {
	startTime := time.Now()
	dataset := c.Param("dataset")

	apiKeyValue, _ := c.Get("apiKey")
	apiKey := apiKeyValue.(model.ApiKey)

	var dataSource model.DataSource
	if err := h.db.First(&dataSource, "id = ?", dataset).Error; err != nil {
		h.logCall(c, apiKey.ID, dataset, "", "", "", http.StatusNotFound, startTime)
		c.JSON(http.StatusNotFound, gin.H{"error": "数据源不存在"})
		return
	}

	if dataSource.Status != "enabled" {
		h.logCall(c, apiKey.ID, dataset, "", "", "", http.StatusForbidden, startTime)
		c.JSON(http.StatusForbidden, gin.H{"error": "数据源已禁用"})
		return
	}

	if _, err := os.Stat(dataSource.Path); os.IsNotExist(err) {
		h.logCall(c, apiKey.ID, dataset, "", "", "", http.StatusNotFound, startTime)
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	h.logCall(c, apiKey.ID, dataset, "", "", "", http.StatusOK, startTime)
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	c.File(dataSource.Path)
}

// logCall 记录调用日志
func (h *Handler) logCall(c *gin.Context, apiKeyID, dataSourceID, z, x, y string, statusCode int, startTime time.Time) {
	tileZ, _ := strconv.Atoi(z)
	tileX, _ := strconv.Atoi(x)
	tileY, _ := strconv.Atoi(y)
	responseTime := int(time.Since(startTime).Milliseconds())

	log := model.CallLog{
		ApiKeyID:     apiKeyID,
		DataSourceID: dataSourceID,
		TileZ:        tileZ,
		TileX:        tileX,
		TileY:        tileY,
		StatusCode:   statusCode,
		ResponseTime: responseTime,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.GetHeader("User-Agent"),
	}

	h.db.Create(&log)

	// 更新 API Key 调用计数
	h.db.Model(&model.ApiKey{}).Where("id = ?", apiKeyID).UpdateColumn("call_count", gorm.Expr("call_count + ?", 1))
}

// ServePreview3DTileFile 内部预览 3D Tiles 文件（包括 tileset.json, *.json, *.glb 等）
func (h *Handler) ServePreview3DTileFile(c *gin.Context) {
	dataset := c.Param("dataset")
	filePath := c.Param("filepath")

	// 移除开头的斜杠
	if len(filePath) > 0 && filePath[0] == '/' {
		filePath = filePath[1:]
	}

	var dataSource model.DataSource
	if err := h.db.First(&dataSource, "id = ?", dataset).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "数据源不存在"})
		return
	}

	if dataSource.Status != "enabled" {
		c.JSON(http.StatusForbidden, gin.H{"error": "数据源已禁用"})
		return
	}

	fullPath := filepath.Join(dataSource.Path, filePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在: " + filePath})
		return
	}

	c.File(fullPath)
}

// ServePreviewTileset 内部预览 3D Tiles tileset.json（JWT 认证，不记录统计）
func (h *Handler) ServePreviewTileset(c *gin.Context) {
	dataset := c.Param("dataset")

	var dataSource model.DataSource
	if err := h.db.First(&dataSource, "id = ?", dataset).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "数据源不存在"})
		return
	}

	if dataSource.Status != "enabled" {
		c.JSON(http.StatusForbidden, gin.H{"error": "数据源已禁用"})
		return
	}

	tilesetPath := filepath.Join(dataSource.Path, "tileset.json")
	if _, err := os.Stat(tilesetPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "tileset.json 不存在"})
		return
	}

	c.File(tilesetPath)
}

// ServePreview3DTiles 内部预览 3D Tiles 其他文件（JWT 认证，不记录统计）
func (h *Handler) ServePreview3DTiles(c *gin.Context) {
	dataset := c.Param("dataset")
	filePath := c.Param("filepath")

	var dataSource model.DataSource
	if err := h.db.First(&dataSource, "id = ?", dataset).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "数据源不存在"})
		return
	}

	if dataSource.Status != "enabled" {
		c.JSON(http.StatusForbidden, gin.H{"error": "数据源已禁用"})
		return
	}

	fullPath := filepath.Join(dataSource.Path, filePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	c.File(fullPath)
}

// ServeTileset 提供 3D Tiles tileset.json
func (h *Handler) ServeTileset(c *gin.Context) {
	startTime := time.Now()
	dataset := c.Param("dataset")

	apiKeyValue, _ := c.Get("apiKey")
	apiKey := apiKeyValue.(model.ApiKey)

	var dataSource model.DataSource
	if err := h.db.First(&dataSource, "id = ?", dataset).Error; err != nil {
		h.logCall(c, apiKey.ID, dataset, "", "", "", http.StatusNotFound, startTime)
		c.JSON(http.StatusNotFound, gin.H{"error": "数据源不存在"})
		return
	}

	if dataSource.Status != "enabled" {
		h.logCall(c, apiKey.ID, dataset, "", "", "", http.StatusForbidden, startTime)
		c.JSON(http.StatusForbidden, gin.H{"error": "数据源已禁用"})
		return
	}

	tilesetPath := filepath.Join(dataSource.Path, "tileset.json")

	if _, err := os.Stat(tilesetPath); os.IsNotExist(err) {
		h.logCall(c, apiKey.ID, dataset, "", "", "", http.StatusNotFound, startTime)
		c.JSON(http.StatusNotFound, gin.H{"error": "Tileset 不存在"})
		return
	}

	h.logCall(c, apiKey.ID, dataset, "", "", "", http.StatusOK, startTime)
	c.File(tilesetPath)
}

// Serve3DTiles 提供 3D Tiles 其他文件
func (h *Handler) Serve3DTiles(c *gin.Context) {
	startTime := time.Now()
	dataset := c.Param("dataset")
	filePath := c.Param("filepath")

	apiKeyValue, _ := c.Get("apiKey")
	apiKey := apiKeyValue.(model.ApiKey)

	var dataSource model.DataSource
	if err := h.db.First(&dataSource, "id = ?", dataset).Error; err != nil {
		h.logCall(c, apiKey.ID, dataset, "", "", "", http.StatusNotFound, startTime)
		c.JSON(http.StatusNotFound, gin.H{"error": "数据源不存在"})
		return
	}

	if dataSource.Status != "enabled" {
		h.logCall(c, apiKey.ID, dataset, "", "", "", http.StatusForbidden, startTime)
		c.JSON(http.StatusForbidden, gin.H{"error": "数据源已禁用"})
		return
	}

	fullPath := filepath.Join(dataSource.Path, filePath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		h.logCall(c, apiKey.ID, dataset, "", "", "", http.StatusNotFound, startTime)
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	h.logCall(c, apiKey.ID, dataset, "", "", "", http.StatusOK, startTime)
	c.File(fullPath)
}
