package gen

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math"
	"sync/atomic"
	"time"

	lg "github.com/impact-eintr/WebKits/loadgenerator"
	gt "github.com/impact-eintr/WebKits/loadgenerator/gotickets"
)

// 调用状态：0-未调用或调用中；1-调用完成；2-调用超时。
const (
	NORMAL uint32 = iota
	SUCCEED
	TIMEOUT
)

type generator struct {
	caller lg.Caller // 调用器。

	timeoutNS  time.Duration // 处理超时时间，单位：纳秒。
	lps        uint32        // 每秒载荷量。
	durationNS time.Duration // 负载持续时间，单位：纳秒。

	concurrency uint32       // 载荷并发量。
	tickets     lg.GoTickets // Goroutine票池。

	ctx        context.Context    // 上下文。
	cancelFunc context.CancelFunc // 取消函数。

	callCount int64  // 调用计数。
	status    uint32 // 状态。

	resultCh chan *lg.CallResult // 调用结果通道
}

func NewGenerator(param ParamSet) (lg.Generator, error) {
	log.Println("新建一个载荷发生器")
	if err := param.Check(); err != nil {
		return nil, err
	}

	gen := &generator{
		caller:     param.Caller,
		timeoutNS:  param.TimeoutNS,
		lps:        param.LPS,
		durationNS: param.DurationNS,
		status:     lg.STATUS_ORIGINAL,
		resultCh:   param.ResultCh,
	}
	if err := gen.init(); err != nil {
		return nil, err
	}

	return gen, nil
}

func (g *generator) init() error {
	var buf bytes.Buffer
	buf.WriteString("初始化一个载荷发生器...\n")
	// 载荷的并发量 ≈ 载荷的响应超时时间 / 载荷的发送间隔时间
	var total64 = int64(g.timeoutNS)/int64(1e9/g.lps) + 1
	if total64 > math.MaxInt32 {
		total64 = math.MaxInt32
	}
	g.concurrency = uint32(total64)

	tickets, err := gt.NewGoTickets(g.concurrency)
	if err != nil {
		return err
	}
	g.tickets = tickets

	buf.WriteString(fmt.Sprintf("结束. (并发量=%d)\n", g.concurrency))
	log.Println(buf.String())
	return nil

}

// 状态机: (可启动) -> (启动)
func (g *generator) Start() bool {
	log.Println("开启载荷中...")
	if !atomic.CompareAndSwapUint32(&g.status, lg.STATUS_ORIGINAL, lg.STATUS_STARTED) {
		if !atomic.CompareAndSwapUint32(&g.status, lg.STATUS_STOPPED, lg.STATUS_STARTED) {
			return false
		}
	}
	// 设定节流阀
	var throttle <-chan time.Time
	if g.lps > 0 {
		interval := time.Duration(1e9 / g.lps)
		log.Printf("设置节流阀(%v)\n", interval)
		throttle = time.Tick(interval)
	}

	// 初始化上下文和取消函数
	g.ctx, g.cancelFunc = context.WithTimeout(context.Background(), g.durationNS)

	// 初始化调用计数
	g.callCount = 0

	atomic.StoreUint32(&g.status, lg.STATUS_STARTED)

	go func() {
		// 生成并发送载荷
		log.Println("载荷发生中。。。")
		g.genLoad(throttle)
		log.Printf("载荷结束. (调用计数: %d)\n", g.callCount)
	}()

	return true
}

// 产生载荷并向承受方发送
func (g *generator) genLoad(throttle <-chan time.Time) {
	for {
		select {
		case <-g.ctx.Done():
			g.prepareToStop(g.ctx.Err())
			return
		default:
		}
		g.asyncCall()
		if g.lps > 0 {
			select {
			case <-throttle:
			case <-g.ctx.Done():
				g.prepareToStop(g.ctx.Err())
				return
			}
		}
	}
}

// 异步调用接口
func (g *generator) asyncCall() {
	// 获取一个gouroutine
	g.tickets.Fetch()
	go func() {

		// 构建调用方请求
		rawReq := g.caller.BuildReq()

		var callStatus uint32

		// 处理超时任务
		timer := time.AfterFunc(g.timeoutNS, func() {
			if !atomic.CompareAndSwapUint32(&callStatus, NORMAL, TIMEOUT) {
				return
			}
			result := &lg.CallResult{
				ID:     rawReq.ID,
				Req:    rawReq,
				Code:   lg.WARNING_CALL_TIMEOUT,
				Msg:    fmt.Sprintf("运行超时! (预计时间: < %v)", g.timeoutNS),
				Elapse: g.timeoutNS,
			}
			g.sendResult(result)
		})

		rawResp := g.callOne(&rawReq)
		if !atomic.CompareAndSwapUint32(&callStatus, NORMAL, SUCCEED) {
			return
		}
		timer.Stop() //到时间后执行超时函数

		//  处理非超时任务
		var result *lg.CallResult
		if rawResp.Err != nil {
			result = &lg.CallResult{
				ID:     rawResp.ID,
				Req:    rawReq,
				Code:   lg.ERROR_CALL,
				Msg:    rawResp.Err.Error(),
				Elapse: rawResp.Elapse}
		} else {
			result = g.caller.CheckResp(rawReq, *rawResp)
			result.Elapse = rawResp.Elapse
		}
		g.sendResult(result)

	}()

}

func (g *generator) prepareToStop(err error) {

}

func (g *generator) sendResult(result *lg.CallResult) {

}

func (g *generator) callOne(rawReq *lg.RawReq) *lg.RawResp {

}

func (g *generator) Stop() bool {
	return true
}

func (g *generator) Status() uint32 {

	return 0
}

func (g *generator) CallCount() int64 {
	return 0
}
