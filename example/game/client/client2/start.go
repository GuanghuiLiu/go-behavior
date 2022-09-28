package main

import (
	_ "net/http/pprof"
	"time"

	server "github.com/GuanghuiLiu/behavior/example/game"
	cli "github.com/GuanghuiLiu/behavior/example/game/client"
)

var name = "client2"

func main() {
	cli.TestGetNode(name)
	time.Sleep(time.Second)
	cli.Start(name, server.ActionRoleGiveAdd, 100, "npc1")
}
