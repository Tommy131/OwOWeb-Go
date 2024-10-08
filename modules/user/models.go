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
 * @LastEditTime : 2024-06-06 02:33:33
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package user

// User 用户结构
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// LoginCredentials 登录凭证结构
type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RecoverPasswordRequest 找回密码结构
type RecoverPasswordRequest struct {
	Email string `json:"email"`
}

// VerifyEmailRequest 验证邮箱结构
type VerifyEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
