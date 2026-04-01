package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"map-tile-system/internal/model"
	"map-tile-system/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Handler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewHandler(db *gorm.DB, rdb *redis.Client) *Handler {
	return &Handler{
		db:  db,
		rdb: rdb,
	}
}

// Login 登录
func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
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

	ctx := context.Background()
	lockKey := fmt.Sprintf("login_lock:%s", req.Username)
	failKey := fmt.Sprintf("login_fail:%s", req.Username)

	// 检查是否被锁定
	locked, err := h.rdb.Get(ctx, lockKey).Result()
	if err == nil && locked == "1" {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "登录失败次数过多，请30分钟后再试",
		})
		return
	}

	// 查询用户
	var user model.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "用户名或密码错误",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "服务器错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证密码（前端已SHA256加密）
	hashedPassword := sha256.Sum256([]byte(user.Password))
	expectedPassword := hex.EncodeToString(hashedPassword[:])

	if req.Password != expectedPassword {
		// 增加失败次数
		failCount, _ := h.rdb.Incr(ctx, failKey).Result()
		h.rdb.Expire(ctx, failKey, 30*time.Minute)

		if failCount >= 5 {
			h.rdb.Set(ctx, lockKey, "1", 30*time.Minute)
			c.JSON(http.StatusOK, gin.H{
				"code":    1,
				"message": "登录失败次数过多，账户已锁定30分钟",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": fmt.Sprintf("用户名或密码错误，还可尝试%d次", 5-failCount),
		})
		return
	}

	// 登录成功，清除失败记录
	h.rdb.Del(ctx, failKey)

	// 生成 token
	token := utils.GenerateToken(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"role":     user.Role,
			},
		},
		"message": "登录成功",
	})
}
