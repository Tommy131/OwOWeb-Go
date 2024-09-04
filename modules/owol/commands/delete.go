/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-09-04 22:19:32
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2024-09-04 22:19:43
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

// deleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete [shortlink ID]",
	Short: "Delete a short link by ID",
	Long:  "This command deletes a short link from the database using its ID.",
	Args:  cobra.ExactArgs(1), // 需要1个参数：短链接ID
	Run: func(cmd *cobra.Command, args []string) {
		deleteLink(args[0])
	},
}

func deleteLink(id string) {
	// 打开数据库连接
	db, err := sql.Open("sqlite", utils.STORAGE_PATH+"url_shortener.db")
	if err != nil {
		fmt.Println("Failed to open the database:", err)
		return
	}
	defer db.Close()

	// 删除短链接
	result, err := db.Exec("DELETE FROM url_map WHERE id = ?", id)
	if err != nil {
		fmt.Println("Failed to delete the short link:", err)
		return
	}

	// 检查删除的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error checking rows affected:", err)
		return
	}

	if rowsAffected > 0 {
		fmt.Println("Successfully deleted short link:", id)
	} else {
		fmt.Println("Short link not found:", id)
	}
}
