package loadgenerator

import "time"

// 载荷发生器接口
type Generator interface {
	Start() bool
	Stop() bool
	Status() uint32
	CallCount() int64 // 获取调用计数 每次启动会重置该计数
}

// 载荷发生器的状态
const (
	// STATUS_ORIGINAL 代表原始。
	STATUS_ORIGINAL uint32 = 0
	// STATUS_STARTING 代表正在启动。
	STATUS_STARTING uint32 = 1
	// STATUS_STARTED 代表已启动。
	STATUS_STARTED uint32 = 2
	// STATUS_STOPPING 代表正在停止。
	STATUS_STOPPING uint32 = 3
	// STATUS_STOPPED 代表已停止。
	STATUS_STOPPED uint32 = 4
)

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

type RetCode int

// 保留1 ～ 1000 给载荷承受方使用
const (
	SUCCESS              RetCode = 0    // 成功。
	WARNING_CALL_TIMEOUT         = 1001 // 调用超时警告。
	ERROR_CALL                   = 2001 // 调用错误。
	ERROR_RESPONSE               = 2002 // 响应内容错误。
	ERROR_CALEE                  = 2003 // 被调用方（被测软件）的内部错误。
	FATAL_CALL                   = 3001 // 调用过程中发生了致命错误！
)

var (
	CodeExplain = map[RetCode]string{
		SUCCESS:              "Success",
		WARNING_CALL_TIMEOUT: "Call Timeout Warning",
		ERROR_CALL:           "Call Error",
		ERROR_RESPONSE:       "Response Error",
		ERROR_CALEE:          "Callee Error",
		FATAL_CALL:           "Call Fatal Error",
	}
)

// GetRetCodePlain 会依据结果代码返回相应的文字解释。
func GetRetCodePlain(code RetCode) string {
	_, ok := CodeExplain[code]
	if !ok {
		return "Unknown result code"
	} else {
		return CodeExplain[code]
	}
}

// CallResult 表示调用结果的结构
type CallResult struct {
	ID     int64         // ID。
	Req    RawReq        // 原生请求。
	Resp   RawResp       // 原生响应。
	Code   RetCode       // 响应代码。
	Msg    string        // 结果成因的简述。
	Elapse time.Duration // 耗时。
}
