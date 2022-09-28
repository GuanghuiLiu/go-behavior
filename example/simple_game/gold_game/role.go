package main

import (
	"fmt"
	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/utils"
	"github.com/spf13/cast"
	"strconv"
)

type Role struct {
	model.SafeHandel
	count     int
	retry     int
	score     int
	token     uint
	scheduler string
}

func NewRole(name string) *Role {
	r := &Role{}
	r.Name = name
	return r
}

func (role *Role) HandlerInfo(msg *model.Message) model.Handler {
	switch msg.Data.ProtoID {
	case StartGame:
		role.token = cast.ToUint(msg.Data.Data)
		role.scheduler = msg.From
	case GameInfo:
		slotGold := make(map[string]int)
		utils.Decode(msg.Data.Data, &slotGold)
		fmt.Println(role.Name, "can choose gold", slotGold)
		id := chooseGold(slotGold)
		fmt.Println(role.Name, "choose", id)
		role.SendInfo(msg.From, RoleChoose, utils.Encode(id))
	case GameResule:
		// role.score += cast.ToInt(*msg.Data.Data)
		var score int
		utils.Decode(msg.Data.Data, &score)
		role.score += score
		fmt.Println(role.Name, "get gold", score, ",all gold", role.score)
	}
	return role

}

func chooseGold(slotGold map[string]int) int {
	var slotWeight []int
	for k, _ := range slotGold {
		s, _ := strconv.Atoi(k)
		slotWeight = append(slotWeight, s)
	}
	return slotWeight[randInt(len(slotWeight))]
}
