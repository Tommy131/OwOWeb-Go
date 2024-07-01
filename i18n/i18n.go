/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-06-11 14:31:27
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2024-07-01 01:52:51
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package i18n

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"owoweb/utils"
	"path/filepath"
)

var Lpk *LanguagePack

// defines a structure for language packs
type LanguagePack struct {
	lang  string
	packs map[string]map[string]string
}

// creates a new LanguagePack
func New() *LanguagePack {
	return &LanguagePack{
		lang:  "en",
		packs: make(map[string]map[string]string),
	}
}

func init() {
	Lpk = New()
	err := Lpk.LoadAllLangs(utils.STORAGE_PATH + "translates")
	if err != nil {
		fmt.Println("Error loading language packs:", err)
		return
	}

	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	err = Lpk.SelectLang(config.TranslateLanguagePack)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 加载所有语言包
func (lp *LanguagePack) LoadAllLangs(translatesDir string) error {
	err := filepath.Walk(translatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			lang := filepath.Base(path[:len(path)-len(filepath.Ext(path))])
			err := lp.LoadLang(lang, path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// 加载指定语言包
func (lp *LanguagePack) LoadLang(lang, path string) error {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return err
	}

	var messages map[string]string
	if err := json.Unmarshal(data, &messages); err != nil {
		return err
	}

	lp.packs[lang] = messages
	return nil
}

// 选择语言包
func (lp *LanguagePack) SelectLang(lang string) error {
	if _, exists := lp.packs[lang]; !exists {
		return errors.New("language not supported")
	}
	lp.lang = lang
	return nil
}

// 获取翻译
func (lp *LanguagePack) GetTranslate(key string) string {
	if messages, exists := lp.packs[lp.lang]; exists {
		if msg, ok := messages[key]; ok {
			return msg
		}
	}
	return lp.lang + "." + key
}

// 格式化语言包翻译
func (lp *LanguagePack) FormatMessage(key string, args ...interface{}) string {
	return fmt.Sprintf(lp.GetTranslate(key), args...)
}
