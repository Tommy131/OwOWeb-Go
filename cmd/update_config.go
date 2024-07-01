/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-06-11 15:57:50
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2024-06-11 15:58:00
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package cmd

import (
	"fmt"
	"log"
	"owoweb/utils"

	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:     "config",
	Short:   "Configuration related commands",
	Aliases: []string{"cfg", "c"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use one of the subcommands: update")
	},
}

var updateConfigCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update the configuration file",
	Aliases: []string{"u"},
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := utils.UpdateConfig(func(cfg *utils.Config) {
			if len(args) > 0 {
				cfg.WebListeningAddress = args[0]
			}
			if len(args) > 1 {
				cfg.TranslateLanguagePack = args[1]
			}
		})
		if err != nil {
			log.Fatalf("Failed to update config: %v", err)
		}

		fmt.Printf("Updated config: %+v\n", cfg)
	},
}

func init() {
	ConfigCmd.AddCommand(updateConfigCmd)

	RootCmd.AddCommand(ConfigCmd)
}
