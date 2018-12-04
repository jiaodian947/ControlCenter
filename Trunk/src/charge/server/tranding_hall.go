package server

import (
	"fmt"
	"time"
)

// 单笔交易
type Trader interface {
	Platform() string
	Check() bool
	Process(chan Trader) error
	Complete(*Context)
}

// 交易大厅，按顺序处理订单
type TradingHall struct {
	ctx      *Context
	pool     map[string]chan Trader
	complete map[string]chan Trader
	size     int
	quit     bool
}

func NewTradingHall(ctx *Context, chan_len int) *TradingHall {
	hall := &TradingHall{}
	hall.ctx = ctx
	hall.pool = make(map[string]chan Trader)
	hall.complete = make(map[string]chan Trader)
	hall.size = chan_len //每个平台最大排队长度
	return hall
}

// 按平台创建处理队列
func (hall *TradingHall) CreatePlatform(name string) {
	if _, dup := hall.pool[name]; dup {
		panic("platform pool dup")
	}
	if _, dup := hall.complete[name]; dup {
		panic("platform complete dup")
	}

	hall.pool[name] = make(chan Trader, hall.size)
	hall.complete[name] = make(chan Trader, hall.size)
}

// 增加一笔交易
func (hall *TradingHall) AddTrader(v Trader) error {
	if p, find := hall.pool[v.Platform()]; find {
		p <- v
		return nil
	}

	return fmt.Errorf("platform not found")
}

// 按平台启动多个处理队列
func (hall *TradingHall) StartAll() {
	for k, v := range hall.pool {
		queuech := v
		donech := hall.complete[k]
		hall.ctx.Server.waitGroup.Wrap(func() {
			hall.worker(queuech, donech)
		})
		hall.ctx.Server.waitGroup.Wrap(func() {
			hall.done(donech)
		})
	}
}

// 退出
func (hall *TradingHall) Shutdown() {
	hall.quit = true
}

// 每个平台的处理队列
func (hall *TradingHall) worker(queue chan Trader, complete chan Trader) {
	for !hall.quit {
		select {
		case v := <-queue:
			vf := v
			if vf.Check() {
				vf.Process(complete)
			}
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

// 每个平台的完成队列
func (hall *TradingHall) done(ch chan Trader) {
	for !hall.quit {
		select {
		case v := <-ch:
			v.Complete(hall.ctx)
		default:
			time.Sleep(time.Millisecond)
		}
	}
}
