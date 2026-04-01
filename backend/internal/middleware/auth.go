package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"map-tile-system/internal/model"
)

var db *gorm.DB

func SetDB(database *gorm.DB) {
	db = database
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			c.Abort()
			return
		}

		// 移除 "Bearer " 前缀
		token = strings.TrimPrefix(token, "Bearer ")

		// TODO: 验证 JWT token
		// 这里暂时使用简单的 mock 验证
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// PreviewAuthMiddleware 预览鉴权中间件，接受 JWT token（query param 或 header），不参与调用统计
func PreviewAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			token = strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func TileAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Query("key")
		if apiKey == "" {
			apiKey = c.GetHeader("X-API-Key")
		}

		if apiKey == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "缺少 API Key"})
			c.Abort()
			return
		}

		var key model.ApiKey
		if err := db.Where("key = ? AND status = ?", apiKey, "enabled").First(&key).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "无效的 API Key"})
			c.Abort()
			return
		}

		c.Set("apiKey", key)
		c.Next()
	}
}
