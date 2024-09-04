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
 * @LastEditTime : 2024-09-05 00:49:21
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package owol

import (
	"database/sql"
	"fmt"
	"os"
	"owoweb/cmd"

	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

var db *sql.DB

var shortLinkCmd = &cobra.Command{
	Use:     "shortlink",
	Aliases: []string{"s", "sl"},
	Short:   "Shortlink CLI is used to manage URL shortening service",
	Long:    "This CLI tool allows you to list, delete, and update short links in the URL shortening service.",
}

// 执行命令
func Execute() {
	if err := shortLinkCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [shortLink ID]",
	Short: "Delete a short link by ID",
	Long:  "This command deletes a short link from the database using its ID.",
	Args:  cobra.ExactArgs(1), // 需要1个参数：短链接ID
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
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
	},
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all short links",
	Long:  "This command lists all short links and their corresponding original URLs.",
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [ShortLink ID] [new original URL]",
	Short: "Update the original URL for a short link",
	Long:  "This command updates the original URL for a given short link ID.",
	Args:  cobra.ExactArgs(2), // 需要2个参数：短链接ID和新的原始URL
	Run: func(cmd *cobra.Command, args []string) {
		id, newURL := args[0], args[1]
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
	},
}

func init() {
	// 注册子命令
	shortLinkCmd.AddCommand(listCmd)
	shortLinkCmd.AddCommand(deleteCmd)
	shortLinkCmd.AddCommand(updateCmd)

	cmd.RootCmd.AddCommand(shortLinkCmd)
}
