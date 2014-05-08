package wechat

import (
	"testing"
)

func estGroup(t *testing.T) {
	x := NewLocalMongo("api")
	wc, err := x.GetWeChat()
	if err != nil {
		t.Error(err)
		return
	}
	m, err := wc.CreateGroup("你好")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(m)
}
