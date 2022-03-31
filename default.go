// Copyright (C) 2019-2022, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package clock

import (
	"sync/atomic"
	"time"
)

const (
	DefaultInterval  = 50 * time.Millisecond
	DefaultSyncTimes = 10
)

type defaultClock struct {
	now  int64
	stop chan struct{}

	// 时间间隔
	interval time.Duration
	// 在syncLoopTime次interval间隔后重新调整时间
	syncLoopTime int
}

func NewClock() *defaultClock {
	ret := &defaultClock{
		now:          time.Now().UnixNano(),
		stop:         make(chan struct{}),
		interval:     DefaultInterval,
		syncLoopTime: DefaultSyncTimes,
	}
	return ret
}

func (c *defaultClock) Now() time.Time {
	return time.Unix(0, atomic.LoadInt64(&c.now))
}

func (c *defaultClock) UnixNano() int64 {
	return atomic.LoadInt64(&c.now)
}

func (c *defaultClock) Sync() {
	atomic.StoreInt64(&c.now, time.Now().UnixNano())
}

func (c *defaultClock) Start() {
	go func() {
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()
		i := c.syncLoopTime
		for {
			select {
			case <-c.stop:
				return
			case <-ticker.C:
				i--
				if i == 0 {
					c.Sync()
					i = c.syncLoopTime
				} else {
					atomic.AddInt64(&c.now, int64(c.interval))
				}
			}
		}
	}()
}

func (c *defaultClock) Stop() {
	close(c.stop)
}
