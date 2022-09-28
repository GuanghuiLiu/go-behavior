package name_server

// MaxPacketSize MaxPacketSize
const MaxPacketSize uint8 = 255

// HeadLen dataLen 1,protoID 1
const HeadLen = 1 + 1

type Message struct {
	DataLen uint8  // 消息的长度 （1 字节）
	ProtoID uint8  // 协议ID （1 字节）
	Data    []byte // 消息的内容
}

const (
	Registyer uint8 = iota
	// GetNode      // todo 删除此协议，只有合法用户才可以获取；在登录到目标节点时，还要使用密码，一次使用，终身有效。
	// 不使用token，因为没必要，只使用一次，不是每次请求都使用
	GetNode // login
	NodeInfo
	Error
)
const (
	ErrDataUnKnow uint32 = iota
)
