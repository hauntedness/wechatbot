package wechat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/hauntedness/httputil"
	"github.com/hauntedness/wechatbot/config"
)

type Messager interface {
	Send(string, []Article, context.Context) error
}

type messager struct {
	config *config.BotConfig
	token  *Token
}

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

func NewMessager(conf *config.BotConfig) Messager {
	return &messager{config: conf, token: &Token{}}
}

var m Messager = &messager{config: config.GetBot(), token: &Token{}}

func GetMessager() Messager {
	return m
}

func (m *messager) GetToken() string {
	m.refreshToken()
	return m.token.AccessToken
}

func (m *messager) refreshToken() {

	if m.token.willExpireAt.After(time.Now()) {
		return
	}
	value := url.Values{}
	value.Add("corpid", m.config.CorpId)
	value.Add("corpsecret", m.config.Secret)
	u := url.URL{
		Scheme:     m.config.Protocol,
		Host:       m.config.Host,
		Path:       m.config.GetTokenUri,
		ForceQuery: true,
		RawQuery:   value.Encode(),
	}
	url_ := u.String()

	data := httputil.Request(http.MethodGet, url_, nil, nil)
	err := json.Unmarshal(data, m.token)
	if err != nil {
		log.Println(err)
		panic("getToken failed")
	}
	m.token.willExpireAt = time.Now().Add(time.Duration(m.token.ExpiresIn - 10))
}

func (m *messager) Send(messages string, articles []Article, ctx context.Context) error {
	select {
	case <-ctx.Done():
		e := "task is cancelled!"
		log.Println(e)
		return errors.New(e)
	default:
		token := m.GetToken()
		query := url.Values{}
		query.Add("access_token", token)
		u := url.URL{
			Scheme:     m.config.Protocol,
			User:       &url.Userinfo{},
			Host:       m.config.Host,
			Path:       m.config.SendMsgUri,
			ForceQuery: true,
			RawQuery:   query.Encode(),
		}
		senderUrl := u.String()
		message := Message{
			Touser:  m.config.UserId,
			Agentid: m.config.Agent,
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
