package main

import (
	"fmt"
	"github.com/GuanghuiLiu/behavior/model"
	nameServer "github.com/GuanghuiLiu/behavior/tcp/short/name_server"
	"github.com/GuanghuiLiu/behavior/tcp/short/name_server/client"
	pb "github.com/GuanghuiLiu/behavior/tcp/short/name_server/pb"
	"github.com/GuanghuiLiu/behavior/utils"
	"google.golang.org/protobuf/proto"
	"net"
	"runtime/debug"
)

const (
	// NameServerIP string =
	Port uint = 9987
)

var ServerAddr = "192.168.1.36:9989"

type serverNameClient struct {
	msg *client.Message
	model.QuickHandel
	conn net.Conn
}

func newNameClient(name string) *serverNameClient {
	addrStr := fmt.Sprintf("%s:%d", utils.GetIP(), Port)
	conn, err := net.Dial("tcp", addrStr)
	if err != nil {
		panic(err.(any))
	}
	client := &serverNameClient{
		conn: conn,
		msg:  client.NewMessage(name),
	}
	client.Name = name
	go client.receive()
	return client
}
func (tc *serverNameClient) receive() bool {
	defer func() {
		if err := recover(); err != any(nil) {
			fmt.Println(tc.Name, "panic:", err, string(debug.Stack()))
		}
	}()

	serverMsg, err := tc.msg.UnpackHead(tc.conn)
	if err != nil {
		fmt.Println(tc.Name, "receive err", err)
		return true
	}
	tc.msg.UnpackData(tc.conn, serverMsg)
	switch serverMsg.ProtoID {
	case nameServer.NodeInfo:
		info := &pb.NodeInfo{}
		_ = proto.Unmarshal(serverMsg.Data, info)
		ServerAddr = info.Addr
		fmt.Println(tc.Name, "receive info", info)
	}
	return true
}

func (tc *serverNameClient) HandlerStop() model.Handler {
	tc.conn.Close()
	return tc
}

func (tc *serverNameClient) HandlerInfo(msg *model.Message) model.Handler {
	switch uint8(msg.Data.ProtoID) {
	case nameServer.GetNode:
		tc.getNode()
		// case nameServer.GetNodeByPass:
		// 	tc.getNodeByPass()
	}
	return tc
}

// func (tc *serverNameClient) getNode() {
// 	msg := &pb.GetNode{
// 		Uid: tc.Name,
// 	}
// 	tc.msg.SendMsg(tc.conn, nameServer.GetNode, msg)
// }
func (tc *serverNameClient) getNode() {
	msg := &pb.GetNode{
		Name:     tc.Name,
		Password: tc.Name,
	}
	tc.msg.SendMsg(tc.conn, nameServer.GetNode, msg)
}
