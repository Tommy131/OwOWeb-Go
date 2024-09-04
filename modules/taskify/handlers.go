/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-02-04 14:10:45
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2024-07-19 17:48:25
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package taskify

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func getUpdateInfo() AppUpdateModel {
	_u := AppUpdateModel{
		VersionCode:   20240212,
		BuildVersion:  20240212,
		VersionName:   "0.0.4",
		UpdateMessage: defaultUpdateMsg,
		DownloadURL:   defaultDownloadURL,
	}
	return _u
}

// Bug report handler
func BugReport(c *gin.Context) {
	var bugReport BugReportModel
	if err := c.BindJSON(&bugReport); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	bugContentStr := "Taskify >> POST: Bug Report from Client [%s]:\nEmail: %s\nTitle: %s\nCategory: %s\nContent: %s\n"
	log.Info(fmt.Sprintf(bugContentStr, c.ClientIP(), bugReport.Email, bugReport.Title, bugReport.Category, bugReport.Content))

	c.JSON(http.StatusOK, map[string]string{"status": "success", "message": "Bug Report received successfully"})
}

// Update checker handler
func Update(c *gin.Context) {
	c.JSON(http.StatusOK, getUpdateInfo())
	log.Info(fmt.Sprintf("Taskify >> GET: Client %s requested to check update", c.ClientIP()))
}
