/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-06-09 23:29:04
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2024-10-01 00:08:00
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package user

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"owoweb/utils"
	"strconv"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

var (
	db           *gorm.DB
	jwtSecret    []byte
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	senderEmail  string
)

func init() {
	var err error
	// 使用 SQLite 作为数据库示例
	db, err = gorm.Open(sqlite.Open(utils.DATABASE_PATH+"user_database.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	// 自动迁移创建表格
	if err := db.AutoMigrate(&User{}, &EmailVerificationToken{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	smtpHost = os.Getenv("SMTP_HOST")
	smtpPort, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))
	smtpUsername = os.Getenv("SMTP_USERNAME")
	smtpPassword = os.Getenv("SMTP_PASSWORD")
	senderEmail = os.Getenv("SENDER_EMAIL")
}

// 发送邮件
func sendEmail(to string, emailTemplate EmailTemplate) error {
	// 读取模板文件内容
	tmplContent, err := os.ReadFile(utils.STORAGE_PATH + "email_template.html")
	if err != nil {
		log.Fatalf("Error reading template file: %v", err)
	}

	// 使用 html/template 解析模板
	t, err := template.New("emailTemplate").Parse(string(tmplContent))
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// 创建一个缓冲区来存储模板执行的结果
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, emailTemplate); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	mail := gomail.NewMessage()

	// 邮件主体
	mail.SetHeader("From", fmt.Sprintf("%s <%s>", smtpUsername, senderEmail))
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", emailTemplate.Subject)
	mail.SetBody("text/html", tpl.String())
	dialer := gomail.NewDialer(smtpHost, smtpPort, senderEmail, smtpPassword)
	if err := dialer.DialAndSend(mail); err != nil {
		return err
	}
	return nil
}

// 生成令牌的函数
func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// 生成JWT令牌
func generateJWT(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(), // 72小时后过期
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 解析JWT令牌
func parseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 检查签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 提取 claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// 检查用户名是否存在
func isUserIdExists(uid uint) (User, error) {
	var existingUser User
	err := db.Raw("SELECT * FROM owoblog_user WHERE id = ? LIMIT 1", uid).Scan(&existingUser).Error
	if err == nil && existingUser.ID != 0 {
		return existingUser, nil
	}
	return User{}, err
}

// 检查用户邮箱是否存在
func isEmailExists(email string) (User, error) {
	var existingUser User
	err := db.Raw("SELECT * FROM owoblog_user WHERE email = ? LIMIT 1", email).Scan(&existingUser).Error
	if err == nil && existingUser.ID != 0 {
		return existingUser, nil
	}
	return User{}, err
}

// 检查临时令牌是否存在
func isTokenExists(token string) (EmailVerificationToken, error) {
	var verificationToken EmailVerificationToken
	err := db.Raw("SELECT * FROM email_verification_tokens WHERE token = ? LIMIT 1", token).Scan(&verificationToken).Error
	if err == nil {
		return verificationToken, nil
	}
	return EmailVerificationToken{}, err
}
