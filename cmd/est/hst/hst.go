package hst

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	lg "github.com/impact-eintr/WebKits/loadgenerator"
	"github.com/impact-eintr/WebKits/loadgenerator/gen"
)

type httpResp struct {
	StatusCode int
	Body       string
}

type HttpST struct {
	Url    string
	Method string
	Header map[string]string
	Body   []byte
}

// 初始化载荷发生器。
var PSet = &gen.ParamSet{}

func init() {
	var url, method, headers, body, filename string
	var testTime, timeOut int64
	var lps uint

	flag.StringVar(&url, "U", "", "http API 地址(仅支持http) 必须添加 http://")
	flag.StringVar(&method, "X", "GET", "http 请求方式")
	flag.StringVar(&headers, "H", "", "http headers")
	flag.StringVar(&body, "d", "", "http body")
	flag.StringVar(&filename, "f", "", "http body from file 输入文件路径 此选项会覆盖 -d")

	flag.Int64Var(&testTime, "t", -1, "测试持续时间 以 S 计时")
	flag.Int64Var(&timeOut, "T", -1, "TimeoutNS 以 MS 计")
	flag.UintVar(&lps, "P", 1, "LPS 默认为1")

	flag.Parse()

	if url == "" {
		log.Fatalln(errors.New("未指定测试API"))
	}
	if testTime == -1 {
		log.Fatalln(errors.New("未指定测试持续时间"))
	}
	if timeOut == -1 {
		log.Fatalln(errors.New("未指定TimeOut"))
	}
	if filename != "" {
		body, err := fileParse(filename)
		if err != nil {
			log.Fatalln(err)
		}
		PSet.Caller = NewHttpST(url, method, headersParse(headers), body)
	} else {
		PSet.Caller = NewHttpST(url, method, headersParse(headers), []byte(body))
	}

	PSet.DurationNS = time.Duration(testTime) * time.Second
	PSet.TimeoutNS = time.Duration(timeOut) * time.Millisecond
	PSet.LPS = uint32(lps)
	PSet.ResultCh = make(chan *lg.CallResult, 50)

}

func fileParse(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	return io.ReadAll(file)
}

func headersParse(raw string) map[string]string {
	headers := make(map[string]string)

	header := strings.Split(raw, ": ")
	headers[header[0]] = header[1]

	return headers
}

func NewHttpST(url string, method string,
	headers map[string]string, body []byte) lg.Caller {
	return &HttpST{
		Url:    url,
		Method: method,
		Header: headers,
		Body:   body,
	}

}

func (h *HttpST) BuildReq() lg.RawReq {
	id := time.Now().UnixNano()

	rawReq := lg.RawReq{
		ID:  id,
		Req: h.Body,
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

	// 添加用户自己加入的headers
	for k, v := range h.Header {
		httpReq.Header.Set(k, v)
	}

	// 发起请求
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Println(err)
	}

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

	var sresp httpResp
	err := json.Unmarshal(rawResp.Resp, &sresp)
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
