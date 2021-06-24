package erouter

import (
	"net/http"
	"sync"
)

type Param struct {
	Key   string
	Value string
}

type Params []Param

type Router struct {
	trees map[string]*node

	paramsPool sync.Pool
	maxParams  uint16
	// 这个参数是否自动处理当访问路径最后带的 /，一般为 true 就行。
	// 例如： 当访问 /foo/ 时， 此时没有定义 /foo/ 这个路由，但是定义了
	// /foo 这个路由，就对自动将 /foo/ 重定向到 /foo (GET 请求
	// 是 http 301 重定向，其他方式的请求是 http 307 重定向）。
	RedirectTrailingSlash bool

	// 是否自动修正路径， 如果路由没有找到时，Router 会自动尝试修复。
	// 首先删除多余的路径，像 ../ 或者 // 会被删除。
	// 然后将清理过的路径再不区分大小写查找，如果能够找到对应的路由， 将请求重定向到
	// 这个路由上 ( GET 是 301， 其他是 307 ) 。
	RedirectFixedPath bool

	HandleMethodNotAllowed bool

	HandleOPTIONS bool

	GlobalOPTIONS http.Handler

	NotFound http.Handler

	MethodNotAllowed http.Handler

	PanicHandler func(http.ResponseWriter, *http.Request, interface{})
}

var _ http.Handler = New()

func New() *Router {
	return &Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.PanicHandler != nil {
		defer r.recv(w, req)
	}

	// 获取当前请求路径
	path := req.URL.Path

	if root := r.trees[req.Method]; root != nil {
		if handle, ps, tsr := root.getVal(path, r.getParams); handle != nil {

		}
	}

}

func (r *Router) recv(w http.ResponseWriter, req *http.Request) {

}

func (r *Router) getParams() *Params {
	ps, _ := r.paramsPool.Get().(*Params)
	*ps = (*ps)[0:0]
	return ps
}
