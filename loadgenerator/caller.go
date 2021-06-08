package loadgenerator

import "time"

type Caller interface {
	BuildReq() RawReq
	Call(req []byte, timeoutNS time.Duration) ([]byte, error)
	CheckResp(rawReq RawResp, rawResp RawResp) *CallResult
}
