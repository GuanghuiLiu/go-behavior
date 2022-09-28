package client

import (
	"fmt"
	"github.com/GuanghuiLiu/behavior/model"
	nameServer "github.com/GuanghuiLiu/behavior/tcp/short/name_server"
	"github.com/GuanghuiLiu/behavior/tcp/short/name_server/client"
	pb "github.com/GuanghuiLiu/behavior/tcp/short/name_server/pb"
	"github.com/GuanghuiLiu/behavior/utils"
	"google.golang.org/protobuf/proto"
	"net"
)

const (
	// NameServerIP string =
	Port uint = 9987
)

var ServerAddr = "192.168.1.36:9989"

type serverNameClient struct {
	isRepair bool
	msg      *client.Message
	model.QuickHandel
	conn net.Conn
}

func NewNameClient(name string) *serverNameClient {
	client := &serverNameClient{
		msg: client.NewMessage(name),
	}
	client.Name = name
	return client
}
func (tc *serverNameClient) receive() {
	defer tc.conn.Close()
	serverMsg, err := tc.msg.UnpackHead(tc.conn)
	if err != nil {
		fmt.Println(tc.Name, "receive err", err)
		return
	}
	tc.msg.UnpackData(tc.conn, serverMsg)
	switch serverMsg.ProtoID {
	case nameServer.NodeInfo:
		info := &pb.NodeInfo{}
		_ = proto.Unmarshal(serverMsg.Data, info)
		ServerAddr = info.Addr
		fmt.Println(tc.Name, "receive info", info)
	}
	return
}

func (tc *serverNameClient) HandlerStop() model.Handler {
	if tc.conn != nil {
		tc.conn.Close()
	}
	return tc
}

func (tc *serverNameClient) HandlerInfo(msg *model.Message) model.Handler {
	switch msg.Data.ProtoID {
	case uint64(nameServer.GetNode):
		tc.getNode(msg.Data.Data)
	case ProtoIDRepair:
		tc.isRepair = true
		tc.getNode(msg.Data.Data)
	}
	return tc
}

// func (tc *serverNameClient) getNode() {
// 	msg := &pb.GetNode{
// 		Uid: tc.Name,
// 	}
// 	tc.msg.SendMsg(tc.conn, nameServer.GetNode, msg)
// }
func (tc *serverNameClient) getNode(name []byte) {
	addrStr := fmt.Sprintf("%s:%d", utils.GetIP(), Port)
	conn, err := net.Dial("tcp", addrStr)
	if err != nil {
		return
	}
	tc.conn = conn
	go tc.receive()
	var n string
	utils.Decode(name, &n)
	msg := &pb.GetNode{
		Name:     n,
		Password: n,
		IsRepair: tc.isRepair,
	}
	tc.msg.SendMsg(conn, nameServer.GetNode, msg)
}
