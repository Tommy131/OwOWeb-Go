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
 * @LastEditTime : 2024-07-19 14:48:04
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package main

import (
	"bufio"
	"os"
	"os/signal"
	"owoweb/a"
	"owoweb/cmd"
	"owoweb/i18n"
	"owoweb/routes"
	"owoweb/utils"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

// 主函数
func main() {
	// 激活头部包
	a.DoNothing()

	if !utils.DEBUG_MODE {
		gin.SetMode(gin.ReleaseMode)
	}

	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if len(os.Args) > 1 {
		cmd.Execute()
		return
	}

	go runWebServer(config.WebListeningAddress)
	log.Info(i18n.Lpk.FormatMessage("main.web_service_listening", aurora.BrightYellow(utils.CreateClickableLink("http://"+utils.WEB_ADDRESS))))

	done := make(chan bool, 1)
	go scanCommands(done)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigs:
		log.Info(i18n.Lpk.GetTranslate("main.forcibly_stopped_terminal"))
	case <-done:
		log.Info(i18n.Lpk.GetTranslate("main.stopping_terminal"))
	}
}

// runWebServer 启动Web服务
func runWebServer(address string) {
	router := gin.Default()
	router.Static("/static", utils.STORAGE_PATH+"static")
	router.LoadHTMLGlob(utils.STORAGE_PATH + "static/*.html")

	// 注册自定义路由
	routes.RegisterRouters(router)

	// 启动Web服务
	router.Run(address)
}

// scanCommands 扫描并执行指令
func scanCommands(done chan bool) {
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
				log.Error("Error:", err)
			}
		}
	}
}
