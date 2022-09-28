package server

import (
	"github.com/GuanghuiLiu/behavior"
	server "github.com/GuanghuiLiu/behavior/example/game"
	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/tcp/long/game"
	game_pb "github.com/GuanghuiLiu/behavior/tcp/long/game/pb"
	"github.com/GuanghuiLiu/behavior/utils"
	"google.golang.org/protobuf/proto"
)

type Role struct {
	*behavior.Actor
	status uint8
	ID     uint64
	UName  string
	golds  uint64
	runTimeInfo
}
type runTimeInfo struct {
	npc       string
	roundID   uint64
	goldLevel uint64
	cards     [2]uint8
}

func NewRouter() behavior.IRouter {
	r := &Role{}
	run := runTimeInfo{
		// isS2S: false,
	}
	r.runTimeInfo = run
	r.status = RoleStatusOver
	return r
}

func (r *Role) SetRouter(actor *behavior.Actor) error {
	r.Actor = actor
	return nil
}

func (r *Role) Router(msg *model.Message) error {
	switch msg.Data.ProtoID {
	// case game.EnterRoom:
	// 	en, err := r.decodeEnterRoom(msg.Data.Data)
	// 	if err != nil {
	// 		r.SentError2client(ErrParam, err.Error())
	// 	}
	// 	if en.RoomID == 0 {
	// 	}

	case RoleAddGolds:
		var golds uint64
		utils.Decode(msg.Data.Data, &golds)
		r.golds += golds
	case StartGame:
		r.status = RoleStatusWait
		r.Sent2clientByte(game.GameStart, msg.Data.Data)
	case Card2Role:
		r.Sent2clientByte(game.FollowAction, msg.Data.Data)
	case game.GameStart:
		if r.status == RoleStatusTurn {
			r.Sent2clientByte(game.GameTurn, nil)
		} else if r.status == RoleStatusOver {
			data := &game_pb.GameStart{}
			proto.Unmarshal(msg.Data.Data, data)
			r.npc = data.Npc
			r.SendInfo(r.npc, game.GameStart, msg.Data.Data)
		}
	case game.GiveCards:
		r.status = RoleStatusWait
		cardsInfo, _ := r.decodeGiveCards(msg.Data.Data)
		if cardsInfo.Action == server.ActionGiveCardSelf {
			r.cards[0] = uint8(cardsInfo.First)
			r.cards[1] = uint8(cardsInfo.Second)
		}
		r.Sent2clientByte(game.GiveCards, msg.Data.Data)
	case game.GameTurn:
		r.status = RoleStatusTurn
		r.Sent2clientByte(game.GameTurn, msg.Data.Data)
	case game.FollowAction:
		r.status = RoleStatusWait
		fa, _ := r.decodeFollowAction(msg.Data.Data)
		r.checkAction(fa.Action, fa.Count)
		r.SendInfo(r.npc, game.FollowAction, msg.Data.Data)
	case game.SettleResult:
		r.status = RoleStatusOver
		r.Sent2clientByte(game.SettleResult, msg.Data.Data)
	case game.ErrorInfo:
		r.Sent2clientByte(game.ErrorInfo, msg.Data.Data)
	}
	return nil
}
func (r *Role) checkAction(action uint32, number uint64) error {
	return nil
}

func (r *Role) decodeFollowAction(data []byte) (*game_pb.FollowAction, error) {

	a := &game_pb.FollowAction{}
	if err := proto.Unmarshal(data, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (r *Role) decodeEnterRoom(data []byte) (*game_pb.EnterRoom, error) {

	en := &game_pb.EnterRoom{}
	if err := proto.Unmarshal(data, en); err != nil {
		return nil, err
	}
	return en, nil
}
func (r *Role) decodeGameStart(data []byte) (*game_pb.GameStart, error) {

	gs := &game_pb.GameStart{}
	if err := proto.Unmarshal(data, gs); err != nil {
		return nil, err
	}
	return gs, nil
}
func (r *Role) decodeGiveCards(data []byte) (*game_pb.GiveCards, error) {

	gc := &game_pb.GiveCards{}
	if err := proto.Unmarshal(data, gc); err != nil {
		return nil, err
	}
	return gc, nil
}
