package model

import (
	s2s "github.com/GuanghuiLiu/behavior/tcp/long/s2s"
	"time"
)

const ChanMaxLen uint8 = 20

// for retry
const (
	RetryCount   uint8         = 3
	RetrySeconds time.Duration = 60 * time.Minute
)

// LiveTime duration
const (
	LiveTime        time.Duration = 20 * time.Minute
	LiveTimeWarning time.Duration = 1 * time.Minute
)

// IsGlobal modify
const IsGlobal = false

type opt func(*model)

func newModel(date Handler, name string, opts ...opt) *model {
	m := &model{
		name: name,
		data: date,
		acceptor: &acceptor{
			sync: make(chan *SyncMessage, ChanMaxLen),
			// async: make(chan *Message, ChanMaxLen),
			info: make(chan *Message, ChanMaxLen),
			stop: make(chan chan struct{}, ChanMaxLen),
		},
		isGlobal: IsGlobal,
	}

	r := retryRange{
		maxRetry: RetryCount,
		duration: RetrySeconds,
	}
	m.retryRange = &r
	tk := time.NewTicker(r.duration)
	m.retry = &retry{
		retried: 0,
		tk:      tk,
	}
	m.liveTime = &liveTime{
		minute: LiveTime,
		tk:     time.NewTicker(LiveTime),
	}
	for _, f := range opts {
		f(m)
	}
	if m.isGlobal {
		node, _ := ClusterCenter.getNodeByModel(m.name)
		if node == localName || node == NullString {
			return m
		}
		sentOtherNode(uint64(s2s.StopProcess), nil, "self", m.name)
	}
	return m
}

func SetRetry(count uint8, seconds uint32) opt {
	if seconds == 0 {
		return func(m *model) {
			m.retryRange = nil
			m.retry = nil
		}
	}
	t := time.Duration(seconds) * time.Second
	r := retryRange{
		maxRetry: count,
		duration: t,
	}
	return func(m *model) {
		m.retryRange = &r
		tk := time.NewTicker(r.duration)
		m.retry = &retry{
			retried: 0,
			tk:      tk,
		}
	}
}

func SetLiveTime(minutes int) opt {
	t := time.Duration(minutes) * time.Minute
	return func(m *model) {
		l := liveTime{
			minute: t,
		}
		if t > 0 {
			tk := time.NewTicker(l.minute)
			l.tk = tk
		}
		m.liveTime = &l
	}
}

func SetIsGlobal(isGlobal bool) opt {
	return func(m *model) {
		m.isGlobal = isGlobal
	}
}

func SetDefaultFunc(defaultFun func() bool) opt {
	return func(m *model) {
		m.defaultFun = defaultFun
	}
}
