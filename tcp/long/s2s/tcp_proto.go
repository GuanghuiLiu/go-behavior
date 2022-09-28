package name_server

// MaxPacketSize MaxPacketSize
const MaxPacketSize uint16 = 65535

// HeadLen dataLen 2,protoID 2
const HeadLen uint8 = 2 + 2

type Message struct {
	DataLen uint16 // 消息的长度 （2 字节）
	ProtoID uint16 // 协议ID （2 字节）
	Data    []byte // 消息的内容
}

const (
	CommonMessage uint16 = iota
	SentInfo
	SyncSent
	SyncResult
	StartProcess
	StopProcess
	BaseResponse
	GetNode
	GetNodeByPass // login
	NodeInfo
)
