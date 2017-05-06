package main

import (
	"github.com/num5/axiom"
	"staff/adapter/wechat"
)

func main() {
	b := axiom.New("Axiom")
	b.AddAdapter(wechat.NewWeChat(b))
	//b.Register()

	b.Run()
}
