/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-09-05 18:09:58
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2024-09-06 00:29:55
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package ocsp

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// SetupDatabase 初始化数据库并自动迁移
func SetupDatabase(dbPath string) {
	var err error
	// 使用 SQLite 作为数据库示例
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	// 自动迁移创建表格
	if err := db.AutoMigrate(&CertificateStatus{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

// AddCertificateStatus 添加证书状态
func AddCertificateStatus(certStatus CertificateStatus) error {
	// 使用 GORM 创建记录
	if err := db.Create(&certStatus).Error; err != nil {
		return err
	}
	return nil
}

// GetCertificateStatus 根据序列号获取证书状态
func GetCertificateStatus(serialNumber string) (*CertificateStatus, error) {
	var certStatus CertificateStatus
	if err := db.First(&certStatus, "serial_number = ?", serialNumber).Error; err != nil {
		return nil, err
	}
	return &certStatus, nil
}
