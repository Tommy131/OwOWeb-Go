/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-07-19 14:20:15
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2024-09-04 23:32:05
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package a

import (
	"fmt"
	"io"
	"os"
	"owoweb/i18n"
	"owoweb/utils"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
)

type CustomLogger struct {
	Logger *lumberjack.Logger
	mu     sync.Mutex
}

// Write method to filter ANSI color codes
func (cl *CustomLogger) Write(p []byte) (n int, err error) {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	cleaned := strings.TrimSpace(string(p))

	ansiEscape := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	linkEscape := regexp.MustCompile(`\x1b\][^\x07]*\x07`)

	cleaned = ansiEscape.ReplaceAllString(cleaned, "")
	cleaned = linkEscape.ReplaceAllString(cleaned, "")

	cl.Logger.Write([]byte(cleaned + "\n"))
	return len(p), err
}

func init() {
	// 在Windows下启用ANSI
	utils.EnableVirtualTerminalProcessing()
	// 设置日志输出文件
	logger := &CustomLogger{
		Logger: &lumberjack.Logger{
			Filename:   utils.LOG_PATH + "app.log",
			MaxSize:    10, // MB
			MaxBackups: 3,
			MaxAge:     28,   // days
			Compress:   true, // Compress old log files
		},
	}

	// 结合输出方法
	log.SetOutput(io.MultiWriter(os.Stdout, logger))

	// 初始化日志记录格式
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
		DisableQuote:    true,
	})

	// 发送欢迎通知
	name, _ := time.Now().Zone()
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println(i18n.Lpk.FormatMessage("main.welcome_message", utils.ColorfulString("OwOWeb-Go")))
	fmt.Println(i18n.Lpk.FormatMessage("main.current_timezone", name, time.Now().Format("2006-01-02 15:04:05")))
	fmt.Println(i18n.Lpk.GetTranslate("main.unique_device_code"), aurora.Bold(aurora.BrightCyan(utils.GetUniqueDeviceCode())))
	fmt.Println(strings.Repeat("-", 50))
}

// DoNothing
func DoNothing() {
}
