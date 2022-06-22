package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type config struct {
	BotConf BotConfig
}

func GetWechatConfig() *config {
	var err error
	path := os.Getenv("WECHAT_CONFIG_PATH")
	if path == "" {
		dir, err := os.UserConfigDir()
		if err != nil {
			panic(err)
		}
		path = dir + `\wechat\.config\wechat.toml`
	}
	conf := config{}
	_, err = toml.DecodeFile(path, &conf)
	if err != nil {
		panic(err)
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
