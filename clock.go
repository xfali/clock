// Copyright (C) 2019-2022, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package clock

import (
	"time"
)

type Clock interface {
	// 获得当前时间
	Now() time.Time

	// 返回从UTC 1970年1月1日至今经过的纳秒时间
	UnixNano() int64

	// 强制同步时钟
	Sync()

	// 开启时钟
	Start()

	// 停止时钟
	Stop()
}
