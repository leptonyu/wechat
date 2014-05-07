package wechat

import (
	"testing"
)

func TestMenu(t *testing.T) {
	wc, err := NewWeChatInMem(`appid`,
		`secret`,
		`token`)

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
