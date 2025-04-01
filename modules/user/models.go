/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-06-06 02:24:56
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2024-09-30 22:20:17
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package user

import (
	"time"
)

// User模型
type User struct {
	ID          uint   `gorm:"primaryKey"`
	Username    string `gorm:"unique;not null" json:"username"`
	Email       string `gorm:"unique;not null" json:"email"`
	Password    string `json:"password"`
	Verified    bool   `json:"verified"`
	RegisterIP  string
	LastLoginIP string
	LastLoginAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName 方法返回自定义的表名
func (User) TableName() string {
	return "owoblog_user"
}

// 用于验证邮箱地址的输入结构体
type EmailVerificationToken struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	UserID    uint      `gorm:"index" json:"userId"`
	Token     string    `gorm:"unique;not null" json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time
}

// 邮件模板结构体
type EmailTemplate struct {
	Subject          string
	Name             string
	Body             string
	VerificationLink string
}

// 用于登录和注册的输入结构体
type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
