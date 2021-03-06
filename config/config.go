package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type config struct {
	BotConf BotConfig
}

func GetWechatConfig() *config {
	var err error
	path_ := os.Getenv("WECHAT_CONFIG_PATH")
	if path_ == "" {
		dir, err := os.UserConfigDir()
		if err != nil {
			panic(err)
		}
		path_ = filepath.Join(dir, "wechat", ".config", "wechat.toml")
	}
	conf := config{}
	_, err = toml.DecodeFile(path_, &conf)
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

func GetBot() *BotConfig {
	c := GetWechatConfig()
	return &c.BotConf
}
