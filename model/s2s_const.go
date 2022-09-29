package model

// conn
const (
	IPVersion string = "tcp4"
)
const RedisAddr = "127.0.0.1:6379"

// default value
const (
	NullString string = ""
	Zero       uint64 = 0
	ZeroUint8  uint8  = 0
	One        uint64 = 1
)

// error code
const (
	ResultOK uint32 = iota
	ErrSocketErr
	ErrUnCodeErr
	ErrParamErr
)
