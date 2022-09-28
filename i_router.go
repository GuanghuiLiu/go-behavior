package behavior

import (
	"github.com/GuanghuiLiu/behavior/model"
)

type IRouter interface {
	Router(*model.Message) error
	SetRouter(*Actor) error
}
type CreateRouter func() IRouter
