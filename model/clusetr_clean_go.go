package model

import (
	"time"
)

const (
	heartbeatTime       = 60 * time.Second
	heartbeatName       = "heartbeat"
	connMaxTimes  uint8 = 3
	connDuration        = 1 * time.Second
)

func initS2SHeartbeat() {
	h := newHeartBeat(heartbeatName)
	h.Run(h, SetDefaultFunc(h.check))
}

type heartBeat struct {
	QuickHandel
	duration time.Duration
	tk       *time.Ticker
}

func newHeartBeat(name string) *heartBeat {
	h := &heartBeat{}
	h.Name = name
	h.duration = heartbeatTime
	h.tk = time.NewTicker(h.duration)
	return h
}

func (h *heartBeat) check() bool {
	select {
	case <-h.tk.C:
		ClusterCenter.localConn.Range(checkConn)
		h.tk.Reset(h.duration)
		return false
	}
}

// func checkConn(node, addr any) bool {
// 	_, err := net.Dial("tcp", addr.(string))
// 	if err != nil {
// 		ClusterCenter.delNodeAddr(S2S, node.(string))
// 	}
// 	return true
// }

// todo 采集云端数据，更新至本地，测试连通性
//  定时拉取云端节点信息，发送其他节点心跳，当发现有节点宕机，清理云端数据和本地的云端数据
func checkConn(addr, conn any) bool {

	if conn != nil {
		return true
	}

	if _, err := repairGetConn(addr.(string)); err != nil {
		ClusterCenter.deleteNodeProcessList(addr.(string))
	}
	return true
}
