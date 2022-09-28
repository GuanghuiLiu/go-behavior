package main

import (
	"time"

	server "github.com/GuanghuiLiu/behavior/example/game"
	cli "github.com/GuanghuiLiu/behavior/example/game/client"
)

var name = "client6"

func main() {
	cli.TestGetNode(name)
	time.Sleep(time.Second)
	cli.Start(name, server.ActionRoleFollow, 100, "npc1")
}
