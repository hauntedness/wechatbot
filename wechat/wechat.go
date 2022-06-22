package wechat

import (
	"context"
)

func Send(messages string, articles []Article, ctx context.Context) {
	m.Send(messages, articles, ctx)
}
