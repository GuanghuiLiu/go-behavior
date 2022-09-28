package model

import (
	"fmt"
	"net"

	s2s "github.com/GuanghuiLiu/behavior/tcp/long/s2s"
	pb "github.com/GuanghuiLiu/behavior/tcp/long/s2s/pb"
	"github.com/GuanghuiLiu/behavior/tcp/long/s2s/receive"
	"google.golang.org/protobuf/proto"
)

type receiver struct {
	QuickHandel
	conn *net.TCPConn
}

func newRouter(name string, conn *net.TCPConn) *receiver {
	r := &receiver{}
	r.Name = "s2s_" + name
	r.conn = conn
	return r
}

func (r *receiver) receive() bool {

	msg, err0 := receive.Unpack(r.conn)
	// 当且仅当 tcp通信出错时，断开连接，终止程序
	if err0 != nil {
		r.sentErr(ErrSocketErr, err0.Error(), nil)
		return true
	}

	switch msg.ProtoID {

	case s2s.CommonMessage:
		data, err := r.unCodeCommon(msg)
		if err != nil {
			fmt.Println("unCode err", err)
			return false
		}

		switch uint16(data.Key) {
		case s2s.SyncSent:
			b, err2 := SentSync(data.Key, data.Message, r.Name, data.Base.To, 0, false)
			if err2 != nil {
				r.sentErr(ErrParamErr, err2.Error(), data.Base)
			}
			r.sentOK(b, data.Base)
		case s2s.SentInfo:
			SentInfo(data.Key, data.Message, r.Name, data.Base.To, false)
		case s2s.StopProcess:
			SentStop(r.Name, data.Base.To)
		}

	case s2s.StartProcess:
		param, err := r.unCodeStarModel(msg)
		if err != nil {
			r.sentErr(ErrUnCodeErr, err.Error(), param.Base)
			return false
		}
		// todo start process
		fmt.Println("param", param)
		r.sentOK(nil, param.Base)

	default:
		// todo

	}

	return false
}

func (r *receiver) HandlerStop() Handler {
	r.conn.Close()
	return r
}

func (r *receiver) sentOK(data []byte, base *pb.Base) {
	base.From, base.To = base.To, base.From
	r.send2client(s2s.SyncResult, &pb.SyncResult{
		Code:   ResultOK,
		Base:   base,
		Result: data,
	})
}
func (r *receiver) sentErr(code uint32, msg string, base *pb.Base) {
	base.From, base.To = base.To, base.From
	r.send2client(s2s.SyncResult, &pb.SyncResult{
		Code:    code,
		Message: msg,
		Base:    base,
	})
}

func (r *receiver) unCodeCommon(msg *s2s.Message) (*pb.CommonS2S, error) {

	param := &pb.CommonS2S{}
	proto.Unmarshal(msg.Data, param)
	return param, nil
}

func (r *receiver) unCodeSyncResult(msg *s2s.Message) (*pb.SyncResult, error) {

	param := &pb.SyncResult{}
	proto.Unmarshal(msg.Data, param)
	return param, nil
}

func (r *receiver) unCodeStarModel(msg *s2s.Message) (*pb.StarModel, error) {

	param := &pb.StarModel{}
	proto.Unmarshal(msg.Data, param)
	return param, nil
}

func (r *receiver) unCodeStopModel(msg *s2s.Message) (*pb.StopModel, error) {

	param := &pb.StopModel{}
	proto.Unmarshal(msg.Data, param)
	return param, nil
}

func (r *receiver) send2client(protoID uint16, data proto.Message) error {
	b, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	sentDate, errP := receive.Pack(protoID, b)
	if errP != nil {
		return errP
	}
	_, errS := r.conn.Write(sentDate)
	if errS != nil {
		return errS
	}
	return nil
}
