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
 * @LastEditTime : 2024-06-10 01:31:56
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package taskify

type AppUpdateModel struct {
	VersionCode   int    `json:"versionCode"`
	BuildVersion  int    `json:"buildVersion"`
	VersionName   string `json:"versionName"`
	UpdateMessage string `json:"updateMessage"`
	DownloadURL   string `json:"downloadUrl"`
}


type BugReportModel struct {
	Email    string `json:"email"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Content  string `json:"content"`
}