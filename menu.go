package wechat

import (
	"encoding/json"
	"fmt"
)

// Use to store QR code
type QRScene struct {
	Ticket        string `json:"ticket"`
	ExpireSeconds int    `json:"expire_seconds"`
}

// Custom Menu
type Menu struct {
	Buttons []MenuButton `json:"button,omitempty"`
}

// Menu Button
type MenuButton struct {
	Name       string       `json:"name"`
	Type       string       `json:"type,omitempty"`
	Key        string       `json:"key,omitempty"`
	Url        string       `json:"url,omitempty"`
	SubButtons []MenuButton `json:"sub_button,omitempty"`
}

// Create QR scene
func (wc *WeChat) CreateQRScene(sceneId int, expires int) (*QRScene, error) {
	data := []byte(fmt.Sprintf(`{"expire_seconds":%d,"action_name":"QR_SCENE","action_info":{"scene":{"scene_id":%d}}}`, expires, sceneId))
	var qr QRScene
	err := wc.post(WeChatQRSceneCreate, data, &qr)
	return &qr, err
}

// Create  QR limit scene
func (wc *WeChat) CreateQRLimitScene(sceneId int) (*QRScene, error) {
	data := []byte(fmt.Sprintf(`{"action_name":"QR_LIMIT_SCENE","action_info":{"scene":{"scene_id":%d}}}`, sceneId))
	var qr QRScene
	err := wc.post(WeChatQRSceneCreate, data, &qr)
	return &qr, err
}

// Custom menu
func (wc *WeChat) CreateMenu(menu *Menu) error {
	if data, err := json.Marshal(menu); err != nil {
		return err
	} else {
		//fmt.Println(string(data))
		return wc.post(WeChatMenuCreate, data, nil)
	}
}

func (wc *WeChat) GetMenu() (*Menu, error) {
	var result struct {
		MenuCtx *Menu `json:"menu"`
	}
	result.MenuCtx = &Menu{}
	err := wc.get(WeChatMenuGet, &result, true)
	if err != nil {
		return nil, err
	}
	return result.MenuCtx, nil
}

// Delete Menu
func (wc *WeChat) DeleteMenu() error {
	return wc.get(WeChatMenuDelete, nil, true)
}
