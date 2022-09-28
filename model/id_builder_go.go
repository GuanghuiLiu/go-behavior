package model

import (
	"github.com/GuanghuiLiu/behavior/utils"
)

const maxID uint64 = 10000000
const processName = "IDBuilder"

type idBuilder struct {
	QuickHandel
	data map[string]uint64
}

func startIDBuild() {
	i := newIDBuilder(processName)
	i.Run(i, SetLiveTime(-1))
}
func newIDBuilder(name string) *idBuilder {
	i := &idBuilder{}
	i.Name = name
	i.data = make(map[string]uint64)
	return i
}
func (i *idBuilder) HandlerInfo(msg *Message) Handler {
	return i
}
func (i *idBuilder) HandlerSync(msg *SyncMessage) ([]byte, Handler) {
	var key string
	utils.Decode(msg.Data.Data, &key)
	if i.data[key] < maxID {
		i.data[key]++
	} else {
		i.data[key] = 1
	}
	return utils.Encode(i.data[key]), i
}
func (i *idBuilder) HandlerStop() Handler {
	return i
}

func GetId(key string) (uint64, error) {
	b, err := SentSync(0, utils.Encode(key), "", processName, 0, false)
	if err != nil {
		return 0, err
	}
	var id uint64
	utils.Decode(b, &id)
	return id, nil
}
