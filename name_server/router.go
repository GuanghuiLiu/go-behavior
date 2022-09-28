package name_server

import (
	"fmt"
	"net"

	"github.com/GuanghuiLiu/behavior/tcp/short/name_server"
	pb "github.com/GuanghuiLiu/behavior/tcp/short/name_server/pb"
	"google.golang.org/protobuf/proto"
)

type router struct {
	conn *net.TCPConn
}

func newRouter(conn *net.TCPConn) *router {
	return &router{
		conn: conn,
	}
}

func (r *router) route() {
	defer r.conn.Close()
	msg, err := name_server.Unpack(r.conn)

	if err != nil {
		return
	}

	switch msg.ProtoID {
	// case name_server.GetNode:
	//
	// 	register, err := r.unCodeGetNode(msg)
	// 	if err != nil {
	// 		r.sentErr(name_server.ErrDataUnKnow, err.Error())
	// 		// return err
	// 	}
	//
	// 	addr, errGN := getNode(register.Uid, register.IsRepair)
	// 	if errGN != nil {
	// 		r.sentErr(name_server.ErrDataUnKnow, errGN.Error())
	// 		// return errGN
	// 	}
	// 	nodeInfo.Addr = addr
	case name_server.Registyer:

		register, err2 := r.unCodeRegister(msg)
		if err2 != nil {
			r.sentErr(name_server.ErrDataUnKnow, err2.Error())
			// return err
		}

		if ok := r.registerUser(register.Name, register.Password); !ok {
			r.sentErr(name_server.ErrDataUnKnow, "register fail")
		}

		r.getAddr(register.Name, register.Password, true)

	case name_server.GetNode:

		gnb, errUP := r.unCodeGetNode(msg)
		if errUP != nil {
			r.sentErr(name_server.ErrDataUnKnow, errUP.Error())
			// return err
		}
		r.getAddr(gnb.Name, gnb.Password, gnb.IsRepair)

	default:
		r.sentErr(name_server.ErrDataUnKnow, "not found proto")
		// return errors.New("not found proto")
	}

	// return nil
}

func (r *router) getAddr(name, passWord string, isRepair bool) {

	addr, err := getNodeByPass(name, passWord, isRepair)
	if err != nil {
		r.sentErr(name_server.ErrDataUnKnow, err.Error())
		return
	}

	nodeInfo := &pb.NodeInfo{Addr: addr}

	if err = r.send2client(name_server.NodeInfo, nodeInfo); err != nil {
		fmt.Println("sent err", err)
	}
}

func (r *router) unCodeGetNode(msg *name_server.Message) (*pb.GetNode, error) {

	gn := &pb.GetNode{}
	proto.Unmarshal(msg.Data, gn)
	return gn, nil
}
func (r *router) unCodeRegister(msg *name_server.Message) (*pb.Register, error) {

	gn := &pb.Register{}
	proto.Unmarshal(msg.Data, gn)
	return gn, nil
}
func (r *router) unCodeGetNodeByPass(msg *name_server.Message) (*pb.GetNode, error) {

	g := &pb.GetNode{}
	proto.Unmarshal(msg.Data, g)
	return g, nil
}
func (r *router) sentErr(code uint32, msg string) error {
	e := &pb.Error{
		Code: code,
		Msg:  msg,
	}
	r.send2client(name_server.Error, e)
	return nil
}

func (r *router) send2client(protoID uint8, data proto.Message) error {
	b, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	sentDate, errP := name_server.Pack(protoID, b)
	if errP != nil {
		return errP
	}
	_, errS := r.conn.Write(sentDate)
	if errS != nil {
		return errS
	}
	return nil
}
func (r *router) registerUser(name, pass string) bool {
	return true
}

func (r *router) checkUser(name, pass string) (string, error) {
	uid := name
	return uid, nil
}
