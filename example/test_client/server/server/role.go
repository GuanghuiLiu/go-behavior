package server

import (
	"github.com/GuanghuiLiu/behavior"
	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/tcp/long/game"
	game_pb "github.com/GuanghuiLiu/behavior/tcp/long/game/pb"
	"google.golang.org/protobuf/proto"
)

type Role struct {
	*behavior.Actor
	status uint8 // 状态 -> 游戏中断线重连，根据状态，限制行为
	ID     uint64
	UName  string
}

func NewRouter() behavior.IRouter {
	r := &Role{}
	return r
}

func (r *Role) SetRouter(actor *behavior.Actor) error {
	r.Actor = actor
	return nil
}
func (r *Role) Router(msg *model.Message) error {
	switch msg.Data.ProtoID {

	case game.Skill:
		skill, _ := r.decodeSkill(msg.Data.Data)
		skill.SkillID++
		skill.TargetID++

		r.Sent2client(game.Skill, skill)

	case game.Heartbeat:

		r.Sent2client(game.Heartbeat, nil)
	}
	return nil
}

func (r *Role) decodeSkill(data []byte) (*game_pb.Skill, error) {

	skill := &game_pb.Skill{}
	if err := proto.Unmarshal(data, skill); err != nil {
		return nil, err
	}
	return skill, nil
}
