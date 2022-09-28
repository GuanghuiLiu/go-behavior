package main

import (
	_ "net/http/pprof"
	"time"

	server "github.com/GuanghuiLiu/behavior/example/game"
	cli "github.com/GuanghuiLiu/behavior/example/game/client"
)

var name = "client3"

func main() {
	cli.TestGetNode(name)
	time.Sleep(time.Second)
	cli.Start(name, server.ActionRoleGiveUP, 0, "npc1")
}
