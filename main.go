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
 * @LastEditTime : 2024-07-01 02:10:50
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"owoweb/cmd"
	"owoweb/i18n"
	"owoweb/modules/test"
	"owoweb/modules/user"
	"owoweb/routes"
	"owoweb/utils"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/logrusorgru/aurora"
)

// 欢迎标题
func init() {
	utils.EnableVirtualTerminalProcessing()
	name, _ := time.Now().Zone()

	fmt.Println(strings.Repeat("-", 50))
	fmt.Println(i18n.Lpk.FormatMessage("main.welcome_message", utils.ColorfulString("OwOWeb-Go")))
	fmt.Println(i18n.Lpk.FormatMessage("main.current_timezone", name, time.Now().Format("2006-01-02 15:04:05")))
	fmt.Println(i18n.Lpk.GetTranslate("main.unique_device_code"), aurora.Bold(aurora.BrightCyan(utils.GetUniqueDeviceCode())))
	fmt.Println(strings.Repeat("-", 50))

	// gin.SetMode(gin.ReleaseMode)
}

// 主函数
func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if len(os.Args) > 1 {
		cmd.Execute()
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		runWebServer(config.WebListeningAddress)
	}()
	fmt.Println(i18n.Lpk.FormatMessage("main.web_service_listening", aurora.BrightYellow(utils.CreateClickableLink("http://"+utils.WEB_ADDRESS))))

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			if !scanner.Scan() {
				break
			}

			command := strings.TrimSpace(scanner.Text())
			if command == "" {
				return
			}

			switch command {
			case "stop", "shutdown":
				done <- true
				return

			default:
				args := strings.Split(command, " ")
				if len(args) == 0 {
					return
				}

				cmd.RootCmd.SetArgs(args)
				if err := cmd.RootCmd.Execute(); err != nil {
					fmt.Println("Error:", err)
				}
			}
		}
	}()

	select {
	case <-sigs:
		fmt.Println("收到强制停止信号, 正在停止OwOWeb相关服务......")
	case <-done:
		fmt.Println("正在停止OwOWeb相关服务......")
	}
}

// 启动Web服务
func runWebServer(address string) {
	router := gin.Default()
	router.Static("/static", utils.STORAGE_PATH+"static")
	router.LoadHTMLFiles(utils.STORAGE_PATH + "static/index.html")

	// 注册各个模块的路由
	test.SetupRoutes(router)
	user.SetupRoutes(router)

	// 注册自定义路由
	routes.RegisterRouters(router)

	// 启动Web服务
	router.Run(address)
}
