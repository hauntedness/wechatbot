package main

import (
	"context"

	"github.com/hauntedness/wechatbot/wechat"
)

//go:generate go run ../cmd/main.go
func main() {
	var articles = []wechat.Article{
		{
			Title:       "中秋节礼品领取",
			Description: "今年中秋节公司有豪礼相送",
			URL:         "URL",
			Picurl:      "http://res.mail.qq.com/node/ww/wwopenmng/images/independent/doc/test_pic_msg1.png",
		},
	}
	wechat.Send("dafsf", articles, context.TODO())

}
