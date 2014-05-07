package db

import (
	"github.com/leptonyu/wechat"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type MongoStorage struct {
	username string
	password string
	host     string
	database string
}

//Query in database
type QueryFunc func(*mgo.Database) error

//Using MongoDB to create WeChat struct
func New(username, password, host, database string) *MongoStorage {
	return &MongoStorage{
		username: username,
		password: password,
		host:     host,
		database: database,
	}
}

//Standard query of mongodb
func (m *MongoStorage) Query(qf QueryFunc) error {
	url := "mongodb://"
	if m.username != "" {
		url += m.username
		if m.password != "" {
			url += ":" + m.password
		}
		url += "@"
	}
	if m.host == "" {
		m.host = "localhost"
	}
	url += m.host
	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}
	defer session.Close()
	return qf(session.DB("wechat_" + m.database))
}
func (m *MongoStorage) GetWeChat() (*wechat.WeChat, error) {
	return wechat.New(m)
}

type storeWeChat struct {
	Name   string
	Appid  string
	Secret string
	Token  string
}

func (m *MongoStorage) Init(appid, secret, token string) error {
	return m.Query(func(d *mgo.Database) error {
		_, err := d.C("wechat").Upsert(bson.M{"name": "wechat"},
			storeWeChat{
				Name:   "wechat",
				Appid:  appid,
				Secret: secret,
				Token:  token,
			})
		return err
	})
}

type access struct {
	Name   string
	Token  string
	Expire time.Time
}

func (m *MongoStorage) ReadAccessToken() (wechat.AccessToken, error) {
	at := wechat.AccessToken{}
	err := m.Query(func(d *mgo.Database) error {
		a := access{}
		if err := d.C("wechat").Find(bson.M{"name": "accesstoken"}).One(&a); err != nil {
			return err
		}
		at.Token = a.Token
		at.ExpireTime = a.Expire
		return nil
	})
	return at, err
}

func (m *MongoStorage) WriteAccessToken(at wechat.AccessToken) error {
	return m.Query(func(d *mgo.Database) error {
		_, err := d.C("wechat").Upsert(bson.M{"name": "accesstoken"},
			access{
				Name:   "accesstoken",
				Token:  at.Token,
				Expire: at.ExpireTime,
			})
		return err
	})
}

func (m *MongoStorage) SaveRequest(r *wechat.Request) {
	m.Query(func(d *mgo.Database) error {
		d.C("request").Insert(r)
		return nil
	})
}

func (m *MongoStorage) WeChatInfo() (appid, secret, token string, err error) {
	err = m.Query(func(d *mgo.Database) error {
		x := storeWeChat{}
		if err := d.C("wechat").Find(bson.M{"name": "wechat"}).One(&x); err != nil {
			return err
		}

		appid = x.Appid
		secret = x.Secret
		token = x.Token
		return nil
	})
	return appid, secret, token, err
}
