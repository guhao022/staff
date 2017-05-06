package wechat

import (
	"github.com/num5/axiom"
	"github.com/num5/webot"
)

// 撤销取回
func (w *WeChatAdapter) revoke(data webot.EventMsgData) {
	if data.MsgType == MSG_WITHDRAW {
		amsg := axiom.Message{
			Text:    data.Content,
			ReplyTo: []interface{}{data.FromUserName},
		}
		w.bot.ReceiveMessage(amsg)
	}
}
