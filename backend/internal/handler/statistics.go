package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"map-tile-system/internal/model"
	"map-tile-system/pkg/scheduler"
)

// GetStatisticsOverview 获取统计概览
func (h *Handler) GetStatisticsOverview(c *gin.Context) {
	var todayTotal int64
	var yesterdayTotal int64
	var monthTotal int64
	var activeKeys int64

	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)
	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

	h.db.Model(&model.CallLog{}).Where("created_at >= ?", today).Count(&todayTotal)
	h.db.Model(&model.CallLog{}).Where("created_at >= ? AND created_at < ?", yesterday, today).Count(&yesterdayTotal)
	h.db.Model(&model.CallLog{}).Where("created_at >= ?", monthStart).Count(&monthTotal)
	h.db.Model(&model.ApiKey{}).Where("status = ?", "enabled").Count(&activeKeys)

	// 全部累计：所有 API Key 的 call_count 求和
	var totalAllTime int64
	h.db.Model(&model.ApiKey{}).Select("COALESCE(SUM(call_count), 0)").Scan(&totalAllTime)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"todayTotal":     todayTotal,
			"yesterdayTotal": yesterdayTotal,
			"monthTotal":     monthTotal,
			"activeKeys":     activeKeys,
			"totalAllTime":   totalAllTime,
		},
	})
}

// GetStatisticsTrend 获取统计趋势
func (h *Handler) GetStatisticsTrend(c *gin.Context) {
	days := 7
	data := make([]gin.H, days)

	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, -days+i+1).Truncate(24 * time.Hour)
		nextDate := date.AddDate(0, 0, 1)

		var count int64
		// 优先从 CallLog 查询（近期数据）
		h.db.Model(&model.CallLog{}).
			Where("created_at >= ? AND created_at < ?", date, nextDate).
			Count(&count)

		// 如果 CallLog 没有数据，从 StatisticsSummary 查询（历史数据）
		if count == 0 {
			var summary model.StatisticsSummary
			if err := h.db.Where("date = ?", date).First(&summary).Error; err == nil {
				count = summary.TotalCalls
			}
		}

		data[i] = gin.H{
			"date":  date.Format("2006-01-02"),
			"count": count,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": data,
	})
}

// BackfillSummary 补跑历史统计汇总
func (h *Handler) BackfillSummary(c *gin.Context) {
	// 默认补跑过去 180 天，可通过 ?days=N 指定
	days := 180
	if d := c.Query("days"); d != "" {
		if n, err := strconv.Atoi(d); err == nil && n > 0 && n <= 365 {
			days = n
		}
	}

	count := scheduler.BackfillSummary(h.db, days)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "补跑完成",
		"data":    gin.H{"days": count},
	})
}

// GetTopKeys 获取调用次数 Top 5 的 API Key
func (h *Handler) GetTopKeys(c *gin.Context) {
	type TopKeyResult struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Key       string `json:"key"`
		CallCount int64  `json:"callCount"`
	}

	// call_count 是历史累计（含已清理的明细），直接排序取 Top 5
	var results []TopKeyResult
	h.db.Model(&model.ApiKey{}).
		Select("id, name, key, call_count as call_count").
		Order("call_count DESC").
		Limit(5).
		Scan(&results)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": results,
	})
}

func (h *Handler) GetStatisticsDetails(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "20")
	apiKeyId := c.Query("apiKeyId")
	dataSourceId := c.Query("dataSourceId")
	ipAddress := c.Query("ipAddress")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	type LogWithNames struct {
		model.CallLog
		ApiKeyName     string `json:"apiKeyName"`
		DataSourceName string `json:"dataSourceName"`
	}

	query := h.db.Model(&model.CallLog{}).
		Select("call_logs.*, api_keys.name as api_key_name, data_sources.name as data_source_name").
		Joins("LEFT JOIN api_keys ON call_logs.api_key_id = api_keys.id").
		Joins("LEFT JOIN data_sources ON call_logs.data_source_id = data_sources.id")

	if apiKeyId != "" {
		query = query.Where("call_logs.api_key_id = ?", apiKeyId)
	}
	if dataSourceId != "" {
		query = query.Where("call_logs.data_source_id = ?", dataSourceId)
	}
	if ipAddress != "" {
		query = query.Where("call_logs.ip_address LIKE ?", "%"+ipAddress+"%")
	}
	if startDate != "" {
		query = query.Where("call_logs.created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("call_logs.created_at <= ?", endDate)
	}

	pageNum, _ := strconv.Atoi(page)
	pageSizeNum, _ := strconv.Atoi(pageSize)
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSizeNum < 1 {
		pageSizeNum = 20
	}
	offset := (pageNum - 1) * pageSizeNum

	var total int64
	query.Count(&total)

	var logs []LogWithNames
	query.Order("call_logs.created_at DESC").Limit(pageSizeNum).Offset(offset).Scan(&logs)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":     logs,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}
