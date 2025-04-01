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
 * @LastEditTime : 2024-09-30 23:31:16
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// JWT验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请求头中未找到令牌"})
			c.Abort()
			return
		}

		claims, err := parseJWT(tokenString)
		// 处理解析令牌的错误
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			c.Abort()
			return
		}

		userId, ok := claims["userId"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			c.Abort()
			return
		}

		c.Set("userId", uint(userId))
		c.Next() // 继续处理下一个请求
	}
}

// 注册功能（包括发送邮箱验证邮件）
func RegisterUser(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户是否已存在
	if _, err := isEmailExists(input.Email); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户已存在"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建新用户
	newUser := User{
		Username:    input.Email,
		Email:       input.Email,
		Password:    string(hashedPassword),
		Verified:    false,
		RegisterIP:  c.ClientIP(),
		LastLoginIP: c.ClientIP(),
		LastLoginAt: time.Now(),
		CreatedAt:   time.Now(),
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	// 生成邮箱验证令牌
	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成验证令牌失败"})
		return
	}

	verificationToken := EmailVerificationToken{
		UserID:    newUser.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 令牌24小时后过期
	}

	if err := db.Create(&verificationToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存验证令牌失败"})
		return
	}

	// 发送验证邮件
	verificationLink := fmt.Sprintf("http://127.0.0.1:8081/user/verify-email?token=%s", token)
	if err := sendEmail(newUser.Email, EmailTemplate{
		Subject:          "邮箱验证",
		Name:             newUser.Email,
		Body:             fmt.Sprintf("如果下方按钮点击后没有反应, 请点击以下链接验证您的邮箱: %s", verificationLink),
		VerificationLink: verificationLink,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送验证邮件失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功，请查收您的邮箱进行验证"})
}

// 邮箱验证功能
func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少验证令牌"})
		return
	}

	verificationToken, err := isTokenExists(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效或已过期的验证令牌"})
		return
	}

	// 检查令牌是否过期
	if verificationToken.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证令牌已过期"})
		return
	}

	// 更新用户为已验证
	user, err := isUserIdExists(verificationToken.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户不存在"})
		return
	}

	user.Verified = true
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户验证状态失败"})
		return
	}

	// 删除验证令牌
	db.Delete(&verificationToken)

	c.JSON(http.StatusOK, gin.H{"message": "邮箱验证成功"})
}

func LoginUser(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	foundUser, err := isEmailExists(input.Email)
	if err != nil {
		// 为了安全考虑, 不直接提示用户不存在
		c.JSON(http.StatusConflict, gin.H{"error": "账号或密码错误"})
		return
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}

	// 生成 JWT 令牌
	token, err := generateJWT(foundUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "userId": foundUser.ID})
}

// 密码找回功能（发送重置密码邮件）
func ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找用户
	user, err := isEmailExists(input.Email)
	if err != nil {
		// 为了安全性，不返回用户是否存在的信息
		c.JSON(http.StatusOK, gin.H{"message": "如果该邮箱存在，我们已发送密码重置链接"})
		return
	}

	// 生成密码重置令牌
	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成重置令牌失败"})
		return
	}

	resetToken := EmailVerificationToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 令牌1小时后过期
	}

	if err := db.Create(&resetToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存重置令牌失败"})
		return
	}

	// 发送重置密码邮件
	resetLink := fmt.Sprintf("http://your-domain.com/user/reset-password?token=%s", token)
	if err := sendEmail(user.Email, EmailTemplate{
		Subject:          "重置密码",
		Name:             user.Email,
		Body:             fmt.Sprintf("如果下方按钮点击后没有反应, 请点击以下链接重置您的密码: %s", resetLink),
		VerificationLink: resetLink,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送重置密码邮件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "如果该邮箱存在，我们已发送密码重置链接"})
}

// 重置密码功能
func ResetPassword(c *gin.Context) {
	var input struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找密码重置令牌
	resetToken, err := isTokenExists(input.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效或已过期的重置令牌"})
		return
	}

	// 检查令牌是否过期
	if resetToken.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "重置令牌已过期"})
		return
	}

	// 查找用户
	user, err := isUserIdExists(resetToken.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户不存在"})
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 更新用户密码
	user.Password = string(hashedPassword)
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
		return
	}

	// 删除重置令牌
	db.Delete(&resetToken)

	c.JSON(http.StatusOK, gin.H{"message": "密码已成功重置"})
}

// 获取用户个人资料
// main.go 中的 GetUserProfile 函数
func GetUserProfile(c *gin.Context) {
	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取用户ID"})
		return
	}

	userIdFloat, ok := userIdInterface.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID"})
		return
	}
	userId := uint(userIdFloat)

	user, err := isUserIdExists(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 返回用户信息，不包括密码
	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"email":     user.Email,
		"verified":  user.Verified,
		"username":  user.Username,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	})
}

// 更新用户个人资料
func UpdateUserProfile(c *gin.Context) {
	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取用户ID"})
		return
	}

	userId, ok := userIdInterface.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID"})
		return
	}

	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := isUserIdExists(uint(userId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 检查新邮箱是否已被使用
	existingUser, err := isEmailExists(input.Email)
	if err == nil && existingUser.ID != user.ID {
		c.JSON(http.StatusConflict, gin.H{"error": "该邮箱已被其他用户使用"})
		return
	}

	user.Email = input.Email
	user.Verified = false // 更新邮箱需要重新验证

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户资料失败"})
		return
	}

	// 生成新的邮箱验证令牌
	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成验证令牌失败"})
		return
	}

	verificationToken := EmailVerificationToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := db.Create(&verificationToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存验证令牌失败"})
		return
	}

	// 发送验证邮件
	verificationLink := fmt.Sprintf("http://your-domain.com/user/verify-email?token=%s", token)
	if err := sendEmail(user.Email, EmailTemplate{
		Subject:          "邮箱验证",
		Name:             user.Email,
		Body:             fmt.Sprintf("如果下方按钮点击后没有反应, 请点击以下链接验证您的邮箱: %s", verificationLink),
		VerificationLink: verificationLink,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送验证邮件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户资料已更新，请查收您的新邮箱进行验证"})
}

// 更新用户设置
func UpdateUserSettings(c *gin.Context) {
	userIdInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取用户ID"})
		return
	}

	userId, ok := userIdInterface.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无效的用户ID"})
		return
	}

	var input struct {
		Username string `json:"username" binding:"required"`
		// 其他设置字段...
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := isUserIdExists(uint(userId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	user.Username = input.Username
	// 更新其他设置字段...

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户设置已更新"})
}
