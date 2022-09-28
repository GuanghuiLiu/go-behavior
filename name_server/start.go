package name_server

import (
	"fmt"
	"net"

	"github.com/GuanghuiLiu/behavior/utils"
)

// var MaxConnID uint64
var Port uint = 9987

func Start(appName string, port uint) {
	initCluster(appName)
	initHeartbeat()
	Port = port
	startUserCenter(UserCenter)
	tcpStart()
}

func tcpStart() {

	addr, err := net.ResolveTCPAddr(IPVersion, fmt.Sprintf("%s:%d", utils.GetIP(), Port))
	if err != nil {
		fmt.Println("resolve tcp addr err: ", err)
		return
	}

	listener, err := net.ListenTCP(IPVersion, addr)
	if err != nil {
		panic(err.(any))
	}

	fmt.Println("listen", Port)

	for {

		conn, errC := listener.AcceptTCP()
		defer conn.Close()

		if errC != nil {
			fmt.Println("Accept err ", errC)
			continue
		}

		// fmt.Println("Get store remote addr = ", store.RemoteAddr().String())

		// if MaxConnID >= MaxConn {
		// 	continue
		// }

		go startRouter(conn)
	}
}
func startRouter(conn *net.TCPConn) {
	r := newRouter(conn)
	r.route()
}
