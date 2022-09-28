package game

// MaxDataSize 拆包时，分段拆包；打包时不需要
const MaxDataSize uint32 = 4

const MaxPackSize uint32 = MaxDataSize + HeadLen

// HeadLen dataLen 4,protoID 8
const LengthDateLen uint32 = 4
const LengthProtoID uint32 = 8
const HeadLen uint32 = LengthDateLen + LengthProtoID

type Message struct {
	DataLen uint32 // 消息的长度 （4 字节）
	ProtoID uint64 // 协议ID （8 字节）
	Data    []byte // 消息的内容
}

// proto：key
// 在跨语言、跨节点交互中，用固定ID方便协调，不用自增ID
const (
	ResetConn uint64 = iota + 10000
	Login
	LoginInfo
	Heartbeat
	ErrorInfo
	Skill uint64 = iota + 10100
	EnterRoom
	GameStart
	GameTurn
	LeaveGame
	FollowAction
	GiveCards
	SettleResult
)
