package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"map-tile-system/internal/model"
)

// ListApiKeys 获取 API Key 列表
func (h *Handler) ListApiKeys(c *gin.Context) {
	var apiKeys []model.ApiKey
	if err := h.db.Find(&apiKeys).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "获取 API Key 列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": apiKeys,
	})
}

// CreateApiKey 创建 API Key
func (h *Handler) CreateApiKey(c *gin.Context) {
	var req model.ApiKey
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	req.ID = uuid.New().String()
	req.Key = generateApiKey()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	if err := h.db.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "创建 API Key 失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    req,
		"message": "创建成功",
	})
}

// GetApiKey 获取单个 API Key
func (h *Handler) GetApiKey(c *gin.Context) {
	id := c.Param("id")
	var apiKey model.ApiKey

	if err := h.db.First(&apiKey, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API Key 不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": apiKey,
	})
}

// UpdateApiKey 更新 API Key
func (h *Handler) UpdateApiKey(c *gin.Context) {
	id := c.Param("id")
	var req model.ApiKey

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	req.ID = id
	req.UpdatedAt = time.Now()

	if err := h.db.Model(&model.ApiKey{}).Where("id = ?", id).Updates(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "更新 API Key 失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    req,
		"message": "更新成功",
	})
}

// DeleteApiKey 删除 API Key
func (h *Handler) DeleteApiKey(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&model.ApiKey{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "删除 API Key 失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

func generateApiKey() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
