package gateway

import (
	"fmt"
	"net"

	"github.com/GuanghuiLiu/behavior"
	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/utils"
)

var port uint = 9989

func Start(toNameServer bool, newPort uint, router behavior.CreateRouter) {
	port = newPort
	startTcp(toNameServer, router)
}

func startTcp(nameServer bool, router behavior.CreateRouter) {

	addr, err := net.ResolveTCPAddr(IPVersion, getAddr())
	if err != nil {
		fmt.Println("resolve tcp addr err: ", err)
		return
	}

	listener, err := net.ListenTCP(IPVersion, addr)
	if err != nil {
		panic(err.(any))
	}

	if nameServer {
		if e := registerNode(getAddr()); e != nil {
			panic(e.(any))
		}
	}

	fmt.Println("listen at", port)

	for {

		conn, err2 := listener.AcceptTCP()

		if err2 != nil {
			fmt.Println("Accept err ", err2)
			conn.Close()
			continue
		}

		fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())

		if MaxConnID >= MaxConn {
			fmt.Println("conn out set ", MaxConnID)
			conn.Close()
			continue
		}

		conn.SetKeepAlive(true)

		runConn(conn, router)
	}
}

func runConn(c *net.TCPConn, router behavior.CreateRouter) {
	reader, err := newReader(c, router)
	if err != nil {
		panic(any(err))
	}
	reader.Run(reader, model.SetRetry(0, 0), model.SetDefaultFunc(reader.receiveFunc))
}

func getAddr() string {
	return fmt.Sprintf("%s:%d", utils.GetIP(), port)
}

func registerNode(addr string) error {
	return model.ClusterCenter.SetNodeAddr(model.C2S, addr)
}
