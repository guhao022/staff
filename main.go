package main

import (
	"github.com/num5/axiom"
	"staff/adapter/wechat"
	"staff/listener"
)

func main() {
	b := axiom.New("Axiom")
	b.AddAdapter(wechat.NewWeChat(b))
	b.Register(listener.NewTime())

	b.Run()
}
