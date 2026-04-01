package database

import (
	"log"
	"map-tile-system/internal/model"
	"map-tile-system/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InitDefaultUser 初始化默认管理员账户
func InitDefaultUser(db *gorm.DB) error {
	// 检查是否已存在管理员账户
	var count int64
	db.Model(&model.User{}).Count(&count)

	if count > 0 {
		log.Println("Admin user already exists")
		return nil
	}

	// 创建默认管理员账户
	// 用户名: admin
	// 密码: admin123
	defaultUser := model.User{
		ID:       uuid.New().String(),
		Username: "admin",
		Password: utils.HashPassword("admin123"),
		Role:     "admin",
	}

	if err := db.Create(&defaultUser).Error; err != nil {
		return err
	}

	log.Println("Default admin user created successfully (username: admin, password: admin123)")
	return nil
}
