package main

import (
	"fmt"
	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/utils"
)

type Role struct {
	model.QuickHandel
	count int
	retry int
	score int
}

func (role *Role) HandlerInfo(msg *model.Message) model.Handler {
	switch msg.Data.ProtoID {
	case StartCompete:
		self := randCompete()
		fmt.Println(role.Name, "show:", self)
		role.SendInfo("npc", NoAction, utils.Encode(self))
	case CompeteResule:
		var result uint8
		utils.Decode(msg.Data.Data, &result)
		if result == Dogfall && role.retry < 30 {
			role.retry++
			self := randCompete()
			fmt.Println(role.Name, "show:", self)
			role.SendInfo("npc", uint64(NoAction), utils.Encode(self))
		}
		if result == Win {
			role.score++
			role.count++
			fmt.Println(role.Name, "in this game: win, score is", role.score, "count", role.count)
		}
		if result == Lose {
			role.score--
			role.count++
			fmt.Println(role.Name, "in this game: lose, score is", role.score, "count", role.count)
		}

	}
	return role

}
