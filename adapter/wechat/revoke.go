package wechat

import (
	"github.com/num5/webot"
	"regexp"
	"encoding/xml"
	"github.com/num5/logger"
	"github.com/go-redis/redis"
)

var char_string = map[string]string{
	"&lt;": "<",
	"&gt;": ">",
}

type sysmsg struct {
	Revokemsg revokemsg `xml:"revokemsg"`
}

type revokemsg struct {
	Session string `xml:"session"`
	Oldmsgid string `xml:"oldmsgid"`
	Msgid string `xml:"msgid"`
	Replacemsg string `xml:"replacemsg"`
}

// 撤销取回
func (w *WeChatAdapter) revoke(data webot.EventMsgData) {
	if data.MsgType == MSG_WITHDRAW {

		var reg *regexp.Regexp
		revokemsg := data.Content
		for k, v := range char_string {
			reg = regexp.MustCompile(`(`+k+`)`)
			revokemsg = reg.ReplaceAllString(revokemsg, v)
		}

		var revoke sysmsg

		err := xml.Unmarshal([]byte(revokemsg), &revoke)
		if err != nil {
			logger.Errorf("xml parse err:%v", err)
		}

		logger.Tracf("revoke:%v", revoke)

		w.Wechat.SendTextMsg(revoke.Revokemsg.Msgid, data.FromUserName)
	}
}

