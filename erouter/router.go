package erouter

import (
	"net/http"
	"sync"
)

type Param struct {
	key   string
	value string
}

// Params 就是一个 Param 的切片，这样就可以看出来， URL 参数可以设置多个了。
// 它是在 tree 的 GetValue() 方法调用的时候设置的
// 这个切片是有顺序的，第一个设置的参数就是切片的第一个值，所以通过索引获取值是安全的
type Params []Param

func (ps Params) ByName(name string) string {
	for _, p := range ps {
		if p.key == name {
			return p.value
		}
	}
	return ""
}

type paramsKey struct{}

var ParamsKey = paramsKey{}

type Handle func(http.ResponseWriter, *http.Request, Params)

type Router struct {
	trees map[string]*node

	paramsPool sync.Pool
	maxParams  uint16

	SaveMatchedRoutePath bool

	// 这个参数是否自动处理当访问路径最后带的 /，一般为 true 就行。
	// 例如： 当访问 /foo/ 时， 此时没有定义 /foo/ 这个路由，但是定义了
	// /foo 这个路由，就对自动将 /foo/ 重定向到 /foo (GET 请求
	// 是 http 301 重定向，其他方式的请求是 http 307 重定向）。
	RedirectTrailingSlash bool

	// 是否自动修正路径， 如果路由没有找到时，Router 会自动尝试修复。
	// 首先删除多余的路径，像 ../ 或者 // 会被删除。
	// 然后将清理过的路径再不区分大小写查找，如果能够找到对应的路由， 将请求重定向到
	// 这个路由上 ( GET 是 301， 其他是 307 ) 。
	RedirectFixedPath bol

	HandleMethodNotAllowed bool

	// 如果为 true ，会自动回复 OPTIONS 方式的请求。
	// 如果自定义了 OPTIONS 路由，会使用自定义的路由，优先级高于这个自动回复。
	HandleOPTHONS bool

	GlobalOPTHONS http.Handler
	GlobalAllowed string

	NotFound         http.Handler
	MethodNotAllowed http.Handler
	PanicHandler     func(http.ResponseWriter, *http.Request, interface{})
}

func New() *Router {
	return &Router{}
}

// 验证Router是否实现了http.Handler
var _ http.Handler = New()

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.PanicHandler != nil {
		defer r.recv(w, req)
	}
	// 请求路径
	path := req.URL.Path

	// 到基数树中去查找匹配的路由
	if root := r.trees[req.Method]; root != nil {
		if handle, ps, tsr := root.getValue(path, r.getParams); handle != nil {

		}
	}

}

func (r *Router) recv(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(w, req, rcv)
	}
}

// Params 是一个 []Param struct{key,value}
func (r *Router) getParams() *Params {
	ps, _ := r.paramsPool.Get().(*Params)
	*ps = (*ps)[0:0]
	return ps

}
