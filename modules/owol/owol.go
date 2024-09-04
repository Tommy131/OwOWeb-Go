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
 * @LastEditTime : 2024-09-05 00:43:09
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */

package owol

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"owoweb/utils"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

var domain = "https://owol.cc/s/"

// 初始化数据库
func init() {
	var err error
	db, err = sql.Open("sqlite", utils.DATABASE_PATH+"owol_database.db")
	if err != nil {
		panic(err)
	}

	// 创建URL映射表
	sqlStmt := `CREATE TABLE IF NOT EXISTS url_map (
		id TEXT PRIMARY KEY,
		original_url TEXT NOT NULL
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		panic(err)
	}

	log.Info("Loaded OwOLink Services.")
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
