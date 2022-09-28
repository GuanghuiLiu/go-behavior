package model

import (
	"errors"
	"fmt"
	"net"
	"runtime/debug"

	s2s "github.com/GuanghuiLiu/behavior/tcp/long/s2s"
	pb "github.com/GuanghuiLiu/behavior/tcp/long/s2s/pb"
	"github.com/GuanghuiLiu/behavior/tcp/long/s2s/sent"
	"google.golang.org/protobuf/proto"
)

const s2sID = "SyncS2S"

var mapIDChan = make(map[uint64]chan []byte)

func sentOtherNode(protoID uint64, data []byte, from, to string) error {
	return catchSent(protoID, data, from, to, false, nil)
}
func sentOtherNodeSync(protoID uint64, data []byte, from, to string, wait chan []byte) error {
	return catchSent(protoID, data, from, to, true, wait)
}

func catchSent(protoID uint64, data []byte, from, to string, sync bool, wait chan []byte) error {
	addr, err := ClusterCenter.getAddr(S2S, to)
	if err != nil {
		return clean(to, NullString)
	}
	conn, err := quickGetConn(addr)
	if err != nil {
		return clean(to, addr)
	}
	if err = sentDate(conn, protoID, data, from, to, sync, wait); err != nil {
		return clean(to, addr)
	}
	return nil
}

func clean(to, addr string) error {
	ClusterCenter.cleanNode(S2S, to)
	if addr != NullString {
		cleanConn(addr)
	}
	return errors.New("other server bad")
}

func repairGetConn(toAddr string) (net.Conn, error) {

	cleanConn(toAddr)
	conn, err := net.Dial("tcp", toAddr)
	if err != nil {
		return nil, err
	}
	ClusterCenter.localConn.Store(toAddr, conn)
	go clientReceive(toAddr, conn)

	return conn, nil
}

func quickGetConn(toAddr string) (net.Conn, error) {

	conn, ok := ClusterCenter.localConn.Load(toAddr)
	if ok {
		return conn.(net.Conn), nil
	}
	return repairGetConn(toAddr)
}

func sentDate(conn net.Conn, protoID uint64, data []byte, from, to string, sync bool, wait chan []byte) error {
	var eventID uint64
	if sync {
		eventID, _ = GetId(s2sID)
		mapIDChan[eventID] = wait
	}
	base := &pb.Base{From: from,
		To:      to,
		EventID: eventID,
	}
	toData := &pb.CommonS2S{
		Base:    base,
		Key:     protoID,
		Message: data,
	}

	e := sent.SendMsg(conn, s2s.CommonMessage, toData)
	if e != nil {
		return e
	}
	return nil
}

func cleanConn(addr string) {
	ClusterCenter.localConn.Delete(addr)
}

func clientReceive(addr string, conn net.Conn) {
	defer func() {
		if err := recover(); err != any(nil) {
			fmt.Println("s2s from", addr, "panic:", err, string(debug.Stack()))
		}
	}()
	defer conn.Close()
	defer cleanConn(addr)
	for {
		msg, err := sent.Unpack(conn)
		if err != nil {
			break
		}
		switch msg.ProtoID {
		// 读取其他节点返回
		case s2s.SyncResult:
			result, err2 := unCodeSyncResult(msg)
			if err2 != nil {
				return
			}
			if result == nil {
				break
			}
			if result.Code != ResultOK {
				break
			}
			if result.Base != nil {
				if result.Base.EventID > 0 {
					mapIDChan[result.Base.EventID] <- result.Result
				}
			}
		}
	}
}
func unCodeSyncResult(msg *s2s.Message) (*pb.SyncResult, error) {
	param := &pb.SyncResult{}
	proto.Unmarshal(msg.Data, param)
	return param, nil
}
