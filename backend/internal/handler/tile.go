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

// ServeTileset 提供 3D Tiles tileset.json
func (h *Handler) ServeTileset(c *gin.Context) {
	dataset := c.Param("dataset")

	tilesetPath := filepath.Join("./data/tiles", dataset, "tileset.json")

	if _, err := os.Stat(tilesetPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tileset 不存在"})
		return
	}

	c.File(tilesetPath)
}

// Serve3DTiles 提供 3D Tiles 其他文件
func (h *Handler) Serve3DTiles(c *gin.Context) {
	dataset := c.Param("dataset")
	filePath := c.Param("filepath")

	fullPath := filepath.Join("./data/tiles", dataset, filePath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	c.File(fullPath)
}
