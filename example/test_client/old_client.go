package main

import (
	"fmt"
	"net"
	"runtime/debug"

	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/tcp/long/game"
	"github.com/GuanghuiLiu/behavior/tcp/long/game/client"
	pb "github.com/GuanghuiLiu/behavior/tcp/long/game/pb"
	"github.com/GuanghuiLiu/behavior/utils"
	"google.golang.org/protobuf/proto"
)

type tcpClientOld struct {
	msg *client.Message
	model.QuickHandel
	conn net.Conn
}

func newTcpClientOld(name string) *TcpClient {
	// ServerAddr = "192.168.1.36:9989"
	conn, err := net.Dial("tcp", ServerAddr)
	if err != nil {
		panic(err.(any))
	}
	client := &TcpClient{
		conn: conn,
		msg:  client.NewMessage(name),
	}
	client.Name = name

	go client.receive()
	return client
}
func (tc *tcpClientOld) receive() {
	defer func() {
		if err := recover(); err != any(nil) {
			fmt.Println(tc.Name, "panic:", err, string(debug.Stack()))
		}
	}()
	for {
		serverMsg, err := tc.msg.UnpackHead(tc.conn)
		if err != nil {
			fmt.Println(tc.Name, "receive err", err)
			break
		}
		tc.msg.UnpackData(tc.conn, serverMsg)
		switch serverMsg.ProtoID {
		case game.LoginInfo:
			loginInfo := &pb.LoginInfo{}
			_ = proto.Unmarshal(serverMsg.Data, loginInfo)
			fmt.Println(tc.Name, "receive loginInfo", loginInfo)
		case game.Skill:
			skill := &pb.Skill{}
			proto.Unmarshal(serverMsg.Data, skill)
			fmt.Println(tc.Name, "receive skill", skill)
		}
	}
}

func (tc *tcpClientOld) HandlerStop() model.Handler {
	// fmt.Println(tc.Name, "stop3")
	tc.conn.Close()
	return tc
}

func (tc *tcpClientOld) HandlerInfo(msg *model.Message) model.Handler {
	// fmt.Println(tc.Name, "handler", msg.Data.ProtoID)
	switch msg.Data.ProtoID {
	case game.Login:
		tc.Login()
	case game.Heartbeat:
		tc.Heartbeat()
	case game.Skill:
		tc.Skill()
	}
	return tc
}

func (tc *tcpClientOld) Login() {
	msg := &pb.Login{
		Name:     tc.Name,
		Password: "testClient",
	}
	tc.msg.SendMsg(tc.conn, game.Login, msg)
}

func (tc *tcpClientOld) Heartbeat() {
	msg := &pb.Heartbeat{}
	tc.msg.SendMsg(tc.conn, game.Heartbeat, msg)
}

func (tc *tcpClientOld) Skill() {
	msg := &pb.Skill{
		SkillID:  utils.RandUint64n(100),
		TargetID: utils.RandUint64n(10000),
	}
	tc.msg.SendMsg(tc.conn, game.Skill, msg)
}
