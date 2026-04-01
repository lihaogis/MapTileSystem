package scheduler

import (
	"log"
	"time"

	"gorm.io/gorm"
	"map-tile-system/internal/model"
)

// Start 启动定时任务
func Start(db *gorm.DB) {
	// 每天凌晨 2 点执行汇总和清理
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 2, 0, 0, 0, now.Location())
			duration := next.Sub(now)

			log.Printf("Next statistics summary and cleanup scheduled at: %s", next.Format("2006-01-02 15:04:05"))
			time.Sleep(duration)

			SummarizeAndCleanup(db)
		}
	}()
}

// SummarizeAndCleanup 汇总昨天数据并清理 180 天前的日志
func SummarizeAndCleanup(db *gorm.DB) {
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	SummarizeDate(db, yesterday)

	// 清理 180 天前的日志
	cutoffDate := time.Now().AddDate(0, 0, -180)
	result := db.Where("created_at < ?", cutoffDate).Delete(&model.CallLog{})
	log.Printf("Cleaned up %d old call logs (before %s)", result.RowsAffected, cutoffDate.Format("2006-01-02"))
}

// SummarizeDate 汇总指定日期的数据，已有记录则更新
func SummarizeDate(db *gorm.DB, date time.Time) {
	date = date.Truncate(24 * time.Hour)
	nextDay := date.AddDate(0, 0, 1)

	var totalCalls, successCalls int64
	db.Model(&model.CallLog{}).Where("created_at >= ? AND created_at < ?", date, nextDay).Count(&totalCalls)
	db.Model(&model.CallLog{}).Where("created_at >= ? AND created_at < ? AND status_code = ?", date, nextDay, 200).Count(&successCalls)

	summary := model.StatisticsSummary{
		Date:         date,
		TotalCalls:   totalCalls,
		SuccessCalls: successCalls,
		FailedCalls:  totalCalls - successCalls,
	}

	// 已存在则更新，不存在则插入
	db.Where(model.StatisticsSummary{Date: date}).Assign(summary).FirstOrCreate(&summary)
	log.Printf("Summarized %d calls for date: %s", totalCalls, date.Format("2006-01-02"))
}

// BackfillSummary 补跑指定天数范围内的历史汇总
func BackfillSummary(db *gorm.DB, days int) int {
	today := time.Now().Truncate(24 * time.Hour)
	count := 0
	for i := 1; i <= days; i++ {
		date := today.AddDate(0, 0, -i)
		SummarizeDate(db, date)
		count++
	}
	return count
}

