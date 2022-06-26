package wechat

import (
	"errors"
	"regexp"
	"sync"
	"time"
)

type Token struct {
	lock         sync.Mutex
	willExpireAt time.Time
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func IsAccessTokenError(res []byte) error {
	match, _ := regexp.Match("\"errcode\":40014", res)
	if match {
		return errors.New(string(res))
	}
	return nil
}
