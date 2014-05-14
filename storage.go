package wechat

import (
	"errors"
	"log"
)

//Store some important data get from wechat server
type Storage interface {
	ReadAccessToken() (AccessToken, error)                // Read access token from storage
	WriteAccessToken(AccessToken) error                   //Write access token to storage
	SaveRequest(*Request)                                 // Save WeChat request
	SaveReply(string)                                     // Save WeChat reply
	WeChatInfo() (appid, secret, token string, err error) //Fetch Basic WeChat Info
	GetUserName(string) (name, admin string, err error)   //Fetch Username of id,if username is not exists, return id
	SetUserName(id, name, admin string)                   //Set Username of id,if username is not exists, return id
}

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
	idname map[string]*user
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
func (s *MemStorage) SaveReply(r string) {
	log.Println(r)
}
func (s *MemStorage) SaveRequest(r *Request) {
	log.Println(r)
}
func (s *MemStorage) WeChatInfo() (appid, secret, token string, err error) {
	return s.appid, s.secret, s.token, nil
}

type user struct {
	Id    string
	Name  string
	Admin string
}

func (s *MemStorage) GetUserName(id string) (name, admin string, err error) {
	u, ok := s.idname[id]
	if ok {
		return u.Name, u.Admin, nil
	} else {
		return "", "", errors.New("Id not have name")
	}
}
func (s *MemStorage) SetUserName(id, name, admin string) {
	s.idname[id] = &user{
		Id:    id,
		Name:  name,
		Admin: admin,
	}
}
