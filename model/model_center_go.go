package model

import (
	"errors"
	"fmt"
)

var allModel map[string]*model

func init() {
	allModel = make(map[string]*model)
	registerChan = make(chan *registerInfo, ChanMaxLen)
	go modelCenter()
	startIDBuild()
}

func OnThisNode(name string) (bool, Handler) {
	if m, ok := allModel[name]; ok {
		return true, m.data
	}
	return false, nil
}

// to comparison sync.map
func modelCenter() {
	for {
		select {
		case info := <-registerChan:
			if info.isSet {
				// deleteModel model
				if info.model == nil {
					delete(allModel, info.name)
					ClusterCenter.sentDelete(info.name)
					continue
				}
				// register model
				if _, ok := allModel[info.name]; ok {
					info.errChan <- "repeated name"
					continue
				}
				allModel[info.name] = info.model
				info.errChan <- ""
				if info.model.isGlobal {
					ClusterCenter.sentRegister(info.name)
				}
				continue
			} else {
				// GetNodeByModel model
				if m, ok := allModel[info.name]; ok {
					info.modelChan <- m
					continue
				}
				info.errChan <- "no this mod"
				continue
			}
		}
	}
}

// register model
func registerModel(m *model) error {
	errChan := make(chan string)
	registerChan <- &registerInfo{
		isSet:   true,
		name:    m.name,
		model:   m,
		errChan: errChan,
	}
	select {
	case e := <-errChan:
		close(errChan)
		if e == "" {
			return nil
		}
		return errors.New(e)
	}
	return nil
}

func quickGetModel(name string) (*model, bool) {
	m, ok := allModel[name]
	return m, ok
}

func getModel(name string) (*model, bool) {
	modeChan := make(chan *model)
	errChan := make(chan string)
	registerChan <- &registerInfo{
		isSet:     false,
		name:      name,
		modelChan: modeChan,
		errChan:   errChan,
	}
	for {
		select {
		case err := <-errChan:
			close(errChan)
			if err != "" {
				fmt.Println("when GetNodeByModel", name, ":", err)
				return nil, false
			}
		case m := <-modeChan:
			close(modeChan)
			if m != nil {
				return m, true
			}
		}
	}

}

func deleteModel(name string) error {
	registerChan <- &registerInfo{
		isSet: true,
		name:  name,
		model: nil,
	}
	return nil
}
