package wechat

import (
	"testing"
)

func TestMenu(t *testing.T) {
	x := NewLocalMongo("api")
	wc, err := x.GetWeChat()

	if err != nil {
		t.Error(err)
		return
	}
	m, err := wc.GetMenu()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(m)
	t.Log(wc.getAccessToken())
}
