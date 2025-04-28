package data

import (
	"encoding/json"
	"os"
)

var Config AppConfig = ReadConfig()

type AppConfig struct {
	Token            string
	Owner            string
	BackupLocation   string
	GuildID          string
	RoleID           string
	GameDataLocation string
	Emote            string
	LogChannelID     string
}

func ReadConfig() AppConfig {
	rawConfig, err := os.ReadFile("config.json")

	if err != nil {
		panic(err)
	}

	var config AppConfig

	err = json.Unmarshal(rawConfig, &config)

	if err != nil {
		panic(err)
	}

	return config
}
