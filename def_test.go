package wechat

import (
	"testing"
	"time"
)

func TestDef(t *testing.T) {
	wc, err := New(&MemStorage{
		at: &AccessToken{
			Token:      "accesstoken",
			ExpireTime: time.Now().Add(1000 * time.Hour),
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(wc.getAccessToken())
	}
}
