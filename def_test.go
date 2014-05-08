package wechat

import (
	"testing"
)

func TestDef(t *testing.T) {
	x := NewLocalMongo("api")
	wc, err := x.GetWeChat()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(wc.getAccessToken())
	}
}
