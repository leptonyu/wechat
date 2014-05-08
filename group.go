package wechat

type Group struct {
	Id   int
	Name string
}

//Create a new Group
func (w *WeChat) CreateGroup(name string) (Group, error) {
	g := map[string]Group{}
	err := w.post(WeChatGroupCreate, []byte(`{"group":{"name":"`+name+`"}}`), &g)
	return g["group"], err
}
