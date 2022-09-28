package main

import (
	"fmt"
	"github.com/GuanghuiLiu/behavior/utils"
	"strconv"

	"github.com/GuanghuiLiu/behavior/model"
)

type slot struct {
	id         int
	isGold     bool
	goldNumber int
	allowRoles []int
	roles      []string
	roleID     int
}

var slots []*slot
var roles []string
var roleGold map[string]int

type Npc struct {
	model.SafeHandel
	Number    int
	count     int
	scheduler string
}

func init() {
	slots = make([]*slot, 9)
	roleGold = make(map[string]int)
	for i := 0; i < 9; i++ {
		s := new(slot)
		s.id = i + 1
		if i%2 == 0 {
			s.isGold = true
			s.allowRoles = getAllowRole(i + 1)
		} else {
			s.roleID = getSlotRole(s.id)
		}
		slots[i] = s
	}
}
func NewNpc(name string) *Npc {
	n := &Npc{}
	n.Name = name
	return n
}
func (npc *Npc) HandlerInfo(msg *model.Message) model.Handler {
	switch msg.Data.ProtoID {
	case StartGame:
		npc.scheduler = msg.From
		npc.count++
		fmt.Println(npc.Name, "start game rand", npc.count)
		for _, s := range slots {
			if s.isGold {
				s.goldNumber = randInt(npc.count*BaseGoldNumbe*len(s.allowRoles)) + 1
			}
		}
		// roles = cast.ToStringSlice(msg.Data.Data)
		utils.Decode(msg.Data.Data, &roles)
		for i, role := range roles {
			slotGold := make(map[string]int)
			for _, s := range slots {
				if s.isGold {
					if isAtInt(i+1, s.allowRoles) {
						slotGold[strconv.Itoa(s.id)] = s.goldNumber
					}
				}
			}
			npc.SendInfo(role, GameInfo, utils.Encode(slotGold))
		}
	case RoleChoose:
		var slotID int
		utils.Decode(msg.Data.Data, &slotID)
		s := slots[slotID-1]
		s.roles = append(s.roles, msg.From)
		npc.SendInfo(npc.scheduler, RoleChoose, nil)

	case GameResule:
		for _, s := range slots {
			if s.isGold {
				if len(s.roles) == 1 {
					roleGold[s.roles[0]] += s.goldNumber
					npc.SendInfo(s.roles[0], GameResule, utils.Encode(s.goldNumber))
				} else {
					for _, role := range s.roles {
						npc.SendInfo(role, GameResule, nil)
					}
				}
				s.roles = []string{}
			}
		}
		npc.SendInfo(npc.scheduler, GameResule, nil)
		fmt.Println(npc.Name, "over game rand", npc.count, roleGold)
	}
	return npc

}

func getAllowRole(i int) []int {
	res := make([]int, 0)
	switch i {
	case 1:
		res = append(res, getSlotRole(2), getSlotRole(4))
	case 3:
		res = append(res, getSlotRole(2), getSlotRole(6))
	case 5:
		res = append(res, getSlotRole(2), getSlotRole(4), getSlotRole(6), getSlotRole(8))
	case 7:
		res = append(res, getSlotRole(4), getSlotRole(8))
	case 9:
		res = append(res, getSlotRole(6), getSlotRole(8))
	}
	return res
}
func getSlotRole(slot int) int {
	switch slot {
	case 2:
		return 1
	case 4:
		return 2
	case 6:
		return 3
	case 8:
		return 4
	}
	return 0
}
