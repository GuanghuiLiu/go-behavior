package name_server

import (
	"net"
	"time"

	"github.com/GuanghuiLiu/behavior/model"
)

const (
	heartbeatTime = 60 * time.Second
	heartbeatName = "heartbeatName"

	connMaxTimes uint8 = 2
	connDuration       = 1 * time.Second
)

func initHeartbeat() {
	h := newHeartBeat(heartbeatName)
	h.Run(h, model.SetDefaultFunc(h.check))
}

type heartBeat struct {
	model.QuickHandel
	duration time.Duration
}

func newHeartBeat(name string) *heartBeat {
	h := &heartBeat{}
	h.Name = name
	h.duration = heartbeatTime
	return h
}

func (h *heartBeat) check() bool {
	time.Sleep(h.duration)
	clusterData.localNodeAddr.Range(checkConn)
	return false
}

func checkConn(node, addr any) bool {

	conn, ok := clusterData.localNodeConn.Load(node)
	if ok && conn != nil {
		return true
	}
	connTimes := ZeroUint8
label:
	conn, err := net.Dial("tcp", addr.(string))
	if connTimes < connMaxTimes && err != nil {
		if err != nil {
			connTimes++
			time.Sleep(connDuration)
			goto label
		}
	}
	if err != nil {
		clusterData.delNodeAddr(node.(string))
		return true
	}
	clusterData.localNodeConn.Store(node, conn)
	return true
}
