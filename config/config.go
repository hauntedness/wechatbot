package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type config struct {
	BotConf BotConfig
}

func GetWechatConfig() *config {
	path := os.Getenv("WECHAT_CONFIG_PATH")
	if path == "" {
		dir, _ := os.UserConfigDir()
		path = dir + `\wechat\.config\wechat.toml`
	}
	conf := config{}
	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		return nil
	}
	return &conf
}

type BotConfig struct {
	CorpId         string
	Agent          string
	Secret         string
	UserAgent      string
	Protocol       string
	Host           string
	Port           int
	GetTokenUri    string
	SendMsgUri     string
	UserId         string
	Token          string
	EncodingAESKey string
}

func GetBot() BotConfig {
	c := GetWechatConfig()
	return c.BotConf
}
