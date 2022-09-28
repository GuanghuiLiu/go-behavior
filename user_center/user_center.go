package user_center

import "github.com/GuanghuiLiu/behavior/model"

type UserCenter struct {
	model.QuickHandel
}

func newUserCenter(name string) *UserCenter {
	uc := &UserCenter{}
	uc.Name = name
	return uc
}

func (uc *UserCenter) HandlerInfo(msg *model.Message) model.Handler {
	// todo register and check password
	return uc
}
