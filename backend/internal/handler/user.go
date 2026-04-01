package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"map-tile-system/internal/model"
	"map-tile-system/pkg/utils"
)

// ListUsers 获取用户列表
func (h *Handler) ListUsers(c *gin.Context) {
	var users []model.User
	if err := h.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "获取用户列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": users,
	})
}

// CreateUser 创建用户
func (h *Handler) CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	var count int64
	h.db.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "用户名已存在",
		})
		return
	}

	user := model.User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Password: utils.HashPassword(req.Password),
		Role:     req.Role,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "创建用户失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    user,
		"message": "创建成功",
	})
}

// GetUser 获取单个用户
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	var user model.User

	if err := h.db.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1,
			"message": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": user,
	})
}

// UpdateUser 更新用户
func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 检查用户是否存在
	var user model.User
	if err := h.db.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1,
			"message": "用户不存在",
		})
		return
	}

	// 如果修改了用户名，检查新用户名是否已存在
	if req.Username != "" && req.Username != user.Username {
		var count int64
		h.db.Model(&model.User{}).Where("username = ? AND id != ?", req.Username, id).Count(&count)
		if count > 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "用户名已存在",
			})
			return
		}
		user.Username = req.Username
	}

	if req.Role != "" {
		user.Role = req.Role
	}

	user.UpdatedAt = time.Now()

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "更新用户失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    user,
		"message": "更新成功",
	})
}

// DeleteUser 删除用户
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// 检查用户是否存在
	var user model.User
	if err := h.db.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1,
			"message": "用户不存在",
		})
		return
	}

	// 不允许删除 admin 用户
	if user.Username == "admin" {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "不允许删除管理员账户",
		})
		return
	}

	if err := h.db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "删除用户失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// UpdatePassword 修改密码
func (h *Handler) UpdatePassword(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 检查用户是否存在
	var user model.User
	if err := h.db.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    1,
			"message": "用户不存在",
		})
		return
	}

	user.Password = utils.HashPassword(req.Password)
	user.UpdatedAt = time.Now()

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "修改密码失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "密码修改成功",
	})
}
