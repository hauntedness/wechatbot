package wechat

import (
	"context"
)

func SendMessage(messages string, ctx context.Context) {
	m.SendMessage(messages, ctx)
}

func SendNewsMessage(articles []Article, ctx context.Context) {
	m.SendNewsMessage(articles, ctx)
}
