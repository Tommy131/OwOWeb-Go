/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-06-08 15:08:35
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2024-07-19 12:59:57
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/logrusorgru/aurora"
)

// GetCPUSerial 获取CPU序列号
func GetCPUSerial() (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("wmic", "cpu", "get", "ProcessorId")
	} else {
		cmd = exec.Command("sh", "-c", "lscpu | grep 'Serial' | awk '{print $2}'")
	}
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetDiskSerial 获取硬盘序列号
func GetDiskSerial() (string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("wmic", "diskdrive", "get", "SerialNumber")
	} else {
		cmd = exec.Command("sh", "-c", "lsblk -o SERIAL | grep -v 'SERIAL'")
	}
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// generateDeviceCode 生成设备唯一识别码
func generateDeviceCode() (string, error) {
	cpuSerial, err := GetCPUSerial()
	if err != nil {
		return "", err
	}

	diskSerial, err := GetDiskSerial()
	if err != nil {
		return "", err
	}

	data := fmt.Sprintf("%s%s", cpuSerial, diskSerial)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:]), nil
}

// GetUniqueDeviceCode 自动获取设备唯一识别码
func GetUniqueDeviceCode() string {
	deviceCode, err := generateDeviceCode()
	if err != nil {
		return "Error generating device code"
	}
	return deviceCode
}

// OpenBrowser 打开外部URL
func OpenBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		fmt.Printf("Unsupported platform\n")
		return
	}

	if err := exec.Command(cmd, args...).Start(); err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
	}
}

// CreateClickableLink 创建一个可点击打开的外部链接
func CreateClickableLink(url string) string {
	// ANSI escape codes for creating hyperlinks
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", url, url)
}

// EnableVirtualTerminalProcessing 启用在Windows环境下的ANSI
func EnableVirtualTerminalProcessing() {
	if runtime.GOOS == "windows" {
		const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004

		var kernel32 = syscall.NewLazyDLL("kernel32.dll")
		var procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
		var procSetConsoleMode = kernel32.NewProc("SetConsoleMode")

		var mode uint32
		handle := syscall.Handle(os.Stdout.Fd())
		procGetConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
		mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING
		procSetConsoleMode.Call(uintptr(handle), uintptr(mode))
	}
}

// GetRandomColor 获取随机颜色
func GetRandomColor() aurora.Color {
	colors := []aurora.Color{
		aurora.RedFg,
		aurora.GreenFg,
		aurora.YellowFg,
		aurora.BlueFg,
		aurora.MagentaFg,
		aurora.CyanFg,
	}
	return colors[rand.Intn(len(colors))]
}

// ColorfulString 赋予字符串随机颜色
func ColorfulString(input string) string {
	rand.Seed(time.Now().UnixNano())
	var result string

	for _, char := range input {
		color := GetRandomColor()
		result += fmt.Sprintf("%s", aurora.Colorize(string(char), color))
	}

	return result
}
