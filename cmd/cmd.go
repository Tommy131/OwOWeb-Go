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
 * @LastEditTime : 2024-07-19 21:30:04
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

var RootCmd = &cobra.Command{
	Use:   "owoweb",
	Short: "OwOWeb CLI",
	Long:  `A command line tool to manage the OwO-Web system.`,
}

func init() {
	var helpCmd = &cobra.Command{
		Use:   "help",
		Short: "自定义帮助命令",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("这是自定义帮助命令的输出!")
			printCommands(RootCmd, 0)
		},
	}
	RootCmd.SetHelpCommand(helpCmd)
}

// printCommands 打印命令及其子命令
func printCommands(cmd *cobra.Command, level int) {
	indent := ""
	for i := 0; i < level; i++ {
		indent += "  "
	}
	fmt.Printf("%s%s\n", indent, cmd.Use)
	for _, c := range cmd.Commands() {
		printCommands(c, level+1)
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
