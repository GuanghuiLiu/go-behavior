package model

type QuickHandel struct {
	Name string
	// DefaultFunc func()
}

func (qh *QuickHandel) HandlerStop() Handler {
	return qh
}
func (qh *QuickHandel) HandlerInfo(msg *Message) Handler {
	return qh

}
func (qh *QuickHandel) HandlerSync(msg *SyncMessage) ([]byte, Handler) {
	return nil, qh
}
func (qh *QuickHandel) HandlerAsync(msg *Message) (Handler, *Message) {
	return qh, msg
}
func (qh *QuickHandel) SendSync(to string, protoID uint64, data []byte) ([]byte, error) {
	return SentSync(protoID, data, qh.Name, to, 0, false)
}

func (qh *QuickHandel) SendInfo(to string, protoID uint64, data []byte) error {
	return SentInfo(protoID, data, qh.Name, to, false)
}

func (qh *QuickHandel) SafeSendInfo(to string, protoID uint64, data []byte) error {
	return SentInfo(protoID, data, qh.Name, to, true)
}

func (qh *QuickHandel) SendStop(to string) (*model, error) {
	return SentStop(qh.Name, to)
}
func (qh *QuickHandel) SafeRun(data Handler, opts ...opt) error {

	return run(data, qh.Name, false, opts...)
}

func (qh *QuickHandel) Run(data Handler, opts ...opt) error {

	return run(data, qh.Name, true, opts...)
}
