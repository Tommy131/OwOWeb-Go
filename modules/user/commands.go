/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-06-08 12:28:02
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2024-09-30 19:03:31
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package user

import (
	"database/sql"
	"fmt"
	"owoweb/cmd"
	"strings"

	"github.com/spf13/cobra"
)

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "User related commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use one of the subcommands: count, last-login, list")
	},
}

var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Get the count of registered users",
	Run: func(cmd *cobra.Command, args []string) {
		var count int
		err := db.Exec("SELECT COUNT(*) FROM users").Scan(&count)
		if err != nil {
			fmt.Printf("Failed to count users: %v\n", err)
			return
		}
		fmt.Printf("Total registered users: %d\n", count)
	},
}

var lastLoginCmd = &cobra.Command{
	Use:   "last-login",
	Short: "Get the last login time of a user",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a username")
			return
		}
		username := args[0]
		var lastLogin string
		err := db.Exec("SELECT last_login FROM users WHERE username = ?", username).Scan(&lastLogin)
		if err != nil {
			if err.Error == sql.ErrNoRows {
				fmt.Println("User not found")
				return
			}
			fmt.Printf("Failed to get last login time: %v\n", err)
			return
		}
		fmt.Printf("Last login time for %s: %s\n", username, lastLogin)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users with detailed information",
	Run: func(cmd *cobra.Command, args []string) {
		var users []User
		result := db.Select("id, username, email, last_login_at").Find(&users)
		if result.Error != nil {
			fmt.Printf("Failed to list users: %v\n", result.Error)
			return
		}

		// 打印表头
		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("%-5s %-20s %-30s %-25s\n", "ID", "Username", "Email", "Last Login")
		fmt.Println(strings.Repeat("-", 80))

		// 打印每个用户的信息
		for _, user := range users {
			fmt.Printf("%-5d %-20s %-30s %-20s\n", user.ID, user.Username, user.Email, user.LastLoginAt.Format("2006-01-02 15:04:05"))
		}

		// 检查迭代中的错误
		if result.Error != nil {
			fmt.Printf("Error occurred during row iteration: %v\n", result.Error)
		}
	},
}

func init() {
	UserCmd.AddCommand(countCmd)
	UserCmd.AddCommand(lastLoginCmd)
	UserCmd.AddCommand(listCmd)

	cmd.RootCmd.AddCommand(UserCmd)
}
