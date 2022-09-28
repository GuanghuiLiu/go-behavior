package main

import (
	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/utils"
)

const maxRound uint8 = 10

type Scheduler struct {
	model.SafeHandel
	tokenList   map[string]uint
	commitCount uint8
	npc         string
	round       uint8
	roleList    []string
}

func newScheduler(name string) *Scheduler {
	s := &Scheduler{tokenList: make(map[string]uint)}
	s.Name = name
	return s
}

func (scheduler *Scheduler) HandlerInfo(msg *model.Message) model.Handler {
	switch msg.Data.ProtoID {
	case NpcName:
		// scheduler.npc = cast.ToString(msg.Data.Data)
		utils.Decode(msg.Data.Data, &scheduler.npc)
	case StartGame:
		// roleList := cast.ToStringSlice(msg.Data.Data)
		var roleList []string
		utils.Decode(msg.Data.Data, &roleList)
		scheduler.roleList = roleList
		for _, role := range roleList {
			token := uint(randInt(100000))
			scheduler.tokenList[role] = token
			scheduler.commitCount++
			scheduler.SendInfo(role, StartGame, utils.Encode(token))
		}
		scheduler.SendInfo(scheduler.npc, StartGame, utils.Encode(roleList))
	// from npc
	case Rolechose:
		scheduler.commitCount--
		if scheduler.commitCount == 0 {
			scheduler.SendInfo(scheduler.npc, GameResule, nil)
		}
	// from role
	case RoleChoose:
		var tk uint
		utils.Decode(msg.Data.Data, &tk)
		if scheduler.tokenList[msg.From] == tk {
			scheduler.commitCount--
			if scheduler.commitCount == 0 {
				scheduler.SendInfo(scheduler.npc, GameResule, nil)
			}
		}
		// fmt.Println(scheduler.Name, "commitCount", scheduler.commitCount, msg.From, "token", cast.ToUint(*msg.Data.Data), "list", scheduler.tokenList)

	case GameResule:
		scheduler.round++
		if scheduler.round < maxRound {
			for _, role := range scheduler.roleList {
				token := uint(randInt(100000))
				scheduler.tokenList[role] = token
				scheduler.commitCount++
				scheduler.SendInfo(role, StartGame, utils.Encode(token))
			}
			scheduler.SendInfo(scheduler.npc, StartGame, utils.Encode(scheduler.roleList))
		}
	}
	return scheduler
}
