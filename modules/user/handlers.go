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
 * @LastEditTime : 2024-06-10 00:22:09
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package user

import (
	"database/sql"
	"net/http"
	"owoweb/utils"

	"github.com/gin-gonic/gin"
)

func RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Email    string `json:"email"`
		}
		if err := c.BindJSON(&req); err != nil {
			utils.RespondJSONWithError(c, http.StatusBadRequest, "Invalid request")
			return
		}

		encryptedPassword, err := utils.EncryptPassword(req.Password)
		if err != nil {
			utils.RespondJSONWithError(c, http.StatusInternalServerError, "Encryption error")
			return
		}

		_, err = UserDb.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)", req.Username, encryptedPassword, req.Email)
		if err != nil {
			utils.RespondJSONWithError(c, http.StatusInternalServerError, "Database error: "+err.Error())
			return
		}

		utils.RespondJSONWithSuccess(c, http.StatusOK, "User registered successfully")
	}
}

func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&req); err != nil {
			utils.RespondJSONWithError(c, http.StatusBadRequest, "Invalid request")
			return
		}

		var storedPassword string
		err := UserDb.QueryRow("SELECT password FROM users WHERE username = ?", req.Username).Scan(&storedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.RespondJSONWithError(c, http.StatusUnauthorized, "User not found")
				return
			}
			utils.RespondJSONWithError(c, http.StatusInternalServerError, "Database error")
			return
		}

		if !utils.CheckPasswordHash(req.Password, storedPassword) {
			utils.RespondJSONWithError(c, http.StatusUnauthorized, "Invalid password")
			return
		}

		// Update last login time
		_, err = UserDb.Exec("UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE username = ?", req.Username)
		if err != nil {
			utils.RespondJSONWithError(c, http.StatusInternalServerError, "Failed to update last login time")
			return
		}

		utils.RespondJSONWithSuccess(c, http.StatusOK, "Login successful")
	}
}

func RecoverHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email"`
		}
		if err := c.BindJSON(&req); err != nil {
			utils.RespondJSONWithError(c, http.StatusBadRequest, "Invalid request")
			return
		}

		var username string
		err := UserDb.QueryRow("SELECT username FROM users WHERE email = ?", req.Email).Scan(&username)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.RespondJSONWithError(c, http.StatusNotFound, "Email not found")
				return
			}
			utils.RespondJSONWithError(c, http.StatusInternalServerError, "Database error")
			return
		}

		// Here you can add code to send a recovery email to the user.
		// For simplicity, we'll just return a success message.

		utils.RespondJSONWithSuccess(c, http.StatusOK, "Recovery email sent")
	}
}

func VerifyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email"`
			Code  string `json:"code"`
		}
		if err := c.BindJSON(&req); err != nil {
			utils.RespondJSONWithError(c, http.StatusBadRequest, "Invalid request")
			return
		}

		// Here you can add code to verify the email with the provided code.
		// For simplicity, we'll just return a success message.

		utils.RespondJSONWithSuccess(c, http.StatusOK, "Email verified")
	}
}
