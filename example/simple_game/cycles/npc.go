package main

import (
	"fmt"
	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/utils"
)

type npc struct {
	model.QuickHandel
	// count int
}

func (npc *npc) HandlerInfo(msg *model.Message) model.Handler {
	var opponent uint8
	utils.Decode(msg.Data.Data, &opponent)
	self := randCompete()
	fmt.Println(npc.Name, "show:", self, msg.From)
	result, ok := rule(opponent, self)
	if ok {
		npc.SendInfo(msg.From, uint64(CompeteResule), utils.Encode(result))
	}
	return npc

}
