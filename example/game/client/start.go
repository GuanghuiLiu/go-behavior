package client

import (
	"fmt"
	"time"

	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/tcp/long/game"
	"github.com/GuanghuiLiu/behavior/tcp/short/name_server"
	"github.com/GuanghuiLiu/behavior/utils"
)

func TestGetNode(clientName string) {
	client := NewNameClient("clientGN1")
	client.Run(client, model.SetLiveTime(1), model.SetRetry(0, 0))
	e := model.SentInfo(uint64(name_server.GetNode), utils.Encode(clientName), "main", client.Name, false)
	if e != nil {
		fmt.Println(client.Name, "sentAction err", e)
	}
}

func Start(name string, action uint32, addCount uint64, npc string) {
	client := newClientWithLogin(name, action, addCount)
	time.Sleep(time.Second)
	// data := &pb.GameStart{
	// 	Npc: npc,
	// }
	// byteData, _ := proto.Marshal(data)
	model.SentInfo(game.GameStart, utils.Encode(npc), "main", client.Name, false)
	tk := time.NewTicker(20 * time.Second)
	defer tk.Stop()
	for {
		select {
		case <-tk.C:
			model.SentInfo(game.Heartbeat, nil, "main", client.Name, false)
			tk.Reset(300 * time.Second)
		}
	}
}
func newClientWithLogin(name string, action uint32, addCount uint64) *TcpClient {
	client := NewTcpClient(name, action, addCount)
	client.Run(client, model.SetLiveTime(10), model.SetRetry(1, 10*60))

	if err := model.SentInfo(game.Login, nil, "main", client.Name, false); err != nil {
		fmt.Println("err:", err)
	}
	return client
}
