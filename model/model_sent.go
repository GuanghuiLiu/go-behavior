package model

import (
	"errors"
	"time"
)

func SentSync(protoID uint64, data []byte, from, to string, chanSize uint8, safe bool) ([]byte, error) {
	if b, err := SentSyncLocal(protoID, data, from, to, chanSize, safe); err == nil {
		return b, nil
	}
	receive := make(chan []byte, chanSize)
	if err := sentOtherNodeSync(protoID, data, from, to, receive); err != nil {
		return nil, err
	}
	tk := time.NewTicker(10 * time.Second)
	select {
	case b := <-receive:
		close(receive)
		return b, nil
	case <-tk.C:
		return nil, errors.New("time out")
	}
}

func SentSyncLocal(protoID uint64, data []byte, from, to string, chanSize uint8, safe bool) ([]byte, error) {
	if from == to {
		return nil, errors.New("can not  sent to self !")
	}
	m, ok := allModel[to]
	if safe {
		m, ok = getModel(to)
	}
	if ok {
		receive := make(chan []byte, chanSize)
		m.acceptor.sync <- &SyncMessage{From: from,
			receiver: receive,
			Data: &Data{
				ProtoID: protoID,
				Data:    data}}

		tk := time.NewTicker(10 * time.Second)
		select {
		case b := <-receive:
			close(receive)
			return b, nil
		case <-tk.C:
			return nil, errors.New("time out")
		}
	}
	return nil, errors.New("no target")
}

func SentInfo(protoID uint64, data []byte, from, to string, safe bool) error {
	if err := doSentInfo(protoID, data, from, to, safe); err != nil {
		if err2 := sentOtherNode(protoID, data, from, to); err2 != nil {
			return err2
		}
		return err
	}
	return nil
}

func doSentInfo(protoID uint64, data []byte, from, to string, safe bool) error {
	if from == to {
		return errors.New("can not  sent to self !")
	}
	m, ok := allModel[to]
	if safe {
		m, ok = getModel(to)
	}
	if ok {
		m.acceptor.info <- &Message{From: from,
			Data: &Data{
				ProtoID: protoID,
				Data:    data}}
		return nil
	}
	return errors.New("no target")
}

func SentStop(from, to string) (*model, error) {
	if from == to {
		return nil, errors.New("can not  sent to self !")
	}
	m, ok := getModel(to)
	if ok {
		stopChan := make(chan struct{})
		m.acceptor.stop <- stopChan
		// todo 互相关闭会造成死锁，去掉同步。后续可能修改，保留同步stop，通过监控进程，收取消息
		// <-stopChan
		close(stopChan)
		return m, nil
	}
	return nil, errors.New("target has stop")
}
