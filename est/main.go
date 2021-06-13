package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	lg "github.com/impact-eintr/WebKits/loadgenerator"
)

type HttpST struct {
	Url    string
	Method string
	Header map[string]string
}

type httpReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type httpResp struct {
	StatusCode int
	Body       string
}

func (h *HttpST) BuildReq() lg.RawReq {
	id := time.Now().UnixNano()

	s := httpReq{
		Username: "yixingwei",
		Password: "123456",
	}
	b, _ := json.Marshal(s)

	rawReq := lg.RawReq{
		ID:  id,
		Req: b,
	}

	return rawReq

}

func (h *HttpST) Call(req []byte, timeoutNS time.Duration) ([]byte, error) {
	client := &http.Client{}

	httpReq, err := http.NewRequest(h.Method, h.Url, strings.NewReader(string(req)))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(httpReq)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	httpResp := httpResp{
		StatusCode: resp.StatusCode,
		Body:       string(body),
	}

	b, err := json.Marshal(httpResp)

	return b, nil

}

func (h *HttpST) CheckResp(rawReq lg.RawReq, rawResp lg.RawResp) *lg.CallResult {
	var commResult lg.CallResult
	commResult.ID = rawResp.ID
	commResult.Req = rawReq
	commResult.Resp = rawResp

	var sreq httpReq
	err := json.Unmarshal(rawReq.Req, &sreq)
	if err != nil {
		commResult.Code = lg.FATAL_CALL
		commResult.Msg =
			fmt.Sprintf("Incorrectly formatted Req: %s!\n", string(rawReq.Req))
		return &commResult

	}

	var sresp httpResp
	err = json.Unmarshal(rawResp.Resp, &sresp)
	if err != nil {
		commResult.Code = lg.ERROR_RESPONSE
		commResult.Msg =
			fmt.Sprintf("Incorrectly formatted Resp: %s!\n", string(rawResp.Resp))
		return &commResult

	}

	if sresp.StatusCode != http.StatusOK {
		commResult.Code = lg.ERROR_RESPONSE
		commResult.Msg = fmt.Sprintf("Incorrect result: %d!\n", sresp.StatusCode)
		return &commResult

	}

	commResult.Code = lg.SUCCESS
	commResult.Msg = fmt.Sprintf("Success. (%d)", sresp.StatusCode)

	return &commResult

}

func NewHttpST(url string, method string, headers map[string]string) lg.Caller {
	return &HttpST{
		Url:    url,
		Method: method,
		Header: headers,
	}
}

func main() {
	url := "http://127.0.0.1:9426/api/v1/login"
	method := "POST"
	// 初始化载荷发生器。
	pset := ParamSet{
		Caller:     NewHttpST(url, method, nil),
		TimeoutNS:  100 * time.Millisecond,
		LPS:        uint32(1500),
		DurationNS: 1 * time.Minute,
		ResultCh:   make(chan *lg.CallResult, 50),
	}

	log.Printf("Initialize load generator (timeoutNS=%v, lps=%d, durationNS=%v)...",
		pset.TimeoutNS, pset.LPS, pset.DurationNS)
	gen, err := NewGenerator(pset)
	if err != nil {
		t.Fatalf("Load generator initialization failing: %s\n", err)
		t.FailNow()
	}

	// 开始！
	log.Println("Start load generator...")
	gen.Start()

	printDetail := false
	// 显示结果。
	countMap := make(map[lg.RetCode]int)
	for r := range pset.ResultCh {
		countMap[r.Code] = countMap[r.Code] + 1
		if printDetail {
			log.Printf("Result: ID=%d, Code=%d, Msg=%s, Elapse=%v.\n",
				r.ID, r.Code, r.Msg, r.Elapse)
		}
	}

	var total int
	log.Println("RetCode Count:")
	for k, v := range countMap {
		codePlain := lg.GetRetCodePlain(k)
		log.Printf("  Code plain: %s (%d), Count: %d.\n",
			codePlain, k, v)
		total += v
	}

	log.Printf("Total: %d.\n", total)
	successCount := countMap[lg.SUCCESS]
	tps := float64(successCount) / float64(pset.DurationNS/1e9)
	log.Printf("Loads per second: %d; Treatments per second: %f.\n", pset.LPS, tps)

}
