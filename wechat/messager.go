package wechat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/hauntedness/httputil"
)

type Messager interface {
	SendMessage(string, context.Context) error
}

type messager struct{}

type Message struct {
	Touser  string `json:"touser"`
	Msgtype string `json:"msgtype"`
	Agentid string `json:"agentid"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

var m Messager = &messager{}

func GetMessager() Messager {
	return m
}

func (m *messager) SendMessage(messages string, ctx context.Context) error {

	select {
	case <-ctx.Done():
		e := "task is cancelled!"
		log.Println(e)
		return errors.New(e)
	default:
		token := GetToken()
		senderUrl := Bot.Protocol + Bot.Host + Bot.SendMsgUri + "?access_token=" + token
		message := Message{Touser: Bot.UserId, Msgtype: "text", Agentid: Bot.Agent}
		message.Text.Content = messages
		json, err := json.Marshal(message)
		if err != nil {
			log.Println("parse message failed")
			log.Println(err)
		}

		res := httputil.Request(http.MethodPost, senderUrl, bytes.NewBuffer([]byte(json)), nil)
		err = IsAccessTokenError(res)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}
}
