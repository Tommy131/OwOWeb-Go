/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-09-04 21:48:39
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2024-09-05 16:28:41
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */

package owol

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"owoweb/utils"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

var shareDomain = "https://owol.cc/s/"

// 预定义的域名列表
var forbiddenDomains = []string{
	"http://localhost/",
	"http://127.0.0.1/",
	"https://owol.cc/",
}

// 初始化数据库
func init() {
	var err error
	db, err = sql.Open("sqlite", utils.DATABASE_PATH+"owol_database.db")
	if err != nil {
		panic(err)
	}

	// 创建 URL 映射表
	urlMapStmt := `CREATE TABLE IF NOT EXISTS url_map (
		id TEXT PRIMARY KEY,
		original_url TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(urlMapStmt)
	if err != nil {
		log.Fatalf("Error creating url_map table: %v", err)
	}

	// 创建访问统计表
	visitStatsStmt := `CREATE TABLE IF NOT EXISTS visit_stats (
		ip_address TEXT PRIMARY KEY,
		visit_count INTEGER DEFAULT 1,
		last_visited DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(visitStatsStmt)
	if err != nil {
		log.Fatalf("Error creating visit_stats table: %v", err)
	}

	log.Println("Loaded OwOLink Services.")
}

// 随机生成安全的短链 ID
func generateRandomID() (string, error) {
	// 生成 6 个随机字节
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Base64 编码生成 ID
	id := base64.URLEncoding.EncodeToString(b)

	// 移除等号
	id = strings.TrimRight(id, "=")
	return id, nil
}

// 检查URL格式是否有效
func isValidURL(url string) bool {
	// 定义URL的正则表达式
	regex := `^(https?|ftp):\/\/[^\s/$.?#].[^\s]*$`

	// 编译正则表达式
	re := regexp.MustCompile(regex)

	// 检查URL是否匹配正则表达式
	return re.MatchString(url)
}

// 检查 URL 是否与禁止的域名匹配
func isForbiddenURL(url string) bool {
	for _, domain := range forbiddenDomains {
		if strings.HasPrefix(url, domain) {
			return true
		}
	}
	return false
}

// 更新独立IP请求次数
func updateVisitCount(ip string) {
	_, err := db.Exec(`INSERT INTO visit_stats (ip_address, visit_count, last_visited)
	VALUES (?, 1, ?)
	ON CONFLICT(ip_address)
	DO UPDATE SET visit_count = visit_count + 1, last_visited = ?`, ip, time.Now(), time.Now())
	if err != nil {
		fmt.Println("Error updating visit count:", err)
	}
}

// 获取总共请求次数
func getTotalVisits() (int, error) {
	var totalVisits int
	rows, err := db.Query("SELECT SUM(visit_count) FROM visit_stats")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&totalVisits)
		if err != nil {
			return 0, err
		}
	}
	return totalVisits, nil
}
