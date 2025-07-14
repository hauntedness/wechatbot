package wechat

import (
	"context"

	"github.com/hauntedness/wechatbot/config"
)

func Send(messages string, articles []Article, ctx context.Context) error {
	if m.config == nil {
		config, err := config.GetBot()
		if err != nil {
			return err
		}
		m.config = config
	}
	return m.Send(messages, articles, ctx)
}
