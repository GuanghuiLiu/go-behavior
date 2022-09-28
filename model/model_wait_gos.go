package model

import (
	"fmt"
	"runtime/debug"
)

func (m *model) wait(isQuick bool) {
	defer func() {
		if err := recover(); err != any(nil) {
			fmt.Println(m.name, "panic:", err, string(debug.Stack()))
			if m.retry != nil && m.retryRange != nil {
				if m.retry.retried < m.retryRange.maxRetry {
					m.retry.retried++
					fmt.Println(m.name, ",restart:", m.retry.retried)
					m.wait(isQuick)
				} else {
					m.stop()
				}
			} else {
				m.stop()
			}
		}
	}()

	// 带WaitFun的进程 是阻塞进程，不计生命周期
	m.switchWait(isQuick)
}
func (m *model) switchWait(isQuick bool) {
	if m.liveTime.minute < 0 {
		m.doWaitPerennially(isQuick)
		return
	}
	if m.retry == nil || m.retryRange == nil {
		if m.defaultFun == nil {
			m.doWaitOnce(isQuick)
			return
		}
		m.doWaitFuncOnce(isQuick)
		return
	}
	if m.defaultFun == nil {
		m.doWait(isQuick)
		return
	}
	m.doWaitFunc(isQuick)
}

func (m *model) doWaitPerennially(isQuick bool) {
	for {
		select {
		case msg := <-m.acceptor.info:
			m.isQuickInfo(isQuick, msg)
		case msg := <-m.acceptor.sync:
			m.isQuickSync(isQuick, msg)
		case stopChan := <-m.acceptor.stop:
			m.stop()
			stopChan <- struct{}{}
			return
		case <-m.retry.tk.C:
			m.retry.retried = 0
			m.retry.tk.Reset(m.retryRange.duration)
		}
	}
}

func (m *model) doWait(isQuick bool) {
	for {
		select {
		case msg := <-m.acceptor.info:
			m.liveTime.tk.Reset(m.liveTime.minute)
			m.isQuickInfo(isQuick, msg)
		case msg := <-m.acceptor.sync:
			m.liveTime.tk.Reset(m.liveTime.minute)
			m.isQuickSync(isQuick, msg)
		case stopChan := <-m.acceptor.stop:
			m.stop()
			stopChan <- struct{}{}
			return
		case <-m.retry.tk.C:
			m.retry.retried = 0
			m.retry.tk.Reset(m.retryRange.duration)
		case <-m.liveTime.tk.C:
			m.stop()
			return
		}
	}
}
func (m *model) doWaitOnce(isQuick bool) {
	for {
		select {
		case msg := <-m.acceptor.info:
			m.liveTime.tk.Reset(m.liveTime.minute)
			m.isQuickInfo(isQuick, msg)
		case msg := <-m.acceptor.sync:
			m.liveTime.tk.Reset(m.liveTime.minute)
			m.isQuickSync(isQuick, msg)
		case stopChan := <-m.acceptor.stop:
			m.stop()
			stopChan <- struct{}{}
			return
		case <-m.liveTime.tk.C:
			m.stop()
			return
		}
	}
}

func (m *model) doWaitFunc(isQuick bool) {
	for {
		select {
		case msg := <-m.acceptor.info:
			m.liveTime.tk.Reset(m.liveTime.minute)
			m.isQuickInfo(isQuick, msg)
		case msg := <-m.acceptor.sync:
			m.liveTime.tk.Reset(m.liveTime.minute)
			m.isQuickSync(isQuick, msg)
		case stopChan := <-m.acceptor.stop:
			m.stop()
			stopChan <- struct{}{}
			return
		case <-m.retry.tk.C:
			m.retry.retried = 0
			m.retry.tk.Reset(m.retryRange.duration)
		// case <-m.liveTime.tk.C:
		// 	m.stop()
		// 	return
		default:
			// m.liveTime.tk.Reset(m.liveTime.minute)
			if m.defaultFun() {
				m.stop()
				return
			}
		}
	}
}
func (m *model) doWaitFuncOnce(isQuick bool) {
	for {
		select {
		case msg := <-m.acceptor.info:
			m.liveTime.tk.Reset(m.liveTime.minute)
			m.isQuickInfo(isQuick, msg)
		case msg := <-m.acceptor.sync:
			m.liveTime.tk.Reset(m.liveTime.minute)
			m.isQuickSync(isQuick, msg)
		case stopChan := <-m.acceptor.stop:
			m.stop()
			stopChan <- struct{}{}
			return
		// case <-m.liveTime.tk.C:
		// 	m.stop()
		// 	return
		default:
			// m.liveTime.tk.Reset(m.liveTime.minute)
			if m.defaultFun() {
				m.stop()
				return
			}
		}
	}
}

func (m *model) isQuickInfo(ok bool, msg *Message) {
	if ok {
		m.data.HandlerInfo(msg)
	} else {
		d := m.data.HandlerInfo(msg)
		m.data = d
	}
}

func (m *model) isQuickSync(ok bool, msg *SyncMessage) {
	if ok {
		data, _ := m.data.HandlerSync(msg)
		msg.receiver <- data
	} else {
		data, mod := m.data.HandlerSync(msg)
		m.data = mod
		msg.receiver <- data
	}
}

// todo 同步消息，还需要chan，可以再使用一个chan，也可以使用监控process（一个接收，一个发送）
// func (m *model) isQuickSync(ok bool, msg *Message) {
//
// 	d, msg := m.data.HandlerSync(msg)
// 	if ok {
//
// 	} else {
// 		m.data = d
// 	}
// }

func (m *model) stop() {
	m.data.HandlerStop()
	deleteModel(m.name)
	if m.retry != nil {
		m.retry.tk.Stop()
	}
	m.liveTime.tk.Stop()
	m.clean()
	m = nil
}

func (m *model) clean() {
	if _, ok := <-m.acceptor.sync; ok {
		close(m.acceptor.sync)
	}
	if _, ok := <-m.acceptor.async; ok {
		close(m.acceptor.async)
	}
	if _, ok := <-m.acceptor.info; ok {
		close(m.acceptor.info)
	}
	if _, ok := <-m.acceptor.stop; ok {
		close(m.acceptor.stop)
	}
}
