package client

import (
	"fmt"
	server "github.com/GuanghuiLiu/behavior/example/game"
	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/tcp/long/game"
	"github.com/GuanghuiLiu/behavior/tcp/long/game/client"
	pb "github.com/GuanghuiLiu/behavior/tcp/long/game/pb"
	"github.com/GuanghuiLiu/behavior/utils"
	"google.golang.org/protobuf/proto"
	"net"
	"runtime/debug"
	"time"
)

type TcpClient struct {
	msg *client.Message
	model.QuickHandel
	conn     net.Conn
	sentType uint32
	addCount uint64
}

func NewTcpClient(name string, sentType uint32, addCount uint64) *TcpClient {
	client := &TcpClient{
		msg: client.NewMessage(name),
	}
	client.Name = name
	client.sentType = sentType
	client.addCount = addCount
	// client.newConn()
	// go client.receive()
	return client
}

func (tc *TcpClient) newConn() {
label:
	fmt.Println(tc.Name, "server addr is", ServerAddr)
	conn, err := net.Dial("tcp", ServerAddr)
	if err != nil {
		tc.SendInfo("clientGN1", ProtoIDRepair, utils.Encode(tc.Name))
		time.Sleep(time.Second * 2)
		goto label
	}
	tc.conn = conn
	go tc.receive()
}

func (tc *TcpClient) receive() {
	defer func() {
		if err := recover(); err != any(nil) {
			fmt.Println(tc.Name, "panic:", err, string(debug.Stack()))
		}
	}()
	for {
		serverMsg, err := tc.msg.UnpackHead(tc.conn)
		if err != nil {
			fmt.Println(tc.Name, "receive err", err)
			break
		}
		tc.msg.UnpackData(tc.conn, serverMsg)
		switch serverMsg.ProtoID {
		case game.LoginInfo:
			loginInfo := &pb.LoginInfo{}
			_ = proto.Unmarshal(serverMsg.Data, loginInfo)
			fmt.Println(tc.Name, "receive loginInfo", loginInfo)
		case game.GameStart:
			other := &pb.GameStart{}
			proto.Unmarshal(serverMsg.Data, other)
			fmt.Println(tc.Name, "role seat", other)
		case game.GiveCards:
			cards := &pb.GiveCards{}
			proto.Unmarshal(serverMsg.Data, cards)
			if cards.Action == 1 {
				fmt.Println(tc.Name, cards.Role, "底牌", tc.formatCards(append([]uint32{}, cards.First, cards.Second)))
			} else {
				fmt.Println(tc.Name, "NPC 发牌,第", cards.First, "张：", tc.formatCards(append([]uint32{}, cards.Second)))
			}
		case game.FollowAction:
			action := &pb.FollowAction{}
			proto.Unmarshal(serverMsg.Data, action)
			fmt.Println(action.Role, tc.formatAction(action.Action), ",奖池:", action.Count)
		case game.GameTurn:
			tc.sentAction()
		case game.SettleResult:
			result := &pb.SettleResult{}
			proto.Unmarshal(serverMsg.Data, result)
			if result.Golds > 0 {
				fmt.Println(tc.Name, result.Role, "赢", result.Golds, ",", tc.formatModel(uint8(result.CardModel)), tc.formatCards(result.Cards))
			} else {
				fmt.Println(tc.Name, result.Role, "result", tc.formatModel(uint8(result.CardModel)), tc.formatCards(result.Cards))
			}
		case game.ErrorInfo:
			err2 := &pb.GameStart{}
			proto.Unmarshal(serverMsg.Data, err2)
			fmt.Println(tc.Name, "err", err2)

		}
	}
}

func (tc *TcpClient) sentAction() {
	time.Sleep(3 * time.Second)
	if tc.addCount == 0 && tc.sentType == server.ActionRoleFollow {
		tc.addCount = 100
	}
	if tc.conn != nil {
		msg := &pb.FollowAction{
			Action: tc.sentType,
			Count:  tc.addCount,
		}
		tc.msg.SendMsg(tc.conn, game.FollowAction, msg)
		// fmt.Println("轮到我了，", tc.formatAction(msg.Action), msg.Count)
	}
}

func (tc *TcpClient) HandlerStop() model.Handler {
	// fmt.Println(tc.Name, "stop3")
	if tc.conn != nil {
		tc.conn.Close()
	}
	return tc
}

func (tc *TcpClient) HandlerInfo(msg *model.Message) model.Handler {
	// fmt.Println(tc.Name, "handler", msg.Data.ProtoID)
	switch msg.Data.ProtoID {
	case game.Login:
		tc.Login()
	case game.Heartbeat:
		tc.Heartbeat()
	case game.GameStart:
		tc.gameStart(msg.Data.Data)
	}
	return tc
}

func (tc *TcpClient) Login() {
	tc.newConn()
	msg := &pb.Login{
		Name:     tc.Name,
		Password: "testClient",
	}
	tc.msg.SendMsg(tc.conn, game.Login, msg)
}

func (tc *TcpClient) Heartbeat() {
	if tc.conn != nil {
		msg := &pb.Heartbeat{}
		tc.msg.SendMsg(tc.conn, game.Heartbeat, msg)
	}

}

func (tc *TcpClient) gameStart(data []byte) {
	if tc.conn != nil {
		var npc string
		utils.Decode(data, &npc)
		msg := &pb.GameStart{
			Seat: 0,
			Npc:  npc,
		}
		tc.msg.SendMsg(tc.conn, game.GameStart, msg)
	}
}
func (tc *TcpClient) formatAction(action uint32) string {
	switch action {
	case server.ActionRoleFollow:
		return "跟注"
	case server.ActionRoleGiveAdd:
		return "加注"
	case server.ActionRoleGiveUP:
		return "弃牌"
	case server.ActionRolePass:
		return "过"
	}
	return "error"
}

func (tc *TcpClient) formatModel(mod uint8) string {
	switch mod {
	case CardTypeHJTHS:
		return "皇家同花顺"
	case CardTypeTHS:
		return "同花顺"
	case CardTypeSIT:
		return "四条"
	case CardTypeHL:
		return "葫芦"
	case CardTypeTH:
		return "同花"
	case CardTypeSZ:
		return "顺子"
	case CardTypeSANT:
		return "三条"
	case CardTypeLD:
		return "两对"
	case CardTypeYD:
		return "一对"
	case CardTypeGP:
		return "高牌"

	}
	return "error"
}
func (tc *TcpClient) formatCards(cards []uint32) []string {
	s := []string{}
	for _, card := range cards {
		switch {
		case card < CardTypeGapLV1:
			s = append(s, fmt.Sprint("♠️", tc.formatCardValue(card)))
		case card > CardTypeGapLV1 && card < CardTypeGapLV2:
			s = append(s, fmt.Sprint("♥️", tc.formatCardValue(card)))
		case card > CardTypeGapLV2 && card < CardTypeGapLV3:
			s = append(s, fmt.Sprint("♣️", tc.formatCardValue(card)))
		case card > CardTypeGapLV3:
			s = append(s, fmt.Sprint("♦️", tc.formatCardValue(card)))
		}
	}
	return s
}
func (tc *TcpClient) formatCardValue(card uint32) string {
	switch {
	case card%CardTypeGapLV1 == 2:
		return "2"
	case card%CardTypeGapLV1 == 3:
		return "3"
	case card%CardTypeGapLV1 == 4:
		return "4"
	case card%CardTypeGapLV1 == 5:
		return "5"
	case card%CardTypeGapLV1 == 6:
		return "6"
	case card%CardTypeGapLV1 == 7:
		return "7"
	case card%CardTypeGapLV1 == 8:
		return "8"
	case card%CardTypeGapLV1 == 9:
		return "9"
	case card%CardTypeGapLV1 == 10:
		return "10"
	case card%CardTypeGapLV1 == 11:
		return "J"
	case card%CardTypeGapLV1 == 12:
		return "Q"
	case card%CardTypeGapLV1 == 13:
		return "K"
	case card%CardTypeGapLV1 == 14:
		return "A"
	}
	return "error"
}
