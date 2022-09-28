package main

import (
	"fmt"

	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/utils"
)

func main() {
	scheduler := newScheduler("scheduler")
	err0 := scheduler.Run(scheduler)
	if err0 != nil {
		fmt.Println(err0)
		return
	}
	npc := NewNpc("npc")
	err := npc.Run(npc)
	if err != nil {
		fmt.Println(err)
		return
	}
	model.SentInfo(NpcName, utils.Encode(npc.Name), "main", scheduler.Name, false)
	role1 := NewRole("role1")
	err = role1.Run(role1)
	if err != nil {
		fmt.Println(err)
		return
	}
	role2 := NewRole("role2")
	err = role2.Run(role2)
	if err != nil {
		fmt.Println(err)
		return
	}
	role3 := NewRole("role3")
	err = role3.Run(role3)
	if err != nil {
		fmt.Println(err)
		return
	}
	role4 := NewRole("role4")
	err = role4.Run(role4)
	if err != nil {
		fmt.Println(err)
		return
	}
	err1 := model.SentInfo(StartGame, utils.Encode([]string{role1.Name, role2.Name, role3.Name, role4.Name}), "main", scheduler.Name, false)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	for {
	}
}
