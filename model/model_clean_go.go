package model

import (
	"time"

	"github.com/GuanghuiLiu/behavior/utils"
)

const (
	MaxModelCount       uint32 = 10000
	SafeModelCount      uint32 = 1000
	WarningModelCountL1 uint32 = 5000
	WarningModelCountL2 uint32 = 6000
	WarningModelCountL3 uint32 = 7000
)

const CleanDuration time.Duration = 30 * time.Second

type processClean struct {
	tk *time.Ticker
}

func init() {
	c := newProcessClean()
	go c.run()
}

func newProcessClean() *processClean {
	return &processClean{
		tk: time.NewTicker(CleanDuration),
	}
}

func (c *processClean) run() {
	for {
		select {
		case <-c.tk.C:
			mc := modelCount()
			switch {
			case mc < SafeModelCount:
				for _, m := range allModel {
					if m.liveTime.minute > 0 {
						m.liveTime.minute = utils.MinTime(m.liveTime.minute*2, LiveTime)
					}
				}
			case mc > WarningModelCountL1 && mc < WarningModelCountL2:
				for _, m := range allModel {
					if m.liveTime.minute > 0 {
						m.liveTime.minute = utils.MaxTime(m.liveTime.minute/2, LiveTimeWarning)
					}
				}
			case mc > WarningModelCountL2 && mc < WarningModelCountL3:
				for _, m := range allModel {
					if m.liveTime.minute > 0 {
						m.liveTime.minute = utils.MaxTime(m.liveTime.minute/3, LiveTimeWarning)
					}
				}
			case mc > WarningModelCountL3:
				for _, m := range allModel {
					if m.liveTime.minute > 0 {
						m.liveTime.minute = LiveTimeWarning
					}
				}
			}
			c.tk.Reset(CleanDuration)
		}
	}
}
