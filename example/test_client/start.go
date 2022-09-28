package main

import (
	"fmt"
	_ "net/http/pprof"
	"time"

	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/tcp/long/game"
	"github.com/GuanghuiLiu/behavior/tcp/short/name_server"
	"github.com/GuanghuiLiu/behavior/utils"
)

func main() {

	// go func() {
	// 	log.Fatal(http.ListenAndServe(":6060", nil))
	// }()

	// model.initCluster("serverClient")
	go testGetNode()
	time.Sleep(time.Second)
	go testLiveTime()
	go testPack()
	go testCommon()
	for {

	}
}

func testGetNode() {
	client := newNameClient("clientGN1")
	client.Run(client, model.SetLiveTime(1), model.SetRetry(0, 1))
	e := model.SentInfo(uint64(name_server.GetNode), utils.Encode(client.Name), "main", client.Name, false)
	if e != nil {
		fmt.Println(client.Name, "sent err", e)
	}
	// model.SentInfo(uint64(name_server.GetNodeByPass), nil, "main", client1.Name)
	for {

	}
}
func testPack() {

	client1 := newClient("clientP1")
	client2 := newClient("clientP2")
	client3 := newClient("clientP3")
	tk := time.NewTicker(3 * time.Second)
	defer tk.Stop()
	model.SentInfo(game.Skill, nil, "main", client1.Name, false)
	time.Sleep(time.Second)
	model.SentInfo(game.Login, nil, "main", client1.Name, false)
	for {
		select {
		case <-tk.C:
			model.SentInfo(game.Heartbeat, nil, "main", client1.Name, false)
			model.SentInfo(game.Skill, nil, "main", client1.Name, false)
			model.SentInfo(game.Heartbeat, nil, "main", client2.Name, false)
			model.SentInfo(game.Skill, nil, "main", client2.Name, false)
			model.SentInfo(game.Heartbeat, nil, "main", client3.Name, false)
			model.SentInfo(game.Skill, nil, "main", client3.Name, false)
			tk.Reset(8 * time.Second)
		}
	}
}
func testCommon() {

	client1 := newClientWithLogin("clientC1")
	client2 := newClientWithLogin("clientC2")
	client3 := newClientWithLogin("clientC3")
	tk := time.NewTicker(3 * time.Second)
	defer tk.Stop()
	for {
		select {
		case <-tk.C:
			model.SentInfo(game.Heartbeat, nil, "main", client1.Name, false)
			model.SentInfo(game.Skill, nil, "main", client1.Name, false)
			model.SentInfo(game.Heartbeat, nil, "main", client2.Name, false)
			model.SentInfo(game.Skill, nil, "main", client2.Name, false)
			model.SentInfo(game.Heartbeat, nil, "main", client3.Name, false)
			model.SentInfo(game.Skill, nil, "main", client3.Name, false)
			tk.Reset(30 * time.Second)
		}
	}
}
func newClientWithLogin(name string) *TcpClient {
	client := newTcpClient(name)
	client.Run(client, model.SetLiveTime(1), model.SetRetry(1, 10*60))

	if err := model.SentInfo(game.Login, nil, "main", client.Name, false); err != nil {
		fmt.Println("err:", err)
	}
	return client
}
func newClient(name string) *TcpClient {
	client := newTcpClient(name)
	client.Run(client, model.SetLiveTime(10), model.SetRetry(0, 0))
	return client
}
func testLiveTime() {

	client1 := newClient("clientL1")
	model.SentInfo(game.Heartbeat, nil, "main", client1.Name, false)
	model.SentInfo(game.Skill, nil, "main", client1.Name, false)
	for {

	}
}
