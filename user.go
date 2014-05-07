package wechat

import (
	"fmt"
)

type User struct {
	Subscribe     int
	Openid        string `json:",omitempty"`
	Nickname      string `json:",omitempty"`
	Sex           int    `json:",omitempty"`
	City          string `json:",omitempty"`
	Country       string `json:",omitempty"`
	Province      string `json:",omitempty"`
	Language      string `json:",omitempty"`
	Headimgurl    string `json:",omitempty"`
	SubscribeTime int64  `json:"subscribe_time,omitempty"`
}

//Get user infomation from wechat
func (w *WeChat) GetUser(openid, lang string) (*User, error) {
	u := &User{}
	if lang == "" {
		lang = LANG_CN
	}
	err := w.get(fmt.Sprintf(WeChatUserGet, openid, lang)+`%v`, u, true)
	return u, err
}

//Get all user from wechat
func (w *WeChat) GetAllUser(firstid string) ([]string, string, error) {
	var a struct {
		Total int
		Count int
		Data  map[string][]string
		Next  string `json:"next_openid"`
	}
	if err := w.get(fmt.Sprintf(WeChatUserGetAll, firstid)+`%v`, &a, true); err != nil {
		return nil, "", err
	}
	return a.Data["openid"], a.Next, nil

}
