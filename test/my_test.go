package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

func TestChan(t *testing.T) {
	// c := make(chan int, 0)
	// result
	// my_test.go:25: 9 5
	// my_test.go:25: 4 9
	// my_test.go:25: 0 1
	// my_test.go:25: 1 2
	// my_test.go:25: 2 3
	// my_test.go:25: 3 6
	// my_test.go:25: 6 5
	// my_test.go:25: 5 7
	// my_test.go:25: 7 8
	// my_test.go:25: 8 1
	c := make(chan int, 10)
	// result
	//    my_test.go:25: 9 9
	//    my_test.go:25: 5 6
	//    my_test.go:25: 4 4
	//    my_test.go:25: 1 1
	//    my_test.go:25: 0 0
	//    my_test.go:25: 6 7
	//    my_test.go:25: 2 2
	//    my_test.go:25: 7 8
	//    my_test.go:25: 3 4
	//    my_test.go:25: 8 9
	go func() {
		for {
			v := <-c
			c <- v + 1
		}
	}()
	for i := 0; i < 10; i++ {
		go func(i int) {
			c <- i
			v := <-c
			t.Log(i, v)
		}(i)
	}
	time.Sleep(1 * time.Second)
}

func TestIP(t *testing.T) {
	ipLocal()
	ip()
	ipAndMac()
	t.Log(getLocalIP())
}

func getLocalIP() []string {
	var ipStr []string
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces error:", err.Error())
		return ipStr
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
						fmt.Println(ipnet.IP.String())
						ipStr = append(ipStr, ipnet.IP.String())

					}
				}
			}
		}
	}
	return ipStr

}

func ipLocal() {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("local", conn.LocalAddr().String())
}
func ipAndMac() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			fmt.Println("ip: ", ip.String(), "mac: ", iface.HardwareAddr.String())
		}
	}
}
func ip() {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						fmt.Println(ipnet.IP.String())
					}
				}
			}
		}
	}
	//结果
	//169.254.86.238
	//192.168.1.2
}

func TestBinary(t *testing.T) {
	circleBinary()
}
func TestString(t *testing.T) {
	var a []string
	a = []string{"s", "d", "f"}
	b, err := json.Marshal(a)
	if err != nil {
		t.Log(err)
	}
	result := string(b)
	fmt.Println(result)
}

type CircleBinary struct {
	maxSize uint32
	array   [5]byte
	head    uint32
	tail    uint32
}

func (this *CircleBinary) IsFull() bool {
	return (this.tail+1)%this.maxSize == this.head
}

func (this *CircleBinary) IsEmpty() bool {
	return this.tail == this.head
}

func (this *CircleBinary) AddQueue(val byte) (err error) {
	if this.IsFull() {
		return errors.New("队列已满")
	}

	this.array[this.tail] = val
	this.tail = (this.tail + 1) % this.maxSize
	return
}

//出队列
func (this *CircleBinary) GetQueue() (val byte, err error) {
	if this.IsEmpty() {
		return 0, errors.New("队列已空")
	}
	val = this.array[this.head]
	this.head = (this.head + 1) % this.maxSize
	return
}

//显示队列元素
func (this *CircleBinary) ListQueue() {
	fmt.Println("队列情况如下：")
	//计算出队列多少元素
	//比较关键的一步
	size := (this.tail + this.maxSize - this.head) % this.maxSize
	if size == 0 {
		fmt.Println("队列已空")
	}
	//定义一个辅助变量 指向head
	tempHead := this.head
	for i := 0; i < int(size); i++ {
		fmt.Printf("arr[%d]=%d\t", tempHead, this.array[tempHead])
		tempHead = (tempHead + 1) % this.maxSize
	}
	fmt.Println()
}

func circleBinary() {
	queue := &CircleBinary{
		maxSize: 5,
		head:    0,
		tail:    0,
	}
	var key string
	var val byte
	for {
		fmt.Println("1. 输入add 表示添加数据到队列")
		fmt.Println("2. 输入get 表示从队列获取数据")
		fmt.Println("3. 输入show 表示显示队列")
		fmt.Println("4. 输入exit 表示显示队列")

		fmt.Print("请输入:")
		fmt.Scanln(&key)
		switch key {
		case "add":
			fmt.Println("输入你要入队列数")
			fmt.Scanln(&val)
			err := queue.AddQueue(val)
			if err != nil {
				fmt.Println(err.Error())
			} else {

				fmt.Println("加入队列ok")
			}
		case "get":
			val, err := queue.GetQueue()
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("从队列中取出了一个数=", val)
			}
		case "show":
			queue.ListQueue()
		case "exit":
			os.Exit(0)
		}
	}
}
