# clock
clock是一个快速获得当前时间的时钟，相较于内置time.Now，Clock使用缓存时间替代查询系统时间减少CPU占用。

## 1.安装
```
go get github.com/xfali/clock
```

## 2.使用
```
c := clock.NewClock()
c.Start()
defer c.Stop()

now := c.Now()
```

## 3.注意事项
由于clock.Now()默认返回基于时间轮询缓存的时间，有较小几率出现一个轮询时间间隔的偏差。

## 4.配置
通过如下方法设置轮询时间间隔：

（间隔越短则时间偏差越小）
```
clock.NewClock(clock.OptSetInterval(interval))
```
通过如下方法设置在多少次轮询时间间隔后校准时间：

（该值影响轮询缓存时间的精度，校准时间越短越精准，但会增加CPU占用)
```
clock.NewClock(clock.OptSetSyncLoopTimes(times))
```