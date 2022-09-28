package main

import (
	"fmt"

	"github.com/GuanghuiLiu/behavior/model"
)

func main() {
	npc := npc{}
	npc.Name = "npc"
	err := npc.Run(&npc)
	if err != nil {
		fmt.Println(err)
		return
	}
	role1 := Role{}
	role1.Name = "role1"
	err = role1.Run(&role1)
	if err != nil {
		fmt.Println(err)
		return
	}
	role2 := Role{}
	role2.Name = "role2"
	err = role2.Run(&role2)
	if err != nil {
		fmt.Println(err)
		return
	}
	role3 := Role{}
	role3.Name = "role3"
	err = role3.Run(&role3)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < 10; i++ {
		err1 := model.SentInfo(StartCompete, nil, "main", role1.Name, false)
		err2 := model.SentInfo(StartCompete, nil, "main", role2.Name, false)
		err3 := model.SentInfo(StartCompete, nil, "main", role3.Name, false)
		if err1 != nil && err2 != nil && err3 != nil {
			fmt.Println(err1, err2, err3)
			return
		}
	}
	for {
	}
}
