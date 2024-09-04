/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-09-04 22:20:56
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2024-09-04 22:20:56
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package commands

import (
	"database/sql"
	"fmt"
	"owoweb/utils"

	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

// updateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:   "update [shortlink ID] [new original URL]",
	Short: "Update the original URL for a short link",
	Long:  "This command updates the original URL for a given short link ID.",
	Args:  cobra.ExactArgs(2), // 需要2个参数：短链接ID和新的原始URL
	Run: func(cmd *cobra.Command, args []string) {
		updateLink(args[0], args[1])
	},
}

func updateLink(id string, newURL string) {
	// 打开数据库连接
	db, err := sql.Open("sqlite", utils.STORAGE_PATH+"url_shortener.db")
	if err != nil {
		fmt.Println("Failed to open the database:", err)
		return
	}
	defer db.Close()

	// 更新短链接对应的原始URL
	result, err := db.Exec("UPDATE url_map SET original_url = ? WHERE id = ?", newURL, id)
	if err != nil {
		fmt.Println("Failed to update the short link:", err)
		return
	}

	// 检查修改的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error checking rows affected:", err)
		return
	}

	if rowsAffected > 0 {
		fmt.Println("Successfully updated short link:", id)
	} else {
		fmt.Println("Short link not found:", id)
	}
}
