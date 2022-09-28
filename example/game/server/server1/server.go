package main

import (
	"fmt"

	"github.com/GuanghuiLiu/behavior/example/game/server"
	"github.com/GuanghuiLiu/behavior/gateway"
	"github.com/GuanghuiLiu/behavior/model"
)

func main() {
	// 向云端存数据
	model.InitCluster("common", "bj_0001_00091", false, 9979)

	npc := server.NewNpc("npc1", 3)
	err := npc.Run(npc)
	if err != nil {
		fmt.Println(err)
		return
	}
	f := server.NewRouter
	gateway.Start(false, 9989, f)
}
