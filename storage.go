package wechat

import (
	"errors"
	"log"
)

//Create WeChat using in memory storage.
func NewWeChatInMem(appid, secret, token string) (*WeChat, error) {
	return New(&MemStorage{
		appid:  appid,
		secret: secret,
		token:  token,
		at:     &AccessToken{},
	})
}

//In memory storage struct
// This storage will not save the request, just print them into log.
type MemStorage struct {
	appid  string
	secret string
	token  string
	at     *AccessToken
}

func (s *MemStorage) ReadAccessToken() (AccessToken, error) {
	if s.at == nil {
		return *s.at, errors.New("No access token was found!")
	} else {
		return *s.at, nil
	}
}
func (s *MemStorage) WriteAccessToken(at AccessToken) error {
	s.at = &at
	return nil
}
func (s *MemStorage) SaveRequest(r *Request) {
	log.Println(r)
}
func (s *MemStorage) WeChatInfo() (appid, secret, token string, err error) {
	return s.appid, s.secret, s.token, nil
}
