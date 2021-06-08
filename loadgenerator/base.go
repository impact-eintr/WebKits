package loadgenerator

import "time"

// 载荷发生器接口
type Generator interface {
	Start() bool
	Stop() bool
	Status() uint32
	CallCount() int64 // 获取调用计数 每次启动会重置该计数
}

// 原生请求
type RawReq struct {
	ID  int64
	Req []byte
}

// 原生响应
type RawResp struct {
	ID     int64
	Resp   []byte
	Err    error
	Elapse time.Duration // 耗时
}

type ResCode int

// 保留1 ～ 1000 给载荷承受方使用
const (
	SUCCESS              RetCode = 0    // 成功。
	WARNING_CALL_TIMEOUT         = 1001 // 调用超时警告。
	ERROR_CALL                   = 2001 // 调用错误。
	ERROR_RESPONSE               = 2002 // 响应内容错误。
	ERROR_CALEE                  = 2003 // 被调用方（被测软件）的内部错误。
	FATAL_CALL                   = 3001 // 调用过程中发生了致命错误！
)

//

//

//
