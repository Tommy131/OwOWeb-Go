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
 * @LastEditTime : 2024-07-01 01:51:47
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package utils

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	WebListeningAddress   string `json:"web-listening-address"`
	TranslateLanguagePack string `json:"translate-language-pack"`
}

var (
	config     *Config
	configPath = STORAGE_PATH + "config.json"
	mutex      sync.Mutex
)

// 读取配置文件
func LoadConfig() (*Config, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if config != nil {
		return config, nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			config = &Config{
				WebListeningAddress:   "127.0.0.1:8080",
				TranslateLanguagePack: "en",
			}
			err = SaveConfig(config)
			if err != nil {
				return nil, err
			}
			return config, nil
		}
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config = &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// 保存配置文件
func SaveConfig(cfg *Config) error {
	file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	config = cfg
	return nil
}

// 更新配置文件
func UpdateConfig(updateFunc func(*Config)) (*Config, error) {
	mutex.Lock()
	defer mutex.Unlock()

	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	updateFunc(cfg)

	err = SaveConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
