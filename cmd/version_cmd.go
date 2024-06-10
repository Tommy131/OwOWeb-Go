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
 * @LastEditTime : 2024-06-09 15:58:20
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package cmd

import (
	"fmt"
	"owoweb/utils"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"ver", "v"},
	Short:   "A command for checking OwOWeb's version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Welcome to use OwOWeb :) Currently version: %s%s\n", aurora.BrightBlue("v"), aurora.BrightGreen(utils.BUILD_VERSION))
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
