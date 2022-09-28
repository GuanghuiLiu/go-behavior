package model

type SafeHandel struct {
	Name string
	// DefaultFunc func()
}

func (h *SafeHandel) HandlerStop() Handler {
	return h
}
func (h *SafeHandel) HandlerInfo(msg *Message) Handler {
	return h

}
func (h *SafeHandel) HandlerSync(msg *SyncMessage) ([]byte, Handler) {
	return nil, h
}
func (h *SafeHandel) HandlerAsync(msg *Message) (Handler, *Message) {
	return h, msg
}

func (h *SafeHandel) SendSync(to string, protoID uint64, data []byte) ([]byte, error) {
	return SentSync(protoID, data, h.Name, to, 0, true)
}

func (h *SafeHandel) SendInfo(to string, protoID uint64, data []byte) error {
	return SentInfo(protoID, data, h.Name, to, true)
}

func (h *SafeHandel) QuickSendInfo(to string, protoID uint64, data []byte) error {
	return SentInfo(protoID, data, h.Name, to, false)
}
func (h *SafeHandel) SendStop(to string) (*model, error) {
	return SentStop(h.Name, to)
}
func (h *SafeHandel) Run(data Handler, opts ...opt) error {

	return run(data, h.Name, false, opts...)
}
func (h *SafeHandel) QuickRun(data Handler, opts ...opt) error {

	return run(data, h.Name, true, opts...)
}
