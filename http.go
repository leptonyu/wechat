package wechat

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"
)

//Check valid from wechat.
func checkSignature(token string, w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()
	var signature string = r.FormValue("signature")
	var timestamp string = r.FormValue("timestamp")
	var nonce string = r.FormValue("nonce")
	strs := sort.StringSlice{token, timestamp, nonce}
	sort.Strings(strs)
	var str string
	for _, s := range strs {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil)) == signature
}

//Handle http request
//implement the http.Handler interface
func (wc *WeChat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !checkSignature(wc.token, w, r) {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	//Valify WeChat request test
	if r.Method == "GET" {
		fmt.Fprintf(w, r.FormValue("echostr"))
		return
	}
	//Read message
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	//Process Message
	msg := &Request{}
	if err := xml.Unmarshal(data, &msg); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	// Storage every valid request
	go wc.atrw.SaveRequest(msg)
	requestPath := msg.MsgType
	if requestPath == msgEvent {
		requestPath += "." + msg.Event
	}
	for _, route := range wc.routes {
		if !route.Regex.MatchString(requestPath) {
			continue
		}
		route.Handle(&Respond{
			wechat:       wc,
			Writer:       w,
			ToUserName:   msg.FromUserName,
			FromUserName: msg.ToUserName,
		}, msg)
		return

	}
}

//Respond to wechat server
type Respond struct {
	wechat       *WeChat
	Writer       http.ResponseWriter
	ToUserName   string
	FromUserName string
}

func (r *Respond) ReplyText(text string) {
	r.reply(
		`<MsgType><![CDATA[text]]></MsgType><Content><![CDATA[` + text + `]]></Content>`)
}

func (r *Respond) reply(message string) {
	head := fmt.Sprintf(`<xml><ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[%s]]></FromUserName>
<CreateTime>%d</CreateTime>%v</xml>`, r.ToUserName, r.FromUserName, time.Now().Unix(), message)
	go r.wechat.atrw.SaveReply(head)
	r.Writer.Write([]byte(head))
}

func (r *Respond) ReplyImage(mediaId string) {
	r.reply(
		`<![CDATA[image]]></MsgType><Image><MediaId><![CDATA[` + mediaId + `]]></MediaId></Image>`)
}
func (r *Respond) ReplyVoice(mediaId string) {
	r.reply(
		`<MsgType><![CDATA[voice]]></MsgType><Voice><MediaId><![CDATA[` + mediaId + `]]></MediaId></Voice>`)
}

func (r *Respond) ReplyVideo(mediaId, title, desp string) {
	r.reply(
		fmt.Sprintf(`<MsgType><![CDATA[video]]></MsgType><Video><MediaId><![CDATA[%s]]></MediaId><Title><![CDATA[%s]]></Title><Description><![CDATA[%s]]></Description></Video>`,
			mediaId, title, desp))
}

type Music struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	MusicUrl     string `json:"musicurl"`
	HQMusicUrl   string `json:"hqmusicurl"`
	ThumbMediaId string `json:"thumb_media_id"`
}

func (r *Respond) ReplyMusic(music *Music) {
	r.reply(
		fmt.Sprintf(`<MsgType><![CDATA[music]]></MsgType><Music><Title><![CDATA[%s]]></Title><Description><![CDATA[%s]]></Description><MusicUrl><![CDATA[%s]]></MusicUrl><HQMusicUrl><![CDATA[%s]]></HQMusicUrl><ThumbMediaId><![CDATA[%s]]></ThumbMediaId></Music>`,
			music.Title, music.Description, music.MusicUrl, music.HQMusicUrl, music.ThumbMediaId))
}

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	PicUrl      string `json:"picurl"`
	Url         string `json:"url"`
}

func (r *Respond) ReplyNews(articles []Article) {
	ctx := ""
	for _, article := range articles {
		ctx += fmt.Sprintf(`<item><Title><![CDATA[%s]]></Title> <Description><![CDATA[%s]]></Description><PicUrl><![CDATA[%s]]></PicUrl><Url><![CDATA[%s]]></Url></item>`,
			article.Title, article.Description, article.PicUrl, article.Url)
	}
	r.reply(
		fmt.Sprintf(`<MsgType><![CDATA[news]]></MsgType><ArticleCount>%d</ArticleCount><Articles>%s</Articles>`,
			len(articles), ctx))
}
func (r *Respond) FromUserId() string {
	return r.ToUserName
}

//Reply messages to wechat
type RespondWriter interface {
	FromUserId() string                            //Get Request
	ReplyText(text string)                         //Reply text message to wechat
	ReplyImage(mediaId string)                     //Reply text message to wechat
	ReplyVoice(mediaId string)                     //Reply text message to wechat
	ReplyVideo(mediaId, title, description string) //Reply text message to wechat
	ReplyMusic(music *Music)                       //Reply text message to wechat
	ReplyNews(articles []Article)                  //Reply text message to wechat
}

type Request struct {
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      string
	MsgId        int64
	Content      string  `json:",omitempty"`
	PicUrl       string  `json:",omitempty"`
	MediaId      string  `json:",omitempty"`
	Format       string  `json:",omitempty"`
	ThumbMediaId string  `json:",omitempty"`
	LocationX    float32 `json:"Location_X,omitempty",xml:"Location_X"`
	LocationY    float32 `json:"Location_Y,omitempty",xml:"Location_Y"`
	Scale        float32 `json:",omitempty"`
	Label        string  `json:",omitempty"`
	Title        string  `json:",omitempty"`
	Description  string  `json:",omitempty"`
	Url          string  `json:",omitempty"`
	Event        string  `json:",omitempty"`
	EventKey     string  `json:",omitempty"`
	Ticket       string  `json:",omitempty"`
	Latitude     float32 `json:",omitempty"`
	Longitude    float32 `json:",omitempty"`
	Precision    float32 `json:",omitempty"`
	Recognition  string  `json:",omitempty"`
	UserName     string  `json:"-"`
}

const (
	msgEvent = "event"
	// Event Type
	EventSubscribe   = "subscribe"
	EventUnsubscribe = "unsubscribe"
	EventScan        = "scan"
	EventClick       = "CLICK"
	EventLocation    = "LOCATION"
	EventView        = "VIEW"
	// Message type
	MsgTypeDefault          = ".*"
	MsgTypeText             = "text"
	MsgTypeImage            = "image"
	MsgTypeVoice            = "voice"
	MsgTypeVideo            = "video"
	MsgTypeLocation         = "location"
	MsgTypeLink             = "link"
	MsgTypeEvent            = msgEvent + ".*"
	MsgTypeEventSubscribe   = msgEvent + "\\." + EventSubscribe
	MsgTypeEventUnsubscribe = msgEvent + "\\." + EventUnsubscribe
	MsgTypeEventScan        = msgEvent + "\\." + EventScan
	MsgTypeEventClick       = msgEvent + "\\." + EventClick
	MsgTypeEventView        = msgEvent + "\\." + EventView
	MsgTypeEventLocation    = msgEvent + "\\." + EventLocation
	// Media type
	MediaTypeImage = "image"
	MediaTypeVoice = "voice"
	MediaTypeVideo = "video"
	MediaTypeThumb = "thumb"
	// Button type
	MenuButtonTypeKey = "click"
	MenuButtonTypeUrl = "view"
)
