package model

import (
	"fmt"
	"net"

	"github.com/GuanghuiLiu/behavior/utils"
)

var portS2S uint = 9979

func startS2S(port uint) {
	if port > 1000 && port < 65535 {
		portS2S = port
	}
	startTcp(IPVersion, getAddr())
}

func startTcp(ipVersion, addr string) {

	TcpAddr, err := net.ResolveTCPAddr(ipVersion, addr)
	if err != nil {
		fmt.Println("resolve tcp addr err: ", err)
		return
	}

	listener, err := net.ListenTCP(ipVersion, TcpAddr)
	if err != nil {
		panic(err.(any))
	}

	fmt.Println("s2s listen at", addr)

	if e := registerNode(getAddr()); e != nil {
		panic(e.(any))
	}

	for {

		conn, err := listener.AcceptTCP()

		if err != nil {
			fmt.Println("Accept err ", err)
			conn.Close()
			continue
		}

		fmt.Println("s2s Get conn remote addr = ", conn.RemoteAddr().String())
		conn.SetKeepAlive(true)
		startRouter(conn.RemoteAddr().String(), conn)
	}
}

func startRouter(node string, conn *net.TCPConn) error {
	r := newRouter(node, conn)
	return r.Run(r, SetRetry(0, 0), SetDefaultFunc(r.receive))

}
func getAddr() string {
	return fmt.Sprintf("%s:%d", utils.GetIP(), portS2S)
}

func registerNode(addr string) error {
	return ClusterCenter.SetNodeAddr(S2S, addr)
}
