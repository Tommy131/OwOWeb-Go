/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-06-06 02:54:05
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2024-09-05 17:01:40
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package owol

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func UrlCheck(c *gin.Context) {
	var urlRequest URLRequest

	// 解析并绑定 JSON 请求数据
	if err := c.ShouldBindJSON(&urlRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid request, please provide a valid 'url' field",
		})
		return
	}

	// 检查是否已经存在相同的原始URL
	var existingID string
	err := db.QueryRow("SELECT id FROM url_map WHERE original_url = ?", urlRequest.URL).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query the database"})
		return
	}

	// 如果已经存在该URL，直接返回对应的短链ID
	if existingID != "" {
		shortURL := fmt.Sprintf(shareDomain+"%s", existingID)
		c.JSON(http.StatusOK, gin.H{
			"result":    true,
			"short_url": shortURL,
		})
		return
	}

	// 检查URL格式有效性
	if !isValidURL(urlRequest.URL) {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid URL format",
		})
		return
	}

	// 进行域名检查
	if isForbiddenURL(urlRequest.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is not allowed"})
		return
	}

	// 使用 curl 发送请求检测 URL
	client := http.Client{
		Timeout: 5 * time.Second, // 设置超时时间
	}

	var resp *http.Response
	resp, err = client.Get(urlRequest.URL)

	if err != nil || resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Unable to access URL",
		})
		return
	}
	defer resp.Body.Close()

	// 生成唯一短链ID
	var shortID string
	for {
		id, err := generateRandomID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": false,
				"error":  "Failed to generate short URL",
			})
			return
		}

		// 检查 ID 是否唯一
		var exists string
		err = db.QueryRow("SELECT id FROM url_map WHERE id = ?", id).Scan(&exists)
		if err == sql.ErrNoRows {
			// 如果不存在重复，则可以使用此 ID
			shortID = id
			break
		}
	}

	// 插入数据库保存原始URL和短链ID
	_, err = db.Exec("INSERT INTO url_map(id, original_url) VALUES(?, ?)", shortID, urlRequest.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": false,
			"error":  "Failed to generate short URL",
		})
		return
	}

	// 返回短链
	shortURL := fmt.Sprintf(shareDomain+"%s", shortID)
	c.JSON(http.StatusOK, gin.H{
		"result":    true,
		"short_url": shortURL,
	})
}

// 重定向到原始URL
func RedirectToOriginalURL(c *gin.Context) {
	id := c.Param("id")

	// 查找原始URL
	var originalURL string
	err := db.QueryRow("SELECT original_url FROM url_map WHERE id = ?", id).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			// 短链接ID不存在，返回自定义404页面
			c.HTML(http.StatusNotFound, "owol-404.html", nil)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query the database"})
		return
	}

	// 重定向到原始URL
	c.Redirect(http.StatusMovedPermanently, originalURL)
}

// 记录独立IP请求并返回全站访问次数
func VisitStats(c *gin.Context) {
	ip := c.ClientIP()
	updateVisitCount(ip)
	totalVisits, err := getTotalVisits()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total visits"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"totalVisits": totalVisits})
}
