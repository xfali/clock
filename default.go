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

type Opt func(c *defaultClock)

func NewClock(opts ...Opt) *defaultClock {
	ret := &defaultClock{
		now:          time.Now().UnixNano(),
		stop:         make(chan struct{}),
		interval:     DefaultInterval,
		syncLoopTime: DefaultSyncTimes,
	}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

func (c *defaultClock) Now() time.Time {
	if c.interval == 0 {
		return time.Now()
	}
	return time.Unix(0, atomic.LoadInt64(&c.now))
}

func (c *defaultClock) UnixNano() int64 {
	if c.interval == 0 {
		return time.Now().UnixNano()
	}
	return atomic.LoadInt64(&c.now)
}

func (c *defaultClock) Since(t time.Time) time.Duration {
	return c.Now().Sub(t)
}

func (c *defaultClock) Until(t time.Time) time.Duration {
	return t.Sub(c.Now())
}

func (c *defaultClock) Sync() {
	atomic.StoreInt64(&c.now, time.Now().UnixNano())
}

func (c *defaultClock) Start() {
	if c.interval == 0 {
		return
	}
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

// 设置递增时间的间隔
func OptSetInterval(interval time.Duration) Opt {
	return func(c *defaultClock) {
		if interval < 0 {
			interval = 0
		}
		c.interval = interval
	}
}

// 设置在多少次递增时间间隔后校准时间
func OptSetSyncLoopTimes(syncTimes int) Opt {
	return func(c *defaultClock) {
		if syncTimes < 1 {
			syncTimes = 1
		}
		c.syncLoopTime = syncTimes
	}
}
