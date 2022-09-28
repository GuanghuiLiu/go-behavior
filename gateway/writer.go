package gateway

import (
	"net"

	"github.com/GuanghuiLiu/behavior/model"
	proto_game "github.com/GuanghuiLiu/behavior/tcp/long/game"
	pb "github.com/GuanghuiLiu/behavior/tcp/long/game/pb"
	"google.golang.org/protobuf/proto"
)

// writer作为actor的原因：net.TCPConn 不可传递，没法重置链接，只能取出、关掉
type writer struct {
	model.QuickHandel
	conn   *net.TCPConn
	doPack proto_game.Packer
}

func newWriter(reader string, conn *net.TCPConn) *writer {
	w := &writer{
		conn: conn,
	}
	w.Name = NameWriterPrefix + reader
	return w
}

func (w *writer) HandlerStop() model.Handler {
	if w.conn != nil {
		w.conn.Close()
	}
	return w
}

func (w *writer) HandlerInfo(msg *model.Message) model.Handler {

	// 只接收来自当前roleProcess的消息
	if e := w.send2Client(msg.Data.ProtoID, msg.Data.Data); e != nil {
		w.conn.Close()
	}

	return w
}

func (w *writer) send2Client(protoID uint64, data []byte) error {

	return w.doPack.Send2Client(w.conn, protoID, data)

}

//
func (w *writer) loginInfo() {
	m := &pb.LoginInfo{
		Id: MaxConnID,
		// Name: e.Name,
		Text: "login success",
	}
	binaryData, _ := proto.Marshal(m)
	w.send2Client(proto_game.LoginInfo, binaryData)
}
func (w *writer) reLogin() {
	m := &pb.LoginInfo{
		Text: "not login,place re login",
	}
	binaryData, _ := proto.Marshal(m)
	w.send2Client(proto_game.LoginInfo, binaryData)
}
