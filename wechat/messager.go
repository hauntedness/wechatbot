package wechat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/hauntedness/httputil"
)

type Messager interface {
	Send(string, []Article, context.Context) error
}

type messager struct{}

type Message struct {
	Touser  string `json:"touser"`
	Toparty string `json:"toparty,omitempty"`
	Totag   string `json:"totag,omitempty"`
	Msgtype string `json:"msgtype"`
	Agentid string `json:"agentid"`
	News    News   `json:"news,omitempty"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text,omitempty"`
}

type News struct {
	Articles []Article `json:"articles"`
}

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Picurl      string `json:"picurl"`
	Appid       string `json:"appid"`
	Pagepath    string `json:"pagepath"`
}

var m Messager = &messager{}

func GetMessager() Messager {
	return m
}

func (m *messager) Send(messages string, articles []Article, ctx context.Context) error {
	select {
	case <-ctx.Done():
		e := "task is cancelled!"
		log.Println(e)
		return errors.New(e)
	default:
		token := GetToken()
		query := url.Values{}
		query.Add("access_token", token)
		u := url.URL{
			Scheme:     Bot.Protocol,
			User:       &url.Userinfo{},
			Host:       Bot.Host,
			Path:       Bot.SendMsgUri,
			ForceQuery: true,
			RawQuery:   query.Encode(),
		}
		senderUrl := u.String()
		message := Message{
			Touser:  Bot.UserId,
			Agentid: Bot.Agent,
		}
		if len(articles) != 0 {
			message.Msgtype = "news"
			message.News = News{
				Articles: articles,
			}
		} else if messages != "" {
			message.Text.Content = messages
			message.Msgtype = "text"
		} else {
			panic("can not send empty message")
		}
		json, err := json.Marshal(message)
		if err != nil {
			log.Println("parse message failed")
			log.Println(err)
			return err
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
