package name_server

import (
	"net"
)

func getIP() ([]string, error) {
	var ipStr []string
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return ipStr, err
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					//获取IPv6
					/*if ipnet.IP.To16() != nil {
					    fmt.Println(ipnet.IP.String())
					    ipStr = append(ipStr, ipnet.IP.String())

					}*/
					//获取IPv4
					if ipnet.IP.To4() != nil {
						ipStr = append(ipStr, ipnet.IP.String())

					}
				}
			}
		}
	}
	return ipStr, nil
}
