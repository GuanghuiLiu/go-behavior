package model

import "time"

type model struct {
	name       string
	data       Handler
	acceptor   *acceptor
	defaultFun func() bool
	retry      *retry
	retryRange *retryRange
	liveTime   *liveTime
	isGlobal   bool
}
type liveTime struct {
	minute time.Duration
	tk     *time.Ticker
}
type retry struct {
	retried uint8
	tk      *time.Ticker
}
type retryRange struct {
	maxRetry uint8
	duration time.Duration
}

type Handler interface {
	HandlerStop() Handler
	HandlerInfo(*Message) Handler
	HandlerSync(*SyncMessage) ([]byte, Handler)
	// HandlerAsync(*Message) Handler
}

type acceptor struct {
	sync  chan *SyncMessage
	async chan *Message
	info  chan *Message
	stop  chan chan struct{}
}

type SyncMessage struct {
	receiver chan []byte
	From     string
	Data     *Data
}
type Message struct {
	From string
	Data *Data
}
type Data struct {
	ProtoID uint64
	Data    []byte
}

var registerChan chan *registerInfo

type registerInfo struct {
	isSet     bool
	name      string
	model     *model
	modelChan chan *model
	errChan   chan string
}

func modelCount() uint32 {
	return uint32(len(allModel))
}
