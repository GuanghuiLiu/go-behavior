package client

const ProtoIDRepair uint64 = iota + 10000

const (
	ErrParam uint32 = iota + 1
	ErrStatus
	ErrTurn
	ErrBaseGoldNumbe
)

// type actionType uint
// const (
// 	ActionGiveCardSelf uint32 = iota + 1
// 	ActionGiveCardGlobal
// 	ActionRoleGiveUP
// 	ActionRolePass
// 	ActionRoleFollow
// 	ActionRoleGiveAdd
// )
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

// proto：key,内部协议ID在10000以内，与客户端协议在10000以上
const (
	Error uint64 = iota + 1
	NpcStart
	NpcGetOne
	StartGame = iota + 100
	// ShuffleCards
	Card2Role
	RoleAddGolds
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
const (
	H2 uint8 = iota + 2
	H3
	H4
	H5
	H6
	H7
	H8
	H9
	H10
	H11
	H12
	H13
	H1
)
const (
	Ho2 uint8 = iota + 2 + CardTypeGapLV1
	Ho3
	Ho4
	Ho5
	Ho6
	Ho7
	Ho8
	Ho9
	Ho10
	Ho11
	Ho12
	Ho13
	Ho1
)
const (
	M2 uint8 = iota + 2 + CardTypeGapLV2
	M3
	M4
	M5
	M6
	M7
	M8
	M9
	M10
	M11
	M12
	M13
	M1
)
const (
	F2 uint8 = iota + 2 + CardTypeGapLV3
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
	F13
	F1
)
