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
 * @LastEditTime : 2024-06-08 12:27:01
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package utils

import (
	"github.com/gin-gonic/gin"
)

// Respond sends a JSON or text response based on the respondType.
func Respond(c *gin.Context, status int, message string, respondType string, data ...interface{}) {
	switch respondType {
	default:
	case "text":
		c.String(status, message)

	case "json":
		if len(data) > 0 {
			c.JSON(status, gin.H{
				"message": message,
				"data":    data[0],
			})
		} else {
			c.JSON(status, gin.H{
				"message": message,
			})
		}
	}
}

// RespondWithError sends a text error response.
func RespondWithError(c *gin.Context, status int, message string) {
	Respond(c, status, message, "text")
}

// RespondWithSuccess sends a text success response.
func RespondWithSuccess(c *gin.Context, status int, message string) {
	Respond(c, status, message, "text")
}

// RespondJSONWithError sends a JSON error response.
func RespondJSONWithError(c *gin.Context, status int, message string) {
	Respond(c, status, message, "json")
}

// RespondJSONWithSuccess sends a JSON success response.
func RespondJSONWithSuccess(c *gin.Context, status int, message string) {
	Respond(c, status, message, "json")
}

// RespondJSONWithData sends a JSON response with additional data.
func RespondJSONWithData(c *gin.Context, status int, message string, data interface{}) {
	Respond(c, status, message, "json", data)
}
