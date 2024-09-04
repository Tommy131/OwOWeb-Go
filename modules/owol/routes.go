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
 * @LastEditTime : 2024-09-04 20:23:04
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package owol

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// 使用 CORS 中间件
	// 配置 CORS 允许前端跨域请求
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},            // 允许来自 http://localhost:3000 的请求
		AllowMethods:     []string{"POST", "GET", "POST", "OPTIONS"},   // 允许的 HTTP 方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},                   // 可以暴露的头
		AllowCredentials: true,                                         // 允许带凭证请求
		MaxAge:           12 * time.Hour,                               // 预检请求缓存12小时
	}))

	globalGroup := router.Group("/s")
	{
		globalGroup.GET("/:id", RedirectToOriginalURL)
		globalGroup.POST("/api/url-check", UrlCheck)
	}
}
