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
 * @LastEditTime : 2024-09-04 22:29:53
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

// Lrepresents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all short links",
	Long:  "This command lists all short links and their corresponding original URLs.",
	Run: func(cmd *cobra.Command, args []string) {
		listAllLinks()
	},
}

func listAllLinks() {
	// 打开数据库连接
	db, err := sql.Open("sqlite", utils.STORAGE_PATH+"url_shortener.db")
	if err != nil {
		fmt.Println("Failed to open the database:", err)
		return
	}
	defer db.Close()

	// 查询所有数据
	rows, err := db.Query("SELECT id, original_url FROM url_map")
	if err != nil {
		fmt.Println("Failed to query the database:", err)
		return
	}
	defer rows.Close()

	fmt.Println("Short Link ID | Original URL")
	fmt.Println("--------------------------------------")

	// 打印所有短链接数据
	for rows.Next() {
		var id, originalURL string
		if err := rows.Scan(&id, &originalURL); err != nil {
			fmt.Println("Error reading data:", err)
			return
		}
		fmt.Printf("%s | %s\n", id, originalURL)
	}
}
