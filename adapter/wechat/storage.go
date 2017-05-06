package wechat

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
	"github.com/num5/ider"
	"github.com/num5/webot"
)

type message struct {
	ID          int64              `json:"id"`
	Msg         webot.EventMsgData `json:"msg"`
	ReceiveTime time.Time    `json:"receive_time"`
}

func Stor(storpath string, data webot.EventMsgData) error {
	msg := new(message)
	id := ider.NewID(1)
	msg.ID = id
	msg.Msg = data
	msg.ReceiveTime = time.Now()

	return write(storpath, msg)
}

func write(storpath string, value interface{}) error {
	content, err := json.Marshal(value)

	if err != nil {
		return err
	}
	return ioutil.WriteFile(storpath, content, os.ModePerm)
}
