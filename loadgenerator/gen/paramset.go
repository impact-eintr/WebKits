package gen

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	lg "github.com/impact-eintr/WebKits/loadgenerator"
)

type ParamSet struct {
	Caller     lg.Caller           // 调用器。
	TimeoutNS  time.Duration       // 响应超时时间，单位：纳秒。
	LPS        uint32              // 每秒载荷量。
	DurationNS time.Duration       // 负载持续时间，单位：纳秒。
	ResultCh   chan *lg.CallResult // 调用结果通道。
}

// Check 会检查当前值的所有字段的有效性。
// 若存在无效字段则返回值非nil。
func (pset *ParamSet) Check() error {
	var errMsgs []string

	if pset.Caller == nil {
		errMsgs = append(errMsgs, "Invalid caller!")
	}
	if pset.TimeoutNS == 0 {
		errMsgs = append(errMsgs, "Invalid timeoutNS!")
	}
	if pset.LPS == 0 {
		errMsgs = append(errMsgs, "Invalid lps(load per second)!")
	}
	if pset.DurationNS == 0 {
		errMsgs = append(errMsgs, "Invalid durationNS!")
	}
	if pset.ResultCh == nil {
		errMsgs = append(errMsgs, "Invalid result channel!")
	}

	var buf bytes.Buffer
	buf.WriteString("Checking the parameters...")
	if errMsgs != nil {
		errMsg := strings.Join(errMsgs, " ")
		buf.WriteString(fmt.Sprintf("NOT passed! (%s)", errMsg))
		log.Println(buf.String())
		return errors.New(errMsg)
	}
	buf.WriteString(
		fmt.Sprintf("Passed. (timeoutNS=%s, lps=%d, durationNS=%s)",
			pset.TimeoutNS, pset.LPS, pset.DurationNS))
	log.Println(buf.String())
	return nil
}
