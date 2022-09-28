package itf

type IReceiver interface {
	Start()
	SetPort()
	SetCookie()
}
