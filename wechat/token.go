package wechat

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/hauntedness/httputil"
)

type Token struct {
	lock         sync.Mutex
	willExpireAt time.Time
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
}

var cached Token = Token{}

func GetToken() string {
	// request new token after it is expired
	if time.Now().After(cached.willExpireAt) {
		cached.lock.Lock()
		updateToken(&cached)
		cached.lock.Unlock()
	}
	return cached.AccessToken
}

func IsAccessTokenError(res []byte) error {
	match, _ := regexp.Match("\"errcode\":40014", res)
	if match {
		return errors.New(string(res))
	}
	return nil
}

func updateToken(token *Token) {
	url := Bot.Protocol + Bot.Host + Bot.GetTokenUri
	data := httputil.Request(http.MethodGet, url, nil, nil)
	err := json.Unmarshal(data, token)
	if err != nil {
		log.Println(err)
		panic("getToken failed")
	}
	cached.willExpireAt = time.Now().Add(time.Duration(cached.ExpiresIn - 10))
}
