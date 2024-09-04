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
 * @LastEditTime : 2024-07-01 02:13:04
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package routes

import (
	"owoweb/modules/owol"
	"owoweb/modules/taskify"
	"owoweb/modules/test"
	"owoweb/modules/user"

	"github.com/gin-gonic/gin"
)

// 注册路由
func RegisterRouters(router *gin.Engine) {
	// 注册各个模块的路由
	test.SetupRoutes(router)
	user.SetupRoutes(router)
	taskify.SetupRoutes(router)
	owol.SetupRoutes(router)

	// an example
	registerCustomRouters(router)
}
