package wechat

import (
	"context"
	"errors"
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
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty,omitempty"`
	ToTag   string `json:"totag,omitempty"`
	MsgType string `json:"msgtype"`
	AgentId string `json:"agentid"`
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

var m = &messager{token: &Token{}}

func GetMessager() Messager {
	return m
}

func (m *messager) GetToken() (string, error) {
	err := m.refreshToken()
	if err != nil {
		return "", err
	}
	return m.token.AccessToken, nil
}

func (m *messager) refreshToken() error {
	if m.token.willExpireAt.After(time.Now()) {
		return nil
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
	token, err := httputil.GetJson[Token](url_, nil, nil)
	if err != nil {
		return err
	}
	if token.Errcode != 0 {
		return &MsgError[Token]{v: token}
	}
	m.token = token
	m.token.willExpireAt = time.Now().Add(time.Second * time.Duration(m.token.ExpiresIn-10))
	return nil
}

func (m *messager) Send(messages string, articles []Article, ctx context.Context) error {
	select {
	case <-ctx.Done():
		return errors.New("task is cancelled!")
	default:
		token, err := m.GetToken()
		if err != nil {
			return err
		}
		query := url.Values{}
		query.Add("access_token", token)
		u := url.URL{
			Scheme:     m.config.Protocol,
			Host:       m.config.Host,
			Path:       m.config.SendMsgUri,
			ForceQuery: true,
			RawQuery:   query.Encode(),
		}
		senderUrl := u.String()
		message := Message{
			ToUser:  m.config.UserId,
			AgentId: m.config.Agent,
		}
		if len(articles) != 0 {
			message.MsgType = "news"
			message.News = News{
				Articles: articles,
			}
		} else if messages != "" {
			message.Text.Content = messages
			message.MsgType = "text"
		} else {
			return errors.New("can not send empty message")
		}
		var h httputil.H
		if m.config.UserAgent != "" {
			h = make(httputil.H)
			h["User-Agent"] = m.config.UserAgent
		}
		res, err := httputil.PostJson[MessageSendResponse](senderUrl, message, h)
		if err != nil {
			return err
		}
		if res.Errcode != 0 {
			return &MsgError[MessageSendResponse]{v: res}
		}
		return nil
	}
}

type MessageSendResponse struct {
	Errcode        int64  `json:"errcode"`
	Errmsg         string `json:"errmsg"`
	Invaliduser    string `json:"invaliduser"`
	Invalidparty   string `json:"invalidparty"`
	Invalidtag     string `json:"invalidtag"`
	Unlicenseduser string `json:"unlicenseduser"`
	Msgid          string `json:"msgid"`
	ResponseCode   string `json:"response_code"`
}
