/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-06-09 23:29:04
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2024-09-05 00:27:47
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package user

import (
	"database/sql"
	"owoweb/i18n"
	"owoweb/utils"

	log "github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

var UserDb *sql.DB

func init() {
	var err error
	UserDb, err = sql.Open("sqlite", utils.DATABASE_PATH+"user_database.db")
	if err != nil {
		log.Fatal(err)
	}
	// defer UserDb.Close()

	createUserTable := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			last_login DATETIME
		);`
	if _, err := UserDb.Exec(createUserTable); err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	if err := UserDb.Ping(); err != nil {
		log.Fatalf("Failed to connect UserModule's Database: %v", err)
	} else {
		log.Info(i18n.Lpk.GetTranslate("module.user.init_database_successful"))
	}
}
