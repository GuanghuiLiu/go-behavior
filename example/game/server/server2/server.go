package main

import (
	"fmt"

	"github.com/GuanghuiLiu/behavior/example/game/server"
	"github.com/GuanghuiLiu/behavior/gateway"
	"github.com/GuanghuiLiu/behavior/model"
)

func main() {
	model.InitCluster("common", "bj_0001_00092", false, 9977)
	npc := server.NewNpc("npc1", 3)

	err := npc.Run(npc)
	if err != nil {
		fmt.Println(err)
		return
	}
	f := server.NewRouter
	gateway.Start(true, 9988, f)
}
