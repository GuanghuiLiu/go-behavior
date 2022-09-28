package name_server

import "github.com/GuanghuiLiu/behavior/model"

type userCenter struct {
	model.QuickHandel
}

func startUserCenter(name string) {
	uc := newUserCenter(name)
	uc.Run(uc)
}

func newUserCenter(name string) *userCenter {
	uc := &userCenter{}
	uc.Name = name
	return uc
}

func (uc *userCenter) HandlerInfo(msg *model.Message) model.Handler {
	// todo register and check password
	return uc
}

func (uc *userCenter) registerUser(name, pass string) bool {
	return true
}

func (uc *userCenter) checkUser(name, pass string) (string, error) {
	uid := name
	return uid, nil
}
