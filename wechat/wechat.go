package wechat

import (
	"context"
)

func Send(messages string, articles []Article, ctx context.Context) error {
	return m.Send(messages, articles, ctx)
}
