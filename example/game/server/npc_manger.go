package server

import (
	"fmt"
	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/utils"
)

type npcManger struct {
	model.QuickHandel
	npcList map[uint32]string
}

func NewNpcManger(name string) *npcManger {
	n := &npcManger{}
	n.Name = name
	n.npcList = make(map[uint32]string)
	if err := n.startNpc(1); err != nil {
		return n
	}
	n.npcList[0] = "npc"

	return n
}

func (m *npcManger) HandlerInfo(msg *model.Message) model.Handler {
	switch msg.Data.ProtoID {
	case NpcStart:
		var n uint32
		utils.Decode(msg.Data.Data, &n)
		m.startNpc(n)
	case NpcGetOne:
		var n uint32
		utils.Decode(msg.Data.Data, &n)
		if npcName, ok := m.npcList[n]; ok {
			m.SendInfo(msg.From, NpcGetOne, utils.Encode(npcName))
			break
		}
		m.SendInfo(msg.From, NpcGetOne, utils.Encode(m.npcList[0]))
	}
	return m
}

func (m *npcManger) startNpc(n uint32) error {
	for i := 0; i < int(n); i++ {
		name := fmt.Sprintf("npc-%d", i)
		npc := NewNpc(name, 12)
		err := npc.Run(npc)
		if err != nil {
			return err
		}
		m.npcList[uint32(i)] = npc.Name
	}
	return nil
}
