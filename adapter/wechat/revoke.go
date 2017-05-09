package wechat

import (
	"encoding/xml"
	"github.com/num5/logger"
	"github.com/num5/webot"
	"regexp"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)

var char_string = map[string]string{
	"&lt;": "<",
	"&gt;": ">",
}

type sysmsg struct {
	Revokemsg revokemsg `xml:"revokemsg"`
}

type revokemsg struct {
	Session    string `xml:"session"`
	Oldmsgid   string `xml:"oldmsgid"`
	Msgid      string `xml:"msgid"`
	Replacemsg string `xml:"replacemsg"`
}

// 撤销取回
func (w *WeChatAdapter) revoke(data webot.EventMsgData) {

	if data.MsgType == MSG_WITHDRAW {

		var reg *regexp.Regexp
		revokemsg := data.Content
		for k, v := range char_string {
			reg = regexp.MustCompile(`(` + k + `)`)
			revokemsg = reg.ReplaceAllString(revokemsg, v)
		}

		var revoke sysmsg

		err := xml.Unmarshal([]byte(revokemsg), &revoke)
		if err != nil {
			logger.Errorf("xml parse err:%v", err)
		}

		bytemsg, err := redis.Bytes(w.conn.Do("GET", revoke.Revokemsg.Msgid))

		if err != nil {
			logger.Errorf("redis GET err:%v", err)
		}

		var msg webot.EventMsgData
		err = json.Unmarshal(bytemsg, &msg)

		if err != nil {
			logger.Errorf("unmarshal bytemsg err:%v", err)
		}

		w.Wechat.SendTextMsg(msg.Content, data.FromUserName)
	}
}
