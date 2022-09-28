package server

const (
	ErrParam uint32 = iota + 1
	ErrStatus
	ErrTurn
	ErrBaseGoldNumbe
)

// type actionType uint
const (
	ActionGiveCardSelf uint32 = iota + 1
	ActionGiveCardGlobal
	ActionRoleGiveUP
	ActionRolePass
	ActionRoleFollow
	ActionRoleGiveAdd
)
const (
	SeatNull uint8 = iota + 1
	SeatGiveup
	SeatRoleFollow
	SeatPass
	SeatGiveAdd
	SeatStartWaiting
	SeatDeal
)

// 皇家同花顺、同花顺、四条、葫芦、同花、顺子、三条、两对、一对、高牌
const (
	CardTypeErr uint8 = iota
	CardTypeHJTHS
	CardTypeTHS
	CardTypeSIT
	CardTypeHL
	CardTypeTH
	CardTypeSZ
	CardTypeSANT
	CardTypeLD
	CardTypeYD
	CardTypeGP
)

const (
	RoundStatusWait uint8 = iota + 1
	RoundStatusStarting
)

const (
	CardTypeGapLV1 = 40
	CardTypeGapLV2 = CardTypeGapLV1 * 2
	CardTypeGapLV3 = CardTypeGapLV1 * 3
)
