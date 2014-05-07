/*
This package implements the WeChat SDK.
*/
package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

//WeChat Language support
const (
	LANG_CN = `zh_CN` // Simplified Chinese
	LANG_TW = `zh_TW` // Traditional Chinese
	LANG_EN = `en`    // English
)

//WeChat URL info
const (
	// WeChat host URL
	WeChatHost = `https://api.weixin.qq.com/cgi-bin/`
	//WeChat Reply
	WeChatPost   = WeChatHost + `message/custom/send?access_token=%v`
	WeChatUpload = WeChatHost + `media/uploadnews?access_token=%v`
	//WeChat User
	WeChatUser       = WeChatHost + `user`
	WeChatUserGet    = WeChatUser + `/info?openid=%v&lang=%v&access_token=`
	WeChatUserGetAll = WeChatUser + `/get?next_openid=%v&access_token=`
	//WeChat Group
	WeChatGroup             = WeChatHost + `groups`
	WeChatGroupCreate       = WeChatGroup + `/create?access_token=%v`
	WeChatGroupGet          = WeChatGroup + `/get?access_token=%v`
	WeChatGroupUpdate       = WeChatGroup + `/update?access_token=%v`
	WeChatGroupMemberUpdate = WeChatGroup + `/members/update?access_token=%v`
	WeChatGroupGetIdByUser  = WeChatGroup + `getid?access_token=%v`
	//WeChat Menu
	WeChatMenu       = WeChatHost + `menu`
	WeChatMenuCreate = WeChatMenu + `/create?access_token=%v`
	WeChatMenuGet    = WeChatMenu + `/get?access_token=%v`
	WeChatMenuDelete = WeChatMenu + `/delete?access_token=%v`
	//WeChat Token
	WeChatToken = WeChatHost + `token?grant_type=client_credential&appid=%v&secret=%v`
	//WeChat QRScene
	WeChatQRScene       = WeChatHost + `qrcode`
	WeChatQRSceneCreate = WeChatQRScene + `/create?access_token=%v`
	WeChatShowQRScene   = "https://mp.weixin.qq.com/cgi-bin/showqrcode"
	//File
	WeChatFileURL = "http://file.api.weixin.qq.com/cgi-bin/media"
)

// Basic struct of wechat.
type WeChat struct {
	appid  string   // Appid of wechat
	secret string   // App secret of wechat
	token  string   // App token of wechat, this is defined by user.
	atrw   Storage  // Storage interface, this interface used to store the limit resource.
	routes []*Route // Route of request handler
}

//Register Route
func (w *WeChat) Register(pattern, keyword string, handler HandleFunc) {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}
	if !reg.MatchString(keyword) {
		panic(errors.New("Pattern " + pattern + " must match the keyword " + keyword + "!"))
	}
	w.routes = append(w.routes, &Route{
		Regex:   reg,
		Keyword: keyword,
		Handle:  handler,
	})
}

//Create wechat struct.
func New(storage Storage) (*WeChat, error) {
	appid, secret, token, err := storage.WeChatInfo()
	if err != nil {
		return nil, err
	}
	return &WeChat{
		appid:  appid,
		secret: secret,
		token:  token,
		atrw:   storage,
	}, nil
}

//Handle Func
type HandleFunc func(RespondWriter, *Request) error

//Route of request handler
type Route struct {
	Regex   *regexp.Regexp //Regexp of words that use this Handle
	Keyword string         //Basic Keyword
	Handle  HandleFunc     // Handle function
}

// Access Token, we need this to verify the identity with WeChat server.
// It is valid for 7200 seconds.
type AccessToken struct {
	Token      string    `json:"access_token"` // Access Token
	ExpireTime time.Time `json:"expires_in"`   // ExpireTime of Access Token
}

//Store some important data get from wechat server
type Storage interface {
	ReadAccessToken() (AccessToken, error)                // Read access token from storage
	WriteAccessToken(AccessToken) error                   //Write access token to storage
	SaveRequest(*Request)                                 // Save WeChat request
	WeChatInfo() (appid, secret, token string, err error) //Fetch Basic WeChat Info
}

//Get Access Token
func (w *WeChat) getAccessToken() (AccessToken, error) {
	at, err := w.atrw.ReadAccessToken()
	if err == nil && time.Since(at.ExpireTime).Seconds() < 0 && at.Token != "" {
		return at, nil
	}
	res := AccessToken{}
	var xxx struct {
		Token  string `json:"access_token"` // Access Token
		Expire int64  `json:"expires_in"`   // ExpireTime of Access Token

	}
	err = w.get(fmt.Sprintf(WeChatToken, w.appid, w.secret), &xxx, false)
	if err == nil {
		res.Token = xxx.Token
		res.ExpireTime = time.Now().Add(time.Duration(xxx.Expire) * time.Second)
	}
	return res, err
}

//WeChat server respond code
type ErrWeChat struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (e *ErrWeChat) Error() string {
	return strconv.Itoa(e.ErrCode) + ":" + e.ErrMsg
}

//Get information from WeChat server.
func (w *WeChat) get(url string, out interface{}, needAccessToken bool) error {
	ewc := &ErrWeChat{}
	for i := 1; i <= 3; i++ {
		ewc.ErrCode = -9999
		urlx := url
		if needAccessToken {
			at, err := w.getAccessToken()
			if err != nil {
				return err
			}
			urlx = fmt.Sprintf(url, at.Token)
		}
		resp, err := http.Get(fmt.Sprintf(urlx))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		//fmt.Println(url, string(body))
		if er := json.Unmarshal(body, ewc); er != nil {
			return er
		} else {
			switch ewc.ErrCode {
			case -9999:
				return json.Unmarshal(body, out)
			case 0:
				return nil
			case 42001:
				continue
			default:
				return ewc
			}
		}
	}
	return ewc
}

//Post json to WeChat server.
func (w *WeChat) post(url string, data []byte, out interface{}) error {
	ewc := &ErrWeChat{}
	for i := 1; i <= 3; i++ {
		ewc.ErrCode = -9999
		at, err := w.getAccessToken()
		if err != nil {
			return err
		}
		resp, err := http.Post(fmt.Sprintf(url, at.Token), "application/json; charset=utf-8", bytes.NewReader(data))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if er := json.Unmarshal(body, ewc); er != nil {
			return err
		} else {
			switch ewc.ErrCode {
			case -9999:
				return json.Unmarshal(body, out)
			case 0:
				return nil
			case 42001:
				continue
			default:
				return ewc
			}
		}
	}
	return ewc
}
