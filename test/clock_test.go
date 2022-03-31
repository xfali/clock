// Copyright (C) 2019-2022, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"clock"
	"testing"
	"time"
)

func TestClock(t *testing.T) {
	c := clock.NewClock()
	c.Start()
	defer c.Stop()

	now := c.Now()
	inval := 50 * time.Millisecond
	ti := time.NewTicker(inval)
	tt := int(10 * time.Second / inval)
	for i := 0; i < tt; i++ {
		select {
		case <-ti.C:
			t.Log(time.Since(now)/time.Second, " second")
			v := time.Since(c.Now()) / time.Millisecond
			t.Log("less ", v)
			if v > clock.DefaultInterval+5 {
				t.Fatal("Cannot be here")
			}
		}
	}
}
