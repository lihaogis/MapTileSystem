package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"map-tile-system/internal/model"
)

// ListDataSources 获取数据源列表
func (h *Handler) ListDataSources(c *gin.Context) {
	var dataSources []model.DataSource
	if err := h.db.Find(&dataSources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "获取数据源列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataSources,
	})
}

// CreateDataSource 创建数据源
func (h *Handler) CreateDataSource(c *gin.Context) {
	var req model.DataSource
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	req.ID = uuid.New().String()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	if err := h.db.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "创建数据源失败",
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

// GetDataSource 获取单个数据源
func (h *Handler) GetDataSource(c *gin.Context) {
	id := c.Param("id")
	var dataSource model.DataSource

	if err := h.db.First(&dataSource, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "数据源不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": dataSource,
	})
}

// UpdateDataSource 更新数据源
func (h *Handler) UpdateDataSource(c *gin.Context) {
	id := c.Param("id")
	var req model.DataSource

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

	if err := h.db.Model(&model.DataSource{}).Where("id = ?", id).Updates(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "更新数据源失败",
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

// DeleteDataSource 删除数据源
func (h *Handler) DeleteDataSource(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&model.DataSource{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "删除数据源失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}
