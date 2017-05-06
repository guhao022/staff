package wechat

import (
	"fmt"
	"github.com/num5/axiom"
	"github.com/num5/webot"
	"log"
	"strings"
)

const (
	// msg types
	MSG_TEXT        = 1     // 文本消息
	MSG_IMG         = 3     // 图片消息
	MSG_ADUP        = 4     // 通讯录更新
	MSG_MONEY       = 6     // 可能是红包
	MSG_VOICE       = 34    // 语音消息
	MSG_FV          = 37    // 朋友验证消息
	MSG_PF          = 40    // POSSIBLEFRIEND_MSG
	MSG_SCC         = 42    // 共享联系人卡片
	MSG_VIDEO       = 43    // 视频信息
	MSG_EMOTION     = 47    // gif
	MSG_LOCATION    = 48    // 位置信息
	MSG_LINK        = 49    // 共享链接的消息
	MSG_VOIP        = 50    // VOIPMSG
	MSG_INIT        = 51    // wechat 初始化消息
	MSG_VOIPNOTIFY  = 52    // VOIPNOTIFY
	MSG_VOIPINVITE  = 53    // VOIPINVITE
	MSG_SHORT_VIDEO = 62    // 短视频信息
	MSG_SYSNOTICE   = 9999  // 系统通知
	MSG_SYS         = 10000 // 系统消息
	MSG_WITHDRAW    = 10002 // 撤回通知消息

)

type WeChatAdapter struct {
	bot    *axiom.Robot
	Wechat *webot.WeChat
}

func NewWeChat(bot *axiom.Robot) *WeChatAdapter {

	wechat, err := webot.AwakenNewBot(nil)
	if err != nil {
		panic(err)
	}

	return &WeChatAdapter{
		bot:    bot,
		Wechat: wechat,
	}
}

var x *xiaoice

// 初始化
func (w *WeChatAdapter) Prepare() error {

	x = newXiaoice(w.Wechat)

	w.Wechat.Handle(`/login`, func(webot.Event) {
		if cs, err := w.Wechat.ContactsByNickName(`小冰`); err == nil {
			for _, c := range cs {
				if c.Type == webot.Offical {
					x.un = c.UserName // 更新小冰的UserName
					break
				}
			}
		}
	})

	return nil
}

func (w *WeChatAdapter) GetName() string {
	return "wechat robot"
}

// 解析
func (w *WeChatAdapter) Process() error {
	//
	w.Wechat.Handle(`/msg`, func(evt webot.Event) {
		msg := evt.Data.(webot.EventMsgData)

		if msg.IsGroupMsg {

			if msg.AtMe {
				realcontent := strings.TrimSpace(strings.Replace(msg.Content, "@"+w.Wechat.MySelf.NickName, "", 1))

				if realcontent == "统计人数" {
					stat, err := w.chatRoomMember(msg.FromUserName)
					if err == nil {
						ans := fmt.Sprintf("群里一共有 %d 人，其中男生 %d 人， 女生 %d 人，未知性别者 %d 人 (ó-ò) ", stat["count"], stat["man"], stat["woman"], stat["none"])

						w.Wechat.SendTextMsg(ans, msg.FromUserName)
					} else {
						w.Wechat.SendTextMsg(err.Error(), msg.FromUserName)
					}
				} else {
					amsg := axiom.Message{
						Text:    realcontent,
						ReplyTo: []interface{}{msg.FromUserName},
					}

					w.bot.ReceiveMessage(amsg)
				}
			}

		} else {

			amsg := axiom.Message{
				Text:    msg.Content,
				ReplyTo: []interface{}{msg.FromUserName},
			}
			x.autoReplay(msg)

			w.bot.ReceiveMessage(amsg)
		}

	})

	w.Wechat.Go()

	return nil
}

// 回应
func (w *WeChatAdapter) Reply(msg axiom.Message, message string) error {

	w.Wechat.SendTextMsg(message, msg.ReplyTo[0].(string))

	return nil
}

// 获取群组用户
func (w *WeChatAdapter) chatRoomMember(room_name string) (map[string]int, error) {

	stats := make(map[string]int)

	RoomContactList, err := w.Wechat.MembersOfGroup(room_name)
	if err != nil {
		return nil, err
	}

	man := 0
	woman := 0
	none := 0
	for _, v := range RoomContactList {

		member, err := w.Wechat.ContactByGGID(v.GGID)

		if err != nil {
			log.Printf("[ERRO] 抓取组群用户 [%s] 信息失败: %s... ", v.NickName)
		} else {
			if member.Sex == 1 {
				man++
			} else if member.Sex == 2 {
				woman++
			} else {
				none++
			}
		}

	}

	stats = map[string]int{
		"count": len(RoomContactList),
		"woman": woman,
		"man":   man,
		"none":  none,
	}

	return stats, nil
}
