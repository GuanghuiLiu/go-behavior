package server

import (
	"errors"
	"time"

	server "github.com/GuanghuiLiu/behavior/example/game"
	"github.com/GuanghuiLiu/behavior/model"
	"github.com/GuanghuiLiu/behavior/tcp/long/game"
	game_pb "github.com/GuanghuiLiu/behavior/tcp/long/game/pb"
	"github.com/GuanghuiLiu/behavior/utils"
	"github.com/golang/protobuf/proto"
)

type Npc struct {
	roundId       uint64
	roundStatus   uint8
	cards         *[52]uint8
	globalCards   []uint8
	global3Open   bool
	cardSite      uint8
	seats         *[]*seatStatus
	seatSite      uint8
	liveRoleCount uint8
	roleCountSet  uint8
	golds         uint64
	goldLevel     uint64
	goldLeveler   int8
	waitRole      time.Ticker
	model.QuickHandel
}
type seatStatus struct {
	id            uint8
	role          string
	putGolds      uint64
	maybeWinGolds uint64
	status        uint8
	firstCard     uint8
	secondCard    uint8
	maxModel      uint8
	maxCards      []uint8
}

func newCards() *[52]uint8 {
	return &[52]uint8{
		H1, H2, H3, H4, H5, H6, H7, H8, H9, H10, H11, H12, H13,
		Ho1, Ho2, Ho3, Ho4, Ho5, Ho6, Ho7, Ho8, Ho9, Ho10, Ho11, Ho12, Ho13,
		M1, M2, M3, M4, M5, M6, M7, M8, M9, M10, M11, M12, M13,
		F1, F2, F3, F4, F5, F6, F7, F8, F9, F10, F11, F12, F13,
	}
}

func NewNpc(name string, roles uint8) *Npc {
	n := &Npc{}
	n.Name = name
	n.cards = newCards()
	n.initNpc(roles)
	return n
}

func (n *Npc) initSeat(long uint8) {
	seats := []*seatStatus{}
	for i := 0; uint8(i) < long; i++ {
		seats = append(seats, &seatStatus{
			status: SeatNull,
		})
	}
	n.seats = &seats
}
func (n *Npc) initNpc(long uint8) {
	n.roundStatus = RoundStatusWait
	n.cardSite = 0
	n.global3Open = false
	n.roundId++
	n.goldLeveler = -1
	n.goldLevel = 10
	n.golds = 0
	n.globalCards = nil
	n.roleCountSet = long
	n.liveRoleCount = 0
	n.initSeat(n.roleCountSet)
}

func (n *Npc) HandlerInfo(msg *model.Message) model.Handler {
	switch msg.Data.ProtoID {
	case game.GameStart:
		if n.roundStatus != RoundStatusWait || n.liveRoleCount >= n.roleCountSet {
			isLive := false
			for _, role := range *n.seats {
				if msg.From == role.role {
					isLive = true
					break
				}
			}
			if !isLive {
				n.sentErr2client(msg.From, ErrStatus, "游戏已经开始了")
				break
			}
		}
		var seat uint32
		utils.Decode(msg.Data.Data, &seat)
		s, err := n.setSeat(msg.From, seat)
		if err != nil {
			break
		}
		s++
		gs := &game_pb.GameStart{
			Role: msg.From,
			Seat: s,
		}
		b, _ := proto.Marshal(gs)
		n.broadcast(StartGame, b)
		if n.liveRoleCount == n.roleCountSet {
			n.roundStatus = RoundStatusStarting
			n.givePersonalCardAnd3Global()
		}
	case game.FollowAction:
		fa, _ := n.decodeFollowAction(msg.Data.Data)
		fa.Role = msg.From
		n.doAction(fa)
	}
	return n

}

func (n *Npc) givePersonalCardAnd3Global() error {
	shuffleCards(n.cards)
	for i := 0; i < len(*n.seats); i++ {
		(*(n.seats))[i].firstCard = n.cards[2*i]
		n.cardSite++
		(*(n.seats))[i].secondCard = n.cards[2*i+1]
		n.cardSite++
		n.giveCards2Client(server.ActionGiveCardSelf, uint32(n.cards[2*i]), uint32(n.cards[2*i+1]), (*(n.seats))[i].role, (*(n.seats))[i].role)
	}
	n.give3GlobalCard()
	n.seatSite = uint8(n.roundId) % n.roleCountSet
	n.goldLeveler = int8(n.seatSite)
	n.SendInfo((*(n.seats))[n.seatSite].role, game.GameTurn, nil)
	return nil
}

func (n *Npc) give3GlobalCard() error {
	n.globalCards = append(n.globalCards, n.cards[n.cardSite], n.cards[n.cardSite+1], n.cards[n.cardSite+2])
	n.cardSite = n.cardSite + 3
	return nil
}
func (n *Npc) turnOff3GlobalCard() error {
	for i, v := range n.globalCards {
		n.globalCard2Client(server.ActionGiveCardGlobal, uint32(i+1), uint32(v))
	}
	n.global3Open = true
	n.goldLevel = 0
	n.SendInfo((*(n.seats))[n.seatSite].role, game.GameTurn, nil)
	return nil
}

func (n *Npc) giveGlobalCard() error {
	n.globalCards = append(n.globalCards, n.cards[n.cardSite])
	n.globalCard2Client(server.ActionGiveCardGlobal, uint32(len(n.globalCards)), uint32(n.cards[n.cardSite]))
	n.cardSite++
	n.goldLevel = 0
	n.SendInfo((*(n.seats))[n.seatSite].role, game.GameTurn, nil)
	return nil
}

func (n *Npc) giveCards2Client(action, first, second uint32, to, role string) (bool, error) {
	cards := &game_pb.GiveCards{
		Action: action,
		First:  first,
		Second: second,
		Role:   role,
	}
	b, _ := proto.Marshal(cards)
	n.SendInfo(to, game.GiveCards, b)
	return true, nil
}

func (n *Npc) globalCard2Client(action, first, second uint32) error {
	cards := &game_pb.GiveCards{
		Action: action,
		First:  first,
		Second: second,
	}
	b, _ := proto.Marshal(cards)
	n.broadcast(game.GiveCards, b)
	return nil
}

// 指定为位置坐下
func (n *Npc) setSeat(role string, seat uint32) (uint32, error) {
	if n.liveRoleCount >= n.roleCountSet || n.roundStatus != RoundStatusWait {
		return 0, errors.New("full , no seat")
	}

	if seat > uint32(n.roleCountSet) || seat == 0 {
		return n.insertSeat(role)
	}
	seat--
	if (*(n.seats))[seat].status == SeatNull {
		(*(n.seats))[seat].role = role
		n.liveRoleCount++
		(*(n.seats))[seat].status = SeatStartWaiting
		return seat, nil
	}
	return n.insertSeat(role)
}

// 自动坐下
func (n *Npc) insertSeat(role string) (uint32, error) {
	for i, v := range *(n.seats) {
		if v.status == SeatNull {
			(*(n.seats))[i].role = role
			(*(n.seats))[i].status = SeatStartWaiting
			n.liveRoleCount++
			return uint32(i), nil
		}
	}
	return 0, errors.New("not insert seat")
}

func (n *Npc) doAction(fa *game_pb.FollowAction) error {
	if (*(n.seats))[n.seatSite].role != fa.Role || n.roundStatus != RoundStatusStarting {
		n.sentErr2client(fa.Role, ErrTurn, "没到你呢，等着")
		return errors.New("error turn")
	}
	if fa.Action == server.ActionRoleFollow {
		(*(n.seats))[n.seatSite].status = SeatRoleFollow
		n.golds += (n.goldLevel - (*(n.seats))[n.seatSite].putGolds)
		(*(n.seats))[n.seatSite].putGolds = n.goldLevel
	}
	if fa.Action == server.ActionRoleGiveUP {
		n.liveRoleCount--
		(*(n.seats))[n.seatSite].status = SeatGiveup
		if n.goldLeveler == int8(n.seatSite) {
			n.goldLeveler = -1
		}
	}
	if fa.Action == server.ActionRoleGiveAdd {
		(*(n.seats))[n.seatSite].status = SeatGiveAdd
		n.goldLevel = fa.Count
		n.golds += fa.Count
		n.goldLeveler = int8(n.seatSite)
	}
	n.broadcastAction(fa)
	return n.nextRoleTurn()
}

func (n *Npc) nextSeatSite() uint8 {
	n.seatSite = (n.seatSite + 1) % n.roleCountSet
	if (*(n.seats))[n.seatSite].status == SeatGiveup {
		return n.nextSeatSite()
	}
	return n.seatSite
}
func (n *Npc) nextRoleTurn() error {
	n.nextSeatSite()
	if n.goldLeveler == int8(n.seatSite) {
		return n.nextNpcAction()
	}
	if n.goldLeveler == -1 {
		n.goldLeveler = int8(n.seatSite)
	}
	n.SendInfo((*(n.seats))[n.seatSite].role, game.GameTurn, nil)
	return nil
}
func (n *Npc) nextNpcAction() error {
	if n.liveRoleCount == 1 {
		return n.doResult()
	}
	if len(n.globalCards) == 5 {
		return n.doResult()
	}
	n.cleanPutGolds()
	if len(n.globalCards) == 3 && !n.global3Open {
		return n.turnOff3GlobalCard()
	}
	return n.giveGlobalCard()
}

// todo 暂时只做简单结算，只给最后唯一赢家；不分堆
func (n *Npc) doResult() error {
	maxSeat := seatStatus{}
	for _, seat := range *(n.seats) {
		personalCards := make([]uint8, 7, 7)
		copy(personalCards, append(n.globalCards, seat.firstCard, seat.secondCard))
		cardModel, cards := judgeCardLv(personalCards)
		if cardModel == CardTypeErr {
			return errors.New("server err")
		}
		seat.maxCards = cards
		seat.maxModel = cardModel
		n.broadcastCard(cardModel, cards, seat.role)
		if seat.status != SeatGiveup {
			if maxSeat.role == "" {
				maxSeat = *seat
				continue
			}
			// model 越小，等级越高
			if maxSeat.maxModel > seat.maxModel {
				maxSeat = *seat
				continue
			}
			if maxSeat.maxModel == seat.maxModel {
				for i := 0; i < len(seat.maxCards); i++ {
					if seat.maxCards[i]%CardTypeGapLV1 < maxSeat.maxCards[i]%CardTypeGapLV1 {
						break
					}
					if seat.maxCards[i]%CardTypeGapLV1 > maxSeat.maxCards[i]%CardTypeGapLV1 {
						maxSeat = *seat
						break
					}
				}
			}
		}
	}
	// 处理多人获胜
	winners := []seatStatus{maxSeat}
	for _, seat := range *(n.seats) {
		if seat.role != maxSeat.role && seat.status != SeatGiveup && seat.maxModel == maxSeat.maxModel {
			equal := true
			for i := 0; i < len(maxSeat.maxCards); i++ {
				if seat.maxCards[i]%CardTypeGapLV1 != maxSeat.maxCards[i]%CardTypeGapLV1 {
					equal = false
					break
				}
			}
			if equal {
				winners = append(winners, *seat)
			}
		}
		seat.status = SeatDeal
	}
	n.broadcastWinners(winners)
	n.initNpc(n.roleCountSet)
	return nil
}

func (n *Npc) broadcastAction(fa *game_pb.FollowAction) error {
	fa.Count = n.golds
	actionBinary, _ := proto.Marshal(fa)
	n.broadcast(Card2Role, actionBinary)
	return nil
}

func (n *Npc) broadcastWinners(seats []seatStatus) {
	for _, w := range seats {
		winner := &game_pb.SettleResult{}
		winner.CardModel = uint32(w.maxModel)
		winner.Role = w.role
		winner.Golds = n.golds / uint64(len(seats))
		winner.Cards = uint8toUint32List(w.maxCards)
		b, _ := proto.Marshal(winner)
		n.broadcast(game.SettleResult, b)
	}
}
func (n *Npc) broadcastCard(cardModel uint8, cards []uint8, role string) {
	personal := &game_pb.SettleResult{}
	personal.CardModel = uint32(cardModel)
	personal.Role = role
	personal.Cards = uint8toUint32List(cards)
	b, _ := proto.Marshal(personal)
	n.broadcast(game.SettleResult, b)
}

func (n *Npc) broadcast(protoID uint64, data []byte) error {
	for _, seat := range *(n.seats) {
		if seat.status != SeatNull {
			n.SendInfo(seat.role, protoID, data)
		}
	}
	return nil
}

func (n *Npc) cleanPutGolds() error {
	for _, s := range *(n.seats) {
		if s.status != SeatGiveup {
			s.putGolds = 0
		}
	}
	return nil
}

func (n *Npc) decodeFollowAction(data []byte) (*game_pb.FollowAction, error) {

	a := &game_pb.FollowAction{}
	if err := proto.Unmarshal(data, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (n *Npc) sentErr2client(role string, code uint32, msg string) error {
	gs := &game_pb.ErrorInfo{
		Code: code,
		Msg:  msg,
	}
	b, _ := proto.Marshal(gs)
	n.SendInfo(role, game.ErrorInfo, b)
	return nil
}
