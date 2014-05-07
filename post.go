package wechat

import (
	"fmt"
)

//Interface to post message to WeChat server
// Only works for Service Account
type PostMessage interface {
	PostText(touser, content string) error                       //Post Text
	PostImage(touser, media_id string) error                     //Post Image
	PostVoice(touser, media_id string) error                     //Post Voice
	PostVideo(touser, media_id, title, description string) error //Post Video
	PostMusic(touser, media_id string) error                     //Post Music
	PostNews(touser string, articles []Article) error            //Post Article
}

func (w *WeChat) PostText(touser, content string) error {
	return w.post(WeChatPost,
		[]byte(fmt.Sprintf(`{"touser":"%v","msgtype":"text","text":{"content":"%v"}}`,
			touser,
			content)),
		nil)
}

func (w *WeChat) PostImage(touser, media_id string) error {
	return w.post(WeChatPost,
		[]byte(fmt.Sprintf(`{"touser":"%v","msgtype":"image","image":{"media_id":"%v"}}`,
			touser,
			media_id)),
		nil)
}

func (w *WeChat) PostVoice(touser, media_id string) error {
	return w.post(WeChatPost,
		[]byte(fmt.Sprintf(`{"touser":"%v","msgtype":"voice","voice":{"media_id":"%v"}}`,
			touser,
			media_id)),
		nil)
}

func (w *WeChat) PostVideo(touser, media_id, title, description string) error {
	return w.post(WeChatPost,
		[]byte(fmt.Sprintf(`{"touser":"%v","msgtype":"video","video":{"media_id":"%v","title":"%v","description":"%v"}}`,
			touser,
			media_id,
			title,
			description)),
		nil)
}

func (w *WeChat) PostMusic(touser string, music Music) error {
	return w.post(WeChatPost,
		[]byte(fmt.Sprintf(`{"touser":"%v","msgtype":"music","music":{"title":"%v","description":"%v","musicurl":"%v","hqmusicurl":"%v","thumb_media_id":"%v"}}`,
			touser,
			music.Title,
			music.Description,
			music.MusicUrl,
			music.HQMusicUrl,
			music.ThumbMediaId)),
		nil)
}

func (w *WeChat) PostNews(touser string, articles []Article) error {
	sas := ""
	for _, a := range articles {
		str := fmt.Sprintf(`{"title":"%v","description":"%v","url":"%v","picurl":"%v"}`, a.Title, a.Description, a.Url, a.PicUrl)
		if sas == "" {
			sas = str
		} else {
			sas += "," + str
		}
	}
	return w.post(WeChatPost,
		[]byte(fmt.Sprintf(`{"touser":"%v","msgtype":"news","news":{"articles":[%v]}}`,
			touser, sas)), nil)

}
