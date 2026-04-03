package main

import (
	"log"
	"map-tile-system/internal/handler"
	"map-tile-system/internal/middleware"
	"map-tile-system/pkg/config"
	"map-tile-system/pkg/database"
	"map-tile-system/pkg/logger"
	"map-tile-system/pkg/scheduler"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化日志
	logger.Init(cfg.Log.Level)

	// 初始化数据库
	db, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化默认管理员账户
	if err := database.InitDefaultUser(db); err != nil {
		log.Printf("Warning: Failed to initialize default user: %v", err)
	}

	// 初始化 Redis
	rdb := database.InitRedis(cfg.Redis)

	// 启动定时任务（每天凌晨 2 点汇总和清理）
	scheduler.Start(db)

	// 设置中间件数据库连接
	middleware.SetDB(db)

	// 创建 Gin 引擎
	r := gin.Default()

	// CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost")
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 日志中间件
	r.Use(middleware.Logger())

	// 初始化处理器
	h := handler.NewHandler(db, rdb)

	// 注册路由
	registerRoutes(r, h)

	// 启动服务器
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	logger.Info("Server starting on " + addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func registerRoutes(r *gin.Engine, h *handler.Handler) {
	api := r.Group("/api")
	{
		// 认证相关
		auth := api.Group("/auth")
		{
			auth.POST("/login", h.Login)
		}

		// 需要认证的路由
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			// 文件浏览
			authorized.GET("/files/drives", h.ListDrives)
			authorized.GET("/files/directories", h.ListDirectories)

			// 数据源管理
			datasource := authorized.Group("/datasources")
			{
				datasource.GET("", h.ListDataSources)
				datasource.POST("", h.CreateDataSource)
				datasource.GET("/:id", h.GetDataSource)
				datasource.PUT("/:id", h.UpdateDataSource)
				datasource.DELETE("/:id", h.DeleteDataSource)
			}

			// API Key 管理
			apikeys := authorized.Group("/apikeys")
			{
				apikeys.GET("", h.ListApiKeys)
				apikeys.POST("", h.CreateApiKey)
				apikeys.GET("/:id", h.GetApiKey)
				apikeys.PUT("/:id", h.UpdateApiKey)
				apikeys.DELETE("/:id", h.DeleteApiKey)
			}

			// 统计数据
			stats := authorized.Group("/statistics")
			{
				stats.GET("/overview", h.GetStatisticsOverview)
				stats.GET("/trend", h.GetStatisticsTrend)
				stats.GET("/details", h.GetStatisticsDetails)
				stats.GET("/top-keys", h.GetTopKeys)
				stats.POST("/backfill", h.BackfillSummary)
			}

			// 用户管理
			users := authorized.Group("/users")
			{
				users.GET("", h.ListUsers)
				users.POST("", h.CreateUser)
				users.GET("/:id", h.GetUser)
				users.PUT("/:id", h.UpdateUser)
				users.DELETE("/:id", h.DeleteUser)
				users.PUT("/:id/password", h.UpdatePassword)
			}
		}
	}

	// 内部预览瓦片路由（JWT 认证，不记录统计）
	// XYZ 瓦片预览
	previewXYZ := r.Group("/api/preview/xyz")
	previewXYZ.Use(middleware.PreviewAuthMiddleware())
	{
		previewXYZ.GET("/:dataset/:z/:x/:y", h.ServePreviewTile)
	}

	// 3D Tiles 预览
	preview3D := r.Group("/api/preview/3dtiles")
	preview3D.Use(middleware.PreviewAuthMiddleware())
	{
		preview3D.GET("/:dataset/*filepath", h.ServePreview3DTileFile)
	}

	// 瓦片服务路由
	tiles := r.Group("/tiles")
	tiles.Use(middleware.TileAuthMiddleware())
	{
		// 3D Tiles 服务（必须在 XYZ 之前注册）
		tiles.GET("/:dataset/tileset.json", h.ServeTileset)
		tiles.GET("/:dataset/3dtiles/*filepath", h.Serve3DTiles)
		// XYZ 瓦片服务
		tiles.GET("/:dataset/:z/:x/:y", h.ServeTile)
	}
}
