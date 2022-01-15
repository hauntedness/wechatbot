package main

import (
	"context"

	"github.com/hauntedness/wechatbot/wechat"
)

func main() {
	wechat.SendMessage("hellow", context.TODO())
}
