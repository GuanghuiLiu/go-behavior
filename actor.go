package behavior

import (
	"fmt"
	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/tcp/long/game"
	game_pb "github.com/GuanghuiLiu/behavior/tcp/long/game/pb"
	s2s "github.com/GuanghuiLiu/behavior/tcp/long/s2s"
	s2s_pb "github.com/GuanghuiLiu/behavior/tcp/long/s2s/pb"
	"github.com/GuanghuiLiu/behavior/utils"
	"google.golang.org/protobuf/proto"
)

type Actor struct {
	writer string
	model.QuickHandel
	Router IRouter
}

func NewActor(name, writer string, r IRouter) *Actor {
	a := &Actor{}
	a.Name = name
	a.writer = writer
	r.SetRouter(a)
	a.Router = r

	return a
}

func (a *Actor) HandlerStop() model.Handler {
	a.SendStop(a.writer)
	return a
}

func (a *Actor) HandlerInfo(msg *model.Message) model.Handler {
	switch msg.Data.ProtoID {

	// 接受返回，异步方式，处理同步事件
	case uint64(s2s.SyncResult):
		param := &s2s_pb.SyncResult{}
		proto.Unmarshal(msg.Data.Data, param)
		if param == nil && param.Base == nil {
			break
		}
		eventID := param.Base.EventID
		if param.Code == model.ResultOK {
			fmt.Println(a.Name, eventID, "ok")
		} else {
			fmt.Println(a.Name, eventID, "fail")
		}

	case game.ResetConn:
		var newConn string
		utils.Decode(msg.Data.Data, &newConn)
		if a.writer != newConn {
			a.SendStop(a.writer)
			a.writer = newConn
			a.loginSuccess()
		}
	case game.Heartbeat:

		a.Sent2client(game.Heartbeat, nil)

	default:

		a.Router.Router(msg)

	}
	return a

}

func (a *Actor) SentError2client(code uint32, msg string) error {

	err := &game_pb.ErrorInfo{
		Code: code,
		Msg:  msg,
	}
	return a.Sent2client(game.ErrorInfo, err)
}

func (a *Actor) Sent2client(protoID uint64, data proto.Message) error {

	binary, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	if errS := a.SendInfo(a.writer, protoID, binary); err != nil {
		return errS
	}
	return nil
}

func (a *Actor) Sent2clientByte(protoID uint64, data []byte) error {
	if err := a.SendInfo(a.writer, protoID, data); err != nil {
		return err
	}
	return nil
}

func (a *Actor) loginSuccess() error {
	m := &game_pb.LoginInfo{
		Text: "login success",
	}
	return a.Sent2client(game.LoginInfo, m)
}
