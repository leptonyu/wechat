package wechat

import (
	"testing"
)

func TestUser(t *testing.T) {
	wc, err := NewWeChatInMem(`wxd`, `848`, `PRH`)
	if err != nil {
		t.Error(err)
		return
	}
	us, _, err := wc.GetAllUser("")
	if err != nil {
		t.Error(err)
	}
	for _, uu := range us {
		u, err := wc.GetUser(uu, LANG_CN)
		if err != nil {
			t.Error(err)
			return
		}

		t.Log(u, err)
	}
}
