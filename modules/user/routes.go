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
 * @LastEditTime : 2024-10-01 00:50:55
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package user

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// 使用 CORS 中间件
	// 配置 CORS 允许前端跨域请求
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:8081"},                             // 允许来自 http://localhost:8081 的请求
		AllowMethods:     []string{"POST", "GET", "PUT", "OPTIONS", "Accept"},   // 允许的 HTTP 方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},                                    // 可以暴露的头
		AllowCredentials: true,                                                          // 允许带凭证请求
		MaxAge:           12 * time.Hour,                                                // 预检请求缓存12小时
	}))
	/* globalGroup := router.Group("/")
	{
		globalGroup.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", nil)
		})
	} */

	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/register", RegisterUser)
		userRoutes.POST("/login", LoginUser)
		userRoutes.GET("/verify-email", VerifyEmail)
		userRoutes.POST("/forgot-password", ForgotPassword)
		userRoutes.POST("/reset-password", ResetPassword)
	}

	// 受保护的路由
	protectedRoutes := router.Group("/user")
	protectedRoutes.Use(AuthMiddleware()) // 添加JWT验证中间件
	{
		userRoutes.GET("/profile", GetUserProfile)
		userRoutes.PUT("/profile", UpdateUserProfile)
		userRoutes.PUT("/settings", UpdateUserSettings)
	}
}
